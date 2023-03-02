package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"attendance-api/common/http/email"
	"attendance-api/common/http/response"
	"attendance-api/common/util/activation"
	"attendance-api/common/util/cryptos"
	"attendance-api/common/util/regex"
	"attendance-api/common/util/token"
	"attendance-api/infra"
	"attendance-api/model"
	"attendance-api/service"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserHandler interface {
	Register(c *gin.Context)
	Refresh(c *gin.Context)
	Login(c *gin.Context)
	Create(c *gin.Context)
	Delete(c *gin.Context)
}

type authUserHandler struct {
	authService service.AuthService
	infra       infra.Infra
}

func NewAuthHandler(authService service.AuthService, infra infra.Infra) AuthUserHandler {
	return &authUserHandler{
		authService: authService,
		infra:       infra,
	}
}

// Register ... Register User
// @Summary Create new user based on paramters
// @Description Create new user
// @Tags Auth
// @Accept json
// @Param user body model.UserForm true "User Data"
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /auth/register [post]
func (h *authUserHandler) Register(c *gin.Context) {
	var data model.User
	c.BindJSON(&data)

	if err := validation.Validate(data.Username, validation.Required, validation.Length(4, 30), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := validation.Validate(data.Password, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("password: %v", err))
		return
	}

	if err := validation.Validate(data.Email, validation.Required, validation.Length(6, 50)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("email: %v", err))
		return
	}

	if err := validation.Validate(data.FirstName, validation.Required, validation.Match(regexp.MustCompile(regex.NAME))); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("first_name: %v", err))
		return
	}

	if !h.authService.CheckHandphone(data.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("handphone: already taken"))
	}

	if !h.authService.CheckEmail(data.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email: already taken"))
	}

	if h.authService.CheckUsername(data.Username) {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("password: %v", err))
			return
		}

		data.Password = string(password)
		if err := h.authService.Register(&data); err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		user, err := h.authService.GetByUsername(data.Username)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		activationString := activation.New(*user).Generate(24)
		// log.Printf("activationString: %v\n", activationString)
		activationToken := cryptos.New(h.infra.Cipher("encrypt")).EncryptAES256(activationString)
		// log.Printf("activationToken: %v\n", activationToken)

		go func(user *model.User) {
			config := h.infra.Config().Sub("server")
			urlActivation := fmt.Sprintf("%s:%s/auth/activation?code=%s", config.GetString("url"), config.GetString("port"), activationToken)

			if err := email.New(h.infra.GoMail(), h.infra.Config()).SendActivation(user.FirstName, user.Email, urlActivation); err != nil {
				log.Printf("Error Send Email E: %v", err)
			}
		}(user)

		response.New(c).Write(http.StatusCreated, "success: user registered")
		return
	}

	response.New(c).Error(http.StatusBadRequest, errors.New("username: already taken"))
}

// Login ... Login User
// @Summary Login user with username and password
// @Description Login User
// @Tags Auth
// @Accept json
// @Param data body model.Login true "Login Data"
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /auth/login [post]
func (h *authUserHandler) Login(c *gin.Context) {
	var data model.Login
	c.BindJSON(&data)

	if err := validation.Validate(data.Username, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := validation.Validate(data.Password, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("password: %v", err))
		return
	}

	if isActive := h.authService.CheckIsActive(data.Username); !isActive {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("account is not activated"))
		return
	}

	hashedPassword, err := h.authService.Login(data.Username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(data.Password)); err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("username or password not match"))
		return
	}

	isSuper, isAdmin, isUser, err := h.authService.GetRole(data.Username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("err: %v", err))
		return
	}

	userData, err := h.authService.GetByUsername(data.Username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("err: %v", err))
		return
	}

	email, err := h.authService.GetEmail(data.Username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("email: %v", err))
		return
	}

	expired, accessToken := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateToken(
		&model.UserTokenPayload{
			UserID:       userData.ID,
			Username:     data.Username,
			IsSuperAdmin: isSuper,
			IsAdmin:      isAdmin,
			IsUser:       isUser,
			Email:        email,
			Expired:      h.infra.Config().GetInt("access_token_expired"),
		},
	)

	refreshExpired, refreshToken := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateRefreshToken(
		&model.UserTokenPayload{
			Username: data.Username,
			Expired:  h.infra.Config().GetInt("refresh_token_expired"),
		},
	)
	dataOutput := map[string]interface{}{
		"access_token":          accessToken,
		"expired_access_token":  expired,
		"refresh_token":         refreshToken,
		"expired_refresh_token": refreshExpired,
	}
	response.New(c).Data(200, "success login", dataOutput)
}

