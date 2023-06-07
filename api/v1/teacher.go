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

type TeacherHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type teacherHandler struct {
	userService            service.UserService
	teacherService         service.TeacherService
	activationTokenService service.ActivationTokenService
	infra                  infra.Infra
	middleware             middleware.Middleware
}

func NewTeacherHandler(userService service.UserService, teacherService service.TeacherService, activationTokenService service.ActivationTokenService, infra infra.Infra, middleware middleware.Middleware) TeacherHandler {
	return &teacherHandler{
		userService:            userService,
		teacherService:         teacherService,
		activationTokenService: activationTokenService,
		infra:                  infra,
		middleware:             middleware,
	}
}

// Create ... Create Teacher
// @Summary Create New Teacher
// @Description Create Teacher
// @Tags Teacher
// @Accept       json
// @Produce      json
// @Param data body model.TeacherForm true "data"
// @Success 200 {object} model.TeacherResponseData
// @Failure 400,500 {object} model.Response
// @Router /teacher/create [post]
// @Security BearerTokenAuth
func (h teacherHandler) Create(c *gin.Context) {
	var data model.Teacher
	c.BindJSON(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	data.GormCustom.CreatedBy = currentUserID

	if err := validation.Validate(data.Nip, validation.Required, validation.Length(1, 20)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nip: %v", err))
		return
	}

	if h.teacherService.CheckIsExistByNip(data.Nip, 0) {
		err := errors.New("nip dosen sudah ada yang menggunakan")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nip: %v", err))
		return
	}

	if err := validation.Validate(data.User.Username, validation.Required, validation.Length(4, 30), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama pengguna: %v", err))
		return
	}

	if err := validation.Validate(data.User.Email, validation.Required, validation.Length(6, 50)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("email: %v", err))
		return
	}

	if err := validation.Validate(data.User.FirstName, validation.Required, validation.Match(regexp.MustCompile(regex.NAME))); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama depan: %v", err))
		return
	}

	if !h.userService.CheckHandphone(data.User.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("no telp sudah digunakan"))
	}

	if !h.userService.CheckEmail(data.User.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email sudah digunakan"))
	}

	if h.userService.CheckUsername(data.User.Username) {

		password, err := bcrypt.GenerateFromPassword([]byte(data.GeneratePassword()), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("kata sandi: %v", err))
			return
		}

		data.User.Password = string(password)
		data.User.IsAdmin = true
		data.User.IsUser = false
		loginDate, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
		data.User.LastLogin = loginDate

		result, err := h.teacherService.CreateTeacher(data)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}

		user, err := h.userService.RetrieveUserByUsername(data.User.Username)
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
			urlActivation := fmt.Sprintf("%s/v1/auth/activation?token=%s", config.GetString("web_url"), activationData.Token)

			if err := email.New(h.infra.GoMail(), h.infra.Config()).SendActivation(user.FirstName, user.Email, urlActivation); err != nil {
				log.Printf("Error Send Email E: %v", err)
			}
		}(user)

		response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
		return
	}
}

