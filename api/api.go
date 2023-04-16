package api

import (
	v1 "attendance-api/api/v1"
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/request"
	docs "attendance-api/docs"
	"attendance-api/infra"
	"attendance-api/manager"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server interface {
	Run()
}

type server struct {
	infra      infra.Infra
	gin        *gin.Engine
	service    manager.ServiceManager
	middleware middleware.Middleware
}

func NewServer(infra infra.Infra) Server {
	return &server{
		infra:      infra,
		gin:        gin.Default(),
		service:    manager.NewServiceManager(infra),
		middleware: middleware.NewMiddleware(infra.Config().GetString("secret.key")),
	}
}

func (c *server) Run() {
	docs.SwaggerInfo.BasePath = "/v1"
	c.gin.Use(c.middleware.CORS())
	c.handlers()
	c.v1()
	c.gin.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	c.gin.Run(c.infra.Port())
}

func (c *server) handlers() {
	h := request.DefaultHandler()

	c.gin.NoRoute(h.NoRoute)
	c.gin.GET("/", h.Index)
}

func (c *server) v1() {
	authHandler := v1.NewAuthHandler(c.service.AuthService(), c.service.ActivationTokenService(), c.infra)
	userHandler := v1.NewUserHandler(c.service.UserService(), c.service.ActivationTokenService(), c.infra, c.middleware)
	profileHandler := v1.NewProfileHandler(c.service.UserService(), c.service.ActivationTokenService(), c.infra, c.middleware)
	subjectHandler := v1.NewSubjectHandler(c.service.SubjectService(), c.infra, c.middleware)
	scheduleHandler := v1.NewScheduleHandler(c.service.ScheduleService(), c.service.SubjectService(), c.infra, c.middleware)
	dailyScheduleHandler := v1.NewDailyScheduleHandler(c.service.DailyScheduleService(), c.infra, c.middleware)
	userScheduleHandler := v1.NewUserScheduleHandler(c.service.UserScheduleService(), c.infra, c.middleware)
	passwordResetTokenHandler := v1.NewPasswordResetTokenHandler(c.service.PasswordResetTokenService(), c.infra, c.middleware)
	activationTokenHandler := v1.NewActivationTokenHandler(c.service.ActivationTokenService(), c.infra, c.middleware)
	attendanceHandler := v1.NewAttendanceHandler(
		c.service.AttendanceService(),
		c.service.AttendanceLogService(),
		c.service.ScheduleService(),
		c.service.UserScheduleService(),
		c.service.DailyScheduleService(),
		c.infra,
		c.middleware,
	)

	v1 := c.gin.Group("v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.GET("/activation", authHandler.Activation)
			auth.Use(c.middleware.AUTH()).PUT("/update-password", userHandler.UpdatePassword)
		}

		user := v1.Group("/user")
		user.Use(c.middleware.SUPERADMIN())
		{
			user.POST("/create", userHandler.Create)
			user.GET("/retrieve", userHandler.Retrieve)
			user.PUT("/update", userHandler.Update)
			user.DELETE("/delete", userHandler.Delete)
			user.GET("/list", userHandler.List)
			user.GET("/drop-down", userHandler.DropDown)
			user.PATCH("/active", userHandler.SetActive)
			user.PATCH("/deactive", userHandler.SetDeactive)
		}

		profile := v1.Group("/profile")
		profile.Use(c.middleware.AUTH())
		{
			profile.GET("/", profileHandler.Retrieve)
			profile.PUT("/update", profileHandler.Update)
			profile.PUT("/update-password", profileHandler.UpdatePassword)
		}

		activationToken := v1.Group("/activation-token")
		activationToken.Use(c.middleware.SUPERADMIN())
		{
			activationToken.POST("/create", activationTokenHandler.Create)
			activationToken.GET("/retrieve", activationTokenHandler.Retrieve)
			activationToken.PUT("/update", activationTokenHandler.Update)
			activationToken.DELETE("/delete", activationTokenHandler.Delete)
			activationToken.GET("/list", activationTokenHandler.List)
			activationToken.GET("/drop-down", activationTokenHandler.DropDown)
		}

		passwordResetToken := v1.Group("/password-reset-token")
		passwordResetToken.Use(c.middleware.SUPERADMIN())
		{
			passwordResetToken.POST("/create", passwordResetTokenHandler.Create)
			passwordResetToken.GET("/retrieve", passwordResetTokenHandler.Retrieve)
			passwordResetToken.PUT("/update", passwordResetTokenHandler.Update)
			passwordResetToken.DELETE("/delete", passwordResetTokenHandler.Delete)
			passwordResetToken.GET("/list", passwordResetTokenHandler.List)
			passwordResetToken.GET("/drop-down", passwordResetTokenHandler.DropDown)
		}

		subject := v1.Group("/subject")
		subject.Use(c.middleware.ADMIN())
		{
			subject.POST("/create", subjectHandler.Create)
			subject.GET("/retrieve", subjectHandler.Retrieve)
			subject.PUT("/update", subjectHandler.Update)
			subject.DELETE("/delete", subjectHandler.Delete)
			subject.GET("/list", subjectHandler.List)
			subject.GET("/drop-down", subjectHandler.DropDown)
		}

		schedule := v1.Group("/schedule")
		schedule.Use(c.middleware.ADMIN())
		{
			schedule.POST("/create", scheduleHandler.Create)
			schedule.GET("/retrieve", scheduleHandler.Retrieve)
			schedule.PUT("/update", scheduleHandler.Update)
			schedule.PUT("/update-qr-code", scheduleHandler.UpdateQRcode)
			schedule.DELETE("/delete", scheduleHandler.Delete)
			schedule.GET("/list", scheduleHandler.List)
			schedule.GET("/drop-down", scheduleHandler.DropDown)
		}

		dailySchedule := v1.Group("/daily-schedule")
		dailySchedule.Use(c.middleware.ADMIN())
		{
			dailySchedule.POST("/create", dailyScheduleHandler.Create)
			dailySchedule.GET("/retrieve", dailyScheduleHandler.Retrieve)
			dailySchedule.PUT("/update", dailyScheduleHandler.Update)
			dailySchedule.DELETE("/delete", dailyScheduleHandler.Delete)
			dailySchedule.GET("/list", dailyScheduleHandler.List)
			dailySchedule.GET("/drop-down", dailyScheduleHandler.DropDown)
		}

		userSchedule := v1.Group("/user-schedule")
		userSchedule.Use(c.middleware.ADMIN())
		{
			userSchedule.POST("/create", userScheduleHandler.Create)
			userSchedule.GET("/retrieve", userScheduleHandler.Retrieve)
			userSchedule.PUT("/update", userScheduleHandler.Update)
			userSchedule.DELETE("/delete", userScheduleHandler.Delete)
			userSchedule.GET("/list", userScheduleHandler.List)
			userSchedule.GET("/drop-down", userScheduleHandler.DropDown)
		}

		attendance := v1.Group("/attendance")
		attendance.Use(c.middleware.AUTH())
		{
			attendance.POST("/create", attendanceHandler.Create)
			attendance.GET("/retrieve", attendanceHandler.Retrieve)
			attendance.PUT("/update", attendanceHandler.Update)
			attendance.DELETE("/delete", attendanceHandler.Delete)
			attendance.GET("/list", attendanceHandler.List)
			attendance.GET("/drop-down", attendanceHandler.DropDown)
			attendance.POST("/clock-in", attendanceHandler.ClockIn)
			attendance.POST("/clock-out", attendanceHandler.ClockOut)
		}
	}

}
