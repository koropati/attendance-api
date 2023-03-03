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
	authHandler := v1.NewAuthHandler(c.service.AuthService(), c.infra)
	userHandler := v1.NewUserHandler(c.service.UserService(), c.infra, c.middleware)
	subjectHandler := v1.NewSubjectHandler(c.service.SubjectService(), c.infra, c.middleware)

	v1 := c.gin.Group("v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
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
	}

}
