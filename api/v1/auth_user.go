package v1

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"attendance-api/common/http/email"
	"attendance-api/common/http/response"
	"attendance-api/common/util/activation"
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
	Activation(c *gin.Context)
	Logout(c *gin.Context)
}

type authUserHandler struct {
	authService            service.AuthService
	activationTokenService service.ActivationTokenService
	infra                  infra.Infra
}

func NewAuthHandler(authService service.AuthService, activationTokenService service.ActivationTokenService, infra infra.Infra) AuthUserHandler {
	return &authUserHandler{
		authService:            authService,
		activationTokenService: activationTokenService,
		infra:                  infra,
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
func (h authUserHandler) Register(c *gin.Context) {
	var data model.User
	c.BindJSON(&data)

	if err := validation.Validate(data.Username, validation.Required, validation.Length(4, 30), is.Alphanumeric); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama pengguna: %v", err))
		return
	}

	if err := validation.Validate(data.Password, validation.Required, validation.Length(6, 40)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kata sandi: %v", err))
		return
	}

	if err := validation.Validate(data.Email, validation.Required, validation.Length(6, 50)); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("email: %v", err))
		return
	}

	if err := validation.Validate(data.FirstName, validation.Required, validation.Match(regexp.MustCompile(regex.NAME))); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama depan: %v", err))
		return
	}

	if !h.authService.CheckHandphone(data.Handphone) {
		response.New(c).Error(http.StatusBadRequest, errors.New("no telp sudah digunakan"))
	}

	if !h.authService.CheckEmail(data.Email) {
		response.New(c).Error(http.StatusBadRequest, errors.New("email sudah digunakan"))
	}

	if h.authService.CheckUsername(data.Username) {
		password, err := bcrypt.GenerateFromPassword([]byte(data.Password), 10)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, fmt.Errorf("kata sandi: %v", err))
			return
		}

		data.Password = string(password)
		if err := h.authService.Register(data); err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		user, err := h.authService.GetByUsername(data.Username)
		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		expiredToken, activationToken := activation.New(user).GenerateSHA1(24)

		// Save Activation token to data base
		activationData, err := h.activationTokenService.CreateActivationToken(model.ActivationToken{
			UserID: user.ID,
			Token:  activationToken,
			Valid:  expiredToken,
		})

		if err != nil {
			response.New(c).Error(http.StatusInternalServerError, err)
			return
		}

		go func(user model.User) {
			config := h.infra.Config().Sub("server")
			urlActivation := fmt.Sprintf("%s/v1/auth/activation?token=%s", config.GetString("base_url"), activationData.Token)

			if err := email.New(h.infra.GoMail(), h.infra.Config()).SendActivation(user.FirstName, user.Email, urlActivation); err != nil {
				log.Printf("Error Send Email E: %v", err)
			}
		}(user)

		response.New(c).Write(http.StatusCreated, "berhasil registrasi pengguna")
		return
	}

	response.New(c).Error(http.StatusBadRequest, errors.New("nama pengguna sudah digunakan"))
}

// Login ... Login User
// @Summary Login user with username and password
// @Description Login User
// @Tags Auth
// @Accept json
// @Param data body model.Login true "Login Data"
// @Success 200 {object} model.AuthDataResponseData
// @Failure 400,500 {object} model.Response
// @Router /auth/login [post]
func (h authUserHandler) Login(c *gin.Context) {
	var data model.Login
	c.BindJSON(&data)

	if err := validation.Validate(data.Username, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama pengguna: %v", err))
		return
	}

	if err := validation.Validate(data.Password, validation.Required); err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("kata sandi: %v", err))
		return
	}

	if isExist := h.authService.CheckUsername(data.Username); !isExist {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama pengguna atau kata sandi salah"))
		return
	}

	if isActive := h.authService.CheckIsActive(data.Username); !isActive {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("akun tidak aktif"))
		return
	}

	hashedPassword, err := h.authService.Login(data.Username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("nama pengguna: %v", err))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(data.Password)); err != nil {
		response.New(c).Error(http.StatusBadRequest, errors.New("nama pengguna atau kata sandi salah"))
		return
	}

	userData, err := h.authService.GetByUsername(data.Username)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("maaf ada error : %v", err))
		return
	}

	expiredTimeAT := time.Now().Add(time.Minute * time.Duration(h.infra.Config().GetInt("access_token_expired"))).Unix()
	authAT, err := h.authService.CreateAuth(userData.ID, expiredTimeAT, "at")
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("error autentikasi: %v", err))
		return
	}

	expired, accessToken := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateToken(
		model.UserTokenPayload{
			UserID:   authAT.UserID,
			AuthUUID: authAT.AuthUUID,
			Expired:  authAT.Expired,
		},
	)

	expiredTimeRT := time.Now().Add(time.Minute * time.Duration(h.infra.Config().GetInt("refresh_token_expired"))).Unix()
	authRT, err := h.authService.CreateAuth(userData.ID, expiredTimeRT, "rt")
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("error autentikasi: %v", err))
		return
	}

	refreshExpired, refreshToken := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateRefreshToken(
		model.UserTokenPayload{
			UserID:   authRT.UserID,
			AuthUUID: authRT.AuthUUID,
			Expired:  authRT.Expired,
		},
	)

	userData.Role = userData.GetRole()
	userData.UserAbilities = userData.GetAbility()
	userData.Avatar = userData.GetAvatar()

	tokenData := model.TokenData{
		AccessToken:         accessToken,
		ExpiredAccessToken:  expired,
		RefreshToken:        refreshToken,
		ExpiredRefreshToken: refreshExpired,
	}
	dataOutput := model.AuthData{
		UserData:  userData,
		TokenData: tokenData,
	}
	response.New(c).Data(200, "berhasil masuk ke dalam sistem", dataOutput)
}

