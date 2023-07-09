package model

// post
type Subject struct {
	GormCustom
	Name    string `json:"name" gorm:"type:varchar(100)" query:"name" form:"name"`
	Code    string `json:"code" gorm:"unique;type:varchar(25)" query:"code" form:"code"`
	Summary string `json:"summary" gorm:"type:text" query:"summary" form:"summary"`
	OwnerID int    `json:"owner_id" gorm:"not null" query:"owner_id" form:"owner_id"`
}
