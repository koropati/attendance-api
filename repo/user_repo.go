package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type UserRepo interface {
	ListUser(user *model.User, pagination *model.Pagination) (*[]model.User, error)
	ListUserMeta(user *model.User, pagination *model.Pagination) (*model.Meta, error)
	CreateUser(user *model.User) (*model.User, error)
	RetrieveUser(id int) (*model.User, error)
	UpdateUser(id int, user *model.User) (*model.User, error)
	DeleteUser(id int) error
	SetActiveUser(id int) (*model.User, error)
	SetDeactiveUser(id int) (*model.User, error)
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) ListUser(user *model.User, pagination *model.Pagination) (*[]model.User, error) {
	var users []model.User
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("users").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterUser(query, user)
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
	if err := r.db.Table("users").Where("id = ?").First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
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