// Refresh ... Refresh Token
// @Summary Get New Access Token using refresh token
// @Description Get New Access Token
// @Tags Auth
// @Accept json
// @Param data body model.Refresh true "Refresh Data"
// @Success 200 {object} model.AuthDataResponseData
// @Failure 400,500 {object} model.Response
// @Router /auth/refresh [post]
func (h authUserHandler) Refresh(c *gin.Context) {
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
	userID := claims["user_id"].(uint)

	user, err := h.authService.GetByID(userID)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}
	expiredTimeAT := time.Now().Add(time.Minute * time.Duration(h.infra.Config().GetInt("access_token_expired"))).Unix()
	authAT, err := h.authService.CreateAuth(user.ID, expiredTimeAT, "at")
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("error autentikasi: %v", err))
		return
	}

	expired, accessToken := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateToken(
		model.UserTokenPayload{
			UserID:   authAT.UserID,
			AuthUUID: authAT.AuthUUID,
			Expired:  authAT.Expired,
		},
	)

	expiredTimeRT := time.Now().Add(time.Minute * time.Duration(h.infra.Config().GetInt("refresh_token_expired"))).Unix()
	authRT, err := h.authService.CreateAuth(user.ID, expiredTimeRT, "rt")
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, fmt.Errorf("error autentikasi: %v", err))
		return
	}

	refreshExpired, refreshToken := token.NewToken(h.infra.Config().GetString("secret.key")).GenerateRefreshToken(
		model.UserTokenPayload{
			UserID:   authRT.UserID,
			AuthUUID: authRT.AuthUUID,
			Expired:  authRT.Expired,
		},
	)

	user.Role = user.GetRole()
	user.UserAbilities = user.GetAbility()
	user.Avatar = user.GetAvatar()

	tokenData := model.TokenData{
		AccessToken:         accessToken,
		ExpiredAccessToken:  expired,
		RefreshToken:        refreshToken,
		ExpiredRefreshToken: refreshExpired,
	}
	dataOutput := model.AuthData{
		UserData:  user,
		TokenData: tokenData,
	}
	response.New(c).Data(200, "sukses menyegarkan data", dataOutput)
}

// Activation ... Activation Account URL
// @Summary Set Active By Click This URL
// @Description Set Active By Click This URL
// @Tags Auth
// @Accept json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /auth/activation [get]
// @param token query string true "token data"
func (h authUserHandler) Activation(c *gin.Context) {
	config := h.infra.Config().Sub("server")
	token := c.Query("token")
	if token == "" {
		data := gin.H{
			"status":  "Gagal",
			"message": "Aktivasi token tidak valid",
			"title":   "Aktifasi Akun",
			"url":     config.GetString("web_url"),
		}
		c.HTML(http.StatusOK, "verify_email.html", data)
		return
	}
	isValid, userID := h.activationTokenService.IsValid(token)
	if userID == 0 && !isValid {

		data := gin.H{
			"status":  "Gagal",
			"message": "Ketika sedang memvalidasi, aktivasi token sudah kedaluarsa",
			"title":   "Aktifasi Akun",
			"url":     config.GetString("web_url"),
		}
		c.HTML(http.StatusOK, "verify_email.html", data)
		return
	}
	user, err := h.authService.SetActiveUser(int(userID))
	if err != nil {
		data := gin.H{
			"status":  "Gagal",
			"message": "Ketika sedang mengaktifkan pengguna aktivasi token sudah kedaluarsa",
			"title":   "Aktifasi Akun",
			"url":     config.GetString("web_url"),
		}
		c.HTML(http.StatusOK, "verify_email.html", data)
		return
	}

	data := gin.H{
		"status":  "Sukses",
		"message": "Halo " + user.FirstName + " " + user.LastName + ", anda telah berhasil memferifikasi email, anda sekarang bisa mealakukan aktifitas pada aplikasi.",
		"title":   "Aktifasi Akun",
		"url":     config.GetString("web_url"),
	}
	c.HTML(http.StatusOK, "verify_email.html", data)
}

// Logout ... Logout System
// @Summary Logout from system
// @Description Logout from system
// @Tags Auth
// @Accept json
// @Success 200 {object} model.Response
// @Failure 400,500 {object} model.Response
// @Router /auth/logout [get]
// @Security BearerTokenAuth
func (h authUserHandler) Logout(c *gin.Context) {

	token := token.NewToken(h.infra.Config().GetString("secret.key"))
	auth, err := token.ExtractTokenAuth(c)
	if err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	if err := h.authService.DeleteAuth(auth.UserID, auth.AuthUUID); err != nil {
		response.New(c).Error(http.StatusBadRequest, err)
		return
	}

	response.New(c).Write(200, "berhasil keluar sistem")
}
