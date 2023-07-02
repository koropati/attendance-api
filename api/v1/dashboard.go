package v1

import (
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/response"
	"attendance-api/infra"
	"attendance-api/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler interface {
	GetDashboardAcademic(c *gin.Context)
	GetDashboardUser(c *gin.Context)
	GetDashboardStudent(c *gin.Context)
	GetDashboardTeacher(c *gin.Context)
}

type dashboardHandler struct {
	dashboardService service.DashboardService
	infra            infra.Infra
	middleware       middleware.Middleware
}

func NewDashboardHandler(dashboardService service.DashboardService, infra infra.Infra, middleware middleware.Middleware) DashboardHandler {
	return &dashboardHandler{
		dashboardService: dashboardService,
		infra:            infra,
		middleware:       middleware,
	}
}

// Retrieve Dashboard Academic ... Retrieve Dashboard Academic
// @Summary Retrieve Dashboard Academic
// @Description Retrieve Dashboard Academic
// @Tags Dashboard
// @Accept       json
// @Produce      json
// @Success 200 {object} model.DashboardAcademicResponseData
// @Failure 400,500 {object} model.Response
// @Router /dashboard/academic [get]
// @Security BearerTokenAuth
func (h dashboardHandler) GetDashboardAcademic(c *gin.Context) {
	result, _ := h.dashboardService.RetrieveDashboardAcademic()

	response.New(c).Data(http.StatusOK, "sukses mengambil dashboard akademik", result)
}

// Retrieve Dashboard Users ... Retrieve Dashboard Users
// @Summary Retrieve Dashboard Users
// @Description Retrieve Dashboard Users
// @Tags Dashboard
// @Accept       json
// @Produce      json
// @Success 200 {object} model.DashboardUserResponseData
// @Failure 400,500 {object} model.Response
// @Router /dashboard/user [get]
// @Security BearerTokenAuth
func (h dashboardHandler) GetDashboardUser(c *gin.Context) {
	result, err := h.dashboardService.RetrieveDashboardUser()
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Data(http.StatusOK, "sukses mengambil dashboard pengguna", result)
}

// Retrieve Dashboard Student ... Retrieve Dashboard Student
// @Summary Retrieve Dashboard Student
// @Description Retrieve Dashboard Student
// @Tags Dashboard
// @Accept       json
// @Produce      json
// @Success 200 {object} model.DashboardStudentResponseData
// @Failure 400,500 {object} model.Response
// @Router /dashboard/student [get]
// @Security BearerTokenAuth
func (h dashboardHandler) GetDashboardStudent(c *gin.Context) {
	result, err := h.dashboardService.RetrieveDashboardStudent()
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Data(http.StatusOK, "sukses mengambil dashboard mahasiswa", result)
}

// Retrieve Dashboard Teacher ... Retrieve Dashboard Teacher
// @Summary Retrieve Dashboard Teacher
// @Description Retrieve Dashboard Teacher
// @Tags Dashboard
// @Accept       json
// @Produce      json
// @Success 200 {object} model.DashboardTeacherResponseData
// @Failure 400,500 {object} model.Response
// @Router /dashboard/teacher [get]
// @Security BearerTokenAuth
func (h dashboardHandler) GetDashboardTeacher(c *gin.Context) {
	result, err := h.dashboardService.RetrieveDashboardTeacher()
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Data(http.StatusOK, "sukses mengambil dashboard dosen", result)
}
