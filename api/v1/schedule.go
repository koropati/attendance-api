package v1

import (
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/response"
	"attendance-api/common/util/myqr"
	"attendance-api/common/util/pagination"
	"attendance-api/infra"
	"attendance-api/model"
	"attendance-api/service"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type ScheduleHandler interface {
	Create(c *gin.Context)
	Retrieve(c *gin.Context)
	Update(c *gin.Context)
	UpdateQRcode(c *gin.Context)
	Delete(c *gin.Context)
	List(c *gin.Context)
	DropDown(c *gin.Context)
}

type scheduleHandler struct {
	scheduleService     service.ScheduleService
	subjectService      service.SubjectService
	userScheduleService service.UserScheduleService
	infra               infra.Infra
	middleware          middleware.Middleware
}

func NewScheduleHandler(
	scheduleService service.ScheduleService,
	subjectService service.SubjectService,
	userScheduleService service.UserScheduleService,
	infra infra.Infra,
	middleware middleware.Middleware,
) ScheduleHandler {
	return &scheduleHandler{
		scheduleService:     scheduleService,
		subjectService:      subjectService,
		userScheduleService: userScheduleService,
		infra:               infra,
		middleware:          middleware,
	}
}

// Create ... Create Schedule
// @Summary Create New Schedule
// @Description Create Schedule
// @Tags Schedule
// @Accept       json
// @Produce      json
// @Param data body model.ScheduleForm true "data"
// @Success 200 {object} model.ScheduleResponseData
// @Failure 400,500 {object} model.Response
// @Router /schedule/create [post]
// @Security BearerTokenAuth
func (h scheduleHandler) Create(c *gin.Context) {
	var data model.Schedule
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

	data.Code = strings.ToUpper(data.Code)

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if err := validation.Validate(data.Code, validation.Required, validation.Length(1, 100)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kode: %v", err))
		return
	}

	if err := validation.Validate(data.SubjectID, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id subjek: %v", "pilihan subjek harus terisi"))
		return
	}

	if exist := h.scheduleService.CheckCodeIsExist(data.Code, 0); exist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kode: %v", "kode tersebut sudah ada yang menggunakan"))
		return
	}

	if isExist := h.subjectService.CheckIsExist(int(data.SubjectID)); !isExist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id subjek: %v", "id subjek tidak tersedia"))
		return
	}

	data.QRCode = myqr.Generate(data.Code, 8)

	result, err := h.scheduleService.CreateSchedule(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	result, err = h.scheduleService.RetrieveSchedule(int(result.ID))
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))

	response.New(c).Data(http.StatusCreated, "sukses membuat data", result)
}

// Retrieve ... Retrieve Schedule
// @Summary Retrieve Single Schedule
// @Description Retrieve Schedule
// @Tags Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.ScheduleResponseData
// @Failure 400,500 {object} model.Response
// @Router /schedule/retrieve [get]
// @Security BearerTokenAuth
// @param id query string true "id schedule"
func (h scheduleHandler) Retrieve(c *gin.Context) {
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

	var result model.Schedule
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.scheduleService.RetrieveSchedule(id)
		result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.scheduleService.RetrieveScheduleByOwner(id, currentUserID)
		result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "sukses mengambil data", result)
}

