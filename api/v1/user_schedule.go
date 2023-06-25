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
	ListUserInRule(c *gin.Context)
	ListUserNotInRule(c *gin.Context)
	DropDown(c *gin.Context)
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

// Create ... Create User Schedule
// @Summary Create New User Schedule
// @Description Create User Schedule
// @Tags User Schedule
// @Accept       json
// @Produce      json
// @Param data body model.UserScheduleForm true "data"
// @Success 200 {object} model.UserScheduleResponseData
// @Failure 400,500 {object} model.Response
// @Router /user-schedule/create [post]
// @Security BearerTokenAuth
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
	response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
}

// Retrieve ... Retrieve User Schedule
// @Summary Retrieve User Schedule
// @Description Retrieve User Schedule
// @Tags User Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserScheduleResponseData
// @Failure 400,500 {object} model.Response
// @Router /user-schedule/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id user schedule"
func (h userScheduleHandler) Retrieve(c *gin.Context) {
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
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update User Schedule
// @Summary Update User Schedule
// @Description Update User Schedule
// @Tags User Schedule
// @Accept       json
// @Produce      json
// @Param data body model.UserScheduleForm true "data"
// @Success 200 {object} model.UserScheduleResponseData
// @Failure 400,500 {object} model.Response
// @Router /user-schedule/update [put]
// @Security BearerTokenAuth
// @param id query string true "id user schedule"
func (h userScheduleHandler) Update(c *gin.Context) {
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
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Delete ... Delete User Schedule
// @Summary Delete User Schedule
// @Description Delete User Schedule
// @Tags User Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /user-schedule/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id user schedule"
func (h userScheduleHandler) Delete(c *gin.Context) {
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

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List All User Schedule
// @Summary List All User Schedule
// @Description List All User Schedule
// @Tags User Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserScheduleResponseList
// @Failure 400,500 {object} model.Response
// @Router /user-schedule/list [get]
// @Security BearerTokenAuth
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

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// List User In Schedule ... List User IN Schedule
// @Summary List All User IN Schedule
// @Description List All User IN Schedule
// @Tags User Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudentResponseList
// @Failure 400,500 {object} model.Response
// @Router /user-schedule/list/user-in-rule [get]
// @Security BearerTokenAuth
func (h userScheduleHandler) ListUserInRule(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Student
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.userScheduleService.ListUserInRule(data.ScheduleID, data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.userScheduleService.ListUserInRuleMeta(data.ScheduleID, data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data mahasiswa dalam aturan", dataList, metaList)
}

// List User Not In Schedule ... List User NOT IN Schedule
// @Summary List All User NOT IN Schedule
// @Description List All User NOT IN Schedule
// @Tags User Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudentResponseList
// @Failure 400,500 {object} model.Response
// @Router /user-schedule/list/user-not-in-rule [get]
// @Security BearerTokenAuth
func (h userScheduleHandler) ListUserNotInRule(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Student
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.userScheduleService.ListUserNotInRule(data.ScheduleID, data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.userScheduleService.ListUserNotInRuleMeta(data.ScheduleID, data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data mahasiswa tidak dalam aturan", dataList, metaList)
}

// Dropdown ... Dropdown All User Schedule
// @Summary Dropdown All User Schedule
// @Description Dropdown All User Schedule
// @Tags User Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserScheduleResponseList
// @Failure 400,500 {object} model.Response
// @Router /user-schedule/drop-down [get]
// @Security BearerTokenAuth
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

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}
