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

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SubjectHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type subjectHandler struct {
	subjectService service.SubjectService
	infra          infra.Infra
	middleware     middleware.Middleware
}

func NewSubjectHandler(subjectService service.SubjectService, infra infra.Infra, middleware middleware.Middleware) SubjectHandler {
	return &subjectHandler{
		subjectService: subjectService,
		infra:          infra,
		middleware:     middleware,
	}
}

func (h *subjectHandler) Create(c *gin.Context) {
	var data model.Subject
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

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}
	result, err := h.subjectService.CreateSubject(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "success create data", result)
}

func (h *subjectHandler) Retrieve(c *gin.Context) {
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

	var result *model.Subject
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.subjectService.RetrieveSubject(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.subjectService.RetrieveSubjectByOwner(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

func (h *subjectHandler) Update(c *gin.Context) {
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

	var data model.Subject
	c.BindJSON(&data)

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	var result *model.Subject
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.subjectService.UpdateSubject(id, &data)
	} else {
		result, err = h.subjectService.UpdateSubjectByOwner(id, currentUserID, &data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "success update data", result)
}

func (h *subjectHandler) Delete(c *gin.Context) {
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
		if err := h.subjectService.DeleteSubject(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.subjectService.DeleteSubjectByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "success delete data")
}

func (h *subjectHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Subject
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.subjectService.ListSubject(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.subjectService.ListSubjectMeta(&data, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "success get list data", dataList, metaList)
}

func (h *subjectHandler) DropDown(c *gin.Context) {
	var data model.Subject
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.subjectService.DropDownSubject(&data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
