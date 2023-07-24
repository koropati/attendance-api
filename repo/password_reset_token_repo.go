package repo

import (
	"attendance-api/model"
	"time"

	"gorm.io/gorm"
)

type PasswordResetTokenRepo interface {
	CreatePasswordResetToken(passwordResetToken model.PasswordResetToken) (model.PasswordResetToken, error)
	RetrievePasswordResetToken(id int) (model.PasswordResetToken, error)
	UpdatePasswordResetToken(id int, passwordResetToken model.PasswordResetToken) (model.PasswordResetToken, error)
	DeletePasswordResetToken(id int) error
	ListPasswordResetToken(passwordResetToken model.PasswordResetToken, pagination model.Pagination) ([]model.PasswordResetToken, error)
	ListPasswordResetTokenMeta(passwordResetToken model.PasswordResetToken, pagination model.Pagination) (model.Meta, error)
	DropDownPasswordResetToken(passwordResetToken model.PasswordResetToken) ([]model.PasswordResetToken, error)
	DeleteExpiredPasswordResetToken(currenTime time.Time) error
}

type passwordResetTokenRepo struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepo(db *gorm.DB) PasswordResetTokenRepo {
	return &passwordResetTokenRepo{db: db}
}

func (r passwordResetTokenRepo) CreatePasswordResetToken(passwordResetToken model.PasswordResetToken) (model.PasswordResetToken, error) {
	if err := r.db.Table("password_reset_tokens").Create(&passwordResetToken).Error; err != nil {
		return model.PasswordResetToken{}, err
	}

	query := r.db.Table("password_reset_tokens").Where("id = ?", passwordResetToken.ID)
	query = PreloadPasswordResetToken(query)
	query = query.First(&passwordResetToken)
	if err := query.Error; err != nil {
		return model.PasswordResetToken{}, err
	}

	return passwordResetToken, nil
}

func (r passwordResetTokenRepo) RetrievePasswordResetToken(id int) (model.PasswordResetToken, error) {
	var passwordResetToken model.PasswordResetToken

	query := r.db.Table("password_reset_tokens").Where("id = ?", id)
	query = PreloadPasswordResetToken(query)

	if err := query.First(&passwordResetToken).Error; err != nil {
		return model.PasswordResetToken{}, err
	}

	return passwordResetToken, nil
}

func (r passwordResetTokenRepo) UpdatePasswordResetToken(id int, passwordResetToken model.PasswordResetToken) (model.PasswordResetToken, error) {
	if err := r.db.Model(&model.PasswordResetToken{}).Where("id = ?", id).Updates(&passwordResetToken).Error; err != nil {
		return model.PasswordResetToken{}, err
	}

	query := r.db.Table("password_reset_tokens").Where("id = ?", id)
	query = PreloadPasswordResetToken(query)

	if err := query.First(&passwordResetToken).Error; err != nil {
		return model.PasswordResetToken{}, err
	}
	return passwordResetToken, nil
}

func (r passwordResetTokenRepo) DeletePasswordResetToken(id int) error {
	if err := r.db.Delete(&model.PasswordResetToken{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r passwordResetTokenRepo) ListPasswordResetToken(passwordResetToken model.PasswordResetToken, pagination model.Pagination) ([]model.PasswordResetToken, error) {
	var passwordResetTokens []model.PasswordResetToken
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("password_reset_tokens").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = PreloadPasswordResetToken(query)
	query = FilterPasswordResetToken(query, passwordResetToken)
	query = SearchPasswordResetToken(query, pagination.Search)
	query = query.Find(&passwordResetTokens)
	if err := query.Error; err != nil {
		return nil, err
	}

	return passwordResetTokens, nil
}

func (r passwordResetTokenRepo) ListPasswordResetTokenMeta(passwordResetToken model.PasswordResetToken, pagination model.Pagination) (model.Meta, error) {
	var passwordResetTokens []model.PasswordResetToken
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.PasswordResetToken{}).Select("count(*)")
	queryTotal = FilterPasswordResetToken(queryTotal, passwordResetToken)
	queryTotal = SearchPasswordResetToken(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("password_reset_tokens").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterPasswordResetToken(query, passwordResetToken)
	query = SearchPasswordResetToken(query, pagination.Search)
	query = query.Find(&passwordResetTokens)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(passwordResetTokens),
	}
	return meta, nil
}

func (r passwordResetTokenRepo) DropDownPasswordResetToken(passwordResetToken model.PasswordResetToken) ([]model.PasswordResetToken, error) {
	var passwordResetTokens []model.PasswordResetToken
	query := r.db.Table("password_reset_tokens").Order("id desc")
	query = PreloadPasswordResetToken(query)
	query = FilterPasswordResetToken(query, passwordResetToken)
	query = query.Find(&passwordResetTokens)
	if err := query.Error; err != nil {
		return nil, err
	}
	return passwordResetTokens, nil
}

func (r passwordResetTokenRepo) DeleteExpiredPasswordResetToken(currenTime time.Time) error {
	if err := r.db.Table("password_reset_tokens").Unscoped().Where("valid < ?", currenTime).Delete(&model.PasswordResetToken{}).Error; err != nil {
		return err
	}
	return nil
}

func FilterPasswordResetToken(query *gorm.DB, passwordResetToken model.PasswordResetToken) *gorm.DB {
	if passwordResetToken.Token != "" {
		query = query.Where("token LIKE ?", "%"+passwordResetToken.Token+"%")
	}
	if passwordResetToken.UserID > 0 {
		query = query.Where("user_id = ?", passwordResetToken.UserID)
	}
	return query
}

func SearchPasswordResetToken(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("token LIKE ? OR valid LIKE ? OR user_id LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

func PreloadPasswordResetToken(query *gorm.DB) *gorm.DB {
	query = query.Preload("User")
	return query
}
