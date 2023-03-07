package repo

import (
	"attendance-api/model"
	"log"
	"time"

	"gorm.io/gorm"
)

type ActivationTokenRepo interface {
	CreateActivationToken(activationToken *model.ActivationToken) (*model.ActivationToken, error)
	RetrieveActivationToken(id int) (*model.ActivationToken, error)
	UpdateActivationToken(id int, activationToken *model.ActivationToken) (*model.ActivationToken, error)
	DeleteActivationToken(id int) error
	ListActivationToken(activationToken *model.ActivationToken, pagination *model.Pagination) (*[]model.ActivationToken, error)
	ListActivationTokenMeta(activationToken *model.ActivationToken, pagination *model.Pagination) (*model.Meta, error)
	DropDownActivationToken(activationToken *model.ActivationToken) (*[]model.ActivationToken, error)
	IsValid(token string) (isValid bool, userID uint)
}

type activationTokenRepo struct {
	db *gorm.DB
}

func NewActivationTokenRepo(db *gorm.DB) ActivationTokenRepo {
	return &activationTokenRepo{db: db}
}

func (r *activationTokenRepo) CreateActivationToken(subject *model.ActivationToken) (*model.ActivationToken, error) {
	if err := r.db.Table("activation_tokens").Create(&subject).Error; err != nil {
		return nil, err
	}

	return subject, nil
}

func (r *activationTokenRepo) RetrieveActivationToken(id int) (*model.ActivationToken, error) {
	var subject model.ActivationToken
	if err := r.db.First(&subject, id).Error; err != nil {
		return nil, err
	}
	return &subject, nil
}

func (r *activationTokenRepo) UpdateActivationToken(id int, subject *model.ActivationToken) (*model.ActivationToken, error) {
	if err := r.db.Model(&model.ActivationToken{}).Where("id = ?", id).Updates(&subject).Error; err != nil {
		return nil, err
	}
	return subject, nil
}

func (r *activationTokenRepo) DeleteActivationToken(id int) error {
	if err := r.db.Delete(&model.ActivationToken{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *activationTokenRepo) ListActivationToken(subject *model.ActivationToken, pagination *model.Pagination) (*[]model.ActivationToken, error) {
	var activationTokens []model.ActivationToken
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("activation_tokens").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterActivationToken(query, subject)
	query = query.Find(&activationTokens)
	if err := query.Error; err != nil {
		return nil, err
	}

	return &activationTokens, nil
}

func (r *activationTokenRepo) ListActivationTokenMeta(subject *model.ActivationToken, pagination *model.Pagination) (*model.Meta, error) {
	var activationTokens []model.ActivationToken
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.ActivationToken{}).Select("count(*)")
	queryTotal = FilterActivationToken(queryTotal, subject)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return nil, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(activationTokens),
	}
	return &meta, nil
}

func (r *activationTokenRepo) DropDownActivationToken(subject *model.ActivationToken) (*[]model.ActivationToken, error) {
	var activationTokens []model.ActivationToken
	query := r.db.Table("activation_tokens")
	query = FilterActivationToken(query, subject)
	query = query.Find(&activationTokens)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &activationTokens, nil
}

func (r *activationTokenRepo) IsValid(token string) (isValid bool, userID uint) {
	var data model.ActivationToken
	if err := r.db.Table("activation_tokens").Where("token = ?", token).First(&data).Error; err != nil {
		log.Printf("[Error] [IsValid] E: %v\n", err)
		isValid = false
		userID = 0
	} else {
		if data.Valid.After(time.Now()) {
			isValid = false
			userID = 0
		} else {
			isValid = true
			userID = data.UserID
		}
	}
	return
}

func FilterActivationToken(query *gorm.DB, subject *model.ActivationToken) *gorm.DB {
	if subject.Token != "" {
		query = query.Where("token LIKE ?", "%"+subject.Token+"%")
	}
	if subject.Valid.String() != "" {
		query = query.Where("valid LIKE ?", "%"+subject.Valid.Local().Format("2006-01-02")+"%")
	}
	if subject.UserID > 0 {
		query = query.Where("user_id = ?", subject.UserID)
	}
	return query
}
