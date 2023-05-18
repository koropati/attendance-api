package repo

import (
	"attendance-api/model"

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
}

type passwordResetTokenRepo struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepo(db *gorm.DB) PasswordResetTokenRepo {
	return &passwordResetTokenRepo{db: db}
}

func (r passwordResetTokenRepo) CreatePasswordResetToken(subject model.PasswordResetToken) (model.PasswordResetToken, error) {
	if err := r.db.Table("password_reset_tokens").Create(&subject).Error; err != nil {
		return model.PasswordResetToken{}, err
	}

	return subject, nil
}

func (r passwordResetTokenRepo) RetrievePasswordResetToken(id int) (model.PasswordResetToken, error) {
	var subject model.PasswordResetToken
	if err := r.db.First(&subject, id).Error; err != nil {
		return model.PasswordResetToken{}, err
	}
	return subject, nil
}

func (r passwordResetTokenRepo) UpdatePasswordResetToken(id int, subject model.PasswordResetToken) (model.PasswordResetToken, error) {
	if err := r.db.Model(&model.PasswordResetToken{}).Where("id = ?", id).Updates(&subject).Error; err != nil {
		return model.PasswordResetToken{}, err
	}
	return subject, nil
}

func (r passwordResetTokenRepo) DeletePasswordResetToken(id int) error {
	if err := r.db.Delete(&model.PasswordResetToken{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r passwordResetTokenRepo) ListPasswordResetToken(subject model.PasswordResetToken, pagination model.Pagination) ([]model.PasswordResetToken, error) {
	var passwordResetTokens []model.PasswordResetToken
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("password_reset_tokens").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterPasswordResetToken(query, subject)
	query = SearchPasswordResetToken(query, pagination.Search)
	query = query.Find(&passwordResetTokens)
	if err := query.Error; err != nil {
		return nil, err
	}

	return passwordResetTokens, nil
}

func (r passwordResetTokenRepo) ListPasswordResetTokenMeta(subject model.PasswordResetToken, pagination model.Pagination) (model.Meta, error) {
	var passwordResetTokens []model.PasswordResetToken
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.PasswordResetToken{}).Select("count(*)")
	queryTotal = FilterPasswordResetToken(queryTotal, subject)
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
	query = FilterPasswordResetToken(query, subject)
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

func (r passwordResetTokenRepo) DropDownPasswordResetToken(subject model.PasswordResetToken) ([]model.PasswordResetToken, error) {
	var passwordResetTokens []model.PasswordResetToken
	query := r.db.Table("password_reset_tokens").Order("id desc")
	query = FilterPasswordResetToken(query, subject)
	query = query.Find(&passwordResetTokens)
	if err := query.Error; err != nil {
		return nil, err
	}
	return passwordResetTokens, nil
}

func FilterPasswordResetToken(query *gorm.DB, subject model.PasswordResetToken) *gorm.DB {
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

func SearchPasswordResetToken(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("token LIKE ? OR valid LIKE ? OR user_id LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}
