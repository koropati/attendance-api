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

type StudyProgramHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type studyProgramHandler struct {
	studyProgramService service.StudyProgramService
	infra               infra.Infra
	middleware          middleware.Middleware
}

func NewStudyProgramHandler(studyProgramService service.StudyProgramService, infra infra.Infra, middleware middleware.Middleware) StudyProgramHandler {
	return &studyProgramHandler{
		studyProgramService: studyProgramService,
		infra:               infra,
		middleware:          middleware,
	}
}

// Create ... Create Study Program
// @Summary Create New Study Program
// @Description Create Study Program
// @Tags Study Program
// @Accept       json
// @Produce      json
// @Param data body model.StudyProgramForm true "data"
// @Success 200 {object} model.StudyProgramResponseData
// @Failure 400,500 {object} model.Response
// @Router /study-program/create [post]
// @Security BearerTokenAuth
func (h studyProgramHandler) Create(c *gin.Context) {
	var data model.StudyProgram
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
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if err := validation.Validate(data.MajorID, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id jurusan: %v", err))
		return
	}

	if h.studyProgramService.CheckIsExistByName(data.Name, int(data.MajorID), 0) {
		err := errors.New("nama program studi sudah ada")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if h.studyProgramService.CheckIsExistByCode(data.Code, 0) {
		err := errors.New("kode program studi sudah ada")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kode: %v", err))
		return
	}

	result, err := h.studyProgramService.CreateStudyProgram(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
}

// Retrieve ... Retrieve Study Program
// @Summary Retrieve New Study Program
// @Description Retrieve Study Program
// @Tags Study Program
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudyProgramResponseData
// @Failure 400,500 {object} model.Response
// @Router /study-program/retrieve [get]
// @Security BearerTokenAuth
func (h studyProgramHandler) Retrieve(c *gin.Context) {
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

	var result model.StudyProgram
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.studyProgramService.RetrieveStudyProgram(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.studyProgramService.RetrieveStudyProgramByOwner(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update Study Program
// @Summary Update single Study Program
// @Description Update Study Program
// @Tags Study Program
// @Accept       json
// @Produce      json
// @Param data body model.StudyProgramForm true "data"
// @Success 200 {object} model.StudyProgramResponseData
// @Failure 400,500 {object} model.Response
// @Router /study-program/update [put]
// @Security BearerTokenAuth
// @param id query string true "id study program"
func (h studyProgramHandler) Update(c *gin.Context) {
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

	var data model.StudyProgram
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if err := validation.Validate(data.MajorID, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id jurusan: %v", err))
		return
	}

	if h.studyProgramService.CheckIsExistByName(data.Name, int(data.MajorID), id) {
		err := errors.New("nama program studi sudah ada")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if h.studyProgramService.CheckIsExistByCode(data.Code, id) {
		err := errors.New("kode program studi sudah ada")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kode: %v", err))
		return
	}

	var result model.StudyProgram
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.studyProgramService.UpdateStudyProgram(id, data)
	} else {
		result, err = h.studyProgramService.UpdateStudyProgramByOwner(id, currentUserID, data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Delete ... Delete Study Program
// @Summary Delete single Study Program
// @Description Delete Study Program
// @Tags Study Program
// @Accept       json
// @Produce      json
// @Param data body model.StudyProgramForm true "data"
// @Success 200 {object} model.StudyProgramResponseData
// @Failure 400,500 {object} model.Response
// @Router /study-program/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id study program"
func (h studyProgramHandler) Delete(c *gin.Context) {
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

	if h.middleware.IsSuperAdmin(c) {
		if err := h.studyProgramService.DeleteStudyProgram(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.studyProgramService.DeleteStudyProgramByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List Study Program
// @Summary List All Study Program
// @Description List Study Program
// @Tags Study Program
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudyProgramResponseList
// @Failure 400,500 {object} model.Response
// @Router /study-program/list [get]
// @Security BearerTokenAuth
func (h studyProgramHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.StudyProgram
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.studyProgramService.ListStudyProgram(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.studyProgramService.ListStudyProgramMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// Dropdown ... Dropdown Study Program
// @Summary Dropdown All Study Program
// @Description Dropdown Study Program
// @Tags Study Program
// @Accept       json
// @Produce      json
// @Success 200 {object} model.StudyProgramResponseList
// @Failure 400,500 {object} model.Response
// @Router /study-program/drop-down [get]
// @Security BearerTokenAuth
func (h studyProgramHandler) DropDown(c *gin.Context) {
	var data model.StudyProgram
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.studyProgramService.DropDownStudyProgram(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}
