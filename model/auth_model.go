package model

type Auth struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	UserID   uint   `json:"user_id" gorm:"not null"`
	AuthUUID string `json:"auth_uuid" gorm:"size:255;not null;"`
	Expired  int64  `json:"expired"`
	TypeAuth string `json:"type_auth" gorm:"type:enum('at','rt');default:'at'"`
}
