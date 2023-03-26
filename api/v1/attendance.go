package v1

import (
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/response"
	"attendance-api/common/util/pagination"
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
}

type attendanceHandler struct {
	attendanceService    service.AttendanceService
	attendanceLogService service.AttendanceLogService
	scheduleService      service.ScheduleService
	dailyScheduleService service.DailyScheduleService
	infra                infra.Infra
	middleware           middleware.Middleware
}

func NewAttendanceHandler(
	attendanceService service.AttendanceService,
	attendanceLogService service.AttendanceLogService,
	scheduleService service.ScheduleService,
	dailyScheduleService service.DailyScheduleService,
	infra infra.Infra,
	middleware middleware.Middleware) AttendanceHandler {
	return &attendanceHandler{
		attendanceService:    attendanceService,
		attendanceLogService: attendanceLogService,
		scheduleService:      scheduleService,
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
