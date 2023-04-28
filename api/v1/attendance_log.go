package v1

import (
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/response"
	"attendance-api/common/util/pagination"
	"attendance-api/infra"
	"attendance-api/model"
	"attendance-api/service"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AttendanceLogHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type attendancelogHandler struct {
	attendancelogService service.AttendanceLogService
	infra                infra.Infra
	middleware           middleware.Middleware
}

func NewAttendanceLogHandler(attendancelogService service.AttendanceLogService, infra infra.Infra, middleware middleware.Middleware) AttendanceLogHandler {
	return &attendancelogHandler{
		attendancelogService: attendancelogService,
		infra:                infra,
		middleware:           middleware,
	}
}

// Create ... Create Attendance Log
// @Summary Create New Attendance Log
// @Description Create Attendance Log
// @Tags Attendance Log
// @Accept       json
// @Produce      json
// @Param data body model.AttendanceLogForm true "data"
// @Success 200 {object} model.AttendanceLogResponseData
// @Failure 400,500 {object} model.Response
// @Router /attendance-log/create [post]
// @Security BearerTokenAuth
func (h attendancelogHandler) Create(c *gin.Context) {
	var data model.AttendanceLog
	c.BindJSON(&data)

	result, err := h.attendancelogService.CreateAttendanceLog(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

// Retrieve ... Retreive Attendance Log
// @Summary Retreive Single Attendance Log
// @Description Retreive Single Attendance Log
// @Tags Attendance Log
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceLogResponseData
// @Failure 400,500 {object} model.Response
// @Router /attendance-log/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id attendance log"
func (h attendancelogHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	result, err := h.attendancelogService.RetrieveAttendanceLog(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

// Update ... Update Attendance Log
// @Summary Update Single Attendance Log
// @Description Update Single Attendance Log
// @Tags Attendance Log
// @Accept       json
// @Produce      json
// @Param data body model.AttendanceLogForm true "data"
// @Success 200 {object} model.AttendanceLogResponseData
// @Failure 400,500 {object} model.Response
// @Router /attendance-log/update [put]
// @Security BearerTokenAuth
// @param id query string true "id attendance log"
func (h attendancelogHandler) Update(c *gin.Context) {
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

	var data model.AttendanceLog
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	result, err := h.attendancelogService.UpdateAttendanceLog(id, data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

// Delete ... Delete Attendance Log
// @Summary Delete Single Attendance Log
// @Description Delete Single Attendance Log
// @Tags Attendance Log
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /attendance-log/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id attendance log"
func (h attendancelogHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	if err := h.attendancelogService.DeleteAttendanceLog(id); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(http.StatusOK, "success delete data")
}

// List ... List All Attendance Log
// @Summary List All Attendance Log
// @Description List All Attendance Log
// @Tags Attendance Log
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceLogResponseData
// @Failure 400,500 {object} model.Response
// @Router /attendance-log/list [get]
// @Security BearerTokenAuth
func (h attendancelogHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.AttendanceLog
	c.BindQuery(&data)

	dataList, err := h.attendancelogService.ListAttendanceLog(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.attendancelogService.ListAttendanceLogMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

// Dropdown ... Dropdown All Attendance Log
// @Summary Dropdown All Attendance Log
// @Description Dropdown All Attendance Log
// @Tags Attendance Log
// @Accept       json
// @Produce      json
// @Success 200 {object} model.AttendanceLogResponseList
// @Failure 400,500 {object} model.Response
// @Router /attendance-log/drop-down [get]
// @Security BearerTokenAuth
func (h attendancelogHandler) DropDown(c *gin.Context) {
	var data model.AttendanceLog
	c.BindQuery(&data)

	dataList, err := h.attendancelogService.DropDownAttendanceLog(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
