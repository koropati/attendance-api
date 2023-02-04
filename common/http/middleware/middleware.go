package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"attendance-api/common/http/response"
	"attendance-api/common/util/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Middleware interface {
	CORS() gin.HandlerFunc
	AUTH() gin.HandlerFunc
	SUPERADMIN() gin.HandlerFunc
}

type middleware struct {
	secretKey string
}

func NewMiddleware(secretKey string) Middleware {
	return &middleware{secretKey: secretKey}
}

func (m *middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	}
}

func ValidateToken(m *middleware, c *gin.Context) (tokenData *jwt.Token, valid bool, err error) {
	authHeader := c.GetHeader("Authorization")
	authBearer := strings.Split(authHeader, " ")

	if len(authBearer) == 2 {
		if token, err := token.NewToken(m.secretKey).ValidateToken(authBearer[1]); token.Valid && err == nil {
			// Validate expired token
			claims, ok := token.Claims.(jwt.MapClaims)
			if ok {
				unixTime, err := strconv.ParseInt(claims["expired"].(string), 10, 64)
				if err != nil {
					return nil, false, fmt.Errorf("invalid expiration token e: %v", err)
				}
				currentDateTime := time.Now()
				expiredDateTime := time.Unix(int64(unixTime), 0)

				if currentDateTime.After(expiredDateTime) {
					return nil, false, fmt.Errorf("token is expired")
				} else {
					return token, true, nil
				}

			} else {
				return nil, false, fmt.Errorf("invalid claim token")
			}

		} else {
			return nil, false, fmt.Errorf("invalid authorization token %v", err)
		}
	} else {
		return nil, false, errors.New("invalid authorization token")
	}
}

func ValidateRole(token *jwt.Token, role string) (valid bool, err error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if claims["role"] == role {
			return true, nil
		} else {
			return false, errors.New("you can't access this")
		}
	} else {
		return false, errors.New("invalid token")
	}
}

func (m *middleware) AUTH() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, validToken, err := ValidateToken(m, c)
		if !validToken && err != nil {
			response.New(c).Error(http.StatusUnauthorized, err)
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func (m *middleware) SUPERADMIN() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, validToken, err := ValidateToken(m, c)
		if !validToken && err != nil {
			response.New(c).Error(http.StatusUnauthorized, err)
			c.Abort()
		} else {
			validRole, err := ValidateRole(token, "super_admin")
			if !validRole && err != nil {
				response.New(c).Error(http.StatusForbidden, err)
				c.Abort()
			} else {
				c.Next()
			}
		}
	}
}

func (m *middleware) ADMIN() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, validToken, err := ValidateToken(m, c)
		if !validToken && err != nil {
			response.New(c).Error(http.StatusUnauthorized, err)
			c.Abort()
		} else {
			validRole, err := ValidateRole(token, "admin")
			if !validRole && err != nil {
				response.New(c).Error(http.StatusForbidden, err)
				c.Abort()
			} else {
				c.Next()
			}
		}
	}
}

func (m *middleware) EDITOR() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, validToken, err := ValidateToken(m, c)
		if !validToken && err != nil {
			response.New(c).Error(http.StatusUnauthorized, err)
			c.Abort()
		} else {
			validRole, err := ValidateRole(token, "editor")
			if !validRole && err != nil {
				response.New(c).Error(http.StatusForbidden, err)
				c.Abort()
			} else {
				c.Next()
			}
		}
	}
}

func (m *middleware) USER() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, validToken, err := ValidateToken(m, c)
		if !validToken && err != nil {
			response.New(c).Error(http.StatusUnauthorized, err)
			c.Abort()
		} else {
			validRole, err := ValidateRole(token, "user")
			if !validRole && err != nil {
				response.New(c).Error(http.StatusForbidden, err)
				c.Abort()
			} else {
				c.Next()
			}
		}
	}
}
