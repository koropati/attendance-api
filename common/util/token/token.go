package token

import (
	"attendance-api/model"
	"fmt"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Token interface {
	GenerateToken(data model.UserTokenPayload) (expiredDate int64, tokenData string)
	GenerateRefreshToken(data model.UserTokenPayload) (expiredDate int64, tokenData string)
	ValidateToken(token string) (*jwt.Token, error)
	ValidateRefreshToken(token string) (*jwt.Token, error)
	ExtractToken(c *gin.Context) string
	ExtractTokenAuth(c *gin.Context) (model.Auth, error)
}

type token struct {
	secretKey string
}

func NewToken(secretKey string) Token {
	return &token{secretKey: secretKey}
}

type authClaims struct {
	UserID   uint   `json:"user_id"`
	AuthUUID string `json:"auth_uuid"`
	Expired  int64  `json:"expired"`
	jwt.StandardClaims
}

type refreshClaims struct {
	UserID   uint   `json:"user_id"`
	AuthUUID string `json:"auth_uuid"`
	Expired  int64  `json:"expired"`
	jwt.StandardClaims
}

func (t *token) GenerateToken(data model.UserTokenPayload) (expiredDate int64, tokenData string) {
	claims := &authClaims{
		data.UserID,
		data.AuthUUID,
		data.Expired,
		jwt.StandardClaims{},
	}

	ctx := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := ctx.SignedString([]byte(t.secretKey))
	if err != nil {
		logrus.Panic(err)
	}

	return data.Expired, token
}

func (t *token) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("token tidak valid : %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
}

func (t *token) GenerateRefreshToken(data model.UserTokenPayload) (expiredDate int64, tokenData string) {
	claims := &refreshClaims{
		data.UserID,
		data.AuthUUID,
		data.Expired,
		jwt.StandardClaims{},
	}

	ctx := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := ctx.SignedString([]byte(t.secretKey))
	if err != nil {
		logrus.Panic(err)
	}

	return data.Expired, token
}

func (t *token) ValidateRefreshToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("token tidak valid : %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
}

func (t *token) ExtractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	authBearer := strings.Split(authHeader, " ")
	//normally Authorization the_token_xxx
	return authBearer[1]
}

func (t *token) ExtractTokenAuth(c *gin.Context) (model.Auth, error) {
	authHeader := c.GetHeader("Authorization")
	authBearer := strings.Split(authHeader, " ")
	if authBearer[1] == "" || authBearer[1] == " " {
		return model.Auth{}, fmt.Errorf("token otorisasi tidak valid")
	}
	if len(authBearer) == 2 {
		if token, err := t.ValidateToken(authBearer[1]); token.Valid && err == nil {
			// Validate expired token
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				authData := model.Auth{
					UserID:   uint(claims["user_id"].(float64)),
					AuthUUID: claims["auth_uuid"].(string),
					Expired:  int64(claims["expired"].(float64)),
				}
				return authData, nil
			} else {
				return model.Auth{}, fmt.Errorf("token klaim tidak valid")
			}
		} else {
			return model.Auth{}, fmt.Errorf("token otorisasi tidak valid E: %v", err)
		}
	} else {
		return model.Auth{}, fmt.Errorf("token otorisasi tidak valid")
	}
}