// Retrieve ... Retrieve Teacher
// @Summary Retrieve Single Teacher
// @Description Retrieve Single Teacher
// @Tags Teacher
// @Accept       json
// @Produce      json
// @Success 200 {object} model.TeacherResponseData
// @Failure 400,500 {object} model.Response
// @Router /teacher/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id teacher"
func (h teacherHandler) Retrieve(c *gin.Context) {
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

	var result model.Teacher
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.teacherService.RetrieveTeacher(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.teacherService.RetrieveTeacherByOwner(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update Teacher
// @Summary Update Single Teacher
// @Description Update Single Teacher
// @Tags Teacher
// @Accept       json
// @Produce      json
// @Param data body model.TeacherForm true "data"
// @Success 200 {object} model.TeacherResponseData
// @Failure 400,500 {object} model.Response
// @Router /teacher/update [put]
// @Security BearerTokenAuth
// @param id query string true "id teacher"
func (h teacherHandler) Update(c *gin.Context) {
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

	var data model.Teacher
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	oldDataTeacher, err := h.teacherService.RetrieveTeacher(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("%v", "data dosen tidak ditemukan"))
		return
	}

	if err := validation.Validate(data.Nip, validation.Required, validation.Length(1, 20)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nip: %v", err))
		return
	}

	if h.teacherService.CheckIsExistByNip(data.Nip, id) {
		err := errors.New("nip dosen sudah ada yang menggunakan")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nip: %v", err))
		return
	}

	data.GormCustom.UpdatedBy = currentUserID
	data.GormCustom.UpdatedAt = time.Now()

	if err := validation.Validate(data.User.Username, validation.Required, validation.Length(4, 30)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama pengguna: %v", err))
		return
	}

	if err := validation.Validate(data.User.Email, validation.Required, validation.Length(6, 50)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("email: %v", err))
		return
	}

	if err := validation.Validate(data.User.FirstName, validation.Required, validation.Match(regexp.MustCompile(regex.NAME))); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama depan: %v", err))
		return
	}

	if !h.userService.CheckUpdateHandphone(int(oldDataTeacher.User.ID), data.User.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("no telp sudah digunakan"))
	}

	if !h.userService.CheckUpdateEmail(int(oldDataTeacher.User.ID), data.User.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email sudah digunakan"))
	}

	if !h.userService.CheckUpdateUsername(int(oldDataTeacher.User.ID), data.User.Username) {
		response.New(c).Error(http.StatusBadRequest, errors.New("nama pengguna sudah digunakan"))
	}

	if data.User.Password != "" {
		password, err := bcrypt.GenerateFromPassword([]byte(data.User.Password), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("kata sandi: %v", err))
			return
		}
		data.User.Password = string(password)
	}

	// get oldUserData
	oldUser, err := h.userService.RetrieveUser(int(oldDataTeacher.User.ID))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	// fill infomartion role
	data.User.IsActive = oldUser.IsActive
	data.User.IsUser = oldUser.IsUser
	data.User.IsAdmin = oldUser.IsAdmin
	data.User.IsSuperAdmin = oldUser.IsSuperAdmin

	_, err = h.userService.UpdateUser(int(oldDataTeacher.User.ID), data.User)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !data.User.IsActive {
		_, err = h.userService.SetDeactiveUser(int(oldDataTeacher.User.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		_, err = h.userService.SetActiveUser(int(oldDataTeacher.User.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	var result model.Teacher
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.teacherService.UpdateTeacher(id, data)
	} else {
		result, err = h.teacherService.UpdateTeacherByOwner(id, currentUserID, data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Delete ... Delete Teacher
// @Summary Delete Single Teacher
// @Description Delete Single Teacher
// @Tags Teacher
// @Accept       json
// @Produce      json
// @Success 200 {object} model.TeacherResponseData
// @Failure 400,500 {object} model.Response
// @Router /teacher/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id teacher"
func (h teacherHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	dataTeacher, err := h.teacherService.RetrieveTeacher(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("data dosen tidak ditemukan"))
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if h.middleware.IsSuperAdmin(c) {
		if err := h.teacherService.DeleteTeacher(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.teacherService.DeleteTeacherByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	if err := h.userService.DeleteUser(int(dataTeacher.UserID)); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List all Teacher
// @Summary List all Teacher
// @Description List all Teacher
// @Tags Teacher
// @Accept       json
// @Produce      json
// @Success 200 {object} model.TeacherResponseList
// @Failure 400,500 {object} model.Response
// @Router /teacher/list [get]
// @Security BearerTokenAuth
func (h teacherHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Teacher
	c.BindQuery(&data)

	dataList, err := h.teacherService.ListTeacher(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.teacherService.ListTeacherMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// Dropdown ... Dropdown all Teacher
// @Summary Dropdown all Teacher
// @Description Dropdown all Teacher
// @Tags Teacher
// @Accept       json
// @Produce      json
// @Success 200 {object} model.TeacherResponseList
// @Failure 400,500 {object} model.Response
// @Router /teacher/drop-down [get]
// @Security BearerTokenAuth
func (h teacherHandler) DropDown(c *gin.Context) {
	var data model.Teacher
	c.BindQuery(&data)

	dataList, err := h.teacherService.DropDownTeacher(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}
