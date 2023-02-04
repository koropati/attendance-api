package api

import (
	v1 "attendance-api/api/v1"
	"attendance-api/common/http/middleware"
	"attendance-api/common/http/request"
	"attendance-api/infra"
	"attendance-api/manager"

	"github.com/gin-gonic/gin"
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
	c.gin.Use(c.middleware.CORS())
	c.handlers()
	c.v1()

	c.gin.Run(c.infra.Port())
}

func (c *server) handlers() {
	h := request.DefaultHandler()

	c.gin.NoRoute(h.NoRoute)
	c.gin.GET("/", h.Index)
}

func (c *server) v1() {
	authHandler := v1.NewAuthHandler(c.service.AuthService(), c.infra)
	userHandler := v1.NewUserHandler(c.service.UserService(), c.infra)

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
			user.GET("/list", userHandler.ListUser)
			user.POST("/create", authHandler.Create)
			user.DELETE("/delete", authHandler.Delete)
		}
	}
}
