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
)

type FacultyHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type facultyHandler struct {
	facultyService service.FacultyService
	infra          infra.Infra
	middleware     middleware.Middleware
}

func NewFacultyHandler(facultyService service.FacultyService, infra infra.Infra, middleware middleware.Middleware) FacultyHandler {
	return &facultyHandler{
		facultyService: facultyService,
		infra:          infra,
		middleware:     middleware,
	}
}

// Create ... Create Faculty
// @Summary Create New Faculty
// @Description Create Faculty
// @Tags Faculty
// @Accept       json
// @Produce      json
// @Param data body model.FacultyForm true "data"
// @Success 200 {object} model.FacultyResponseData
// @Failure 400,500 {object} model.Response
// @Router /faculty/create [post]
// @Security BearerTokenAuth
func (h facultyHandler) Create(c *gin.Context) {
	var data model.Faculty
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

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if h.facultyService.CheckIsExistByName(data.Name, 0) {
		err := errors.New("faculty name is already exists")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if h.facultyService.CheckIsExistByCode(data.Code, 0) {
		err := errors.New("faculty code is already exists")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", err))
		return
	}

	result, err := h.facultyService.CreateFaculty(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

// Retrieve ... Retrieve Faculty
// @Summary Retrieve Single Faculty
// @Description Retrieve Single Faculty
// @Tags Faculty
// @Accept       json
// @Produce      json
// @Success 200 {object} model.FacultyResponseData
// @Failure 400,500 {object} model.Response
// @Router /faculty/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id faculty"
func (h facultyHandler) Retrieve(c *gin.Context) {
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

	var result model.Faculty
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.facultyService.RetrieveFaculty(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.facultyService.RetrieveFacultyByOwner(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

// Update ... Update Faculty
// @Summary Update Single Faculty
// @Description Update Single Faculty
// @Tags Faculty
// @Accept       json
// @Produce      json
// @Param data body model.FacultyForm true "data"
// @Success 200 {object} model.FacultyResponseData
// @Failure 400,500 {object} model.Response
// @Router /faculty/update [put]
// @Security BearerTokenAuth
// @param id query string true "id faculty"
func (h facultyHandler) Update(c *gin.Context) {
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

	var data model.Faculty
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if h.facultyService.CheckIsExistByName(data.Name, id) {
		err := errors.New("faculty name is already exists")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if h.facultyService.CheckIsExistByCode(data.Code, id) {
		err := errors.New("faculty code is already exists")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", err))
		return
	}

	var result model.Faculty
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.facultyService.UpdateFaculty(id, data)
	} else {
		result, err = h.facultyService.UpdateFacultyByOwner(id, currentUserID, data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

// Delete ... Delete Faculty
// @Summary Delete Single Faculty
// @Description Delete Single Faculty
// @Tags Faculty
// @Accept       json
// @Produce      json
// @Success 200 {object} model.FacultyResponseData
// @Failure 400,500 {object} model.Response
// @Router /faculty/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id faculty"
func (h facultyHandler) Delete(c *gin.Context) {
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
		if err := h.facultyService.DeleteFaculty(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.facultyService.DeleteFacultyByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "success delete data")
}

// List ... List all Faculty
// @Summary List all Faculty
// @Description List all Faculty
// @Tags Faculty
// @Accept       json
// @Produce      json
// @Success 200 {object} model.FacultyResponseList
// @Failure 400,500 {object} model.Response
// @Router /faculty/list [get]
// @Security BearerTokenAuth
func (h facultyHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Faculty
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.facultyService.ListFaculty(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.facultyService.ListFacultyMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

// Dropdown ... Dropdown all Faculty
// @Summary Dropdown all Faculty
// @Description Dropdown all Faculty
// @Tags Faculty
// @Accept       json
// @Produce      json
// @Success 200 {object} model.FacultyResponseList
// @Failure 400,500 {object} model.Response
// @Router /faculty/drop-down [get]
// @Security BearerTokenAuth
func (h facultyHandler) DropDown(c *gin.Context) {
	var data model.Faculty
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.facultyService.DropDownFaculty(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
