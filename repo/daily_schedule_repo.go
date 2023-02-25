package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type DailyScheduleRepo interface {
	CreateDailySchedule(dailyschedule *model.DailySchedule) (*model.DailySchedule, error)
	RetrieveDailySchedule(id int) (*model.DailySchedule, error)
	UpdateDailySchedule(id int, dailyschedule *model.DailySchedule) (*model.DailySchedule, error)
	DeleteDailySchedule(id int) error
	ListDailySchedule(dailyschedule *model.DailySchedule, pagination *model.Pagination) (*[]model.DailySchedule, error)
	ListDailyScheduleMeta(dailyschedule *model.DailySchedule, pagination *model.Pagination) (*model.Meta, error)
	DropDownDailySchedule(dailyschedule *model.DailySchedule) (*[]model.DailySchedule, error)
}

type dailyScheduleRepo struct {
	db *gorm.DB
}

func NewDailyScheduleRepo(db *gorm.DB) DailyScheduleRepo {
	return &dailyScheduleRepo{db: db}
}

func (r *dailyScheduleRepo) CreateDailySchedule(dailyschedule *model.DailySchedule) (*model.DailySchedule, error) {
	if err := r.db.Table("daily_schedules").Create(&dailyschedule).Error; err != nil {
		return nil, err
	}
	return dailyschedule, nil
}

func (r *dailyScheduleRepo) RetrieveDailySchedule(id int) (*model.DailySchedule, error) {
	var dailyschedule model.DailySchedule
	if err := r.db.First(&dailyschedule, id).Error; err != nil {
		return nil, err
	}
	return &dailyschedule, nil
}

func (r *dailyScheduleRepo) UpdateDailySchedule(id int, dailyschedule *model.DailySchedule) (*model.DailySchedule, error) {
	if err := r.db.Model(&model.DailySchedule{}).Where("id = ?", id).Updates(&dailyschedule).Error; err != nil {
		return nil, err
	}
	return dailyschedule, nil
}

func (r *dailyScheduleRepo) DeleteDailySchedule(id int) error {
	if err := r.db.Delete(&model.DailySchedule{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *dailyScheduleRepo) ListDailySchedule(dailyschedule *model.DailySchedule, pagination *model.Pagination) (*[]model.DailySchedule, error) {
	var dailyschedules []model.DailySchedule
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("daily_schedules").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterDailySchedule(query, dailyschedule)
	query = query.Find(&dailyschedules)
	if err := query.Error; err != nil {
		return nil, err
	}

	return &dailyschedules, nil
}

func (r *dailyScheduleRepo) ListDailyScheduleMeta(dailyschedule *model.DailySchedule, pagination *model.Pagination) (*model.Meta, error) {
	var dailyschedules []model.DailySchedule
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.DailySchedule{}).Select("count(*)")
	queryTotal = FilterDailySchedule(queryTotal, dailyschedule)
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
		CurrentRecord: len(dailyschedules),
	}
	return &meta, nil
}

func (r *dailyScheduleRepo) DropDownDailySchedule(dailyschedule *model.DailySchedule) (*[]model.DailySchedule, error) {
	var dailyschedules []model.DailySchedule
	query := r.db.Table("daily_schedules")
	query = FilterDailySchedule(query, dailyschedule)
	query = query.Find(&dailyschedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &dailyschedules, nil
}

func FilterDailySchedule(query *gorm.DB, dailyschedule *model.DailySchedule) *gorm.DB {
	if dailyschedule.Name != "" {
		query = query.Where("name LIKE ?", "%"+dailyschedule.Name+"%")
	}
	if dailyschedule.StartTime != "" {
		query = query.Where("start_time LIKE ?", "%"+dailyschedule.StartTime+"%")
	}
	if dailyschedule.EndTime != "" {
		query = query.Where("end_time LIKE ?", "%"+dailyschedule.EndTime+"%")
	}
	if dailyschedule.OwnerID > 0 {
		query = query.Where("owner_id = ?", dailyschedule.OwnerID)
	}
	return query
}
