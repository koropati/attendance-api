package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type SubjectRepo interface {
	CreateSubject(subject *model.Subject) (*model.Subject, error)
}

type subjectRepo struct {
	db *gorm.DB
}

func NewSubjectRepo(db *gorm.DB) SubjectRepo {
	return &subjectRepo{db: db}
}

func (r *subjectRepo) CreateSubject(subject *model.Subject) (*model.Subject, error) {
	if err := r.db.Table("subjects").Create(&subject).Error; err != nil {
		return nil, err
	}

	return subject, nil
}
