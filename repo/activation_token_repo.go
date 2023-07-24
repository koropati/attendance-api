package repo

import (
	"attendance-api/model"
	"log"
	"time"

	"gorm.io/gorm"
)

type ActivationTokenRepo interface {
	CreateActivationToken(activationToken model.ActivationToken) (model.ActivationToken, error)
	RetrieveActivationToken(id int) (model.ActivationToken, error)
	UpdateActivationToken(id int, activationToken model.ActivationToken) (model.ActivationToken, error)
	DeleteActivationToken(id int) error
	ListActivationToken(activationToken model.ActivationToken, pagination model.Pagination) ([]model.ActivationToken, error)
	ListActivationTokenMeta(activationToken model.ActivationToken, pagination model.Pagination) (model.Meta, error)
	DropDownActivationToken(activationToken model.ActivationToken) ([]model.ActivationToken, error)
	IsValid(token string) (isValid bool, userID uint)
	DeleteExpiredActivationToken(currentTime time.Time) error
}

type activationTokenRepo struct {
	db *gorm.DB
}

func NewActivationTokenRepo(db *gorm.DB) ActivationTokenRepo {
	return &activationTokenRepo{db: db}
}

func (r activationTokenRepo) CreateActivationToken(activationToken model.ActivationToken) (model.ActivationToken, error) {

	if err := r.db.Table("activation_tokens").Create(&activationToken).Error; err != nil {
		return model.ActivationToken{}, err
	}

	query := r.db.Table("activation_tokens").Where("id = ?", activationToken.ID)
	query = PreloadActivationToken(query)
	query = query.First(&activationToken)
	if err := query.Error; err != nil {
		return model.ActivationToken{}, err
	}

	return activationToken, nil
}

func (r activationTokenRepo) RetrieveActivationToken(id int) (model.ActivationToken, error) {
	var activationToken model.ActivationToken
	query := r.db.Table("activation_tokens").Where("id = ?", id)
	query = PreloadActivationToken(query)

	if err := query.First(&activationToken).Error; err != nil {
		return model.ActivationToken{}, err
	}
	return activationToken, nil
}

func (r activationTokenRepo) UpdateActivationToken(id int, activationToken model.ActivationToken) (model.ActivationToken, error) {
	if err := r.db.Model(&model.ActivationToken{}).Where("id = ?", id).Updates(&activationToken).Error; err != nil {
		return model.ActivationToken{}, err
	}

	query := r.db.Table("activation_tokens").Where("id = ?", id)
	query = PreloadActivationToken(query)

	if err := query.First(&activationToken).Error; err != nil {
		return model.ActivationToken{}, err
	}

	return activationToken, nil
}

func (r activationTokenRepo) DeleteActivationToken(id int) error {
	if err := r.db.Delete(&model.ActivationToken{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r activationTokenRepo) ListActivationToken(activationToken model.ActivationToken, pagination model.Pagination) ([]model.ActivationToken, error) {
	var activationTokens []model.ActivationToken
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("activation_tokens").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = PreloadActivationToken(query)
	query = FilterActivationToken(query, activationToken)
	query = SearchActivationToken(query, pagination.Search)
	query = query.Find(&activationTokens)
	if err := query.Error; err != nil {
		return nil, err
	}

	return activationTokens, nil
}

func (r activationTokenRepo) ListActivationTokenMeta(activationToken model.ActivationToken, pagination model.Pagination) (model.Meta, error) {
	var activationTokens []model.ActivationToken
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.ActivationToken{}).Select("count(*)")
	queryTotal = FilterActivationToken(queryTotal, activationToken)
	queryTotal = SearchActivationToken(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("activation_tokens").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterActivationToken(query, activationToken)
	query = SearchActivationToken(query, pagination.Search)
	query = query.Find(&activationTokens)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(activationTokens),
	}
	return meta, nil
}

func (r activationTokenRepo) DropDownActivationToken(activationToken model.ActivationToken) ([]model.ActivationToken, error) {
	var activationTokens []model.ActivationToken
	query := r.db.Table("activation_tokens").Order("id desc")
	query = PreloadActivationToken(query)
	query = FilterActivationToken(query, activationToken)
	query = query.Find(&activationTokens)
	if err := query.Error; err != nil {
		return nil, err
	}
	return activationTokens, nil
}

func (r activationTokenRepo) IsValid(token string) (isValid bool, userID uint) {
	var data model.ActivationToken
	if err := r.db.Table("activation_tokens").Where("token = ?", token).First(&data).Error; err != nil {
		log.Printf("[Error] [IsValid] E: %v\n", err)
		isValid = false
		userID = 0
	} else {
		if time.Now().After(data.Valid) {
			isValid = false
			userID = 0
		} else {
			isValid = true
			userID = data.UserID
		}
	}
	return
}

func (r activationTokenRepo) DeleteExpiredActivationToken(currentTime time.Time) error {
	if err := r.db.Table("activation_tokens").Unscoped().Where("valid < ?", currentTime).Delete(&model.ActivationToken{}).Error; err != nil {
		return err
	}
	return nil
}

func FilterActivationToken(query *gorm.DB, activationToken model.ActivationToken) *gorm.DB {
	if activationToken.Token != "" {
		query = query.Where("token LIKE ?", "%"+activationToken.Token+"%")
	}
	if activationToken.UserID > 0 {
		query = query.Where("user_id = ?", activationToken.UserID)
	}
	return query
}

func SearchActivationToken(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("token LIKE ? OR valid LIKE ? OR user_id LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

func PreloadActivationToken(query *gorm.DB) *gorm.DB {
	query = query.Preload("User")
	return query
}
