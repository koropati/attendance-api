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

func (h *attendancelogHandler) Create(c *gin.Context) {
	var data model.AttendanceLog
	c.BindJSON(&data)

	result, err := h.attendancelogService.CreateAttendanceLog(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

func (h *attendancelogHandler) Retrieve(c *gin.Context) {
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

func (h *attendancelogHandler) Update(c *gin.Context) {
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

	result, err := h.attendancelogService.UpdateAttendanceLog(id, &data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

func (h *attendancelogHandler) Delete(c *gin.Context) {
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

func (h *attendancelogHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.AttendanceLog
	c.BindQuery(&data)

	dataList, err := h.attendancelogService.ListAttendanceLog(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.attendancelogService.ListAttendanceLogMeta(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

func (h *attendancelogHandler) DropDown(c *gin.Context) {
	var data model.AttendanceLog
	c.BindQuery(&data)

	dataList, err := h.attendancelogService.DropDownAttendanceLog(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
