package model

import (
	"database/sql"
	"time"
)

type GormCustom struct {
	ID        uint         `json:"id" gorm:"primary_key"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" sql:"index"`
	CreatedBy int          `json:"created_by"`
	UpdatedBy int          `json:"updated_by"`
	DeletedBy int          `json:"deleted_by"`
}

type User struct {
	GormCustom
	Username     string    `json:"username" gorm:"unique"`
	Password     string    `json:"password"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Handphone    string    `json:"handphone" gorm:"unique"`
	Email        string    `json:"email" gorm:"unique"`
	Intro        string    `json:"intro" gorm:"type:varchar(255)"`
	Profile      string    `json:"profile" gorm:"type:varchar(255)"`
	IsActive     bool      `json:"is_active"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	IsAdmin      bool      `json:"is_admin"`
	IsUser       bool      `json:"is_user"`
	LastLogin    time.Time `json:"last_login"`
}

type UserDropDown struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Handphone string `json:"handphone"`
	Email     string `json:"email"`
}

// username string, email string, isSuperAdmin bool, isAdmin bool, isUser bool, expired int
type UserTokenPayload struct {
	UserID   uint   `json:"user_id"`
	AuthUUID string `json:"auth_uuid"`
	Expired  int64  `json:"expired"`
}

type TokenData struct {
	AccessToken         string `json:"access_token"`
	ExpiredAccessToken  int64  `json:"expired_access_token"`
	RefreshToken        string `json:"refresh_token"`
	ExpiredRefreshToken int64  `json:"expired_refresh_token"`
}

type AuthData struct {
	UserData  User      `json:"user_data"`
	TokenData TokenData `json:"token_data"`
}
type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token"`
}
