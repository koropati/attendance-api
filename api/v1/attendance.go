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
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type AttendanceHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
	ClockIn(c *gin.Context)
	ClockOut(c *gin.Context)
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

func (h *attendanceHandler) Create(c *gin.Context) {
	var data model.Attendance
	c.BindJSON(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	data.GormCustom.CreatedBy = currentUserID
	if !h.middleware.IsSuperAdmin(c) {
		data.UserID = currentUserID
	}

	if err := validation.Validate(data.LatitudeIn, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude_in: %v", err))
		return
	}

	if err := validation.Validate(data.LongitudeIn, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude_in: %v", err))
		return
	}

	if err := validation.Validate(data.LatitudeOut, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude_out: %v", err))
		return
	}

	if err := validation.Validate(data.LongitudeOut, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude_out: %v", err))
		return
	}

	result, err := h.attendanceService.CreateAttendance(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

func (h *attendanceHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	var result *model.Attendance
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
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

func (h *attendanceHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
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

	if err := validation.Validate(data.LatitudeIn, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude_in: %v", err))
		return
	}

	if err := validation.Validate(data.LongitudeIn, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude_in: %v", err))
		return
	}

	if err := validation.Validate(data.LatitudeOut, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude_out: %v", err))
		return
	}

	if err := validation.Validate(data.LongitudeOut, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("longitude_out: %v", err))
		return
	}

	var result *model.Attendance
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.attendanceService.UpdateAttendance(id, &data)
	} else {
		result, err = h.attendanceService.UpdateAttendanceByUserID(id, currentUserID, &data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

func (h *attendanceHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
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

	response.New(c).Write(http.StatusOK, "success delete data")
}

func (h *attendanceHandler) List(c *gin.Context) {
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

	dataList, err := h.attendanceService.ListAttendance(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.attendanceService.ListAttendanceMeta(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

func (h *attendanceHandler) DropDown(c *gin.Context) {
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

	dataList, err := h.attendanceService.DropDownAttendance(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}

func (h *attendanceHandler) ClockIn(c *gin.Context) {

	var dataClockIn model.CheckInData
	c.BindJSON(&dataClockIn)

	currentCheckIn := presence.GetCurrentMillis()
	toDay := time.Now()

	if err := validation.Validate(dataClockIn.Latitude, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude: %v", err))
		return
	}

	if err := validation.Validate(dataClockIn.Longitude, validation.Required, is.Float); err != nil {
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
		err = errors.New("sorry only role user can clock-in")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// check data schedule dari scan
	schedule, err := h.scheduleService.RetrieveScheduleByQRcode(dataClockIn.QRCode)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check In Radius
	if inRadius := schedule.InRange(dataClockIn.Latitude, dataClockIn.Longitude); !inRadius {
		err = errors.New("sorry you are out of radius")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check daily Schedule
	isExistDailySchedule, dailyScheduleID, err := h.dailyScheduleService.CheckHaveDailySchedule(int(schedule.ID), converter.GetDayName(toDay))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isExistDailySchedule {
		err = errors.New("attendance is not possible on this day")
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
		err = errors.New("the user is not on this schedule")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check Attendance is Exist or not
	isExistAttendance := h.attendanceService.CheckIsExistByDate(currentUserID, int(schedule.ID), toDay.Format("2006-01-02"))
	if isExistAttendance {
		// Get
		attendance, err := h.attendanceService.RetrieveAttendanceByDate(currentUserID, int(schedule.ID), toDay.Format("2006-01-02"))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		attendanceNew := attendance
		attendanceNew.ClockIn = currentCheckIn
		attendanceNew.LateIn = calculation.CalculateLateDuration(dailySchedule.StartTime, currentCheckIn, dataClockIn.TimeZone, schedule.LateDuration)
		attendanceNew.StatusPresence = attendanceNew.GenerateStatusPresence()
		attendanceNew.Status = attendanceNew.GenerateStatus()

		if attendance.ClockIn <= 0 {
			// Update attendance
			attendance, err = h.attendanceService.UpdateAttendance(int(attendance.ID), attendanceNew)
			if err != nil {
				response.New(c).Error(http.StatusBadRequest, err)
				return
			}
		}

		// Add Log
		h.attendanceLogService.CreateAttendanceLog(&model.AttendanceLog{
			AttendanceID: attendance.ID,
			LogType:      "clock_in",
			CheckIn:      currentCheckIn,
			Status:       attendanceNew.GenerateStatus(),
			Latitude:     dataClockIn.Latitude,
			Longitude:    dataClockIn.Longitude,
			TimeZone:     dataClockIn.TimeZone,
			Location:     dataClockIn.Location,
		})

		response.New(c).Data(http.StatusCreated, "success clock in", attendance)

	} else {

		newAttendance := model.Attendance{
			GormCustom: model.GormCustom{
				CreatedBy: currentUserID,
				CreatedAt: toDay,
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
		attendance, err := h.attendanceService.CreateAttendance(&newAttendance)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		// Add Log
		h.attendanceLogService.CreateAttendanceLog(&model.AttendanceLog{
			AttendanceID: attendance.ID,
			LogType:      "clock_in",
			CheckIn:      attendance.ClockIn,
			Status:       attendance.Status,
			Latitude:     attendance.LatitudeIn,
			Longitude:    attendance.LongitudeIn,
			TimeZone:     attendance.TimeZoneIn,
			Location:     attendance.LocationIn,
		})

		response.New(c).Data(http.StatusCreated, "success clock in", attendance)
	}

}

func (h *attendanceHandler) ClockOut(c *gin.Context) {

	var dataClockOut model.CheckInData
	c.BindJSON(&dataClockOut)

	currentCheckIn := presence.GetCurrentMillis()
	toDay := time.Now()

	if err := validation.Validate(dataClockOut.Latitude, validation.Required, is.Float); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("latitude: %v", err))
		return
	}

	if err := validation.Validate(dataClockOut.Longitude, validation.Required, is.Float); err != nil {
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
		err = errors.New("sorry only role user can clock-out")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// check data schedule dari scan
	schedule, err := h.scheduleService.RetrieveScheduleByQRcode(dataClockOut.QRCode)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check In Radius
	if inRadius := schedule.InRange(dataClockOut.Latitude, dataClockOut.Longitude); !inRadius {
		err = errors.New("sorry you are out of radius")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check daily Schedule
	isExistDailySchedule, dailyScheduleID, err := h.dailyScheduleService.CheckHaveDailySchedule(int(schedule.ID), converter.GetDayName(toDay))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !isExistDailySchedule {
		err = errors.New("attendance is not possible on this day")
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
		err = errors.New("the user is not on this schedule")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Check Attendance is Exist or not
	isExistAttendance := h.attendanceService.CheckIsExistByDate(currentUserID, int(schedule.ID), toDay.Format("2006-01-02"))
	if isExistAttendance {
		// Get
		attendance, err := h.attendanceService.RetrieveAttendanceByDate(currentUserID, int(schedule.ID), toDay.Format("2006-01-02"))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		attendanceNew := attendance
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
		h.attendanceLogService.CreateAttendanceLog(&model.AttendanceLog{
			AttendanceID: attendance.ID,
			LogType:      "clock_out",
			CheckIn:      currentCheckIn,
			Status:       attendanceNew.GenerateStatus(),
			Latitude:     dataClockOut.Latitude,
			Longitude:    dataClockOut.Longitude,
			TimeZone:     dataClockOut.TimeZone,
			Location:     dataClockOut.Location,
		})

		response.New(c).Data(http.StatusCreated, "success clock out", attendance)

	} else {

		newAttendance := model.Attendance{
			GormCustom: model.GormCustom{
				CreatedBy: currentUserID,
				CreatedAt: toDay,
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
		attendance, err := h.attendanceService.CreateAttendance(&newAttendance)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		// Add Log
		h.attendanceLogService.CreateAttendanceLog(&model.AttendanceLog{
			AttendanceID: attendance.ID,
			LogType:      "clock_out",
			CheckIn:      attendance.ClockOut,
			Status:       attendance.Status,
			Latitude:     attendance.LatitudeOut,
			Longitude:    attendance.LongitudeOut,
			TimeZone:     attendance.TimeZoneOut,
			Location:     attendance.LocationOut,
		})

		response.New(c).Data(http.StatusCreated, "success clock out", attendance)
	}

}
