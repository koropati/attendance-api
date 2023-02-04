package model

import "time"

type Comment struct {
	GormCustom
	PostID      int       `json:"post_id"`
	ParentID    int       `json:"parent_id"`
	Parent      *Comment  `json:"parent" gorm:"foreignkey:ParentID"`
	Title       string    `json:"title" gorm:"type:varchar(100)"`
	Content     string    `json:"content" gorm:"type:text"`
	Published   bool      `json:"published"`
	PublishedAt time.Time `json:"published_at"`
}