// Refresh ... Refresh Token
// @Summary Get New Access Token using refresh token
// @Description Get New Access Token
// @Tags Auth
// @Accept json
// @Param data body model.Refresh true "Refresh Data"
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /auth/refresh [post]
func (h *authUserHandler) Refresh(c *gin.Context) {
	var data model.Refresh
	c.BindJSON(&data)
	dataToken, err := token.NewToken(h.infra.Config().GetString("secret.key")).ValidateRefreshToken(data.RefreshToken)
	if err != nil {
		response.New(c).Error(http.StatusUnauthorized, err)
		return
	}
	claims, ok := dataToken.Claims.(jwt.MapClaims)
	if !ok {
		response.New(c).Error(http.StatusUnauthorized, err)
		return
	}
	username := claims["username"].(string)
	user, err := h.authService.GetByUsername(username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	expired, accessToken := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateToken(
		&model.UserTokenPayload{
			UserID:       user.ID,
			Username:     user.Username,
			IsSuperAdmin: user.IsSuperAdmin,
			IsAdmin:      user.IsAdmin,
			IsUser:       user.IsUser,
			Email:        user.Email,
			Expired:      h.infra.Config().GetInt("access_token_expired"),
		},
	)
	refreshExpired, refreshToken := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateRefreshToken(
		&model.UserTokenPayload{
			UserID:   user.ID,
			Username: user.Username,
			Expired:  h.infra.Config().GetInt("refresh_token_expired"),
		},
	)
	dataOutput := map[string]interface{}{
		"access_token":          accessToken,
		"expired_access_token":  expired,
		"refresh_token":         refreshToken,
		"expired_refresh_token": refreshExpired,
	}
	response.New(c).Data(200, "success refresh", dataOutput)
}

func (h *authUserHandler) Create(c *gin.Context) {
	var data model.User
	c.BindJSON(&data)

	if err := validation.Validate(data.Username, validation.Required, validation.Length(4, 30), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("username: %v", err))
		return
	}

	if err := validation.Validate(data.Password, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("password: %v", err))
		return
	}

	if err := validation.Validate(data.Email, validation.Required, validation.Length(6, 50)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("email: %v", err))
		return
	}

	if err := validation.Validate(data.FirstName, validation.Required, validation.Match(regexp.MustCompile(regex.NAME))); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("first_name: %v", err))
		return
	}

	if !h.authService.CheckHandphone(data.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("handphone: already taken"))
	}

	if !h.authService.CheckEmail(data.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email: already taken"))
	}

	if h.authService.CheckUsername(data.Username) {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("password: %v", err))
			return
		}

		data.Password = string(password)
		if err := h.authService.Create(&data); err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		response.New(c).Write(http.StatusCreated, "success: user created")
		return
	}

	response.New(c).Error(http.StatusBadRequest, errors.New("username: already taken"))
}

func (h *authUserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("id must be filled and valid number"))
		return
	}

	if h.authService.CheckID(id) {
		if err := h.authService.Delete(id); err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		response.New(c).Write(http.StatusOK, "success: user deleted")
		return
	}

	response.New(c).Error(http.StatusBadRequest, errors.New("id: not found"))
}
