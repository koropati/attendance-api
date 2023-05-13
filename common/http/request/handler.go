package request

import (
	"errors"
	"net/http"

	"attendance-api/common/http/response"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	NoRoute(c *gin.Context)
	Index(c *gin.Context)
}

type handler struct {
	// Stuff maybe needed for handler
}

func DefaultHandler() Handler {
	return &handler{}
}

func (h *handler) NoRoute(c *gin.Context) {
	response.New(c).Error(http.StatusNotFound, errors.New("rute tidak di temukan"))
}

func (h *handler) Index(c *gin.Context) {
	response.New(c).Write(http.StatusOK, "aplikasi berjalan")
}
