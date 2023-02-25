package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type UserScheduleRepo interface {
	CreateUserSchedule(userschedule *model.UserSchedule) (*model.UserSchedule, error)
	RetrieveUserSchedule(id int) (*model.UserSchedule, error)
	UpdateUserSchedule(id int, userschedule *model.UserSchedule) (*model.UserSchedule, error)
	DeleteUserSchedule(id int) error
	ListUserSchedule(userschedule *model.UserSchedule, pagination *model.Pagination) (*[]model.UserSchedule, error)
	ListUserScheduleMeta(userschedule *model.UserSchedule, pagination *model.Pagination) (*model.Meta, error)
	DropDownUserSchedule(userschedule *model.UserSchedule) (*[]model.UserSchedule, error)
}

type userScheduleRepo struct {
	db *gorm.DB
}

func NewUserScheduleRepo(db *gorm.DB) UserScheduleRepo {
	return &userScheduleRepo{db: db}
}

func (r *userScheduleRepo) CreateUserSchedule(userschedule *model.UserSchedule) (*model.UserSchedule, error) {
	if err := r.db.Table("user_schedules").Create(&userschedule).Error; err != nil {
		return nil, err
	}
	return userschedule, nil
}

func (r *userScheduleRepo) RetrieveUserSchedule(id int) (*model.UserSchedule, error) {
	var userschedule model.UserSchedule
	if err := r.db.First(&userschedule, id).Error; err != nil {
		return nil, err
	}
	return &userschedule, nil
}

func (r *userScheduleRepo) UpdateUserSchedule(id int, userschedule *model.UserSchedule) (*model.UserSchedule, error) {
	if err := r.db.Model(&model.UserSchedule{}).Where("id = ?", id).Updates(&userschedule).Error; err != nil {
		return nil, err
	}
	return userschedule, nil
}

func (r *userScheduleRepo) DeleteUserSchedule(id int) error {
	if err := r.db.Delete(&model.UserSchedule{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *userScheduleRepo) ListUserSchedule(userschedule *model.UserSchedule, pagination *model.Pagination) (*[]model.UserSchedule, error) {
	var userschedules []model.UserSchedule
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("user_schedules").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterUserSchedule(query, userschedule)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return nil, err
	}

	return &userschedules, nil
}

func (r *userScheduleRepo) ListUserScheduleMeta(userschedule *model.UserSchedule, pagination *model.Pagination) (*model.Meta, error) {
	var userschedules []model.UserSchedule
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.UserSchedule{}).Select("count(*)")
	queryTotal = FilterUserSchedule(queryTotal, userschedule)
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
		CurrentRecord: len(userschedules),
	}
	return &meta, nil
}

func (r *userScheduleRepo) DropDownUserSchedule(userschedule *model.UserSchedule) (*[]model.UserSchedule, error) {
	var userschedules []model.UserSchedule
	query := r.db.Table("user_schedules")
	query = FilterUserSchedule(query, userschedule)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &userschedules, nil
}

func FilterUserSchedule(query *gorm.DB, userschedule *model.UserSchedule) *gorm.DB {
	if userschedule.UserID > 0 {
		query = query.Where("user_id = ?", userschedule.UserID)
	}
	if userschedule.ScheduleID > 0 {
		query = query.Where("schedule_id = ?", userschedule.ScheduleID)
	}
	if userschedule.OwnerID > 0 {
		query = query.Where("owner_id = ?", userschedule.OwnerID)
	}
	return query
}
