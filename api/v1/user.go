package v1

import (
	"attendance-api/common/http/email"
	"attendance-api/common/http/response"
	"attendance-api/common/util/pagination"
	"attendance-api/infra"
	"attendance-api/model"
	"attendance-api/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	ListUser(c *gin.Context)
}

type userHandler struct {
	userService service.UserService
	infra       infra.Infra
}

func NewUserHandler(userService service.UserService, infra infra.Infra) UserHandler {
	return &userHandler{
		userService: userService,
		infra:       infra,
	}
}

func (h *userHandler) ListUser(c *gin.Context) {
	pagination := pagination.GeneratePaginationFromRequest(c)
	var user model.User

	userList, err := h.userService.ListUser(&user, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	metaList, err := h.userService.ListUserMeta(&user, &pagination)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
	}

	// if err := email.New(h.infra.GoMail()).SendActivation("windowsdewa@gmail.com", "xarafaefdfas"); err != nil {
	// 	log.Printf("Error Send Email E: %v", err)
	// }

	if err := email.NewSendGrid(h.infra.SendGrid(), h.infra.Config()).SendActivation("Dewok Satria", "windowsdewa@gmail.com", "ahikjdfahsjdha"); err != nil {
		log.Printf("Error Send Email E: %v", err)
	}

	response.New(c).List(http.StatusOK, "success get list user", userList, metaList)
}
