package v1

import (
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/response"
	"attendance-api/common/util/regex"
	"attendance-api/infra"
	"attendance-api/model"
	"attendance-api/service"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type ProfileHandler interface {
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

type profileHandler struct {
	userService            service.UserService
	activationTokenService service.ActivationTokenService
	infra                  infra.Infra
	middleware             middleware.Middleware
}

func NewProfileHandler(userService service.UserService, activationTokenService service.ActivationTokenService, infra infra.Infra, middleware middleware.Middleware) ProfileHandler {
	return &profileHandler{
		userService:            userService,
		activationTokenService: activationTokenService,
		infra:                  infra,
		middleware:             middleware,
	}
}

// Retrieve ... Retrieve Profile
// @Summary Retrieve Profile
// @Description Retrieve Profile
// @Tags Profile
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserResponseData
// @Failure 400,500 {object} model.Response
// @Router /profile/ [get]
// @Security BearerTokenAuth
func (h profileHandler) Retrieve(c *gin.Context) {
	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	result, err := h.userService.RetrieveUser(currentUserID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

// Update ... Update Profile
// @Summary Update Profile
// @Description Update Profile
// @Tags Profile
// @Accept       json
// @Produce      json
// @Param data body model.UserForm true "data"
// @Success 200 {object} model.UserResponseData
// @Failure 400,500 {object} model.Response
// @Router /profile/update [put]
// @Security BearerTokenAuth
func (h profileHandler) Update(c *gin.Context) {
	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// Get User / valid exist data
	_, err = h.userService.RetrieveUser(currentUserID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	var data model.User
	c.BindJSON(&data)

	data.GormCustom.UpdatedBy = currentUserID
	data.GormCustom.UpdatedAt = time.Now()

	if err := validation.Validate(data.Username, validation.Required, validation.Length(4, 30), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := validation.Validate(data.Email, validation.Required, validation.Length(6, 50)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("email: %v", err))
		return
	}

	if err := validation.Validate(data.FirstName, validation.Required, validation.Match(regexp.MustCompile(regex.NAME))); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("first_name: %v", err))
		return
	}

	if !h.userService.CheckUpdateHandphone(currentUserID, data.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("handphone: already taken"))
	}

	if !h.userService.CheckUpdateEmail(currentUserID, data.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email: already taken"))
	}

	if !h.userService.CheckUpdateUsername(currentUserID, data.Username) {
		response.New(c).Error(http.StatusBadRequest, errors.New("username: already taken"))
	}

	if data.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("password: %v", err))
			return
		}
		data.Password = string(password)
	}

	result, err := h.userService.UpdateUser(currentUserID, data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

// Update Password ... Update Password
// @Summary Update Password
// @Description Update Password
// @Tags Profile
// @Accept       json
// @Produce      json
// @Param data body model.UserUpdatePasswordForm true "data"
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /profile/update-password [put]
// @Security BearerTokenAuth
func (h profileHandler) UpdatePassword(c *gin.Context) {
	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	var data model.UserUpdatePasswordForm
	c.BindJSON(&data)
	if err := validation.Validate(data.CurrentPassword, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("current_password: %v", err))
		return
	}

	if err := validation.Validate(data.NewPassword, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("new_password: %v", err))
		return
	}

	if err := validation.Validate(data.ConfirmPassword, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("confirm_password: %v", err))
		return
	}

	if data.NewPassword != data.ConfirmPassword {
		err := fmt.Errorf("confirm_password: password confirmation not match")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	data.ID = uint(currentUserID)
	hashPassword, err := h.userService.GetPassword(currentUserID)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("user: %v", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(data.CurrentPassword)); err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("password: password not match"))
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 10)
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("password: %v", err))
		return
	}

	data.NewPassword = string(password)
	err = h.userService.UpdatePassword(data)
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("user: %v", err))
		return
	}
	response.New(c).Write(http.StatusOK, "success update password")
}
