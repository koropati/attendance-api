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

type PasswordResetTokenHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type passwordResetTokenHandler struct {
	passwordResetTokenService service.PasswordResetTokenService
	infra                     infra.Infra
	middleware                middleware.Middleware
}

func NewPasswordResetTokenHandler(passwordResetTokenService service.PasswordResetTokenService, infra infra.Infra, middleware middleware.Middleware) PasswordResetTokenHandler {
	return &passwordResetTokenHandler{
		passwordResetTokenService: passwordResetTokenService,
		infra:                     infra,
		middleware:                middleware,
	}
}

// Create ... Create Password Reset Token
// @Summary Create New Password Reset Token
// @Description Create Password Reset Token
// @Tags Password Reset Token
// @Accept       json
// @Produce      json
// @Param data body model.PasswordResetTokenForm true "data"
// @Success 200 {object} model.PasswordResetTokenResponseData
// @Failure 400,500 {object} model.Response
// @Router /password-reset-token/create [post]
// @Security BearerTokenAuth
func (h passwordResetTokenHandler) Create(c *gin.Context) {
	var data model.PasswordResetToken
	c.BindJSON(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	data.GormCustom.CreatedBy = currentUserID

	if err := validation.Validate(data.Token, validation.Required, validation.Length(1, 64), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("token: %v", err))
		return
	}

	if err := validation.Validate(data.UserID, validation.Required, is.UTFNumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id pengguna: %v", err))
		return
	}

	data.Valid = time.Now().Add((time.Hour * 2))

	result, err := h.passwordResetTokenService.CreatePasswordResetToken(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

// Retrieve ... Retrieve Password Reset Token
// @Summary Retrieve Single Password Reset Token
// @Description Retrieve Password Reset Token
// @Tags Password Reset Token
// @Accept       json
// @Produce      json
// @Success 200 {object} model.PasswordResetTokenResponseData
// @Failure 400,500 {object} model.Response
// @Router /password-reset-token/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id"
func (h passwordResetTokenHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	result, err := h.passwordResetTokenService.RetrievePasswordResetToken(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

// Update ... Update Password Reset Token
// @Summary Update Single Password Reset Token
// @Description Update Password Reset Token
// @Tags Password Reset Token
// @Accept       json
// @Produce      json
// @Param data body model.PasswordResetTokenForm true "data"
// @Success 200 {object} model.PasswordResetTokenResponseData
// @Failure 400,500 {object} model.Response
// @Router /password-reset-token/update [put]
// @Security BearerTokenAuth
// @param id query string true "id"
func (h passwordResetTokenHandler) Update(c *gin.Context) {
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

	var data model.PasswordResetToken
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Token, validation.Required, validation.Length(1, 64), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("token: %v", err))
		return
	}

	if err := validation.Validate(data.UserID, validation.Required, is.UTFNumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id pengguna: %v", err))
		return
	}

	result, err := h.passwordResetTokenService.UpdatePasswordResetToken(id, data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

// Delete ... Delete Password Reset Token
// @Summary Delete Single Password Reset Token
// @Description Delete Password Reset Token
// @Tags Password Reset Token
// @Accept       json
// @Produce      json
// @Success 200 {object} model.PasswordResetTokenResponseData
// @Failure 400,500 {object} model.Response
// @Router /password-reset-token/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id"
func (h passwordResetTokenHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	if err := h.passwordResetTokenService.DeletePasswordResetToken(id); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(http.StatusOK, "success delete data")
}

// List ... List Password Reset Token
// @Summary List all Password Reset Token
// @Description List all Password Reset Token
// @Tags Password Reset Token
// @Accept       json
// @Produce      json
// @Success 200 {object} model.PasswordResetTokenResponseList
// @Failure 400,500 {object} model.Response
// @Router /password-reset-token/list [get]
// @Security BearerTokenAuth
func (h passwordResetTokenHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.PasswordResetToken
	c.BindQuery(&data)

	dataList, err := h.passwordResetTokenService.ListPasswordResetToken(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.passwordResetTokenService.ListPasswordResetTokenMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

// Dropdown ... Dropdown Password Reset Token
// @Summary Dropdown all Password Reset Token
// @Description Dropdown all Password Reset Token
// @Tags Password Reset Token
// @Accept       json
// @Produce      json
// @Success 200 {object} model.PasswordResetTokenResponseList
// @Failure 400,500 {object} model.Response
// @Router /password-reset-token/drop-down [get]
// @Security BearerTokenAuth
func (h passwordResetTokenHandler) DropDown(c *gin.Context) {
	var data model.PasswordResetToken
	c.BindQuery(&data)

	dataList, err := h.passwordResetTokenService.DropDownPasswordResetToken(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
