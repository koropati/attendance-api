package model

import (
	"time"
)

type PasswordResetToken struct {
	GormCustom
	UserID uint      `json:"user_id"`
	User   User      `json:"user"`
	Token  string    `json:"token" gorm:"type:varchar(64);unique"`
	Valid  time.Time `json:"valid"`
}
