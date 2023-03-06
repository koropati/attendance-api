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

func (h *activationTokenHandler) Create(c *gin.Context) {
	var data model.ActivationToken
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
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("user_id: %v", err))
		return
	}

	data.Valid = time.Now().Add((time.Hour * 2))

	result, err := h.activationTokenService.CreateActivationToken(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

func (h *activationTokenHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	result, err := h.activationTokenService.RetrieveActivationToken(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

func (h *activationTokenHandler) Update(c *gin.Context) {
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

	var data model.ActivationToken
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Token, validation.Required, validation.Length(1, 64), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("token: %v", err))
		return
	}

	if err := validation.Validate(data.UserID, validation.Required, is.UTFNumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("user_id: %v", err))
		return
	}

	result, err := h.activationTokenService.UpdateActivationToken(id, &data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

func (h *activationTokenHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	if err := h.activationTokenService.DeleteActivationToken(id); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(http.StatusOK, "success delete data")
}

func (h *activationTokenHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.ActivationToken
	c.BindQuery(&data)

	dataList, err := h.activationTokenService.ListActivationToken(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.activationTokenService.ListActivationTokenMeta(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

func (h *activationTokenHandler) DropDown(c *gin.Context) {
	var data model.ActivationToken
	c.BindQuery(&data)

	dataList, err := h.activationTokenService.DropDownActivationToken(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
