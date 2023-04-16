package repo

import (
	"attendance-api/model"
	"strconv"

	"gorm.io/gorm"
)

type AttendanceLogRepo interface {
	CreateAttendanceLog(attendancelog *model.AttendanceLog) (*model.AttendanceLog, error)
	RetrieveAttendanceLog(id int) (*model.AttendanceLog, error)
	UpdateAttendanceLog(id int, attendancelog *model.AttendanceLog) (*model.AttendanceLog, error)
	DeleteAttendanceLog(id int) error
	ListAttendanceLog(attendancelog *model.AttendanceLog, pagination *model.Pagination) (*[]model.AttendanceLog, error)
	ListAttendanceLogMeta(attendancelog *model.AttendanceLog, pagination *model.Pagination) (*model.Meta, error)
	DropDownAttendanceLog(attendancelog *model.AttendanceLog) (*[]model.AttendanceLog, error)
}

type attendanceLogRepo struct {
	db *gorm.DB
}

func NewAttendanceLogRepo(db *gorm.DB) AttendanceLogRepo {
	return &attendanceLogRepo{db: db}
}

func (r *attendanceLogRepo) CreateAttendanceLog(attendancelog *model.AttendanceLog) (*model.AttendanceLog, error) {
	if err := r.db.Table("attendance_logs").Create(&attendancelog).Error; err != nil {
		return nil, err
	}
	return attendancelog, nil
}

func (r *attendanceLogRepo) RetrieveAttendanceLog(id int) (*model.AttendanceLog, error) {
	var attendancelog model.AttendanceLog
	if err := r.db.First(&attendancelog, id).Error; err != nil {
		return nil, err
	}
	return &attendancelog, nil
}

func (r *attendanceLogRepo) UpdateAttendanceLog(id int, attendancelog *model.AttendanceLog) (*model.AttendanceLog, error) {
	if err := r.db.Model(&model.AttendanceLog{}).Where("id = ?", id).Updates(&attendancelog).Error; err != nil {
		return nil, err
	}
	return attendancelog, nil
}

func (r *attendanceLogRepo) DeleteAttendanceLog(id int) error {
	if err := r.db.Delete(&model.AttendanceLog{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *attendanceLogRepo) ListAttendanceLog(attendancelog *model.AttendanceLog, pagination *model.Pagination) (*[]model.AttendanceLog, error) {
	var attendance_logs []model.AttendanceLog
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("attendance_logs").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterAttendanceLog(query, attendancelog)
	query = SearchAttendanceLog(query, pagination.Search)
	query = query.Find(&attendance_logs)
	if err := query.Error; err != nil {
		return nil, err
	}

	return &attendance_logs, nil
}

func (r *attendanceLogRepo) ListAttendanceLogMeta(attendancelog *model.AttendanceLog, pagination *model.Pagination) (*model.Meta, error) {
	var attendance_logs []model.AttendanceLog
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.AttendanceLog{}).Select("count(*)")
	queryTotal = FilterAttendanceLog(queryTotal, attendancelog)
	queryTotal = SearchAttendanceLog(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return nil, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("attendance_logs").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterAttendanceLog(query, attendancelog)
	query = SearchAttendanceLog(query, pagination.Search)
	query = query.Find(&attendance_logs)
	if err := query.Error; err != nil {
		return nil, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(attendance_logs),
	}
	return &meta, nil
}

func (r *attendanceLogRepo) DropDownAttendanceLog(attendancelog *model.AttendanceLog) (*[]model.AttendanceLog, error) {
	var attendance_logs []model.AttendanceLog
	query := r.db.Table("attendance_logs")
	query = FilterAttendanceLog(query, attendancelog)
	query = query.Find(&attendance_logs)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &attendance_logs, nil
}

func FilterAttendanceLog(query *gorm.DB, attendancelog *model.AttendanceLog) *gorm.DB {
	if attendancelog.AttendanceID > 0 {
		query = query.Where("attendance_id = ?", attendancelog.AttendanceID)
	}
	if attendancelog.LogType != "" {
		query = query.Where("log_type = ?", attendancelog.LogType)
	}
	if attendancelog.Latitude > 0 {
		query = query.Where("latitude LIKE ?", "%"+strconv.Itoa(int(attendancelog.Latitude))+"%")
	}
	if attendancelog.Longitude > 0 {
		query = query.Where("longitude LIKE ?", "%"+strconv.Itoa(int(attendancelog.Longitude))+"%")
	}
	if attendancelog.TimeZone > 0 {
		query = query.Where("time_zone = ?", attendancelog.TimeZone)
	}
	if attendancelog.Location != "" {
		query = query.Where("location LIKE ?", "%"+attendancelog.Location+"%")
	}
	if attendancelog.Status != "" {
		query = query.Where("status = ?", attendancelog.Status)
	}
	return query
}

func SearchAttendanceLog(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("attendance_id LIKE ? OR log_type LIKE ? OR latitude LIKE ? OR longitude LIKE ? OR time_zone LIKE ? OR location LIKE ? OR status LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}
