package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type AuthRepo interface {
	CheckUsername(username string) bool
	CheckEmail(email string) bool
	CheckHandphone(handphone string) bool
	CheckIsActive(username string) bool
	GetRole(username string) (string, error)
	GetEmail(username string) (string, error)
	Register(user *model.User) error
	Login(username string) (string, error)
	CheckID(id int) bool
	GetByUsername(username string) (user *model.User, err error)
	GetByEmail(email string) (user *model.User, err error)
	Create(user *model.User) error
	Delete(id int) error
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepo {
	return &authRepo{db: db}
}

func (r *authRepo) CheckUsername(username string) bool {
	var count int64
	if err := r.db.Table("users").Where("username = ?", username).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *authRepo) CheckEmail(email string) bool {
	var count int64
	if err := r.db.Table("users").Where("email = ?", email).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *authRepo) CheckHandphone(handphone string) bool {
	var count int64
	if err := r.db.Table("users").Where("handphone = ?", handphone).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *authRepo) CheckIsActive(username string) bool {
	var user model.User
	if err := r.db.Table("users").Select("is_active").Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return user.IsActive
}

func (r *authRepo) GetRole(username string) (string, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}

	return user.Role, nil
}

func (r *authRepo) GetEmail(username string) (string, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}

	return user.Email, nil
}

func (r *authRepo) Register(user *model.User) error {
	if err := r.db.Table("users").Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *authRepo) Login(username string) (string, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}

	return user.Password, nil
}

func (r *authRepo) CheckID(id int) bool {
	var count int64
	if err := r.db.Table("users").Where("id = ?", id).Count(&count).Error; err != nil {
		return false
	}

	if count < 1 {
		return false
	}

	return true
}

func (r *authRepo) GetByUsername(username string) (user *model.User, err error) {
	var userData model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&userData).Error; err != nil {
		return nil, err
	}

	return &userData, nil
}

func (r *authRepo) GetByEmail(email string) (user *model.User, err error) {
	var userData model.User
	if err := r.db.Table("users").Where("email = ?", email).First(&userData).Error; err != nil {
		return nil, err
	}

	return &userData, nil
}

func (r *authRepo) Create(user *model.User) error {
	if err := r.db.Table("users").Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r *authRepo) Delete(id int) error {
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}

	return nil
}
