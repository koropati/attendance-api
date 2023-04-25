package repo

import (
	"attendance-api/model"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type UserScheduleRepo interface {
	CreateUserSchedule(userschedule model.UserSchedule) (model.UserSchedule, error)
	RetrieveUserSchedule(id int) (model.UserSchedule, error)
	RetrieveUserScheduleByOwner(id int, ownerID int) (model.UserSchedule, error)
	ListMySchedule(userID int) ([]model.MySchedule, error)
	ListUserInRule(scheduleID int, user model.User, pagination model.Pagination) ([]model.User, error)
	ListUserInRuleMeta(scheduleID int, user model.User, pagination model.Pagination) (model.Meta, error)
	UpdateUserSchedule(id int, userschedule model.UserSchedule) (model.UserSchedule, error)
	UpdateUserScheduleByOwner(id int, ownerID int, userschedule model.UserSchedule) (model.UserSchedule, error)
	DeleteUserSchedule(id int) error
	DeleteUserScheduleByOwner(id int, ownerID int) error
	ListUserSchedule(userschedule model.UserSchedule, pagination model.Pagination) ([]model.UserSchedule, error)
	ListUserScheduleMeta(userschedule model.UserSchedule, pagination model.Pagination) (model.Meta, error)
	DropDownUserSchedule(userschedule model.UserSchedule) ([]model.UserSchedule, error)
	CheckHaveSchedule(userID int, date time.Time) (isHaveSchedule bool, scheduleID int, err error)
	CheckUserInSchedule(scheduleID int, userID int) bool
	CountByScheduleID(scheduleID int) (total int)
}

type userScheduleRepo struct {
	db *gorm.DB
}

func NewUserScheduleRepo(db *gorm.DB) UserScheduleRepo {
	return userScheduleRepo{db: db}
}

func (r userScheduleRepo) CreateUserSchedule(userschedule model.UserSchedule) (model.UserSchedule, error) {
	if err := r.db.Table("user_schedules").Create(&userschedule).Error; err != nil {
		return model.UserSchedule{}, err
	}
	return userschedule, nil
}

func (r userScheduleRepo) RetrieveUserSchedule(id int) (model.UserSchedule, error) {
	var userschedule model.UserSchedule
	if err := r.db.First(&userschedule, id).Error; err != nil {
		return model.UserSchedule{}, err
	}
	return userschedule, nil
}

func (r userScheduleRepo) RetrieveUserScheduleByOwner(id int, ownerID int) (model.UserSchedule, error) {
	var userschedule model.UserSchedule
	if err := r.db.Model(&model.UserSchedule{}).Where("id = ? AND owner_id = ?", id, ownerID).First(&userschedule).Error; err != nil {
		return model.UserSchedule{}, err
	}
	return userschedule, nil
}

func (r userScheduleRepo) UpdateUserSchedule(id int, userschedule model.UserSchedule) (model.UserSchedule, error) {
	if err := r.db.Model(&model.UserSchedule{}).Where("id = ?", id).Updates(&userschedule).Error; err != nil {
		return model.UserSchedule{}, err
	}
	return userschedule, nil
}

func (r userScheduleRepo) UpdateUserScheduleByOwner(id int, ownerID int, userschedule model.UserSchedule) (model.UserSchedule, error) {
	if err := r.db.Model(&model.UserSchedule{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(&userschedule).Error; err != nil {
		return model.UserSchedule{}, err
	}
	return userschedule, nil
}

func (r userScheduleRepo) DeleteUserSchedule(id int) error {
	if err := r.db.Delete(&model.UserSchedule{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r userScheduleRepo) DeleteUserScheduleByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.UserSchedule{}).Error; err != nil {
		return err
	}
	return nil
}

func (r userScheduleRepo) ListMySchedule(userID int) (results []model.MySchedule, err error) {
	today := time.Now().Format("2006-01-02")
	query := fmt.Sprintf(`
	SELECT 
	us.schedule_id as schedule_id, 
	s.name as schedule_name, 
	s.code as schedule_code, 
	s.start_date as start_date, 
	s.end_date as end_date, 
	s.subject_id as subject_id, 
	sbj.name as subject_name, 
	sbj.code as subject_code, 
	s.late_duration as late_duration, 
	s.latitude as latitude, 
	s.longitude as longitude, 
	s.radius as radius 
	FROM user_schedules us 
	LEFT JOIN schedules s ON us.schedule_id = s.id 
	LEFT JOIN subjects sbj ON s.subject_id = sbj.id 
	WHERE us.user_id = %d AND '%s' BETWEEN DATE(s.start_date) AND DATE(s.end_date) AND us.deleted_at IS NULL`, userID, today)
	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	return
}

func (r userScheduleRepo) ListUserSchedule(userschedule model.UserSchedule, pagination model.Pagination) ([]model.UserSchedule, error) {
	var userschedules []model.UserSchedule
	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("user_schedules").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterUserSchedule(query, userschedule)
	query = SearchUserSchedule(query, pagination.Search)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return nil, err
	}

	return userschedules, nil
}

func (r userScheduleRepo) ListUserInRule(scheduleID int, user model.User, pagination model.Pagination) ([]model.User, error) {
	var userID []int
	if err := r.db.Table("user_schedules").Select("user_id").Where("schedule_id = ?", scheduleID).Find(&userID).Error; err != nil {
		return nil, err
	}
	var users []model.User
	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("users").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = query.Where("id IN ?", userID)
	query = FilterUser(query, user)
	query = SearchUser(query, pagination.Search)
	query = query.Find(&users)
	if err := query.Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r userScheduleRepo) ListUserScheduleMeta(userschedule model.UserSchedule, pagination model.Pagination) (model.Meta, error) {
	var userschedules []model.UserSchedule
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.UserSchedule{}).Select("count(*)")
	queryTotal = FilterUserSchedule(queryTotal, userschedule)
	queryTotal = SearchUserSchedule(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("user_schedules").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterUserSchedule(query, userschedule)
	query = SearchUserSchedule(query, pagination.Search)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(userschedules),
	}
	return meta, nil
}

func (r userScheduleRepo) ListUserInRuleMeta(scheduleID int, user model.User, pagination model.Pagination) (model.Meta, error) {
	var users []model.User
	var totalRecord int
	var totalPage int

	var userID []int
	if err := r.db.Table("user_schedules").Select("user_id").Where("schedule_id = ?", scheduleID).Find(&userID).Error; err != nil {
		return model.Meta{}, err
	}

	queryTotal := r.db.Model(&model.User{}).Select("count(*)")
	queryTotal = queryTotal.Where("id IN ?", userID)
	queryTotal = FilterUser(queryTotal, user)
	queryTotal = SearchUser(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("users").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = query.Where("id IN ?", userID)
	query = FilterUser(query, user)
	query = SearchUser(query, pagination.Search)
	query = query.Find(&users)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(users),
	}
	return meta, nil
}

func (r userScheduleRepo) DropDownUserSchedule(userschedule model.UserSchedule) ([]model.UserSchedule, error) {
	var userschedules []model.UserSchedule
	query := r.db.Table("user_schedules")
	query = FilterUserSchedule(query, userschedule)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	return userschedules, nil
}

func (r userScheduleRepo) CheckHaveSchedule(userID int, date time.Time) (isHaveSchedule bool, scheduleID int, err error) {
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

func (r userScheduleRepo) CheckUserInSchedule(scheduleID int, userID int) (isHave bool) {
	rawQuery := fmt.Sprintf(`SELECT COUNT(*) > 0 FROM user_schedules us WHERE us.user_id = %d AND us.schedule_id = %d`, userID, scheduleID)
	if err := r.db.Raw(rawQuery).Scan(&isHave).Error; err != nil {
		return false
	}
	return
}

func (r userScheduleRepo) CountByScheduleID(scheduleID int) (total int) {
	if err := r.db.Table("user_schedules").Select("count(*)").Where("schedule_id = ? AND user_id != ?", scheduleID, 0).Find(&total).Error; err != nil {
		total = 0
	}
	return
}

func FilterUserSchedule(query *gorm.DB, userschedule model.UserSchedule) *gorm.DB {
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

func SearchUserSchedule(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("user_id LIKE ? OR schedule_id LIKE ? OR owner_id LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}
