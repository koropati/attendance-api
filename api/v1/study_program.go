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
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if err := validation.Validate(data.MajorID, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("major_id: %v", err))
		return
	}

	if h.studyProgramService.CheckIsExistByName(data.Name, int(data.MajorID), 0) {
		err := errors.New("study program name is already exists")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if h.studyProgramService.CheckIsExistByCode(data.Code, 0) {
		err := errors.New("study program code is already exists")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", err))
		return
	}

	result, err := h.studyProgramService.CreateStudyProgram(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

func (h studyProgramHandler) Retrieve(c *gin.Context) {
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
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

func (h studyProgramHandler) Update(c *gin.Context) {
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

	var data model.StudyProgram
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if err := validation.Validate(data.MajorID, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("major_id: %v", err))
		return
	}

	if h.studyProgramService.CheckIsExistByName(data.Name, int(data.MajorID), id) {
		err := errors.New("study program name is already exists")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if h.studyProgramService.CheckIsExistByCode(data.Code, id) {
		err := errors.New("study program code is already exists")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", err))
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
	response.New(c).Data(http.StatusOK, "success update data", result)
}

func (h studyProgramHandler) Delete(c *gin.Context) {
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

	response.New(c).Write(http.StatusOK, "success delete data")
}

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

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

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

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
