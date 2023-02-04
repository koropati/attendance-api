package token

import (
	"attendance-api/model"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type Token interface {
	GenerateToken(data *model.UserTokenPayload) (expiredDate int64, tokenData string)
	GenerateRefreshToken(data *model.UserTokenPayload) (expiredDate int64, tokenData string)
	ValidateToken(token string) (*jwt.Token, error)
	ValidateRefreshToken(token string) (*jwt.Token, error)
}

type token struct {
	secretKey string
}

func NewToken(secretKey string) Token {
	return &token{secretKey: secretKey}
}

type authClaims struct {
	Username     string `json:"username"`
	IsSuperAdmin bool   `json:"is_Super_admin"`
	IsAdmin      bool   `json:"is_admin"`
	IsUser       bool   `json:"is_user"`
	Email        string `json:"email"`
	Expired      string `json:"expired"`
	jwt.StandardClaims
}

type refreshClaims struct {
	Username string `json:"username"`
	Expired  string `json:"expired"`
	jwt.StandardClaims
}

func (t *token) GenerateToken(data *model.UserTokenPayload) (expiredDate int64, tokenData string) {
	expiredTime := time.Now().Add(time.Minute * time.Duration(data.Expired)).Unix()
	claims := &authClaims{
		data.Username,
		data.IsSuperAdmin,
		data.IsAdmin,
		data.IsUser,
		data.Email,
		strconv.FormatInt(expiredTime, 10),
		jwt.StandardClaims{},
	}

	ctx := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := ctx.SignedString([]byte(t.secretKey))
	if err != nil {
		logrus.Panic(err)
	}

	return expiredTime, token
}

func (t *token) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("invalid token %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
}

func (t *token) GenerateRefreshToken(data *model.UserTokenPayload) (expiredDate int64, tokenData string) {
	expiredTime := time.Now().Add(time.Minute * time.Duration(data.Expired)).Unix()
	claims := &refreshClaims{
		data.Username,
		strconv.FormatInt(expiredTime, 10),
		jwt.StandardClaims{},
	}

	ctx := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := ctx.SignedString([]byte(t.secretKey))
	if err != nil {
		logrus.Panic(err)
	}

	return expiredTime, token
}

func (t *token) ValidateRefreshToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, valid := token.Method.(*jwt.SigningMethodHMAC); !valid {
			return nil, fmt.Errorf("invalid token %v", token.Header["alg"])
		}
		return []byte(t.secretKey), nil
	})
}
