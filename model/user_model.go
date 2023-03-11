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

type UserForm struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Handphone string `json:"handphone"`
	Email     string `json:"email"`
	Intro     string `json:"intro"`
	Profile   string `json:"profile"`
}

type UserForgotPasswordForm struct {
	Email           string `json:"email"`
	Activation      string `json:"activation"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserUpdatePasswordForm struct {
	ID              uint   `json:"id" query:"id"`
	CurrentPassword string `json:"current_password" query:"current_password"`
	NewPassword     string `json:"new_password" query:"new_password"`
	ConfirmPassword string `json:"confirm_password" query:"confirm_password"`
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
	UserID       uint   `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	IsSuperAdmin bool   `json:"is_super_admin"`
	IsAdmin      bool   `json:"is_admin"`
	IsUser       bool   `json:"is_user"`
	Expired      int    `json:"expired"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token"`
}
