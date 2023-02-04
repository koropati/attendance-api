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
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	Name      string    `json:"name"`
	Handphone string    `json:"handphone" gorm:"unique"`
	Email     string    `json:"email" gorm:"unique"`
	Intro     string    `json:"intro" gorm:"type:varchar(255)"`
	Profile   string    `json:"profile" gorm:"type:varchar(255)"`
	LastLogin time.Time `json:"last_login"`
	Role      string    `json:"role" gorm:"type:enum('super_admin','admin','editor','user');default:'user'"`
	IsActive  bool      `json:"is_active"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token"`
}
