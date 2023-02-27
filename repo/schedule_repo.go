package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type ScheduleRepo interface {
	CreateSchedule(schedule *model.Schedule) (*model.Schedule, error)
	RetrieveSchedule(id int) (*model.Schedule, error)
	RetrieveScheduleByOwner(id int, ownerID int) (*model.Schedule, error)
	UpdateSchedule(id int, schedule *model.Schedule) (*model.Schedule, error)
	UpdateScheduleByOwner(id int, ownerID int, schedule *model.Schedule) (*model.Schedule, error)
	DeleteSchedule(id int) error
	DeleteScheduleByOwner(id int, ownerID int) error
	ListSchedule(schedule *model.Schedule, pagination *model.Pagination) (*[]model.Schedule, error)
	ListScheduleMeta(schedule *model.Schedule, pagination *model.Pagination) (*model.Meta, error)
	DropDownSchedule(schedule *model.Schedule) (*[]model.Schedule, error)
}

type scheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) ScheduleRepo {
	return &scheduleRepo{db: db}
}

func (r *scheduleRepo) CreateSchedule(schedule *model.Schedule) (*model.Schedule, error) {
	if err := r.db.Table("schedules").Create(&schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

func (r *scheduleRepo) RetrieveSchedule(id int) (*model.Schedule, error) {
	var schedule model.Schedule
	if err := r.db.First(&schedule, id).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepo) RetrieveScheduleByOwner(id int, ownerID int) (*model.Schedule, error) {
	var schedule model.Schedule
	if err := r.db.Model(&model.Schedule{}).Where("id = ? AND owner_id = ?", id, ownerID).First(&schedule).Error; err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepo) UpdateSchedule(id int, schedule *model.Schedule) (*model.Schedule, error) {
	if err := r.db.Model(&model.Schedule{}).Where("id = ?", id).Updates(&schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

func (r *scheduleRepo) UpdateScheduleByOwner(id int, ownerID int, schedule *model.Schedule) (*model.Schedule, error) {
	if err := r.db.Model(&model.Schedule{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(&schedule).Error; err != nil {
		return nil, err
	}
	return schedule, nil
}

func (r *scheduleRepo) DeleteSchedule(id int) error {
	if err := r.db.Delete(&model.Schedule{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *scheduleRepo) DeleteScheduleByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Schedule{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *scheduleRepo) ListSchedule(schedule *model.Schedule, pagination *model.Pagination) (*[]model.Schedule, error) {
	var schedules []model.Schedule
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("schedules").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterSchedule(query, schedule)
	query = query.Find(&schedules)
	if err := query.Error; err != nil {
		return nil, err
	}

	return &schedules, nil
}

func (r *scheduleRepo) ListScheduleMeta(schedule *model.Schedule, pagination *model.Pagination) (*model.Meta, error) {
	var schedules []model.Schedule
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Schedule{}).Select("count(*)")
	queryTotal = FilterSchedule(queryTotal, schedule)
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
		CurrentRecord: len(schedules),
	}
	return &meta, nil
}

func (r *scheduleRepo) DropDownSchedule(schedule *model.Schedule) (*[]model.Schedule, error) {
	var schedules []model.Schedule
	query := r.db.Table("schedules")
	query = FilterSchedule(query, schedule)
	query = query.Find(&schedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &schedules, nil
}

func FilterSchedule(query *gorm.DB, schedule *model.Schedule) *gorm.DB {
	if schedule.Name != "" {
		query = query.Where("name LIKE ?", "%"+schedule.Name+"%")
	}
	if schedule.Code != "" {
		query = query.Where("code LIKE ?", "%"+schedule.Code+"%")
	}
	if schedule.StartDate.String() != "" {
		query = query.Where("start_date LIKE ?", "%"+schedule.StartDate.String()+"%")
	}
	if schedule.EndDate.String() != "" {
		query = query.Where("end_date LIKE ?", "%"+schedule.EndDate.String()+"%")
	}
	if schedule.OwnerID > 0 {
		query = query.Where("owner_id = ?", schedule.OwnerID)
	}
	return query
}