// Update ... Update Schedule
// @Summary Update Single Schedule
// @Description Update Schedule
// @Tags Schedule
// @Accept       json
// @Produce      json
// @Param data body model.ScheduleForm true "data"
// @Success 200 {object} model.ScheduleResponseData
// @Failure 400,500 {object} model.Response
// @Router /schedule/update [put]
// @Security BearerTokenAuth
// @param id query string true "id schedule"
func (h scheduleHandler) Update(c *gin.Context) {
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

	var data model.Schedule
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama: %v", err))
		return
	}

	if err := validation.Validate(data.Code, validation.Required, validation.Length(1, 100)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kode: %v", err))
		return
	}

	if err := validation.Validate(data.SubjectID, validation.Required, validation.Min(1), is.Int); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id subjek: %v", "pilihan subjek harus terisi"))
		return
	}

	if exist := h.scheduleService.CheckCodeIsExist(data.Code, int(data.ID)); exist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kode: %v", "kode tersebut sudah ada yang menggunakan"))
		return
	}

	if isExist := h.subjectService.CheckIsExist(int(data.SubjectID)); !isExist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("id subjek: %v", "id subjek tidak tersedia"))
		return
	}

	var result model.Schedule
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.scheduleService.UpdateSchedule(id, data)
		result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
		result, err = h.scheduleService.RetrieveSchedule(int(result.ID))
		result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.scheduleService.UpdateScheduleByOwner(id, currentUserID, data)
		result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
		result, err = h.scheduleService.RetrieveScheduleByOwner(int(result.ID), currentUserID)
		result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Data(http.StatusOK, "sukses memperbaharui data", result)
}

// Update QR Code ... Update QR Code
// @Summary Update QR Code
// @Description Update QR Code
// @Tags Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.ScheduleResponseData
// @Failure 400,500 {object} model.Response
// @Router /schedule/update-qr-code [put]
// @Security BearerTokenAuth
// @param id query string true "id schedule"
func (h scheduleHandler) UpdateQRcode(c *gin.Context) {
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

	schedule, err := h.scheduleService.RetrieveSchedule(id)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	qrCode := myqr.Generate(schedule.Code, 8)

	var result model.Schedule
	if h.middleware.IsSuperAdmin(c) {
		result, err = h.scheduleService.UpdateQRcode(id, qrCode)
		result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		result, err = h.scheduleService.UpdateQRcodeByOwner(id, currentUserID, qrCode)
		result.UserInRule = h.userScheduleService.CountByScheduleID(int(result.ID))
		if err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}
	response.New(c).Data(http.StatusCreated, "berhasil memperbaharui kode qr", result)
}

// Delete ... Delete Schedule
// @Summary Delete Schedule
// @Description Delete Schedule
// @Tags Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /schedule/delete [delete]
// @Security BearerTokenAuth
// @param id query string true "id schedule"
func (h scheduleHandler) Delete(c *gin.Context) {
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
		if err := h.scheduleService.DeleteSchedule(id); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	} else {
		if err := h.scheduleService.DeleteScheduleByOwner(id, currentUserID); err != nil {
			response.New(c).Error(http.StatusBadRequest, err)
			return
		}
	}

	response.New(c).Write(http.StatusOK, "sukses menghapus data")
}

// List ... List Schedule
// @Summary List Schedule
// @Description List Schedule
// @Tags Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.ScheduleResponseList
// @Failure 400,500 {object} model.Response
// @Router /schedule/list [get]
// @Security BearerTokenAuth
func (h scheduleHandler) List(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var data model.Schedule
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataLists, err := h.scheduleService.ListSchedule(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	for i, dataList := range dataLists {
		dataLists[i].UserInRule = h.userScheduleService.CountByScheduleID(int(dataList.ID))
	}

	metaLists, err := h.scheduleService.ListScheduleMeta(data, pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).List(http.StatusOK, "sukses mengambil list data", dataLists, metaLists)
}

// Dropdown ... Dropdown Schedule
// @Summary Dropdown Schedule
// @Description Dropdown Schedule
// @Tags Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.ScheduleResponseList
// @Failure 400,500 {object} model.Response
// @Router /schedule/drop-down [get]
// @Security BearerTokenAuth
func (h scheduleHandler) DropDown(c *gin.Context) {
	var data model.Schedule
	c.BindQuery(&data)

	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if !h.middleware.IsSuperAdmin(c) {
		data.OwnerID = currentUserID
	}

	dataList, err := h.scheduleService.DropDownSchedule(data)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mendapatkan data drop down", dataList)
}
