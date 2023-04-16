package repo

import (
	"attendance-api/model"
	"time"

	"gorm.io/gorm"
)

type UserRepo interface {
	CheckID(id int) bool
	CheckUsername(username string) bool
	CheckEmail(email string) bool
	CheckHandphone(handphone string) bool
	CheckIsActive(username string) bool
	CheckUpdateUsername(id int, username string) bool
	CheckUpdateEmail(id int, email string) bool
	CheckUpdateHandphone(id int, handphone string) bool
	ListUser(user *model.User, pagination *model.Pagination) (*[]model.User, error)
	ListUserMeta(user *model.User, pagination *model.Pagination) (*model.Meta, error)
	CreateUser(user *model.User) (*model.User, error)
	RetrieveUser(id int) (*model.User, error)
	RetrieveUserByUsername(username string) (user *model.User, err error)
	RetrieveUserByEmail(email string) (user *model.User, err error)
	UpdateUser(id int, user *model.User) (*model.User, error)
	DeleteUser(id int) error
	SetActiveUser(id int) (*model.User, error)
	SetDeactiveUser(id int) (*model.User, error)
	DropDownUser(user *model.User) (*[]model.UserDropDown, error)
	GetPassword(id int) (hashPassword string, err error)
	UpdatePassword(updatePasswordData *model.UserUpdatePasswordForm) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CheckID(id int) bool {
	var count int64
	if err := r.db.Table("users").Where("id = ?", id).Count(&count).Error; err != nil {
		return false
	}

	if count < 1 {
		return false
	}

	return true
}

func (r *userRepo) CheckUsername(username string) bool {
	var count int64
	if err := r.db.Table("users").Where("username = ?", username).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *userRepo) CheckEmail(email string) bool {
	var count int64
	if err := r.db.Table("users").Where("email = ?", email).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *userRepo) CheckHandphone(handphone string) bool {
	var count int64
	if err := r.db.Table("users").Where("handphone = ?", handphone).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *userRepo) CheckIsActive(username string) bool {
	var user model.User
	if err := r.db.Table("users").Select("is_active").Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return user.IsActive
}

func (r *userRepo) CheckUpdateUsername(id int, username string) bool {
	var count int64
	if err := r.db.Table("users").Where("username = ? AND id != ?", username, id).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *userRepo) CheckUpdateEmail(id int, email string) bool {
	var count int64
	if err := r.db.Table("users").Where("email = ? AND id != ?", email, id).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *userRepo) CheckUpdateHandphone(id int, handphone string) bool {
	var count int64
	if err := r.db.Table("users").Where("handphone = ? AND id != ?", handphone, id).Count(&count).Error; err != nil {
		return false
	}

	if count > 0 {
		return false
	}

	return true
}

func (r *userRepo) ListUser(user *model.User, pagination *model.Pagination) (*[]model.User, error) {
	var users []model.User
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("users").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterUser(query, user)
	query = SearchUser(query, pagination.Search)
	query = query.Find(&users)
	if err := query.Error; err != nil {
		return nil, err
	}

	return &users, nil
}

func (r *userRepo) ListUserMeta(user *model.User, pagination *model.Pagination) (*model.Meta, error) {
	var users []model.User
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.User{}).Select("count(*)")
	queryTotal = FilterUser(queryTotal, user)
	queryTotal = SearchUser(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return nil, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("users").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterUser(query, user)
	query = SearchUser(query, pagination.Search)
	query = query.Find(&users)
	if err := query.Error; err != nil {
		return nil, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(users),
	}

	return &meta, nil
}

func (r *userRepo) CreateUser(user *model.User) (*model.User, error) {
	if err := r.db.Table("users").Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) RetrieveUser(id int) (user *model.User, err error) {
	if err := r.db.Table("users").Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) RetrieveUserByUsername(username string) (user *model.User, err error) {
	var userData model.User
	if err := r.db.Table("users").Where("username = ?", username).First(&userData).Error; err != nil {
		return nil, err
	}

	return &userData, nil
}

func (r *userRepo) RetrieveUserByEmail(email string) (user *model.User, err error) {
	var userData model.User
	if err := r.db.Table("users").Where("email = ?", email).First(&userData).Error; err != nil {
		return nil, err
	}

	return &userData, nil
}

func (r *userRepo) UpdateUser(id int, user *model.User) (*model.User, error) {
	if err := r.db.Model(&model.User{}).Where("id = ?", id).Updates(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepo) DeleteUser(id int) error {
	if err := r.db.Delete(&model.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepo) SetActiveUser(id int) (*model.User, error) {
	var user model.User
	if err := r.db.Model(&user).Where("id = ?", id).Update("is_active", true).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) SetDeactiveUser(id int) (*model.User, error) {
	var user model.User
	if err := r.db.Model(&user).Where("id = ?", id).Update("is_active", false).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) DropDownUser(user *model.User) (results *[]model.UserDropDown, err error) {
	query := r.db.Table("users")
	query = FilterUser(query, user)
	query = query.Find(&results)
	if err := query.Error; err != nil {
		return nil, err
	}
	return
}

func (r *userRepo) GetPassword(id int) (hashPassword string, err error) {
	if err := r.db.Select("password").Model(&model.User{}).Where("id = ?", id).First(&hashPassword).Error; err != nil {
		return "", err
	}

	return
}

func (r *userRepo) UpdatePassword(userPasswordData *model.UserUpdatePasswordForm) (err error) {
	query := r.db.Model(&model.User{}).Where("id = ?", userPasswordData.ID)

	if err := query.Updates(map[string]interface{}{
		"password":   userPasswordData.NewPassword,
		"updated_at": time.Now(),
	}).Error; err != nil {
		return err
	}
	return
}

func FilterUser(query *gorm.DB, user *model.User) *gorm.DB {
	if user.Username != "" {
		query = query.Where("username LIKE ?", "%"+user.Username+"%")
	}
	if user.FirstName != "" {
		query = query.Where("first_name LIKE ?", "%"+user.FirstName+"%")
	}
	if user.LastName != "" {
		query = query.Where("last_name LIKE ?", "%"+user.LastName+"%")
	}
	if user.Handphone != "" {
		query = query.Where("handphone LIKE ?", "%"+user.Handphone+"%")
	}
	if user.Email != "" {
		query = query.Where("email LIKE ?", "%"+user.Email+"%")
	}

	return query
}

func SearchUser(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("username LIKE ? OR first_name LIKE ? OR last_name LIKE ? OR handphone LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}
