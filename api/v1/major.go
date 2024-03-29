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

type MajorHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
	DropDownByFaculty(c *gin.Context)
}

type majorHandler struct {
	majorService service.MajorService
	infra        infra.Infra
	middleware   middleware.Middleware
}

func NewMajorHandler(majorService service.MajorService, infra infra.Infra, middleware middleware.Middleware) MajorHandler {
	return &majorHandler{
		majorService: majorService,
		infra:        infra,
		middleware:   middleware,
	}
}

// Create ... Create Major
// @Summary Create New Major
// @Description Create Major
// @Tags Major
// @Accept       json
// @Produce      json
// @Param data body model.MajorForm true "data"
// @Success 200 {object} model.MajorResponseData
// @Failure 400,500 {object} model.Response
// @Router /major/create [post]
// @Security BearerTokenAuth
func (h majorHandler) Create(c *gin.Context) {
	var data model.Major
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

	if h.majorService.CheckIsExistByName(data.Name, 0) {
		err := errors.New("nama jurusan sudah ada")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if h.majorService.CheckIsExistByCode(data.Code, 0) {
		err := errors.New("kode jurusan sudah ada")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kode: %v", err))
		return
	}

	result, err := h.majorService.CreateMajor(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
}

// Retrieve ... Retrieve Major
// @Summary Retrieve Single Major
// @Description Retrieve Single Major
// @Tags Major
// @Accept       json
// @Produce      json
// @Success 200 {object} model.MajorResponseData
// @Failure 400,500 {object} model.Response
// @Router /major/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id major"
func (h majorHandler) Retrieve(c *gin.Context) {
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

	var result model.Major
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.majorService.RetrieveMajor(id)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.majorService.RetrieveMajorByOwner(id, currentUserID)
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update Major
// @Summary Update Single Major
// @Description Update Single Major
// @Tags Major
// @Accept       json
// @Produce      json
// @Param data body model.MajorForm true "data"
// @Success 200 {object} model.MajorResponseData
// @Failure 400,500 {object} model.Response
// @Router /major/update [put]
// @Security BearerTokenAuth
// @param id query string true "id major"
func (h majorHandler) Update(c *gin.Context) {
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

	var data model.Major
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if h.majorService.CheckIsExistByName(data.Name, id) {
		err := errors.New("nama jurusan sudah ada")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if h.majorService.CheckIsExistByCode(data.Code, id) {
		err := errors.New("kode jurusan sudah ada")
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kode: %v", err))
		return
	}

	var result model.Major
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.majorService.UpdateMajor(id, data)
	} else {
		result, err = h.majorService.UpdateMajorByOwner(id, currentUserID, data)
	}

	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Delete ... Delete Major
// @Summary Delete Single Major
// @Description Delete Single Major
// @Tags Major
// @Accept       json
// @Produce      json
// @Success 200 {object} model.MajorResponseData
// @Failure 400,500 {object} model.Response
// @Router /major/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id major"
func (h majorHandler) Delete(c *gin.Context) {
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
		if err := h.majorService.DeleteMajor(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.majorService.DeleteMajorByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List all Major
// @Summary List all Major
// @Description List all Major
// @Tags Major
// @Accept       json
// @Produce      json
// @Success 200 {object} model.MajorResponseList
// @Failure 400,500 {object} model.Response
// @Router /major/list [get]
// @Security BearerTokenAuth
func (h majorHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Major
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.majorService.ListMajor(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.majorService.ListMajorMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataList, metaList)
}

// Dropdown ... Dropdown all Major
// @Summary Dropdown all Major
// @Description Dropdown all Major
// @Tags Major
// @Accept       json
// @Produce      json
// @Success 200 {object} model.MajorResponseList
// @Failure 400,500 {object} model.Response
// @Router /major/drop-down [get]
// @Security BearerTokenAuth
func (h majorHandler) DropDown(c *gin.Context) {
	var data model.Major
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.majorService.DropDownMajor(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}

// Dropdown ... Dropdown all Major By Faculty ID
// @Summary Dropdown all Major By Faculty ID
// @Description Dropdown all Major By Faculty ID
// @Tags Major
// @Accept       json
// @Produce      json
// @Success 200 {object} model.MajorResponseList
// @Failure 400,500 {object} model.Response
// @Router /major/drop-down-by-faculty [get]
// @Security BearerTokenAuth
// @param id query string true "id faculty"
func (h majorHandler) DropDownByFaculty(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		var result []model.Major
		result = append(result, model.Major{
			GormCustom: model.GormCustom{
				ID: 0,
			},
			Name: "Pilih Fakultas Terlebih Dahulu",
		})
		response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down by faculty id", result)
	}

	dataList, err := h.majorService.DropDownByFaculty(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down by faculty id", dataList)
}
