package model

import "time"

type Post struct {
	GormCustom
	AuthorID    int        `json:"author_id"`
	Author      User       `json:"author" gorm:"foreignKey:AuthorID"`
	ParentID    int        `json:"parent_id"`
	Parent      *Post      `json:"parent" gorm:"foreignkey:ParentID"`
	Category    []Category `json:"category" gorm:"many2many:post_categories;"`
	Tag         []Tag      `json:"tag" gorm:"many2many:post_tags;"`
	PostMeta    []PostMeta `json:"post_meta" gorm:"foreignKey:PostID;references:ID"`
	Title       string     `json:"title" gorm:"type:varchar(75)"`
	MetaTitle   string     `json:"meta_title" gorm:"type:varchar(100)"`
	Slug        string     `json:"slug" gorm:"type:varchar(100)"`
	Summary     string     `json:"summary" gorm:"type:text"`
	Content     string     `json:"content" gorm:"type:text"`
	Published   bool       `json:"published"`
	PublishedAt time.Time  `json:"published_at"`
}
