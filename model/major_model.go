package model

type Major struct {
	GormCustom
	Name    string `json:"name" gorm:"type:varchar(100)"`
	Code    string `json:"code" gorm:"unique;type:varchar(25)"`
	Summary string `json:"summary" gorm:"type:text"`
	OwnerID int    `json:"owner_id" gorm:"not null"`
}
