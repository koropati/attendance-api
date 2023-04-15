package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type AttendanceRepo interface {
	CreateAttendance(attendance *model.Attendance) (*model.Attendance, error)
	RetrieveAttendance(id int) (*model.Attendance, error)
	RetrieveAttendanceByUserID(id int, userID int) (*model.Attendance, error)
	RetrieveAttendanceByDate(userID int, scheduleID int, date string) (*model.Attendance, error)
	UpdateAttendance(id int, attendance *model.Attendance) (*model.Attendance, error)
	UpdateAttendanceByUserID(id int, userID int, attendance *model.Attendance) (*model.Attendance, error)
	DeleteAttendance(id int) error
	DeleteAttendanceByUserID(id int, userID int) error
	ListAttendance(attendance *model.Attendance, pagination *model.Pagination) (*[]model.Attendance, error)
	ListAttendanceMeta(attendance *model.Attendance, pagination *model.Pagination) (*model.Meta, error)
	DropDownAttendance(attendance *model.Attendance) (*[]model.Attendance, error)
	CheckIsExist(id int) (isExist bool, err error)
	CheckIsExistByDate(userID int, scheduleID int, date string) bool
}

type attendanceRepo struct {
	db *gorm.DB
}

func NewAttendanceRepo(db *gorm.DB) AttendanceRepo {
	return &attendanceRepo{db: db}
}

func (r *attendanceRepo) CreateAttendance(attendance *model.Attendance) (*model.Attendance, error) {
	if err := r.db.Table("attendances").Create(&attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}

func (r *attendanceRepo) RetrieveAttendance(id int) (*model.Attendance, error) {
	var attendance model.Attendance
	if err := r.db.First(&attendance, id).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepo) RetrieveAttendanceByUserID(id int, userID int) (*model.Attendance, error) {
	var attendance model.Attendance
	if err := r.db.Model(&model.Attendance{}).Where("id = ? AND user_id = ?", id, userID).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepo) RetrieveAttendanceByDate(userID int, scheduleID int, date string) (*model.Attendance, error) {
	var attendance model.Attendance
	if err := r.db.Model(&model.Attendance{}).Where("user_id = ? AND schedule_id = ? AND DATE(date) = ?", userID, scheduleID, date).First(&attendance).Error; err != nil {
		return nil, err
	}
	return &attendance, nil
}

func (r *attendanceRepo) UpdateAttendance(id int, attendance *model.Attendance) (*model.Attendance, error) {
	if err := r.db.Model(&model.Attendance{}).Where("id = ?", id).Updates(&attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}

func (r *attendanceRepo) UpdateAttendanceByUserID(id int, userID int, attendance *model.Attendance) (*model.Attendance, error) {
	if err := r.db.Model(&model.Attendance{}).Where("id = ? AND user_Id = ?", id, userID).Updates(&attendance).Error; err != nil {
		return nil, err
	}
	return attendance, nil
}

func (r *attendanceRepo) DeleteAttendance(id int) error {
	if err := r.db.Delete(&model.Attendance{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *attendanceRepo) DeleteAttendanceByUserID(id int, userID int) error {
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&model.Attendance{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *attendanceRepo) ListAttendance(attendance *model.Attendance, pagination *model.Pagination) (*[]model.Attendance, error) {
	var attendances []model.Attendance
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("attendances").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterAttendance(query, attendance)
	query = query.Find(&attendances)
	if err := query.Error; err != nil {
		return nil, err
	}

	return &attendances, nil
}

func (r *attendanceRepo) ListAttendanceMeta(attendance *model.Attendance, pagination *model.Pagination) (*model.Meta, error) {
	var attendances []model.Attendance
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Attendance{}).Select("count(*)")
	queryTotal = FilterAttendance(queryTotal, attendance)
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
		CurrentRecord: len(attendances),
	}
	return &meta, nil
}

func (r *attendanceRepo) DropDownAttendance(attendance *model.Attendance) (*[]model.Attendance, error) {
	var attendances []model.Attendance
	query := r.db.Table("attendances")
	query = FilterAttendance(query, attendance)
	query = query.Find(&attendances)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &attendances, nil
}

func (r *attendanceRepo) CheckIsExist(id int) (isExist bool, err error) {
	if err := r.db.Table("attendances").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false, err
	}
	return
}

func (r *attendanceRepo) CheckIsExistByDate(userID int, scheduleID int, date string) (isExist bool) {
	if err := r.db.Table("attendances").Select("count(*) > 0").Where("user_id = ? AND schedule_id = ? AND DATE(date) = ?", userID, scheduleID, date).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func FilterAttendance(query *gorm.DB, attendance *model.Attendance) *gorm.DB {
	if attendance.UserID > 0 {
		query = query.Where("user_id = ?", attendance.UserID)
	}
	if attendance.ScheduleID > 0 {
		query = query.Where("schedule_id = ?", attendance.ScheduleID)
	}
	if attendance.Date.Format("2006-01-02") != "" {
		query = query.Where("DATE(date) = ?", attendance.Date.Format("2006-01-02"))
	}
	if attendance.Status != "" {
		query = query.Where("status = ?", attendance.Status)
	}
	return query
}
