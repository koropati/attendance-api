package model

type Tag struct {
	GormCustom
	Title     string `json:"title" gorm:"type:varchar(100)"`
	MetaTitle string `json:"meta_title" gorm:"type:varchar(100)"`
	Slug      string `json:"slug" gorm:"type:varchar(100)"`
	Content   string `json:"content" gorm:"type:text"`
}
