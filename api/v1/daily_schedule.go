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

type DailyScheduleHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type dailyScheduleHandler struct {
	dailyScheduleService service.DailyScheduleService
	infra                infra.Infra
	middleware           middleware.Middleware
}

func NewDailyScheduleHandler(dailyScheduleService service.DailyScheduleService, infra infra.Infra, middleware middleware.Middleware) DailyScheduleHandler {
	return &dailyScheduleHandler{
		dailyScheduleService: dailyScheduleService,
		infra:                infra,
		middleware:           middleware,
	}
}

func (h *dailyScheduleHandler) Create(c *gin.Context) {
	var data model.DailySchedule
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

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	result, err := h.dailyScheduleService.CreateDailySchedule(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

func (h *dailyScheduleHandler) Retrieve(c *gin.Context) {
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

	var result *model.DailySchedule
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.dailyScheduleService.RetrieveDailySchedule(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.dailyScheduleService.RetrieveDailyScheduleByOwner(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

func (h *dailyScheduleHandler) Update(c *gin.Context) {
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

	var data model.DailySchedule
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	var result *model.DailySchedule
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.dailyScheduleService.UpdateDailySchedule(id, &data)
	} else {
		result, err = h.dailyScheduleService.UpdateDailyScheduleByOwner(id, currentUserID, &data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

func (h *dailyScheduleHandler) Delete(c *gin.Context) {
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
		if err := h.dailyScheduleService.DeleteDailySchedule(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.dailyScheduleService.DeleteDailyScheduleByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "success delete data")
}

func (h *dailyScheduleHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.DailySchedule
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.dailyScheduleService.ListDailySchedule(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.dailyScheduleService.ListDailyScheduleMeta(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

func (h *dailyScheduleHandler) DropDown(c *gin.Context) {
	var data model.DailySchedule
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.dailyScheduleService.DropDownDailySchedule(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
