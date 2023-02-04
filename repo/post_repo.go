package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type PostRepo interface {
	CreatePost(post *model.Post) (*model.Post, error)
}

type postRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) PostRepo {
	return &postRepo{db: db}
}

func (r *postRepo) CreatePost(post *model.Post) (*model.Post, error) {
	if err := r.db.Table("post").Create(&post).Error; err != nil {
		return nil, err
	}

	return post, nil
}
