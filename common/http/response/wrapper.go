package response

import (
	"net/http"

	"attendance-api/model"

	"github.com/gin-gonic/gin"
)

type Wrapper interface {
	Write(code int, message string)
	Data(code int, message string, data interface{})
	List(code int, message string, data interface{}, meta interface{})
	Error(code int, err error)
	Token(expired string, token string)
}

type wrapper struct {
	c *gin.Context
}

func New(c *gin.Context) Wrapper {
	return &wrapper{c: c}
}

func (w *wrapper) Write(code int, message string) {
	w.c.JSON(code, model.Response{Code: code, Message: message})
}

func (w *wrapper) Data(code int, message string, data interface{}) {
	w.c.JSON(code, model.ResponseList{Code: code, Data: data, Message: message})
}

func (w *wrapper) Error(code int, err error) {
	w.c.JSON(code, model.Response{Code: code, Message: err.Error()})
}

func (w *wrapper) Token(expired string, token string) {
	w.c.JSON(http.StatusOK, model.Token{Expired: expired, Token: token})
}

func (w *wrapper) List(code int, message string, data interface{}, meta interface{}) {
	w.c.JSON(code, model.ResponseList{Code: code, Data: data, Meta: meta, Message: message})
}
