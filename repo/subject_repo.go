package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type SubjectRepo interface {
	CreateSubject(subject model.Subject) (model.Subject, error)
	RetrieveSubject(id int) (model.Subject, error)
	RetrieveSubjectByOwner(id int, ownerID int) (model.Subject, error)
	UpdateSubject(id int, subject model.Subject) (model.Subject, error)
	UpdateSubjectByOwner(id int, ownerID int, subject model.Subject) (model.Subject, error)
	DeleteSubject(id int) error
	DeleteSubjectByOwner(id int, ownerID int) error
	ListSubject(subject model.Subject, pagination model.Pagination) ([]model.Subject, error)
	ListSubjectMeta(subject model.Subject, pagination model.Pagination) (model.Meta, error)
	DropDownSubject(subject model.Subject) ([]model.Subject, error)
	CheckIsExist(id int) (isExist bool)
}

type subjectRepo struct {
	db *gorm.DB
}

func NewSubjectRepo(db *gorm.DB) SubjectRepo {
	return &subjectRepo{db: db}
}

func (r subjectRepo) CreateSubject(subject model.Subject) (model.Subject, error) {
	if err := r.db.Table("subjects").Create(&subject).Error; err != nil {
		return model.Subject{}, err
	}

	return subject, nil
}

func (r subjectRepo) RetrieveSubject(id int) (model.Subject, error) {
	var subject model.Subject
	if err := r.db.First(&subject, id).Error; err != nil {
		return model.Subject{}, err
	}
	return subject, nil
}

func (r subjectRepo) RetrieveSubjectByOwner(id int, ownerID int) (model.Subject, error) {
	var subject model.Subject
	if err := r.db.Model(&model.Subject{}).Where("id = ? AND owner_id = ?", id, ownerID).First(&subject).Error; err != nil {
		return model.Subject{}, err
	}
	return subject, nil
}

func (r subjectRepo) UpdateSubject(id int, subject model.Subject) (model.Subject, error) {
	if err := r.db.Model(&model.Subject{}).Where("id = ?", id).Updates(&subject).Error; err != nil {
		return model.Subject{}, err
	}
	return subject, nil
}

func (r subjectRepo) UpdateSubjectByOwner(id int, ownerID int, subject model.Subject) (model.Subject, error) {
	if err := r.db.Model(&model.Subject{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(&subject).Error; err != nil {
		return model.Subject{}, err
	}
	return subject, nil
}

func (r subjectRepo) DeleteSubject(id int) error {
	if err := r.db.Delete(&model.Subject{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r subjectRepo) DeleteSubjectByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Subject{}).Error; err != nil {
		return err
	}
	return nil
}

func (r subjectRepo) ListSubject(subject model.Subject, pagination model.Pagination) ([]model.Subject, error) {
	var subjects []model.Subject
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("subjects").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterSubject(query, subject)
	query = SearchSubject(query, pagination.Search)
	query = query.Find(&subjects)
	if err := query.Error; err != nil {
		return nil, err
	}

	return subjects, nil
}

func (r subjectRepo) ListSubjectMeta(subject model.Subject, pagination model.Pagination) (model.Meta, error) {
	var subjects []model.Subject
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Subject{}).Select("count(*)")
	queryTotal = FilterSubject(queryTotal, subject)
	queryTotal = SearchSubject(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("subjects").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterSubject(query, subject)
	query = SearchSubject(query, pagination.Search)
	query = query.Find(&subjects)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(subjects),
	}
	return meta, nil
}

func (r subjectRepo) DropDownSubject(subject model.Subject) ([]model.Subject, error) {
	var subjects []model.Subject
	query := r.db.Table("subjects").Order("id desc")
	query = FilterSubject(query, subject)
	query = query.Find(&subjects)
	if err := query.Error; err != nil {
		return nil, err
	}
	return subjects, nil
}

func (r subjectRepo) CheckIsExist(id int) (isExist bool) {
	if err := r.db.Table("subjects").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func FilterSubject(query *gorm.DB, subject model.Subject) *gorm.DB {
	if subject.Name != "" {
		query = query.Where("name LIKE ?", "%"+subject.Name+"%")
	}
	if subject.Code != "" {
		query = query.Where("code LIKE ?", "%"+subject.Code+"%")
	}
	if subject.OwnerID > 0 {
		query = query.Where("owner_id = ?", subject.OwnerID)
	}
	return query
}

func SearchSubject(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR summary LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}
