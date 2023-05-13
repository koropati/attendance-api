package repo

import (
	"attendance-api/model"
	"time"

	"github.com/twinj/uuid"
	"gorm.io/gorm"
)

type AuthRepo interface {
	CheckID(id int) bool
	CheckUsername(username string) bool
	CheckEmail(email string) bool
	CheckHandphone(handphone string) bool
	CheckIsActive(username string) bool
	IsSuperAdmin(username string) (bool, error)
	IsAdmin(username string) (bool, error)
	IsUser(username string) (bool, error)
	GetRole(username string) (bool, bool, bool, error)
	GetEmail(username string) (string, error)
	Register(user model.User) error
	Login(username string) (string, error)
	GetByUsername(username string) (user model.User, err error)
	GetByEmail(email string) (user model.User, err error)
	GetByID(id uint) (user model.User, err error)
	Create(user model.User) error
	Delete(id int) error
	SetActiveUser(id int) (model.User, error)
	SetDeactiveUser(id int) (model.User, error)
	FetchAuth(userID uint, authUUID string) (model.Auth, error)
	DeleteAuth(userID uint, authUUID string) error
	CreateAuth(userID uint, expired int64, typeAuth string) (model.Auth, error)
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepo {
	return &authRepo{db: db}
}

func (r authRepo) CheckID(id int) bool {
	var count int64
	if err := r.db.Table("users").Where("id = ?", id).Count(&count).Error; err != nil {
		return false
	}

	if count < 1 {
		return false
	}

	return true
}

func (r authRepo) CheckUsername(username string) (isExist bool) {
	if err := r.db.Table("users").Select("count(*) > 0").Where("username = ?", username).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r authRepo) CheckEmail(email string) bool {
	var count int64
	if err := r.db.Table("users").Where("email = ?", email).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r authRepo) CheckHandphone(handphone string) bool {
	var count int64
	if err := r.db.Table("users").Where("handphone = ?", handphone).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r authRepo) CheckIsActive(username string) bool {
	var user model.User
	if err := r.db.Table("users").Select("is_active").Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return user.IsActive
}

func (r authRepo) IsSuperAdmin(username string) (bool, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return false, err
	}

	return user.IsSuperAdmin, nil
}

func (r authRepo) IsAdmin(username string) (bool, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return false, err
	}

	return user.IsAdmin, nil
}

func (r authRepo) IsUser(username string) (bool, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return false, err
	}

	return user.IsUser, nil
}

func (r authRepo) GetRole(username string) (bool, bool, bool, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return false, false, false, err
	}

	return user.IsSuperAdmin, user.IsAdmin, user.IsUser, nil
}

func (r authRepo) GetEmail(username string) (string, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}

	return user.Email, nil
}

func (r authRepo) Register(user model.User) error {
	if err := r.db.Table("users").Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r authRepo) Login(username string) (string, error) {
	var user model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&user).Error; err != nil {
		return "", err
	}
	if err := r.db.Table("users").Where("username = ?", username).Update("last_login", time.Now()).Error; err != nil {
		return "", err
	}

	return user.Password, nil
}

func (r authRepo) GetByUsername(username string) (user model.User, err error) {
	var userData model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&userData).Error; err != nil {
		return model.User{}, err
	}

	return userData, nil
}

func (r authRepo) GetByEmail(email string) (user model.User, err error) {
	var userData model.User
	if err := r.db.Table("users").Where("email = ?", email).First(&userData).Error; err != nil {
		return model.User{}, err
	}

	return userData, nil
}

func (r authRepo) GetByID(id uint) (user model.User, err error) {
	var userData model.User
	if err := r.db.Table("users").Where("id = ?", id).First(&userData).Error; err != nil {
		return model.User{}, err
	}

	return userData, nil
}

func (r authRepo) Create(user model.User) error {
	if err := r.db.Table("users").Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (r authRepo) Delete(id int) error {
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}

	return nil
}

func (r authRepo) SetActiveUser(id int) (model.User, error) {
	var user model.User
	if err := r.db.Model(&user).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return model.User{}, err
	}
	query := r.db.Table("users")
	if err := query.Where("id = ?", id).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r authRepo) SetDeactiveUser(id int) (model.User, error) {
	var user model.User
	if err := r.db.Model(&user).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return model.User{}, err
	}
	query := r.db.Table("users")
	if err := query.Where("id = ?", id).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r authRepo) FetchAuth(userID uint, authUUID string) (model.Auth, error) {
	var auth model.Auth
	if err := r.db.Table("auths").Where("user_id = ? AND auth_uuid = ?", userID, authUUID).First(&auth).Error; err != nil {
		return model.Auth{}, err
	}
	return auth, nil
}

func (r authRepo) DeleteAuth(userID uint, authUUID string) error {
	if err := r.db.Table("auths").Unscoped().Where("user_id = ? AND auth_uuid = ?", userID, authUUID).Delete(&model.Auth{}).Error; err != nil {
		return err
	}
	return nil
}

func (r authRepo) CreateAuth(userID uint, expired int64, typeAuth string) (model.Auth, error) {
	var auth model.Auth
	auth.UserID = userID
	auth.AuthUUID = uuid.NewV4().String()
	auth.Expired = expired
	auth.TypeAuth = typeAuth

	if err := r.db.Create(&auth).Error; err != nil {
		return model.Auth{}, err
	}
	return auth, nil
}
