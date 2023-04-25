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

type UserScheduleHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
	MySchedule(s *gin.Context)
}

type userScheduleHandler struct {
	userScheduleService service.UserScheduleService
	infra               infra.Infra
	middleware          middleware.Middleware
}

func NewUserScheduleHandler(userScheduleService service.UserScheduleService, infra infra.Infra, middleware middleware.Middleware) UserScheduleHandler {
	return &userScheduleHandler{
		userScheduleService: userScheduleService,
		infra:               infra,
		middleware:          middleware,
	}
}

func (h userScheduleHandler) Create(c *gin.Context) {
	var data model.UserSchedule
	c.BindJSON(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	data.GormCustom.CreatedBy = currentUserID
	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	result, err := h.userScheduleService.CreateUserSchedule(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

func (h userScheduleHandler) Retrieve(c *gin.Context) {
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

	var result model.UserSchedule
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.userScheduleService.RetrieveUserSchedule(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.userScheduleService.RetrieveUserScheduleByOwner(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

func (h userScheduleHandler) Update(c *gin.Context) {
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

	var data model.UserSchedule
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	var result model.UserSchedule
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.userScheduleService.UpdateUserSchedule(id, data)
	} else {
		result, err = h.userScheduleService.UpdateUserScheduleByOwner(id, currentUserID, data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

func (h userScheduleHandler) Delete(c *gin.Context) {
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
		if err := h.userScheduleService.DeleteUserSchedule(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.userScheduleService.DeleteUserScheduleByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "success delete data")
}

func (h userScheduleHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.UserSchedule
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.userScheduleService.ListUserSchedule(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.userScheduleService.ListUserScheduleMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

func (h userScheduleHandler) DropDown(c *gin.Context) {
	var data model.UserSchedule
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.userScheduleService.DropDownUserSchedule(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}

func (h userScheduleHandler) MySchedule(c *gin.Context) {
	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	results, err := h.userScheduleService.ListMySchedule(currentUserID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success list data my schedule", results)
}
