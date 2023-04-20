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

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if err := validation.Validate(data.Code, validation.Required, validation.Length(1, 100), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", err))
		return
	}

	if err := validation.Validate(data.SubjectID, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("subject_id: %v", "subject choice must be filled in"))
		return
	}

	if exist := h.scheduleService.CheckCodeIsExist(data.Code, 0); exist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", "the code is already in use"))
		return
	}

	if isExist := h.subjectService.CheckIsExist(int(data.SubjectID)); !isExist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("subject_id: %v", "subject id is not exist"))
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

	response.New(c).Data(http.StatusCreated, "success create data", result)
}

func (h scheduleHandler) Retrieve(c *gin.Context) {
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
	response.New(c).Data(http.StatusCreated, "success retrieve data", result)
}

func (h scheduleHandler) Update(c *gin.Context) {
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

	var data model.Schedule
	c.BindJSON(&data)

	data.UpdatedBy = currentUserID
	data.UpdatedAt = time.Now()

	if err := validation.Validate(data.Name, validation.Required, validation.Length(1, 255), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("name: %v", err))
		return
	}

	if err := validation.Validate(data.Code, validation.Required, validation.Length(1, 100), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", err))
		return
	}

	if err := validation.Validate(data.SubjectID, validation.Required, validation.Min(1), is.Int); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", "subject choice must be filled in"))
		return
	}

	if exist := h.scheduleService.CheckCodeIsExist(data.Code, int(data.ID)); exist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("code: %v", "the code is already in use"))
		return
	}

	if isExist := h.subjectService.CheckIsExist(int(data.SubjectID)); !isExist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("subject_id: %v", "subject id is not exist"))
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

	response.New(c).Data(http.StatusOK, "success update data", result)
}

func (h scheduleHandler) UpdateQRcode(c *gin.Context) {
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
	response.New(c).Data(http.StatusCreated, "success update qr code data", result)
}

func (h scheduleHandler) Delete(c *gin.Context) {
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

	response.New(c).Write(http.StatusOK, "success delete data")
}

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

	response.New(c).List(http.StatusOK, "success get list data", dataLists, metaLists)
}

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

	response.New(c).Data(http.StatusOK, "success get drop down data", dataList)
}
