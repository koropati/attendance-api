package model

type Category struct {
	GormCustom
	ParentID  int       `json:"parent_id"`
	Parent    *Category `json:"parent" gorm:"foreignkey:ParentID"`
	Title     string    `json:"title" gorm:"type:varchar(100)"`
	MetaTitle string    `json:"meta_title" gorm:"type:varchar(100)"`
	Slug      string    `json:"slug" gorm:"type:varchar(100)"`
	Content   string    `json:"content" gorm:"type:text"`
}
