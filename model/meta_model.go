package model

type PostMeta struct {
	GormCustom
	PostID  int    `json:"post_id"`
	Key     string `json:"key" gorm:"type:varchar(50)"`
	Content string `json:"content" gorm:"type:text"`
}
