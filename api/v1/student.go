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

type StudentHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type studentHandler struct {
	userService            service.UserService
	studentService         service.StudentService
	activationTokenService service.ActivationTokenService
	infra                  infra.Infra
	middleware             middleware.Middleware
}

func NewStudentHandler(userService service.UserService, studentService service.StudentService, activationTokenService service.ActivationTokenService, infra infra.Infra, middleware middleware.Middleware) StudentHandler {
	return &studentHandler{
		userService:            userService,
		studentService:         studentService,
		activationTokenService: activationTokenService,
		infra:                  infra,
		middleware:             middleware,
	}
}

// Create ... Create Student
// @Summary Create New Student
// @Description Create Student
// @Tags Student
// @Accept       json
// @Produce      json
// @Param data body model.StudentForm true "data"
// @Success 200 {object} model.StudentResponseData
// @Failure 400,500 {object} model.Response
// @Router /student/create [post]
// @Security BearerTokenAuth
func (h studentHandler) Create(c *gin.Context) {
	var data model.Student
	c.BindJSON(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	data.GormCustom.CreatedBy = currentUserID

	if err := validation.Validate(data.NIM, validation.Required, validation.Length(1, 20)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nim: %v", err))
		return
	}

	if h.studentService.CheckIsExistByNIM(data.NIM, 0) {
		err := errors.New("nim mahasiswa sudah ada yang menggunakan")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nim: %v", err))
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
		data.User.IsUser = true
		data.User.IsAdmin = false
		loginDate, _ := time.Parse("2006-01-02 15:04:05", "0001-01-01 00:00:00")
		data.User.LastLogin = loginDate

		result, err := h.studentService.CreateStudent(data)
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

// Retrieve ... Retrieve Student
// @Summary Retrieve Single Student
// @Description Retrieve Single Student
// @Tags Student
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudentResponseData
// @Failure 400,500 {object} model.Response
// @Router /student/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id student"
func (h studentHandler) Retrieve(c *gin.Context) {
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

	var result model.Student
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.studentService.RetrieveStudent(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.studentService.RetrieveStudentByOwner(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update Student
// @Summary Update Single Student
// @Description Update Single Student
// @Tags Student
// @Accept       json
// @Produce      json
// @Param data body model.StudentForm true "data"
// @Success 200 {object} model.StudentResponseData
// @Failure 400,500 {object} model.Response
// @Router /student/update [put]
// @Security BearerTokenAuth
// @param id query string true "id student"
func (h studentHandler) Update(c *gin.Context) {
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

	var data model.Student
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	oldDataStudent, err := h.studentService.RetrieveStudent(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("%v", "data mahasiswa tidak ditemukan"))
		return
	}

	if err := validation.Validate(data.NIM, validation.Required, validation.Length(1, 20)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nim: %v", err))
		return
	}

	if h.studentService.CheckIsExistByNIM(data.NIM, id) {
		err := errors.New("nim mahasiswa sudah ada yang menggunakan")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nim: %v", err))
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

	if !h.userService.CheckUpdateHandphone(int(oldDataStudent.User.ID), data.User.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("no telp sudah digunakan"))
	}

	if !h.userService.CheckUpdateEmail(int(oldDataStudent.User.ID), data.User.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email sudah digunakan"))
	}

	if !h.userService.CheckUpdateUsername(int(oldDataStudent.User.ID), data.User.Username) {
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

	_, err = h.userService.UpdateUser(int(oldDataStudent.User.ID), data.User)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !data.User.IsActive {
		_, err = h.userService.SetDeactiveUser(int(oldDataStudent.User.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		_, err = h.userService.SetActiveUser(int(oldDataStudent.User.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	var result model.Student
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.studentService.UpdateStudent(id, data)
	} else {
		result, err = h.studentService.UpdateStudentByOwner(id, currentUserID, data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Delete ... Delete Student
// @Summary Delete Single Student
// @Description Delete Single Student
// @Tags Student
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudentResponseData
// @Failure 400,500 {object} model.Response
// @Router /student/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id student"
func (h studentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	dataStudent, err := h.studentService.RetrieveStudent(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("data mahasiswa tidak ditemukan"))
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if h.middleware.IsSuperAdmin(c) {
		if err := h.studentService.DeleteStudent(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.studentService.DeleteStudentByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	if err := h.userService.DeleteUser(int(dataStudent.UserID)); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List all Student
// @Summary List all Student
// @Description List all Student
// @Tags Student
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudentResponseList
// @Failure 400,500 {object} model.Response
// @Router /student/list [get]
// @Security BearerTokenAuth
func (h studentHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Student
	c.BindQuery(&data)

	dataList, err := h.studentService.ListStudent(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.studentService.ListStudentMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// Dropdown ... Dropdown all Student
// @Summary Dropdown all Student
// @Description Dropdown all Student
// @Tags Student
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudentResponseList
// @Failure 400,500 {object} model.Response
// @Router /student/drop-down [get]
// @Security BearerTokenAuth
func (h studentHandler) DropDown(c *gin.Context) {
	var data model.Student
	c.BindQuery(&data)

	dataList, err := h.studentService.DropDownStudent(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}
