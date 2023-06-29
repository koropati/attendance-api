package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"attendance-api/common/http/response"
	"attendance-api/common/util/token"
	"attendance-api/service"

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
	LOGOUT(c *gin.Context) error
}

type middleware struct {
	secretKey   string
	authService service.AuthService
}

func NewMiddleware(secretKey string, authService service.AuthService) Middleware {
	return &middleware{
		secretKey:   secretKey,
		authService: authService,
	}
}

func (m *middleware) CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Request-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		c.Writer.Header().Set("Accept-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Accept-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Accept-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Content-Type", "application/json")

		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Methods", "*")
			c.Header("Access-Control-Allow-Headers", "*")
			c.Status(200)
			c.Abort()
			// c.AbortWithStatus(200)
			return
		}
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
				if err != nil {
					return nil, false, fmt.Errorf("token kedaluarsa tidak valid e: %v", err)
				}
				currentDateTime := time.Now()
				expiredDateTime := time.Unix(int64(claims["expired"].(float64)), 0)

				if currentDateTime.After(expiredDateTime) {
					return nil, false, fmt.Errorf("token sudah kedaluarsa")
				} else {
					// check in DB
					_, err := m.authService.FetchAuth(uint(claims["user_id"].(float64)), claims["auth_uuid"].(string))
					if err != nil {
						return nil, false, errors.New("akses tidak sah ditolak")
					}
					return token, true, nil
				}

			} else {
				return nil, false, fmt.Errorf("klaim token tidak valid")
			}

		} else {
			return nil, false, fmt.Errorf("token otorisasi tidak valid : %v", err)
		}
	} else {
		return nil, false, errors.New("token otorisasi tidak valid")
	}
}

func ValidateRole(m *middleware, token *jwt.Token, roles ...string) (valid bool, err error) {

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		authUser, err := m.authService.FetchAuth(uint(claims["user_id"].(float64)), claims["auth_uuid"].(string))
		if err != nil {
			valid = false
		} else {
			user, err := m.authService.GetByID(authUser.UserID)
			if err != nil {
				valid = false
			} else {
				for _, role := range roles {
					if role == "user" && user.IsUser {
						valid = true
						break
					}
					if role == "admin" && user.IsAdmin {
						valid = true
						break
					}
					if role == "super_admin" && user.IsSuperAdmin {
						valid = true
						break
					}
				}
			}

		}

		if valid {
			return true, nil
		} else {
			return false, errors.New("kamu tidak bisa mengakses ini")
		}
	} else {
		return false, errors.New("token tidak valid")
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
			validRole, err := ValidateRole(m, token, "super_admin")
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
			validRole, err := ValidateRole(m, token, "admin", "super_admin")
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
			validRole, err := ValidateRole(m, token, "editor")
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
			validRole, err := ValidateRole(m, token, "user", "admin", "super_admin")
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
			return 0, fmt.Errorf("token tidak valid")
		}
	}
}

func (m *middleware) IsSuperAdmin(c *gin.Context) bool {
	token, validToken, err := ValidateToken(m, c)
	if !validToken && err != nil {
		return false
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			authUser, err := m.authService.FetchAuth(uint(claims["user_id"].(float64)), claims["auth_uuid"].(string))
			if err != nil {
				return false
			}
			user, err := m.authService.GetByID(authUser.UserID)
			if err != nil {
				return false
			}
			return user.IsSuperAdmin
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
			authUser, err := m.authService.FetchAuth(uint(claims["user_id"].(float64)), claims["auth_uuid"].(string))
			if err != nil {
				return false
			}
			user, err := m.authService.GetByID(authUser.UserID)
			if err != nil {
				return false
			}
			return user.IsUser
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
			authUser, err := m.authService.FetchAuth(uint(claims["user_id"].(float64)), claims["auth_uuid"].(string))
			if err != nil {
				return false
			}
			user, err := m.authService.GetByID(authUser.UserID)
			if err != nil {
				return false
			}
			return user.IsAdmin
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
				authUser, err := m.authService.FetchAuth(uint(claims["user_id"].(float64)), claims["auth_uuid"].(string))
				if err != nil {
					response.New(c).Error(http.StatusUnauthorized, err)
					c.Abort()
				}
				user, err := m.authService.GetByID(authUser.UserID)
				if err != nil {
					response.New(c).Error(http.StatusUnauthorized, err)
					c.Abort()
				}

				if user.IsSuperAdmin {
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

func (m *middleware) LOGOUT(c *gin.Context) error {
	token, validToken, err := ValidateToken(m, c)
	if !validToken && err != nil {
		return fmt.Errorf("token tidak valid")
	} else {
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			err := m.authService.DeleteAuth(uint(claims["user_id"].(float64)), claims["auth_uuid"].(string))
			if err != nil {
				return err
			} else {
				return nil
			}

		} else {
			return fmt.Errorf("token tidak valid")
		}
	}
}
