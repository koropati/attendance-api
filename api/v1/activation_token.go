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

type ActivationTokenHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type activationTokenHandler struct {
	activationTokenService service.ActivationTokenService
	infra                  infra.Infra
	middleware             middleware.Middleware
}

func NewActivationTokenHandler(activationTokenService service.ActivationTokenService, infra infra.Infra, middleware middleware.Middleware) ActivationTokenHandler {
	return &activationTokenHandler{
		activationTokenService: activationTokenService,
		infra:                  infra,
		middleware:             middleware,
	}
}

// Create ... Create Activation Token
// @Summary Create Activation Token
// @Description Create Activation Token
// @Tags Activation Token
// @Accept       json
// @Produce      json
// @Param data body model.ActivationTokenForm true "data"
// @Success 200 {object} model.ActivationTokenResponseData
// @Failure 400,500 {object} model.Response
// @Router /activation-token/create [post]
// @Security BearerTokenAuth
func (h activationTokenHandler) Create(c *gin.Context) {
	var data model.ActivationToken
	c.BindJSON(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	data.GormCustom.CreatedBy = currentUserID

	if err := validation.Validate(data.Token, validation.Required, validation.Length(1, 64)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("token: %v", err))
		return
	}

	if err := validation.Validate(data.UserID, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id pengguna: %v", err))
		return
	}

	data.Valid = time.Now().Add((time.Hour * 2))

	result, err := h.activationTokenService.CreateActivationToken(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
}

// Retrieve ... Retrieve Activation Token
// @Summary Retrieve Single Activation Token
// @Description Retrieve Single Activation Token
// @Tags Activation Token
// @Accept       json
// @Produce      json
// @Success 200 {object} model.ActivationTokenResponseData
// @Failure 400,500 {object} model.Response
// @Router /activation-token/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id activation token"
func (h activationTokenHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	result, err := h.activationTokenService.RetrieveActivationToken(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update Activation Token
// @Summary Update Single Activation Token
// @Description Update Single Activation Token
// @Tags Activation Token
// @Accept       json
// @Produce      json
// @Param data body model.ActivationTokenForm true "data"
// @Success 200 {object} model.ActivationTokenResponseData
// @Failure 400,500 {object} model.Response
// @Router /activation-token/update [put]
// @Security BearerTokenAuth
// @param id query string true "id activation token"
func (h activationTokenHandler) Update(c *gin.Context) {
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

	var data model.ActivationToken
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Token, validation.Required, validation.Length(1, 64)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("token: %v", err))
		return
	}

	if err := validation.Validate(data.UserID, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id pengguna: %v", err))
		return
	}

	result, err := h.activationTokenService.UpdateActivationToken(id, data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Delete ... Delete Activation Token
// @Summary Delete Single Activation Token
// @Description Delete Single Activation Token
// @Tags Activation Token
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /activation-token/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id activation token"
func (h activationTokenHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	if err := h.activationTokenService.DeleteActivationToken(id); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List All Activation Token
// @Summary List All Activation Token
// @Description List All Activation Token
// @Tags Activation Token
// @Accept       json
// @Produce      json
// @Success 200 {object} model.ActivationTokenResponseList
// @Failure 400,500 {object} model.Response
// @Router /activation-token/list [get]
// @Security BearerTokenAuth
func (h activationTokenHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.ActivationToken
	c.BindQuery(&data)

	dataList, err := h.activationTokenService.ListActivationToken(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.activationTokenService.ListActivationTokenMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// Dropdown ... Dropdown All Activation Token
// @Summary Dropdown All Activation Token
// @Description Dropdown All Activation Token
// @Tags Activation Token
// @Accept       json
// @Produce      json
// @Success 200 {object} model.ActivationTokenResponseList
// @Failure 400,500 {object} model.Response
// @Router /activation-token/drop-down [get]
// @Security BearerTokenAuth
func (h activationTokenHandler) DropDown(c *gin.Context) {
	var data model.ActivationToken
	c.BindQuery(&data)

	dataList, err := h.activationTokenService.DropDownActivationToken(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}
