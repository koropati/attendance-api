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

// Create ... Create Subject
// @Summary Create New Subject
// @Description Create Subject
// @Tags Subject
// @Accept       json
// @Produce      json
// @Param data body model.SubjectForm true "data"
// @Success 200 {object} model.SubjectResponseData
// @Failure 400,500 {object} model.Response
// @Router /subject/create [post]
// @Security BearerTokenAuth
func (h subjectHandler) Create(c *gin.Context) {
	var data model.Subject
	c.BindJSON(&data)

	if !h.middleware.IsSuperAdmin(c) {
		err := errors.New("anda tidak memiliki akses untuk melakukan proses ini")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}

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
	result, err := h.subjectService.CreateSubject(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
}

// Retrieve ... Retrieve Subject
// @Summary Retrieve Single Subject
// @Description Retrieve Subject
// @Tags Subject
// @Accept       json
// @Produce      json
// @Success 200 {object} model.SubjectResponseData
// @Failure 400,500 {object} model.Response
// @Router /subject/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id subejct"
func (h subjectHandler) Retrieve(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	// currentUserID, err := h.middleware.GetUserID(c)
	// if err != nil {
	// 	response.New(c).Error(http.StatusBadRequest, err)
	// 	return
	// }

	var result model.Subject
	// if h.middleware.IsSuperAdmin(c) {
	// 	result, err = h.subjectService.RetrieveSubject(id)
	// 	if err != nil {
	// 		response.New(c).Error(http.StatusBadRequest, err)
	// 		return
	// 	}
	// } else {
	// 	result, err = h.subjectService.RetrieveSubjectByOwner(id, currentUserID)
	// 	if err != nil {
	// 		response.New(c).Error(http.StatusBadRequest, err)
	// 		return
	// 	}
	// }

	result, err = h.subjectService.RetrieveSubject(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update Subject
// @Summary Update Single Subject
// @Description Update Subject
// @Tags Subject
// @Accept       json
// @Produce      json
// @Param data body model.SubjectForm true "data"
// @Success 200 {object} model.SubjectResponseData
// @Failure 400,500 {object} model.Response
// @Router /subject/update [put]
// @Security BearerTokenAuth
// @param id query string true "id subject"
func (h subjectHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		err := errors.New("anda tidak memiliki akses untuk melakukan proses ini")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("%v", err))
		return
	}

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	var data model.Subject
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	var result model.Subject
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.subjectService.UpdateSubject(id, data)
	} else {
		result, err = h.subjectService.UpdateSubjectByOwner(id, currentUserID, data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Delete ... Delete Subject
// @Summary Delete Single Subject
// @Description Delete Subject
// @Tags Subject
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /subject/delete [put]
// @Security BearerTokenAuth
// @param id query string true "id subject"
func (h subjectHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id harus diisi dengan nomor yang valid"))
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		err := errors.New("anda tidak memiliki akses untuk melakukan proses ini")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("%v", err))
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

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List Subject
// @Summary List All Subject
// @Description List Subject
// @Tags Subject
// @Accept       json
// @Produce      json
// @Success 200 {object} model.SubjectResponseList
// @Failure 400,500 {object} model.Response
// @Router /subject/list [get]
// @Security BearerTokenAuth
func (h subjectHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Subject
	c.BindQuery(&data)

	// currentUserID, err := h.middleware.GetUserID(c)
	// if err != nil {
	// 	response.New(c).Error(http.StatusBadRequest, err)
	// 	return
	// }

	// if !h.middleware.IsSuperAdmin(c) {
	// 	data.OwnerID = currentUserID
	// }

	dataList, err := h.subjectService.ListSubject(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.subjectService.ListSubjectMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// Dropdown ... Dropdown Subject
// @Summary Dropdown All Subject
// @Description Dropdown Subject
// @Tags Subject
// @Accept       json
// @Produce      json
// @Success 200 {object} model.SubjectResponseList
// @Failure 400,500 {object} model.Response
// @Router /subject/drop-down [get]
// @Security BearerTokenAuth
func (h subjectHandler) DropDown(c *gin.Context) {
	var data model.Subject
	c.BindQuery(&data)

	// currentUserID, err := h.middleware.GetUserID(c)
	// if err != nil {
	// 	response.New(c).Error(http.StatusBadRequest, err)
	// 	return
	// }

	// if !h.middleware.IsSuperAdmin(c) {
	// 	data.OwnerID = currentUserID
	// }

	dataList, err := h.subjectService.DropDownSubject(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}
