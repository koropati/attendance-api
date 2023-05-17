package v1_test

import (
	v1 "attendance-api/api/v1"
	"attendance-api/common/http/middleware"
	"attendance-api/infra"
	"attendance-api/manager"
	"attendance-api/mocks"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListUser(t *testing.T) {
	t.Run("test normal case list user", func(t *testing.T) {
		userServiceMock := new(mocks.UserServiceMock)
		userServiceMock.On("ListUser", mock.AnythingOfType("model.User"), mock.AnythingOfType("model.Pagination")).Return(nil)
		userServiceMock.On("ListUserMeta", mock.AnythingOfType("model.User"), mock.AnythingOfType("model.Pagination")).Return(nil)

		activationTokenServiceMoc := new(mocks.ActivationTokenServiceMock)

		gin := gin.New()
		rec := httptest.NewRecorder()
		infra := infra.New("../../config/config.json")
		UserHandler := v1.NewUserHandler(userServiceMock, activationTokenServiceMoc, infra, middleware.NewMiddleware(infra.Config().GetString("secret.key"), manager.NewServiceManager(infra).AuthService()))
		gin.GET("/user/list", UserHandler.List)

		req := httptest.NewRequest(http.MethodGet, "/user/list", strings.NewReader(""))
		gin.ServeHTTP(rec, req)

		// convert Json Result to map string
		mapStringResult := map[string]interface{}{}
		json.Unmarshal([]byte(rec.Body.Bytes()), &mapStringResult)

		expMessage := `sukses mendapatkan list pengguna`
		expResponse := `{"code":200,"data":[{"id":1,"created_at":"2006-01-02T15:04:05-07:00","updated_at":"2006-01-02T15:04:05-07:00","deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false},"username":"windowsdewa1","password":"Password123","name":"Dewok Satria 1","handphone":"081222333441","email":"windowsdewa1.com","intro":"Hay guysss","profile":"My Name is Dewok 1","last_login":"2006-01-02T15:04:05-07:00","role":"user","is_active":true},{"id":2,"created_at":"2006-01-02T15:04:05-07:00","updated_at":"2006-01-02T15:04:05-07:00","deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false},"username":"windowsdewa2","password":"Password123","name":"Dewok Satria 2","handphone":"081222333442","email":"windowsdewa2.com","intro":"Hay guysss","profile":"My Name is Dewok 2","last_login":"2006-01-02T15:04:05-07:00","role":"user","is_active":true},{"id":3,"created_at":"2006-01-02T15:04:05-07:00","updated_at":"2006-01-02T15:04:05-07:00","deleted_at":{"Time":"0001-01-01T00:00:00Z","Valid":false},"username":"windowsdewa3","password":"Password123","name":"Dewok Satria 3","handphone":"081222333443","email":"windowsdewa3.com","intro":"Hay guysss","profile":"My Name is Dewok 3","last_login":"2006-01-02T15:04:05-07:00","role":"user","is_active":true}],"meta":{"total_page":1,"current_page":1,"total_record":3,"current_record":3},"message":"sukses mendapatkan list pengguna"}`
		// log.Printf("Response : %v", rec.Body.String())
		t.Run("test status code and response body", func(t *testing.T) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, expMessage, mapStringResult["message"])
			assert.Equal(t, expResponse, rec.Body.String())
		})
	})
}
