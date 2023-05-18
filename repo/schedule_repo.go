package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type ScheduleRepo interface {
	CreateSchedule(schedule model.Schedule) (model.Schedule, error)
	RetrieveSchedule(id int) (model.Schedule, error)
	RetrieveScheduleByOwner(id int, ownerID int) (model.Schedule, error)
	RetrieveScheduleByQRcode(QRcode string) (model.Schedule, error)
	UpdateSchedule(id int, schedule model.Schedule) (model.Schedule, error)
	UpdateScheduleByOwner(id int, ownerID int, schedule model.Schedule) (model.Schedule, error)
	UpdateQRcode(id int, QRcode string) (model.Schedule, error)
	UpdateQRcodeByOwner(id int, ownerID int, QRcode string) (model.Schedule, error)
	DeleteSchedule(id int) error
	DeleteScheduleByOwner(id int, ownerID int) error
	ListSchedule(schedule model.Schedule, pagination model.Pagination) ([]model.Schedule, error)
	ListScheduleMeta(schedule model.Schedule, pagination model.Pagination) (model.Meta, error)
	DropDownSchedule(schedule model.Schedule) ([]model.Schedule, error)
	CheckIsExist(id int) (isExist bool, err error)
	CheckCodeIsExist(code string, exceptID int) bool
}

type scheduleRepo struct {
	db *gorm.DB
}

func NewScheduleRepo(db *gorm.DB) ScheduleRepo {
	return &scheduleRepo{db: db}
}

func (r scheduleRepo) CreateSchedule(schedule model.Schedule) (model.Schedule, error) {

	query := r.db.Table("schedules")
	query = PreloadSchedule(query)
	if err := query.Create(&schedule).Error; err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}

func (r scheduleRepo) RetrieveSchedule(id int) (model.Schedule, error) {
	var schedule model.Schedule
	query := r.db.Table("schedules")
	query = PreloadSchedule(query)
	if err := query.First(schedule, id).Error; err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}

func (r scheduleRepo) RetrieveScheduleByOwner(id int, ownerID int) (model.Schedule, error) {
	var schedule model.Schedule
	query := r.db.Model(&model.Schedule{})
	query = PreloadSchedule(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).First(&schedule).Error; err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}

func (r scheduleRepo) RetrieveScheduleByQRcode(QRcode string) (model.Schedule, error) {
	var schedule model.Schedule
	query := r.db.Model(&model.Schedule{})
	query = PreloadSchedule(query)
	if err := query.Where("qr_code = ?", QRcode).First(&schedule).Error; err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}

func (r scheduleRepo) UpdateSchedule(id int, schedule model.Schedule) (model.Schedule, error) {
	query := r.db.Model(&model.Schedule{})
	query = PreloadSchedule(query)
	if err := query.Where("id = ?", id).Updates(&schedule).Error; err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}

func (r scheduleRepo) UpdateScheduleByOwner(id int, ownerID int, schedule model.Schedule) (model.Schedule, error) {
	if err := r.db.Model(&model.Schedule{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(&schedule).Error; err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}

func (r scheduleRepo) UpdateQRcode(id int, QRcode string) (schedule model.Schedule, err error) {
	if err := r.db.Model(&model.Schedule{}).Where("id = ?", id).Update("qr_code", QRcode).Find(&schedule).Error; err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}

func (r scheduleRepo) UpdateQRcodeByOwner(id int, ownerID int, QRcode string) (schedule model.Schedule, err error) {
	if err := r.db.Model(&model.Schedule{}).Where("id = ? AND owner_id = ?", id, ownerID).Update("qr_code", QRcode).Find(&schedule).Error; err != nil {
		return model.Schedule{}, err
	}
	return schedule, nil
}

func (r scheduleRepo) DeleteSchedule(id int) error {
	if err := r.db.Delete(&model.Schedule{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r scheduleRepo) DeleteScheduleByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Schedule{}).Error; err != nil {
		return err
	}
	return nil
}

func (r scheduleRepo) ListSchedule(schedule model.Schedule, pagination model.Pagination) ([]model.Schedule, error) {
	var schedules []model.Schedule
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("schedules").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = PreloadSchedule(query)
	query = FilterSchedule(query, schedule)
	query = SearchSchedule(query, pagination.Search)
	query = query.Find(&schedules)
	if err := query.Error; err != nil {
		return nil, err
	}

	return schedules, nil
}

func (r scheduleRepo) ListScheduleMeta(schedule model.Schedule, pagination model.Pagination) (model.Meta, error) {
	var schedules []model.Schedule
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Schedule{}).Select("count(*)")
	queryTotal = FilterSchedule(queryTotal, schedule)
	queryTotal = SearchSchedule(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("schedules").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterSchedule(query, schedule)
	query = SearchSchedule(query, pagination.Search)
	query = query.Find(&schedules)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(schedules),
	}
	return meta, nil
}

func (r scheduleRepo) DropDownSchedule(schedule model.Schedule) ([]model.Schedule, error) {
	var schedules []model.Schedule
	query := r.db.Table("schedules").Order("id desc")
	query = PreloadSchedule(query)
	query = FilterSchedule(query, schedule)
	query = query.Find(&schedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	return schedules, nil
}

func (r scheduleRepo) CheckIsExist(id int) (isExist bool, err error) {
	if err := r.db.Table("schedules").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false, err
	}
	return
}

func (r scheduleRepo) CheckCodeIsExist(code string, exceptID int) bool {
	var count int64
	if err := r.db.Table("schedules").Where("code = ? AND id != ?", code, exceptID).Count(&count).Error; err != nil {
		return true
	}

	if count > 0 {
		return true
	}

	return false
}

func FilterSchedule(query *gorm.DB, schedule model.Schedule) *gorm.DB {
	if schedule.Name != "" {
		query = query.Where("name LIKE ?", "%"+schedule.Name+"%")
	}
	if schedule.Code != "" {
		query = query.Where("code LIKE ?", "%"+schedule.Code+"%")
	}
	if schedule.StartDate != "" {
		query = query.Where("start_date LIKE ?", "%"+schedule.StartDate+"%")
	}
	if schedule.EndDate != "" {
		query = query.Where("end_date LIKE ?", "%"+schedule.EndDate+"%")
	}
	if schedule.OwnerID > 0 {
		query = query.Where("owner_id = ?", schedule.OwnerID)
	}
	return query
}

func SearchSchedule(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR qr_code LIKE ? OR start_date LIKE ? OR end_date LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

func PreloadSchedule(query *gorm.DB) *gorm.DB {
	query = query.Preload("Subject")
	query = query.Preload("DailySchedule")
	return query
}
