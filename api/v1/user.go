package v1

import (
	"attendance-api/common/http/email"
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/response"
	"attendance-api/common/util/activation"
	"attendance-api/common/util/pagination"
	"attendance-api/common/util/regex"
	"attendance-api/infra"
	"attendance-api/model"
	"attendance-api/service"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
	SetActive(c *gin.Context)
	SetDeactive(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

type userHandler struct {
	userService            service.UserService
	activationTokenService service.ActivationTokenService
	infra                  infra.Infra
	middleware             middleware.Middleware
}

func NewUserHandler(userService service.UserService, activationTokenService service.ActivationTokenService, infra infra.Infra, middleware middleware.Middleware) UserHandler {
	return &userHandler{
		userService:            userService,
		activationTokenService: activationTokenService,
		infra:                  infra,
		middleware:             middleware,
	}
}

// Create ... Create User
// @Summary Create New User
// @Description Create user
// @Tags User
// @Accept       json
// @Produce      json
// @Param data body model.UserForm true "data"
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /user/create [post]
// @Security BearerTokenAuth
func (h userHandler) Create(c *gin.Context) {
	var data model.User
	c.BindJSON(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	data.GormCustom.CreatedBy = currentUserID

	if err := validation.Validate(data.Username, validation.Required, validation.Length(4, 30), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := validation.Validate(data.Password, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("password: %v", err))
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

	if !h.userService.CheckHandphone(data.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("handphone: already taken"))
	}

	if !h.userService.CheckEmail(data.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email: already taken"))
	}

	if h.userService.CheckUsername(data.Username) {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("password: %v", err))
			return
		}

		data.Password = string(password)
		if _, err := h.userService.CreateUser(data); err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		user, err := h.userService.RetrieveUserByUsername(data.Username)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		expiredToken, activationToken := activation.New(user).GenerateSHA1(24)

		// Save Activation token to data base
		activationData, err := h.activationTokenService.CreateActivationToken(model.ActivationToken{
			UserID: user.ID,
			Token:  activationToken,
			Valid:  expiredToken,
		})

		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		go func(user model.User) {
			config := h.infra.Config().Sub("server")
			urlActivation := fmt.Sprintf("%s:%s/v1/auth/activation?token=%s", config.GetString("url"), config.GetString("port"), activationData.Token)

			if err := email.New(h.infra.GoMail(), h.infra.Config()).SendActivation(user.FirstName, user.Email, urlActivation); err != nil {
				log.Printf("Error Send Email E: %v", err)
			}
		}(user)

		response.New(c).Write(http.StatusCreated, "success: user registered")
		return
	}
}

// Retrieve ... Retrieve User
// @Summary Retreive Single User
// @Description Retreive Single User
// @Tags User
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserResponseData
// @Failure 400,500 {object} model.Response
// @Router /user/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id user"
func (h userHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	result, err := h.userService.RetrieveUser(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

// Update ... Update User
// @Summary Update Single User
// @Description Update Single User
// @Tags User
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserResponseData
// @Failure 400,500 {object} model.Response
// @Router /user/update [put]
// @Security BearerTokenAuth
// @param id query string true "id user"
// @Param data body model.UserForm true "data"
func (h userHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	// Get User / valid exist data
	_, err = h.userService.RetrieveUser(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
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

	if !h.userService.CheckUpdateHandphone(id, data.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("handphone: already taken"))
	}

	if !h.userService.CheckUpdateEmail(id, data.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email: already taken"))
	}

	if !h.userService.CheckUpdateUsername(id, data.Username) {
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

	result, err := h.userService.UpdateUser(id, data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

// Delete ... Delete User
// @Summary Delete Single User
// @Description Delete Single User
// @Tags User
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /user/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id user"
func (h userHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	if err := h.userService.DeleteUser(id); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(http.StatusOK, "success delete data")
}

// List ... List User
// @Summary List Data User
// @Description List Data User
// @Tags User
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserResponseList
// @Failure 400,500 {object} model.Response
// @Router /user/list [get]
// @Security BearerTokenAuth
func (h userHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var user model.User
	c.BindQuery(&user)

	userList, err := h.userService.ListUser(user, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.userService.ListUserMeta(user, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list user", userList, metaList)
}

// Dropdown ... Dropdown User
// @Summary List Data User
// @Description List Data User
// @Tags User
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserResponseList
// @Failure 400,500 {object} model.Response
// @Router /user/drop-down [get]
// @Security BearerTokenAuth
func (h userHandler) DropDown(c *gin.Context) {
	var data model.User
	c.BindQuery(&data)

	dataList, err := h.userService.DropDownUser(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}

// Set Active ... Set Active User
// @Summary Set Active Data User
// @Description Set Active Data User
// @Tags User
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserResponseData
// @Failure 400,500 {object} model.Response
// @Router /user/active [patch]
// @Security BearerTokenAuth
// @param id query string true "id user"
func (h userHandler) SetActive(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	result, err := h.userService.SetActiveUser(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}
	response.New(c).Data(http.StatusOK, "success set active data", result)
}

// Set Deactive ... Set Deactive User
// @Summary Set Deactive Data User
// @Description Set Deactive Data User
// @Tags User
// @Accept       json
// @Produce      json
// @Success 200 {object} model.UserResponseData
// @Failure 400,500 {object} model.Response
// @Router /user/deactive [patch]
// @Security BearerTokenAuth
// @param id query string true "id user"
func (h userHandler) SetDeactive(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	result, err := h.userService.SetDeactiveUser(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}
	response.New(c).Data(http.StatusOK, "success set deactive data", result)
}

func (h userHandler) UpdatePassword(c *gin.Context) {
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

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
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
