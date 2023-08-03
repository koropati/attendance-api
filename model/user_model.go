package model

import (
	"database/sql"
	"time"
)

type GormCustom struct {
	ID        uint         `json:"id" gorm:"primary_key" query:"id" form:"id"`
	CreatedAt time.Time    `json:"created_at" query:"created_at" form:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" query:"updated_at" form:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" sql:"index" query:"deleted_at" form:"deleted_at" gorm:"index"`
	CreatedBy int          `json:"created_by" query:"created_by" form:"created_by"`
	UpdatedBy int          `json:"updated_by" query:"updated_by" form:"updated_by"`
	DeletedBy int          `json:"deleted_by" query:"deleted_by" form:"deleted_by"`
}

type User struct {
	GormCustom
	Username      string    `json:"username" gorm:"unique" query:"username" form:"username"`
	Password      string    `json:"password" query:"password" form:"password"`
	FirstName     string    `json:"first_name" query:"first_name" form:"first_name"`
	LastName      string    `json:"last_name" query:"last_name" form:"last_name"`
	Handphone     string    `json:"handphone" gorm:"unique" query:"handphone" form:"handphone"`
	Email         string    `json:"email" gorm:"unique" query:"email" form:"email"`
	Intro         string    `json:"intro" gorm:"type:varchar(255)" query:"intro" form:"intro"`
	Profile       string    `json:"profile" gorm:"type:varchar(255)" query:"profile" form:"profile"`
	IsActive      bool      `json:"is_active" query:"is_active" form:"is_active"`
	IsSuperAdmin  bool      `json:"is_super_admin" query:"is_super_admin" form:"is_super_admin"`
	IsAdmin       bool      `json:"is_admin" query:"is_admin" form:"is_admin"`
	IsUser        bool      `json:"is_user" query:"is_user" form:"is_user"`
	LastLogin     time.Time `json:"last_login" gorm:"default:'0001-01-01 11:11:11.111'" query:"last_login" form:"last_login"`
	Role          string    `json:"role" gorm:"-" query:"role" form:"role"`
	UserAbilities []Ability `json:"user_abilities" gorm:"-" query:"user_abilities" form:"user_abilities"`
	Avatar        string    `json:"avatar" gorm:"-" query:"avatar" form:"avatar"`
	ScheduleID    int       `json:"schedule_id" gorm:"-" query:"schedule_id" form:"schedule_id"`
	OwnerID       int       `json:"owner_id" gorm:"-" query:"owner_id" form:"owner_id"`
}

func (data User) GetRole() (role string) {
	if data.IsSuperAdmin {
		role = "Super Admin"
	} else if data.IsAdmin {
		role = "Dosen"
	} else if data.IsUser {
		role = "Mahasiswa"
	} else {
		role = "-"
	}
	return
}

func (data User) GetAbility() (abilities []Ability) {
	if data.IsSuperAdmin {
		return GetSuperAdminAbility()
	} else if data.IsAdmin {
		return GetAdminAbility()
	} else if data.IsUser {
		return GetUserAbility()
	} else {
		return GetDefaultAbility()
	}
}

func (data User) GetAvatar() (url string) {
	if data.IsSuperAdmin {
		return "https://cdn-icons-png.flaticon.com/512/1535/1535835.png"
	} else if data.IsAdmin {
		return "https://cdn-icons-png.flaticon.com/512/8443/8443259.png"
	} else if data.IsUser {
		return "https://cdn-icons-png.flaticon.com/512/201/201818.png"
	} else {
		return "https://cdn-icons-png.flaticon.com/512/7878/7878622.png"
	}
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
