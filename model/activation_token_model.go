package model

import (
	"time"
)

type ActivationToken struct {
	GormCustom
	UserID uint      `json:"user_id" query:"user_id" form:"user_id"`
	User   User      `json:"user" query:"user" form:"user"`
	Token  string    `json:"token" gorm:"type:varchar(64);unique" query:"token" form:"token"`
	Valid  time.Time `json:"valid" query:"valid" form:"valid"`
}
