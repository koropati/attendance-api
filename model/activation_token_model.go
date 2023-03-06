package model

import "time"

type ActivationToken struct {
	GormCustom
	UserID uint      `json:"user_id"`
	Token  string    `json:"token" gorm:"type:varchar(64);unique"`
	Valid  time.Time `json:"valid"`
}
