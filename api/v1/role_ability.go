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

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
)

type RoleAbilityHandler interface {
	Create(c *gin.Context)
	Generate(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type roleAbilityHandler struct {
	roleAbilityService service.RoleAbilityService
	infra              infra.Infra
	middleware         middleware.Middleware
}

func NewRoleAbilityHandler(roleAbilityService service.RoleAbilityService, infra infra.Infra, middleware middleware.Middleware) RoleAbilityHandler {
	return &roleAbilityHandler{
		roleAbilityService: roleAbilityService,
		infra:              infra,
		middleware:         middleware,
	}
}

// Create ... Create RoleAbility
// @Summary Create New RoleAbility
// @Description Create RoleAbility
// @Tags RoleAbility
// @Accept       json
// @Produce      json
// @Param data body model.RoleAbilityForm true "data"
// @Success 200 {object} model.RoleAbilityResponseData
// @Failure 400,500 {object} model.Response
// @Router /role-ability/create [post]
// @Security BearerTokenAuth
func (h roleAbilityHandler) Create(c *gin.Context) {
	var data model.RoleAbility
	c.BindJSON(&data)

	if !h.middleware.IsSuperAdmin(c) {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("anda tidak memiliki akses untuk fitur ini"))
		return
	}

	if err := validation.Validate(data.Action, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("action: %v", err))
		return
	}

	if err := validation.Validate(data.Subject, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("subject: %v", err))
		return
	}

	if isExist := h.roleAbilityService.CheckIsExistByActionAndSubject(data.Action, data.Subject, 0); isExist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("tidak bisa membuat data baru, data aksi dan subject sudah ada"))
		return
	}

	result, err := h.roleAbilityService.CreateRoleAbility(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
}

// Generate ... Generate Default RoleAbility
// @Summary Generate Default RoleAbility
// @Description Generate RoleAbility
// @Tags RoleAbility
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /role-ability/generate [get]
// @Security BearerTokenAuth
func (h roleAbilityHandler) Generate(c *gin.Context) {

	if !h.middleware.IsSuperAdmin(c) {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("anda tidak memiliki akses untuk fitur ini"))
		return
	}
	data1 := model.RoleAbility{
		IsSuperAdmin: true,
		IsAdmin:      true,
		IsUser:       true,
		Action:       "read",
		Subject:      "auth",
	}

	data2 := model.RoleAbility{
		IsSuperAdmin: true,
		IsAdmin:      false,
		IsUser:       false,
		Action:       "manage",
		Subject:      "all",
	}

	if isExist := h.roleAbilityService.CheckIsExistByActionAndSubject(data1.Action, data1.Subject, 0); !isExist {
		_, err := h.roleAbilityService.CreateRoleAbility(data1)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	if isExist := h.roleAbilityService.CheckIsExistByActionAndSubject(data2.Action, data2.Subject, 0); !isExist {
		_, err := h.roleAbilityService.CreateRoleAbility(data2)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "sukses generate default ability")
}

// Retrieve ... Retrieve RoleAbility
// @Summary Retrieve Single RoleAbility
// @Description Retrieve RoleAbility
// @Tags RoleAbility
// @Accept       json
// @Produce      json
// @Success 200 {object} model.RoleAbilityResponseData
// @Failure 400,500 {object} model.Response
// @Router /role-ability/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id subejct"
func (h roleAbilityHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("anda tidak memiliki akses untuk fitur ini"))
		return
	}

	result, err := h.roleAbilityService.RetrieveRoleAbility(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update RoleAbility
// @Summary Update Single RoleAbility
// @Description Update RoleAbility
// @Tags RoleAbility
// @Accept       json
// @Produce      json
// @Param data body model.RoleAbilityForm true "data"
// @Success 200 {object} model.RoleAbilityResponseData
// @Failure 400,500 {object} model.Response
// @Router /role-ability/update [put]
// @Security BearerTokenAuth
// @param id query string true "id roleability"
func (h roleAbilityHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("anda tidak memiliki akses untuk fitur ini"))
		return
	}

	var data model.RoleAbility
	c.BindJSON(&data)

	if err := validation.Validate(data.Action, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("action: %v", err))
		return
	}

	if err := validation.Validate(data.Subject, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("subject: %v", err))
		return
	}

	if isExist := h.roleAbilityService.CheckIsExistByActionAndSubject(data.Action, data.Subject, id); isExist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("tidak bisa membuat data baru, data aksi dan subject sudah ada"))
		return
	}

	var result model.RoleAbility

	result, err = h.roleAbilityService.UpdateRoleAbility(id, data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Delete ... Delete RoleAbility
// @Summary Delete Single RoleAbility
// @Description Delete RoleAbility
// @Tags RoleAbility
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /role-ability/delete [put]
// @Security BearerTokenAuth
// @param id query string true "id roleability"
func (h roleAbilityHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("anda tidak memiliki akses untuk fitur ini"))
		return
	}

	if err := h.roleAbilityService.DeleteRoleAbility(id); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List RoleAbility
// @Summary List All RoleAbility
// @Description List RoleAbility
// @Tags RoleAbility
// @Accept       json
// @Produce      json
// @Success 200 {object} model.RoleAbilityResponseList
// @Failure 400,500 {object} model.Response
// @Router /role-ability/list [get]
// @Security BearerTokenAuth
func (h roleAbilityHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.RoleAbility
	c.BindQuery(&data)

	if !h.middleware.IsSuperAdmin(c) {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("anda tidak memiliki akses untuk fitur ini"))
		return
	}

	dataList, err := h.roleAbilityService.ListRoleAbility(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.roleAbilityService.ListRoleAbilityMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// Dropdown ... Dropdown RoleAbility
// @Summary Dropdown All RoleAbility
// @Description Dropdown RoleAbility
// @Tags RoleAbility
// @Accept       json
// @Produce      json
// @Success 200 {object} model.RoleAbilityResponseList
// @Failure 400,500 {object} model.Response
// @Router /role-ability/drop-down [get]
// @Security BearerTokenAuth
func (h roleAbilityHandler) DropDown(c *gin.Context) {
	var data model.RoleAbility
	c.BindQuery(&data)

	if !h.middleware.IsSuperAdmin(c) {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("anda tidak memiliki akses untuk fitur ini"))
		return
	}

	dataList, err := h.roleAbilityService.DropDownRoleAbility(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}
