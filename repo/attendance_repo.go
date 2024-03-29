package repo

import (
	"attendance-api/model"
	"errors"
	"sync"

	"gorm.io/gorm"
)

type AttendanceRepo interface {
	CreateAttendance(attendance model.Attendance) (model.Attendance, error)
	RetrieveAttendance(id int) (model.Attendance, error)
	RetrieveAttendanceByUserID(id int, userID int) (model.Attendance, error)
	RetrieveAttendanceByDate(userID int, scheduleID int, date string) (model.Attendance, error)
	UpdateAttendance(id int, attendance model.Attendance) (model.Attendance, error)
	UpdateAttendanceByUserID(id int, userID int, attendance model.Attendance) (model.Attendance, error)
	UpdateStatusAttendance(id int, statusPresence string, userID int) (model.Attendance, error)
	DeleteAttendance(id int) error
	DeleteAttendanceByUserID(id int, userID int) error
	ListAttendance(attendance model.Attendance, pagination model.Pagination) ([]model.Attendance, error)
	ListAttendanceMeta(attendance model.Attendance, pagination model.Pagination) (model.Meta, error)
	DropDownAttendance(attendance model.Attendance) ([]model.Attendance, error)
	CheckIsExist(id int) (isExist bool, err error)
	CheckIsExistByDate(userID int, scheduleID int, date string) bool
	CountAttendanceByStatus(userID int, statusAttendance string, startDate string, endDate string) (result int)
}

type attendanceRepo struct {
	db *gorm.DB
}

func NewAttendanceRepo(db *gorm.DB) AttendanceRepo {
	return &attendanceRepo{db: db}
}

func (r attendanceRepo) CreateAttendance(attendance model.Attendance) (result model.Attendance, err error) {
	if err := r.db.Table("attendances").Create(&attendance).Error; err != nil {
		return model.Attendance{}, err
	}

	if err := PreloadAttendance(r.db.Table("attendances")).Where("id = ?", attendance.ID).First(&result).Error; err != nil {
		return model.Attendance{}, err
	}

	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility(r.db)
	return
}

func (r attendanceRepo) RetrieveAttendance(id int) (result model.Attendance, err error) {
	if err := PreloadAttendance(r.db.Table("attendances")).Where("id = ?", id).First(&result).Error; err != nil {
		return model.Attendance{}, err
	}

	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility(r.db)
	return
}

func (r attendanceRepo) RetrieveAttendanceByUserID(id int, userID int) (result model.Attendance, err error) {
	if err := PreloadAttendance(r.db.Table("attendances")).Joins("JOIN schedules ON attendances.schedule_id = schedules.id").Where("attendances.id = ? AND schedules.owner_id = ?", id, userID).First(&result).Error; err != nil {
		return model.Attendance{}, err
	}

	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility(r.db)
	return
}

func (r attendanceRepo) RetrieveAttendanceByDate(userID int, scheduleID int, date string) (result model.Attendance, err error) {
	if err := PreloadAttendance(r.db.Table("attendances")).Where("user_id = ? AND schedule_id = ? AND DATE(date) = ?", userID, scheduleID, date).First(&result).Error; err != nil {
		return model.Attendance{}, err
	}
	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility(r.db)
	return
}

func (r attendanceRepo) UpdateAttendance(id int, attendance model.Attendance) (result model.Attendance, err error) {
	if err := r.db.Table("attendances").Where("id = ?", id).Updates(&attendance).Error; err != nil {
		return model.Attendance{}, err
	}
	if err := PreloadAttendance(r.db.Table("attendances")).Where("id = ?", id).First(&result).Error; err != nil {
		return model.Attendance{}, err
	}

	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility(r.db)
	return
}

func (r attendanceRepo) UpdateStatusAttendance(id int, statusPresence string, userID int) (result model.Attendance, err error) {
	if statusPresence != "" {
		if err := r.db.Table("attendances").Where("id = ?", id).Updates(map[string]interface{}{
			"status_presence": statusPresence,
			"status":          "-",
			"updated_by":      userID,
		}).Error; err != nil {
			return model.Attendance{}, err
		}
		if err := PreloadAttendance(r.db.Table("attendances")).Where("id = ?", id).First(&result).Error; err != nil {
			return model.Attendance{}, err
		}

		result.User.Role = result.User.GetRole()
		result.User.Avatar = result.User.GetAvatar()
		result.User.UserAbilities = result.User.GetAbility(r.db)
		return
	} else {
		err = errors.New("status presensi tidak boleh kosong")
		return
	}

}

func (r attendanceRepo) UpdateAttendanceByUserID(id int, userID int, attendance model.Attendance) (result model.Attendance, err error) {
	if err := r.db.Model(&model.Attendance{}).Where("id = ? AND user_Id = ?", id, userID).Updates(&attendance).Error; err != nil {
		return model.Attendance{}, err
	}

	if err := PreloadAttendance(r.db.Table("attendances")).Where("id = ?", id).First(&result).Error; err != nil {
		return model.Attendance{}, err
	}

	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility(r.db)
	return
}

