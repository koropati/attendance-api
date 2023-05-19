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
	Student(c *gin.Context)
	Teacher(c *gin.Context)
	Update(c *gin.Context)
	UpdatePassword(c *gin.Context)
}

type profileHandler struct {
	userService            service.UserService
	studentService         service.StudentService
	teacherService         service.TeacherService
	activationTokenService service.ActivationTokenService
	infra                  infra.Infra
	middleware             middleware.Middleware
}

func NewProfileHandler(
	userService service.UserService,
	studentService service.StudentService,
	teacherService service.TeacherService,
	activationTokenService service.ActivationTokenService,
	infra infra.Infra,
	middleware middleware.Middleware,
) ProfileHandler {
	return &profileHandler{
		userService:            userService,
		studentService:         studentService,
		teacherService:         teacherService,
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
// @Success 200 {object} model.AuthDataResponseData
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
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Student ... Student Profile
// @Summary Student Profile
// @Description Student Profile
// @Tags Profile
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudentResponseData
// @Failure 400,500 {object} model.Response
// @Router /profile/student [get]
// @Security BearerTokenAuth
func (h profileHandler) Student(c *gin.Context) {
	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	result, err := h.studentService.RetrieveStudentByUserID(currentUserID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Teacher ... Teacher Profile
// @Summary Teacher Profile
// @Description Teacher Profile
// @Tags Profile
// @Accept       json
// @Produce      json
// @Success 200 {object} model.TeacherResponseData
// @Failure 400,500 {object} model.Response
// @Router /profile/teacher [get]
// @Security BearerTokenAuth
func (h profileHandler) Teacher(c *gin.Context) {
	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	result, err := h.teacherService.RetrieveTeacherByUserID(currentUserID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
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
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama pengguna: %v", err))
		return
	}

	if err := validation.Validate(data.Email, validation.Required, validation.Length(6, 50)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("email: %v", err))
		return
	}

	if err := validation.Validate(data.FirstName, validation.Required, validation.Match(regexp.MustCompile(regex.NAME))); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama depan: %v", err))
		return
	}

	if !h.userService.CheckUpdateHandphone(currentUserID, data.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("no telp sudah digunakan"))
	}

	if !h.userService.CheckUpdateEmail(currentUserID, data.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email sudah digunakan"))
	}

	if !h.userService.CheckUpdateUsername(currentUserID, data.Username) {
		response.New(c).Error(http.StatusBadRequest, errors.New("nama pengguna sudah digunakan"))
	}

	if data.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("kata sandi: %v", err))
			return
		}
		data.Password = string(password)
	}

	result, err := h.userService.UpdateUser(currentUserID, data)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
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
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kata sandi saat ini: %v", err))
		return
	}

	if err := validation.Validate(data.NewPassword, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kata sandi baru: %v", err))
		return
	}

	if err := validation.Validate(data.ConfirmPassword, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("konfirmasi kata sandi: %v", err))
		return
	}

	if data.NewPassword != data.ConfirmPassword {
		err := fmt.Errorf("konfirmasi kata sandi tidak sesuai")
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	data.ID = uint(currentUserID)
	hashPassword, err := h.userService.GetPassword(currentUserID)

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("pengguna: %v", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(data.CurrentPassword)); err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("kata sandi tidak sesuai"))
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), 10)
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("kata sandi: %v", err))
		return
	}

	data.NewPassword = string(password)
	err = h.userService.UpdatePassword(data)
	if err != nil {
		response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("pengguna: %v", err))
		return
	}
	response.New(c).Write(http.StatusOK, "berhasil memperbaharui kata sandi")
}
