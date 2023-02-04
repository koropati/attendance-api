package token

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
)

type Token interface {
	GenerateToken(username string, email string, role string, expired int) (expiredDate int64, tokenData string)
	GenerateRefreshToken(username string, expired int) (expiredDate int64, tokenData string)
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
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Expired  string `json:"expired"`
	jwt.StandardClaims
}

type refreshClaims struct {
	Username string `json:"username"`
	Expired  string `json:"expired"`
	jwt.StandardClaims
}

func (t *token) GenerateToken(username string, email string, role string, expired int) (expiredDate int64, tokenData string) {
	expiredTime := time.Now().Add(time.Minute * time.Duration(expired)).Unix()
	claims := &authClaims{
		username,
		role,
		email,
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

func (t *token) GenerateRefreshToken(username string, expired int) (expiredDate int64, tokenData string) {
	expiredTime := time.Now().Add(time.Minute * time.Duration(expired)).Unix()
	claims := &refreshClaims{
		username,
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
