package model

import (
	"database/sql"
	"time"
)

type Auth struct {
	ID       uint   `json:"id" gorm:"primary_key" query:"id" form:"id"`
	UserID   uint   `json:"user_id" gorm:"not null" query:"user_id" form:"user_id"`
	AuthUUID string `json:"auth_uuid" gorm:"size:255;not null;" query:"auth_uuid" form:"auth_uuid"`
	Expired  int64  `json:"expired" query:"expired" form:"expired"`
	TypeAuth string `json:"type_auth" gorm:"type:enum('at','rt');default:'at'" query:"type_auth" form:"type_auth"`
}

type RoleAbility struct {
	ID           uint         `json:"id" gorm:"primary_key" query:"id" form:"id"`
	CreatedAt    time.Time    `json:"created_at" query:"created_at" form:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" query:"updated_at" form:"updated_at"`
	DeletedAt    sql.NullTime `json:"deleted_at" sql:"index" query:"deleted_at" form:"deleted_at" gorm:"index"`
	IsSuperAdmin bool         `json:"is_super_admin" query:"is_super_admin" form:"is_super_admin" gorm:"type:boolean;default:1"`
	IsAdmin      bool         `json:"is_admin" query:"is_admin" form:"is_admin" gorm:"type:boolean;default:0"`
	IsUser       bool         `json:"is_user" query:"is_user" form:"is_user" gorm:"type:boolean;default:0"`
	Action       string       `json:"action" query:"action" form:"action" gorm:"type:enum('create','read','update','delete','manage');default:'read'"`
	Subject      string       `json:"subject" query:"subject" form:"subject" gorm:"type:enum('activation_token','teacher','faculty','schedule','major','student','subject','user','attendance','study_program','reset_token','auth','all');default:'auth'"`
}

type Ability struct {
	Action  string `json:"action"`
	Subject string `json:"subject"`
}

type ForgotPassword struct {
	Email string `json:"email" query:"email"`
}

type ConfirmForgotPassword struct {
	Token           string `json:"token" query:"token" form:"token"`
	Password        string `json:"password" query:"password" form:"password"`
	ConfirmPassword string `json:"confirm_password" query:"confirm_password" form:"confirm_password"`
}

func GetDefaultAbility() (results []Ability) {
	results = append(results, Ability{
		Action:  "auth",
		Subject: "read",
	})
	return
}

func GetSuperAdminAbility() (results []Ability) {
	results = append(results, Ability{
		Action:  "manage",
		Subject: "all",
	})
	return
}

func GetAdminAbility() (results []Ability) {
	menus := []string{"student", "teacher", "attendance", "schedule", "daily_schedule", "subject", "user_schedule"}
	actions := []string{"read", "read", "manage", "manage", "manage", "manage", "manage"}
	for i, menu := range menus {
		results = append(results, Ability{
			Action:  actions[i],
			Subject: menu,
		})
	}
	return
}

func GetUserAbility() (results []Ability) {
	menus := []string{"student", "teacher", "attendance", "schedule", "daily_schedule", "subject", "user_schedule"}
	actions := []string{"read", "read", "create", "read", "read", "read", "read"}
	for i, menu := range menus {
		results = append(results, Ability{
			Action:  actions[i],
			Subject: menu,
		})
	}
	return
}
