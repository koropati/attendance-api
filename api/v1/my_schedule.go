package v1

import (
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/response"
	"attendance-api/common/util/converter"
	"attendance-api/infra"
	"attendance-api/model"
	"attendance-api/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type MyScheduleHandler interface {
	List(c *gin.Context)
	Today(c *gin.Context)
}

type myScheduleHandler struct {
	userScheduleService service.UserScheduleService
	attendanceService   service.AttendanceService
	infra               infra.Infra
	middleware          middleware.Middleware
}

func NewMyScheduleHandler(userScheduleService service.UserScheduleService, attendanceService service.AttendanceService, infra infra.Infra, middleware middleware.Middleware) MyScheduleHandler {
	return &myScheduleHandler{
		userScheduleService: userScheduleService,
		attendanceService:   attendanceService,
		infra:               infra,
		middleware:          middleware,
	}
}

// List ... My List Schedule
// @Summary My List Schedule
// @Description My List Schedule
// @Tags My Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.MyScheduleResponseList
// @Failure 400,500 {object} model.Response
// @Router /my-schedule/list [get]
// @Security BearerTokenAuth
func (h myScheduleHandler) List(c *gin.Context) {
	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	var filter model.MyScheduleFilter
	errFilter := c.BindQuery(&filter)
	if errFilter != nil {
		response.New(c).Error(http.StatusBadRequest, errFilter)
	}

	results, err := h.userScheduleService.ListMySchedule(currentUserID, filter)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	response.New(c).Data(http.StatusOK, "sukses mengambil data jadwal", results)
}

// Today ... List Today Schedule
// @Summary List Today Schedule
// @Description List Today Schedule
// @Tags My Schedule
// @Accept       json
// @Produce      json
// @Success 200 {object} model.TodayScheduleResponseList
// @Failure 400,500 {object} model.Response
// @Router /my-schedule/today [get]
// @Security BearerTokenAuth
func (h myScheduleHandler) Today(c *gin.Context) {
	currentUserID, err := h.middleware.GetUserID(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	dayName := converter.GetDayName(time.Now())
	todayDate := time.Now().Format("2006-01-02")

	results, err := h.userScheduleService.ListTodaySchedule(currentUserID, dayName)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	for i, result := range results {
		if isExistAttendance := h.attendanceService.CheckIsExistByDate(currentUserID, int(result.ScheduleID), todayDate); isExistAttendance {
			attendance, err := h.attendanceService.RetrieveAttendanceByDate(currentUserID, int(result.ScheduleID), todayDate)
			if err != nil {
				log.Printf("[Err] RetrieveAttendanceByDate E: %v\n", err)
				results[i].AttendanceID = 0
				results[i].ClockInMillis = 0
				results[i].ClockOutMillis = 0
				results[i].ClockIn = "--:--"
				results[i].ClockOut = "--:--"
				results[i].TimeZoneIn = 0
				results[i].TimeZoneOut = 0
				results[i].LocationIn = "-"
				results[i].LocationOut = "-"
			} else {
				results[i].AttendanceID = attendance.ID
				results[i].ClockInMillis = attendance.ClockIn
				results[i].ClockOutMillis = attendance.ClockOut
				results[i].ClockIn = converter.MillisToTimeString(attendance.ClockIn, attendance.TimeZoneIn)
				results[i].ClockOut = converter.MillisToTimeString(attendance.ClockOut, attendance.TimeZoneOut)
				results[i].TimeZoneIn = attendance.TimeZoneIn
				results[i].TimeZoneOut = attendance.TimeZoneOut
				results[i].LocationIn = attendance.LocationIn
				results[i].LocationOut = attendance.LocationOut
			}
		} else {
			results[i].AttendanceID = 0
			results[i].ClockInMillis = 0
			results[i].ClockOutMillis = 0
			results[i].ClockIn = "--:--"
			results[i].ClockOut = "--:--"
			results[i].TimeZoneIn = 0
			results[i].TimeZoneOut = 0
			results[i].LocationIn = "-"
			results[i].LocationOut = "-"
		}

	}

	response.New(c).Data(http.StatusOK, "sukses mengambil data jadwal", results)
}
