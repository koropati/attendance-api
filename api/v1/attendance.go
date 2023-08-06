package v1

import (
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/response"
	"attendance-api/common/util/calculation"
	"attendance-api/common/util/converter"
	"attendance-api/common/util/pagination"
	"attendance-api/common/util/presence"
	"attendance-api/infra"
	"attendance-api/model"
	"attendance-api/service"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type AttendanceHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	UpdateStatusPresence(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
	ClockIn(c *gin.Context)
	ClockOut(c *gin.Context)
	Summary(s *gin.Context)
	AutoGenerate(s *gin.Context)
}

type attendanceHandler struct {
	attendanceService    service.AttendanceService
	attendanceLogService service.AttendanceLogService
	scheduleService      service.ScheduleService
	userScheduleService  service.UserScheduleService
	dailyScheduleService service.DailyScheduleService
	infra                infra.Infra
	middleware           middleware.Middleware
}

func NewAttendanceHandler(
	attendanceService service.AttendanceService,
	attendanceLogService service.AttendanceLogService,
	scheduleService service.ScheduleService,
	userScheduleService service.UserScheduleService,
	dailyScheduleService service.DailyScheduleService,
	infra infra.Infra,
	middleware middleware.Middleware) AttendanceHandler {
	return &attendanceHandler{
		attendanceService:    attendanceService,
		attendanceLogService: attendanceLogService,
		scheduleService:      scheduleService,
		userScheduleService:  userScheduleService,
		dailyScheduleService: dailyScheduleService,
		infra:                infra,
		middleware:           middleware,
	}
}

// Create ... Create Attendance
// @Summary Create New Attendance
// @Description Create Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Param data body model.AttendanceForm true "data"
// @Success 200 {object} model.AttendanceResponseData
// @Failure 400,500 {object} model.Response
// @Router /attendance/create [post]
// @Security BearerTokenAuth
func (h attendanceHandler) Create(c *gin.Context) {
	var data model.Attendance
	errBind := c.BindJSON(&data)
	if errBind != nil {
		log.Printf("ERROR BIND : %v\n", errBind)
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	data.GormCustom.CreatedBy = currentUserID
	if !h.middleware.IsSuperAdmin(c) {
		data.UserID = currentUserID
	}

	if err := validation.Validate(data.LatitudeIn, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude_in: %v", err))
		return
	}

	if err := validation.Validate(data.LongitudeIn, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude_in: %v", err))
		return
	}

	if err := validation.Validate(data.LatitudeOut, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude_out: %v", err))
		return
	}

	if err := validation.Validate(data.LongitudeOut, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude_out: %v", err))
		return
	}

	schedule, err := h.scheduleService.RetrieveSchedule(int(data.ScheduleID))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	//check tanggal dalam range aturan jadwal?
	isDateInRange, err := presence.IsDateInRange(data.Date, schedule.StartDate, schedule.EndDate)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isDateInRange {
		err = errors.New("tidak bisa membuat absensi pada tanggal tersebut")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check In Radius
	if inRadius := schedule.InRange(data.LatitudeIn, data.LongitudeIn); !inRadius {
		err = errors.New("data jam masuk berada di luar radius")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check In Radius
	if inRadius := schedule.InRange(data.LatitudeOut, data.LongitudeOut); !inRadius {
		err = errors.New("data jam keluar berada di luar radius")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check daily Schedule
	isExistDailySchedule, dailyScheduleID, err := h.dailyScheduleService.CheckHaveDailySchedule(int(schedule.ID), converter.GetDayNameFromDateString(data.Date))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isExistDailySchedule {
		err = errors.New("absensi tidak bisa dilakukan pada tanggal tersebut")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	dailySchedule, err := h.dailyScheduleService.RetrieveDailySchedule(dailyScheduleID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	data.LateIn = calculation.CalculateLateDuration(dailySchedule.StartTime, data.ClockIn, data.TimeZoneIn, schedule.LateDuration)
	data.EarlyOut = calculation.CalculateEarlyDuration(dailySchedule.EndTime, data.ClockOut, data.TimeZoneOut)
	data.StatusPresence = data.GenerateStatusPresence()
	data.Status = data.GenerateStatus()

	result, err := h.attendanceService.CreateAttendance(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
}

// Retrieve ... Retrieve Attendance
// @Summary Retrieve Single Attendance
// @Description Retrieve Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceResponseData
// @Failure 400,500 {object} model.Response
// @Router /attendance/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id attendance"
func (h attendanceHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	var result model.Attendance
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.attendanceService.RetrieveAttendance(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.attendanceService.RetrieveAttendanceByUserID(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update Attendance
// @Summary Update Single Attendance
// @Description Update Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceResponseData
// @Failure 400,500 {object} model.Response
// @Router /attendance/update [put]
// @Security BearerTokenAuth
// @param id query string true "id attendance"
func (h attendanceHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	var data model.Attendance
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.LatitudeIn, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude_in: %v", err))
		return
	}

	if err := validation.Validate(data.LongitudeIn, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude_in: %v", err))
		return
	}

	if err := validation.Validate(data.LatitudeOut, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude_out: %v", err))
		return
	}

	if err := validation.Validate(data.LongitudeOut, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude_out: %v", err))
		return
	}

	schedule, err := h.scheduleService.RetrieveSchedule(int(data.ScheduleID))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	//check tanggal dalam range aturan jadwal?
	isDateInRange, err := presence.IsDateInRange(data.Date, schedule.StartDate, schedule.EndDate)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isDateInRange {
		err = errors.New("tidak bisa membuat absensi pada tanggal tersebut")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check In Radius
	if inRadius := schedule.InRange(data.LatitudeIn, data.LongitudeIn); !inRadius {
		err = errors.New("data jam masuk berada di luar radius")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check In Radius
	if inRadius := schedule.InRange(data.LatitudeOut, data.LongitudeOut); !inRadius {
		err = errors.New("data jam keluar berada di luar radius")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check daily Schedule
	isExistDailySchedule, dailyScheduleID, err := h.dailyScheduleService.CheckHaveDailySchedule(int(schedule.ID), converter.GetDayNameFromDateString(data.Date))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isExistDailySchedule {
		err = errors.New("absensi tidak bisa dilakukan pada tanggal tersebut")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	dailySchedule, err := h.dailyScheduleService.RetrieveDailySchedule(dailyScheduleID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	data.LateIn = calculation.CalculateLateDuration(dailySchedule.StartTime, data.ClockIn, data.TimeZoneIn, schedule.LateDuration)
	data.EarlyOut = calculation.CalculateEarlyDuration(dailySchedule.EndTime, data.ClockOut, data.TimeZoneOut)
	data.StatusPresence = data.GenerateStatusPresence()
	data.Status = data.GenerateStatus()

	var result model.Attendance
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.attendanceService.UpdateAttendance(id, data)
	} else {
		result, err = h.attendanceService.UpdateAttendanceByUserID(id, currentUserID, data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// UpdateStatusPresence ... Update Status Attendance
// @Summary Quick Update Status Attendance
// @Description Quick Update Status Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceResponseData
// @Failure 400,500 {object} model.Response
// @Router /attendance/update-status [put]
// @Security BearerTokenAuth
// @param id query string true "id attendance"
func (h attendanceHandler) UpdateStatusPresence(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	var data model.QuickUpdateAttendance
	c.BindJSON(&data)

	var result model.Attendance
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.attendanceService.UpdateStatusAttendance(id, data.StatusPresence, 0)
	} else {
		result, err = h.attendanceService.UpdateStatusAttendance(id, data.StatusPresence, currentUserID)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui status presensi", result)
}

// Delete ... Delete Attendance
// @Summary Delete Single Attendance
// @Description Delete Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /attendance/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id attendance"
func (h attendanceHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if h.middleware.IsSuperAdmin(c) {
		if err := h.attendanceService.DeleteAttendance(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.attendanceService.DeleteAttendanceByUserID(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List All Attendance
// @Summary List All Attendance
// @Description List All Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceResponseList
// @Failure 400,500 {object} model.Response
// @Router /attendance/list [get]
// @Security BearerTokenAuth
func (h attendanceHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Attendance
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.UserID = currentUserID
	}

	dataList, err := h.attendanceService.ListAttendance(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.attendanceService.ListAttendanceMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// Dropdown ... Dropdown All Attendance
// @Summary Dropdown All Attendance
// @Description Dropdown All Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceResponseList
// @Failure 400,500 {object} model.Response
// @Router /attendance/drop-down [get]
// @Security BearerTokenAuth
func (h attendanceHandler) DropDown(c *gin.Context) {
	var data model.Attendance
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.UserID = currentUserID
	}

	dataList, err := h.attendanceService.DropDownAttendance(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}

// Clock In ... Clock In Attendance
// @Summary Clock In Attendance
// @Description Clock In Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.CheckInData
// @Failure 400,500 {object} model.Response
// @Router /attendance/clock-in [post]
// @Security BearerTokenAuth
func (h attendanceHandler) ClockIn(c *gin.Context) {

	var dataClockIn model.CheckInData
	c.BindJSON(&dataClockIn)

	currentCheckIn := presence.GetCurrentMillis()
	toDay := time.Now().Format("2006-01-02")

	if err := validation.Validate(dataClockIn.Latitude, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude: %v", err))
		return
	}

	if err := validation.Validate(dataClockIn.Longitude, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude: %v", err))
		return
	}

	if dataClockIn.TimeZone == 0 {
		dataClockIn.TimeZone = converter.GetTimeZone(dataClockIn.Latitude, dataClockIn.Longitude)
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsUser(c) {
		err = errors.New("maaf hanya role user yang bisa melakukan clock-in")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// check data schedule dari scan
	schedule, err := h.scheduleService.RetrieveScheduleByQRcode(dataClockIn.QRCode)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	//check tanggal dalam range aturan jadwal?
	isDateInRange, err := presence.IsDateInRange(toDay, schedule.StartDate, schedule.EndDate)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isDateInRange {
		err = errors.New("tidak bisa membuat absensi pada tanggal tersebut")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check In Radius
	if inRadius := schedule.InRange(dataClockIn.Latitude, dataClockIn.Longitude); !inRadius {
		err = errors.New("maaf anda berada di luar radius")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check daily Schedule
	isExistDailySchedule, dailyScheduleID, err := h.dailyScheduleService.CheckHaveDailySchedule(int(schedule.ID), converter.GetDayNameFromDateString(toDay))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isExistDailySchedule {
		err = errors.New("absensi tidak bisa dilakukan pada hari ini")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	dailySchedule, err := h.dailyScheduleService.RetrieveDailySchedule(dailyScheduleID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check Employee dalam schedule kah?
	if isValid := h.userScheduleService.CheckUserInSchedule(int(schedule.ID), currentUserID); !isValid {
		err = errors.New("user tersebut tidak berada dalam jadwal ini")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check Attendance is Exist or not
	isExistAttendance := h.attendanceService.CheckIsExistByDate(currentUserID, int(schedule.ID), toDay)
	if isExistAttendance {
		// Get
		attendance, err := h.attendanceService.RetrieveAttendanceByDate(currentUserID, int(schedule.ID), toDay)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		attendanceNew := attendance

		if attendance.GormCustom.CreatedBy == 0 {
			attendanceNew.GormCustom = model.GormCustom{
				CreatedBy: currentUserID,
				UpdatedBy: currentUserID,
				UpdatedAt: time.Now(),
			}
		} else {
			attendanceNew.GormCustom = model.GormCustom{
				UpdatedBy: currentUserID,
				UpdatedAt: time.Now(),
			}
		}

		attendanceNew.ClockIn = currentCheckIn
		attendanceNew.LateIn = calculation.CalculateLateDuration(dailySchedule.StartTime, currentCheckIn, dataClockIn.TimeZone, schedule.LateDuration)

		if attendance.LatitudeIn == 0 {
			attendanceNew.LatitudeIn = dataClockIn.Latitude
		}
		if attendance.LongitudeIn == 0 {
			attendanceNew.LongitudeIn = dataClockIn.Longitude
		}
		if attendance.TimeZoneIn == 0 {
			attendanceNew.TimeZoneIn = dataClockIn.TimeZone
		}
		if attendance.LocationIn == "" {
			attendanceNew.LocationIn = dataClockIn.Location
		}

		attendanceNew.StatusPresence = attendanceNew.GenerateStatusPresence()
		attendanceNew.Status = attendanceNew.GenerateStatus()

		if attendance.ClockIn == 0 {
			// Update attendance
			attendance, err = h.attendanceService.UpdateAttendance(int(attendance.ID), attendanceNew)
			if err != nil {
				response.New(c).Error(http.StatusBadRequest, err)
				return
			}
		}

		// Add Log
		h.attendanceLogService.CreateAttendanceLog(model.AttendanceLog{
			GormCustom: model.GormCustom{
				CreatedBy: currentUserID,
				UpdatedBy: currentUserID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			AttendanceID: attendance.ID,
			LogType:      "clock_in",
			CheckIn:      currentCheckIn,
			Status:       attendanceNew.GenerateStatus(),
			Latitude:     dataClockIn.Latitude,
			Longitude:    dataClockIn.Longitude,
			TimeZone:     dataClockIn.TimeZone,
			Location:     dataClockIn.Location,
		})

		response.New(c).Data(http.StatusCreated, "berhasil absen masuk", attendance)

	} else {

		newAttendance := model.Attendance{
			GormCustom: model.GormCustom{
				CreatedBy: currentUserID,
				UpdatedBy: currentUserID,
				CreatedAt: time.Now(),
			},
			UserID:      currentUserID,
			ScheduleID:  schedule.ID,
			Date:        toDay,
			ClockIn:     currentCheckIn,
			LateIn:      calculation.CalculateLateDuration(dailySchedule.StartTime, currentCheckIn, dataClockIn.TimeZone, schedule.LateDuration),
			LatitudeIn:  dataClockIn.Latitude,
			LongitudeIn: dataClockIn.Longitude,
			TimeZoneIn:  dataClockIn.TimeZone,
			LocationIn:  dataClockIn.Location,
		}
		newAttendance.StatusPresence = newAttendance.GenerateStatusPresence()
		newAttendance.Status = newAttendance.GenerateStatus()
		// Create attendance
		attendance, err := h.attendanceService.CreateAttendance(newAttendance)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		// Add Log
		h.attendanceLogService.CreateAttendanceLog(model.AttendanceLog{
			GormCustom: model.GormCustom{
				CreatedBy: currentUserID,
				UpdatedBy: currentUserID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			AttendanceID: attendance.ID,
			LogType:      "clock_in",
			CheckIn:      attendance.ClockIn,
			Status:       attendance.Status,
			Latitude:     attendance.LatitudeIn,
			Longitude:    attendance.LongitudeIn,
			TimeZone:     attendance.TimeZoneIn,
			Location:     attendance.LocationIn,
		})

		response.New(c).Data(http.StatusCreated, "berhasil absen masuk", attendance)
	}

}

// Clock Out ... Clock Out Attendance
// @Summary Clock Out Attendance
// @Description Clock Out Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.CheckInData
// @Failure 400,500 {object} model.Response
// @Router /attendance/clock-out [post]
// @Security BearerTokenAuth
func (h attendanceHandler) ClockOut(c *gin.Context) {

	var dataClockOut model.CheckInData
	c.BindJSON(&dataClockOut)

	currentCheckIn := presence.GetCurrentMillis()
	toDay := time.Now().Format("2006-01-02")

	if err := validation.Validate(dataClockOut.Latitude, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude: %v", err))
		return
	}

	if err := validation.Validate(dataClockOut.Longitude, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude: %v", err))
		return
	}

	if dataClockOut.TimeZone == 0 {
		dataClockOut.TimeZone = converter.GetTimeZone(dataClockOut.Latitude, dataClockOut.Longitude)
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsUser(c) {
		err = errors.New("maaf hanya role user yang bisa melakukan clock out")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// check data schedule dari scan
	schedule, err := h.scheduleService.RetrieveScheduleByQRcode(dataClockOut.QRCode)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	//check tanggal dalam range aturan jadwal?
	isDateInRange, err := presence.IsDateInRange(toDay, schedule.StartDate, schedule.EndDate)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isDateInRange {
		err = errors.New("tidak bisa membuat absensi pada tanggal tersebut")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check In Radius
	if inRadius := schedule.InRange(dataClockOut.Latitude, dataClockOut.Longitude); !inRadius {
		err = errors.New("maaf anda berada di luar radius")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check daily Schedule
	isExistDailySchedule, dailyScheduleID, err := h.dailyScheduleService.CheckHaveDailySchedule(int(schedule.ID), converter.GetDayNameFromDateString(toDay))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isExistDailySchedule {
		err = errors.New("absensi tidak bisa dilakukan pada hari ini")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	dailySchedule, err := h.dailyScheduleService.RetrieveDailySchedule(dailyScheduleID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check Employee dalam schedule kah?
	if isValid := h.userScheduleService.CheckUserInSchedule(int(schedule.ID), currentUserID); !isValid {
		err = errors.New("user tersebut tidak berada dalam jadwal ini")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check Attendance is Exist or not
	isExistAttendance := h.attendanceService.CheckIsExistByDate(currentUserID, int(schedule.ID), toDay)
	if isExistAttendance {
		// Get
		attendance, err := h.attendanceService.RetrieveAttendanceByDate(currentUserID, int(schedule.ID), toDay)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		attendanceNew := attendance
		if attendance.GormCustom.CreatedBy == 0 {
			attendanceNew.GormCustom = model.GormCustom{
				CreatedBy: currentUserID,
				UpdatedBy: currentUserID,
				UpdatedAt: time.Now(),
			}
		} else {
			attendanceNew.GormCustom = model.GormCustom{
				UpdatedBy: currentUserID,
				UpdatedAt: time.Now(),
			}
		}
		attendanceNew.ClockOut = currentCheckIn
		attendanceNew.EarlyOut = calculation.CalculateEarlyDuration(dailySchedule.EndTime, currentCheckIn, dataClockOut.TimeZone)
		attendanceNew.StatusPresence = attendanceNew.GenerateStatusPresence()
		attendanceNew.Status = attendanceNew.GenerateStatus()
		attendanceNew.LatitudeOut = dataClockOut.Latitude
		attendanceNew.LongitudeOut = dataClockOut.Longitude
		attendanceNew.TimeZoneOut = dataClockOut.TimeZone
		attendanceNew.LocationOut = dataClockOut.Location

		// Update attendance
		attendance, err = h.attendanceService.UpdateAttendance(int(attendance.ID), attendanceNew)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		// Add Log
		h.attendanceLogService.CreateAttendanceLog(model.AttendanceLog{
			GormCustom: model.GormCustom{
				CreatedBy: currentUserID,
				UpdatedBy: currentUserID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			AttendanceID: attendance.ID,
			LogType:      "clock_out",
			CheckIn:      currentCheckIn,
			Status:       attendanceNew.GenerateStatus(),
			Latitude:     dataClockOut.Latitude,
			Longitude:    dataClockOut.Longitude,
			TimeZone:     dataClockOut.TimeZone,
			Location:     dataClockOut.Location,
		})

		response.New(c).Data(http.StatusCreated, "berhasil absen keluar", attendance)

	} else {

		newAttendance := model.Attendance{
			GormCustom: model.GormCustom{
				CreatedBy: currentUserID,
				UpdatedBy: currentUserID,
				CreatedAt: time.Now(),
			},
			UserID:       currentUserID,
			ScheduleID:   schedule.ID,
			Date:         toDay,
			ClockOut:     currentCheckIn,
			EarlyOut:     calculation.CalculateEarlyDuration(dailySchedule.EndTime, currentCheckIn, dataClockOut.TimeZone),
			LatitudeOut:  dataClockOut.Latitude,
			LongitudeOut: dataClockOut.Longitude,
			TimeZoneOut:  dataClockOut.TimeZone,
			LocationOut:  dataClockOut.Location,
		}
		newAttendance.StatusPresence = newAttendance.GenerateStatusPresence()
		newAttendance.Status = newAttendance.GenerateStatus()
		// Create attendance
		attendance, err := h.attendanceService.CreateAttendance(newAttendance)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		// Add Log
		h.attendanceLogService.CreateAttendanceLog(model.AttendanceLog{
			GormCustom: model.GormCustom{
				CreatedBy: currentUserID,
				UpdatedBy: currentUserID,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			AttendanceID: attendance.ID,
			LogType:      "clock_out",
			CheckIn:      attendance.ClockOut,
			Status:       attendance.Status,
			Latitude:     attendance.LatitudeOut,
			Longitude:    attendance.LongitudeOut,
			TimeZone:     attendance.TimeZoneOut,
			Location:     attendance.LocationOut,
		})

		response.New(c).Data(http.StatusCreated, "berhasil absen keluar", attendance)
	}

}

// Summary ... Summary Attendance
// @Summary Summary Attendance
// @Description Summary Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceSummary
// @Failure 400,500 {object} model.Response
// @Router /attendance/summary [get]
// @Security BearerTokenAuth
func (h attendanceHandler) Summary(c *gin.Context) {

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	y, m, _ := time.Now().Date()

	first, last := converter.MonthInterval(y, m)

	presence := h.attendanceService.CountAttendanceByStatus(currentUserID, "presence", first.Format("2006-01-02"), last.Format("2006-01-02"))
	notPresence := h.attendanceService.CountAttendanceByStatus(currentUserID, "not_presence", first.Format("2006-01-02"), last.Format("2006-01-02"))
	sick := h.attendanceService.CountAttendanceByStatus(currentUserID, "sick", first.Format("2006-01-02"), last.Format("2006-01-02"))
	leaveAttendance := h.attendanceService.CountAttendanceByStatus(currentUserID, "leave_attendance", first.Format("2006-01-02"), last.Format("2006-01-02"))
	data := model.AttendanceSummary{
		Presence:        presence,
		NotPresence:     notPresence,
		Sick:            sick,
		LeaveAttendance: leaveAttendance,
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan rangkuman data", data)
}

// AutoGenerate ... Auto Generate Attendance
// @Summary Auto Generate Attendance
// @Description Auto Generate Attendance
// @Tags Attendance
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /attendance/auto-generate [get]
// @Security BearerTokenAuth
func (h attendanceHandler) AutoGenerate(c *gin.Context) {

	if !h.middleware.IsSuperAdmin(c) {
		err := errors.New("maaf hanya super user yang bisa melakukan tindakan ini")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	userSchedules, err := h.userScheduleService.GetAll()
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	wg := sync.WaitGroup{}
	for _, userSchedule := range userSchedules {
		startDate := converter.GetOnlyDateString(userSchedule.Schedule.StartDate)
		// endDate := converter.GetOnlyDateString(userSchedule.Schedule.EndDate)
		endDate := time.Now().Format("2006-01-02")

		listDates, err := converter.GetDatesArrayFromStartEndDate(startDate, endDate)
		if err != nil {
			log.Printf("Error Get List Date E: %v\n", err)
			break
		}
		for j, date := range listDates {
			wg.Add(1)
			go func(j int, date string, userSchedule model.UserSchedule) {
				if !h.attendanceService.CheckIsExistByDate(userSchedule.UserID, int(userSchedule.ScheduleID), date) {
					// Create Attendance with default (alpa)
					// Buat Data Presensi kosong / tidak hadir secara default terlebih dahulu
					dataAttendance := model.Attendance{
						UserID:         userSchedule.UserID,
						ScheduleID:     userSchedule.ScheduleID,
						Date:           date,
						ClockIn:        0,
						ClockOut:       0,
						Status:         "-",
						StatusPresence: "not_presence",
					}
					_, err := h.attendanceService.CreateAttendance(dataAttendance)
					if err != nil {
						log.Printf("[Error] [Attendance-CreateAttendance] E: %v\n", err)
					}
				}
				wg.Done()
			}(j, date, userSchedule)

		}
	}
	wg.Wait()

	response.New(c).Write(http.StatusOK, "sukses melakukan generate data presensi")
}
