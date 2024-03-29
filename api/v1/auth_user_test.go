package v1_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	v1 "attendance-api/api/v1"
	"attendance-api/infra"
	"attendance-api/mocks"
	"attendance-api/model"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockUser = model.User{
	Email:     "dewok@gmail.com",
	Username:  "dewok",
	Password:  "password",
	FirstName: "Dewok",
	LastName:  "Satria",
	Handphone: "082340803646",
}

func TestRegister(t *testing.T) {
	t.Run("test normal case register", func(t *testing.T) {
		authServiceMock := new(mocks.AuthServiceMock)
		authServiceMock.On("CheckUsername", mock.AnythingOfType("string")).Return(nil)
		authServiceMock.On("CheckHandphone", mock.AnythingOfType("string")).Return(nil)
		authServiceMock.On("CheckEmail", mock.AnythingOfType("string")).Return(nil)
		authServiceMock.On("Register", mock.AnythingOfType("model.User")).Return(nil)

		activationTokenServiceMoc := new(mocks.ActivationTokenServiceMock)
		passwordResetTokenServiceMoc := new(mocks.PasswordResetTokenServiceMock)
		userServiceMoc := new(mocks.UserServiceMock)

		gin := gin.New()
		rec := httptest.NewRecorder()

		authHandler := v1.NewAuthHandler(authServiceMock, userServiceMoc, activationTokenServiceMoc, passwordResetTokenServiceMoc, infra.New("../../config/config.json"))
		gin.POST("/register", authHandler.Register)

		body, err := json.Marshal(mockUser)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(string(body)))
		gin.ServeHTTP(rec, req)

		exp := `{"code":201,"message":"berhasil menambah data"}`

		t.Run("test status code and response body", func(t *testing.T) {
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Equal(t, exp, rec.Body.String())
		})
	})
}

func TestLogin(t *testing.T) {
	t.Run("test normal case login", func(t *testing.T) {
		authServiceMock := new(mocks.AuthServiceMock)
		authServiceMock.On("Login", mock.AnythingOfType("string")).Return(nil)

		activationTokenServiceMoc := new(mocks.ActivationTokenServiceMock)
		passwordResetTokenServiceMoc := new(mocks.PasswordResetTokenServiceMock)
		userServiceMoc := new(mocks.UserServiceMock)

		gin := gin.New()
		rec := httptest.NewRecorder()

		authHandler := v1.NewAuthHandler(authServiceMock, userServiceMoc, activationTokenServiceMoc, passwordResetTokenServiceMoc, infra.New("../../config/config.json"))
		gin.POST("/login", authHandler.Login)

		body, err := json.Marshal(mockUser)
		assert.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(string(body)))
		gin.ServeHTTP(rec, req)

		var mockResponse model.Token
		err = json.Unmarshal(rec.Body.Bytes(), &mockResponse)
		assert.NoError(t, err)

		exp := string(time.Now().Add(time.Hour * 2).Format(time.RFC3339))

		t.Run("test status code and token expiration", func(t *testing.T) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, exp, mockResponse.Expired)
		})
	})
}