func (r attendanceRepo) DeleteAttendance(id int) error {
	if err := r.db.Delete(&model.Attendance{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r attendanceRepo) DeleteAttendanceByUserID(id int, userID int) error {
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Attendance{}).Error; err != nil {
		return err
	}
	return nil
}

func (r attendanceRepo) ListAttendance(attendance model.Attendance, pagination model.Pagination) ([]model.Attendance, error) {
	var attendances []model.Attendance
	offset := (pagination.Page - 1) * pagination.Limit

	query := PreloadAttendance(r.db.Table("attendances")).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterAttendance(query, attendance)
	query = SearchAttendance(query, pagination.Search)
	query = query.Find(&attendances)
	if err := query.Error; err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	for i, data := range attendances {
		wg.Add(1)
		go func(i int, data model.Attendance) {
			attendances[i].User.Role = data.User.GetRole()
			attendances[i].User.Avatar = data.User.GetAvatar()
			attendances[i].User.UserAbilities = data.User.GetAbility(r.db)
			wg.Done()
		}(i, data)
	}
	wg.Wait()
	return attendances, nil
}

func (r attendanceRepo) ListAttendanceMeta(attendance model.Attendance, pagination model.Pagination) (model.Meta, error) {
	var attendances []model.Attendance
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Attendance{}).Select("count(*)")
	queryTotal = FilterAttendance(queryTotal, attendance)
	queryTotal = SearchAttendance(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("attendances").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterAttendance(query, attendance)
	query = SearchAttendance(query, pagination.Search)
	query = query.Find(&attendances)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(attendances),
	}
	return meta, nil
}

func (r attendanceRepo) DropDownAttendance(attendance model.Attendance) ([]model.Attendance, error) {
	var attendances []model.Attendance
	query := PreloadAttendance(r.db.Table("attendances")).Order("id desc")
	query = FilterAttendance(query, attendance)
	query = query.Find(&attendances)
	if err := query.Error; err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for i, data := range attendances {
		wg.Add(1)
		go func(i int, data model.Attendance) {
			attendances[i].User.Role = data.User.GetRole()
			attendances[i].User.Avatar = data.User.GetAvatar()
			attendances[i].User.UserAbilities = data.User.GetAbility(r.db)
			wg.Done()
		}(i, data)
	}
	wg.Wait()
	return attendances, nil
}

func (r attendanceRepo) CheckIsExist(id int) (isExist bool, err error) {
	if err := r.db.Table("attendances").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false, err
	}
	return
}

func (r attendanceRepo) CheckIsExistByDate(userID int, scheduleID int, date string) (isExist bool) {
	if err := r.db.Table("attendances").Select("count(*) > 0").Where("user_id = ? AND schedule_id = ? AND DATE(date) = ?", userID, scheduleID, date).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r attendanceRepo) CountAttendanceByStatus(userID int, statusAttendance string, startDate string, endDate string) (result int) {
	if err := r.db.Table("attendances").Select("count(*)").Where("user_id = ? AND status_presence = ? AND DATE(date) BETWEEN ? AND ?", userID, statusAttendance, startDate, endDate).Find(&result).Error; err != nil {
		return 0
	}
	return
}

func FilterAttendance(query *gorm.DB, attendance model.Attendance) *gorm.DB {
	if attendance.UserID > 0 {
		query = query.Where("user_id = ?", attendance.UserID)
	}
	if attendance.ScheduleID > 0 {
		query = query.Where("schedule_id = ?", attendance.ScheduleID)
	}
	if attendance.Date != "" {
		query = query.Where("DATE(date) = ?", attendance.Date)
	}
	if attendance.Status != "" {
		query = query.Where("status = ?", attendance.Status)
	}
	if attendance.Schedule.OwnerID > 0 {
		query = query.Joins("JOIN schedules ON attendances.schedule_id = schedules.id").Where("schedules.owner_id = ?", attendance.Schedule.OwnerID)
	}
	return query
}

func SearchAttendance(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		searchQuery := "%" + search + "%"
		query = query.Joins("JOIN users ON attendances.user_id = users.id").Where("users.first_name LIKE ? OR users.last_name LIKE ? OR users.email LIKE ? OR attendances.date LIKE ? OR attendances.status_presence LIKE ?", searchQuery, searchQuery, searchQuery, searchQuery, searchQuery)
	}
	return query
}

func PreloadAttendance(query *gorm.DB) *gorm.DB {
	query = query.Preload("User")
	query = query.Preload("Schedule")
	query = query.Preload("Schedule.Subject")
	query = query.Preload("Schedule.DailySchedule")
	query = query.Preload("AttendanceLog")
	return query
}
