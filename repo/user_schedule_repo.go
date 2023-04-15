package repo

import (
	"attendance-api/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserScheduleRepo interface {
	CreateUserSchedule(userschedule *model.UserSchedule) (*model.UserSchedule, error)
	RetrieveUserSchedule(id int) (*model.UserSchedule, error)
	RetrieveUserScheduleByOwner(id int, ownerID int) (*model.UserSchedule, error)
	UpdateUserSchedule(id int, userschedule *model.UserSchedule) (*model.UserSchedule, error)
	UpdateUserScheduleByOwner(id int, ownerID int, userschedule *model.UserSchedule) (*model.UserSchedule, error)
	DeleteUserSchedule(id int) error
	DeleteUserScheduleByOwner(id int, ownerID int) error
	ListUserSchedule(userschedule *model.UserSchedule, pagination *model.Pagination) (*[]model.UserSchedule, error)
	ListUserScheduleMeta(userschedule *model.UserSchedule, pagination *model.Pagination) (*model.Meta, error)
	DropDownUserSchedule(userschedule *model.UserSchedule) (*[]model.UserSchedule, error)
	CheckHaveSchedule(userID int, date time.Time) (isHaveSchedule bool, scheduleID int, err error)
	CheckUserInSchedule(scheduleID int, userID int) bool
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

func (r *userScheduleRepo) RetrieveUserScheduleByOwner(id int, ownerID int) (*model.UserSchedule, error) {
	var userschedule model.UserSchedule
	if err := r.db.Model(&model.UserSchedule{}).Where("id = ? AND owner_id = ?", id, ownerID).First(&userschedule).Error; err != nil {
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

func (r *userScheduleRepo) UpdateUserScheduleByOwner(id int, ownerID int, userschedule *model.UserSchedule) (*model.UserSchedule, error) {
	if err := r.db.Model(&model.UserSchedule{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(&userschedule).Error; err != nil {
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

func (r *userScheduleRepo) DeleteUserScheduleByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.UserSchedule{}).Error; err != nil {
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

func (r *userScheduleRepo) CheckHaveSchedule(userID int, date time.Time) (isHaveSchedule bool, scheduleID int, err error) {
	type DataSchedule struct {
		IsHaveSchedule bool `json:"is_have_schedule" query:"is_have_schedule"`
		ScheduleID     int  `json:"schedule_id" query:"schedule_id"`
	}
	var data DataSchedule
	rawQuery := fmt.Sprintf(`SELECT COUNT(*) > 0 as is_have_schedule, us.schedule_id as schedule_id 
	FROM user_schedules us 
	LEFT JOIN schedules s ON us.schedule_id = s.id 
	WHERE us.user_id = %d AND %v BETWEEN s.start_date AND s.end_date`, userID, date)

	if err := r.db.Raw(rawQuery).Scan(&data).Error; err != nil {
		return data.IsHaveSchedule, data.ScheduleID, err
	}
	return data.IsHaveSchedule, data.ScheduleID, nil
}

func (r *userScheduleRepo) CheckUserInSchedule(scheduleID int, userID int) (isHave bool) {
	rawQuery := fmt.Sprintf(`SELECT COUNT(*) > 0 FROM user_schedules us WHERE us.user_id = %d AND us.schedule_id = %d`, userID, scheduleID)
	if err := r.db.Raw(rawQuery).Scan(&isHave).Error; err != nil {
		return false
	}
	return
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
