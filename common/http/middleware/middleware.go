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
	ADMIN() gin.HandlerFunc
	USER() gin.HandlerFunc
	GetUserID(c *gin.Context) (userID int, err error)
	HaveAccess(c *gin.Context, ownerID int) gin.HandlerFunc
	IsSuperAdmin(c *gin.Context) bool
	IsUser(c *gin.Context) bool
	IsAdmin(c *gin.Context) bool
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

func ValidateRole(token *jwt.Token, roles ...string) (valid bool, err error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		for _, role := range roles {
			if role == "user" && claims["is_user"].(bool) {
				valid = true
				break
			}
			if role == "admin" && claims["is_admin"].(bool) {
				valid = true
				break
			}
			if role == "super_admin" && claims["is_super_admin"].(bool) {
				valid = true
				break
			}
		}
		if valid {
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
			validRole, err := ValidateRole(token, "admin", "super_admin")
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
			validRole, err := ValidateRole(token, "user", "admin", "super_admin")
			if !validRole && err != nil {
				response.New(c).Error(http.StatusForbidden, err)
				c.Abort()
			} else {
				c.Next()
			}
		}
	}
}

func (m *middleware) GetUserID(c *gin.Context) (userId int, err error) {
	token, validToken, err := ValidateToken(m, c)
	if !validToken && err != nil {
		return 0, err
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return int(claims["user_id"].(float64)), nil
		} else {
			return 0, fmt.Errorf("token is invalid")
		}
	}
}

func (m *middleware) IsSuperAdmin(c *gin.Context) bool {
	token, validToken, err := ValidateToken(m, c)
	if !validToken && err != nil {
		return false
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims["is_super_admin"].(bool)
		} else {
			return false
		}
	}
}

func (m *middleware) IsUser(c *gin.Context) bool {
	token, validToken, err := ValidateToken(m, c)
	if !validToken && err != nil {
		return false
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims["is_user"].(bool)
		} else {
			return false
		}
	}
}

func (m *middleware) IsAdmin(c *gin.Context) bool {
	token, validToken, err := ValidateToken(m, c)
	if !validToken && err != nil {
		return false
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			return claims["is_admin"].(bool)
		} else {
			return false
		}
	}
}
func (m *middleware) HaveAccess(c *gin.Context, ownerID int) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, validToken, err := ValidateToken(m, c)
		if !validToken && err != nil {
			response.New(c).Error(http.StatusUnauthorized, err)
			c.Abort()
		} else {
			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				if claims["is_super_admin"].(bool) {
					c.Next()
				} else {
					if claims["user_id"].(int) == ownerID {
						c.Next()
					} else {
						response.New(c).Error(http.StatusUnauthorized, err)
						c.Abort()
					}
				}
			} else {
				response.New(c).Error(http.StatusUnauthorized, err)
				c.Abort()
			}
		}
	}
}
