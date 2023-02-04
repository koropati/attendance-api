package model

// post
type Subject struct {
	GormCustom
	Name    string `json:"name" gorm:"type:varchar(100)"`
	Code    string `json:"code" gorm:"unique;type:varchar(25)"`
	Summary string `json:"summary" gorm:"type:text"`
}
