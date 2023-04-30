package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type AttendanceService interface {
	CreateAttendance(attendance model.Attendance) (model.Attendance, error)
	RetrieveAttendance(id int) (model.Attendance, error)
	RetrieveAttendanceByUserID(id int, userID int) (model.Attendance, error)
	RetrieveAttendanceByDate(userID int, scheduleID int, date string) (model.Attendance, error)
	UpdateAttendance(id int, attendance model.Attendance) (model.Attendance, error)
	UpdateAttendanceByUserID(id int, userID int, attendance model.Attendance) (model.Attendance, error)
	DeleteAttendance(id int) error
	DeleteAttendanceByUserID(id int, userID int) error
	ListAttendance(attendance model.Attendance, pagination model.Pagination) ([]model.Attendance, error)
	ListAttendanceMeta(attendance model.Attendance, pagination model.Pagination) (model.Meta, error)
	DropDownAttendance(attendance model.Attendance) ([]model.Attendance, error)
	CheckIsExist(id int) (isExist bool, err error)
	CheckIsExistByDate(userID int, scheduleID int, date string) bool
	CountAttendanceByStatus(userID int, statusAttendance string, startDate string, endDate string) (result int)
}

type attendanceService struct {
	attendanceRepo repo.AttendanceRepo
}

func NewAttendanceService(attendanceRepo repo.AttendanceRepo) AttendanceService {
	return &attendanceService{attendanceRepo: attendanceRepo}
}

func (s attendanceService) CreateAttendance(attendance model.Attendance) (model.Attendance, error) {
	data, err := s.attendanceRepo.CreateAttendance(attendance)
	if err != nil {
		return model.Attendance{}, err
	}
	return data, nil
}

func (s attendanceService) RetrieveAttendance(id int) (model.Attendance, error) {
	data, err := s.attendanceRepo.RetrieveAttendance(id)
	if err != nil {
		return model.Attendance{}, err
	}
	return data, nil
}

func (s attendanceService) RetrieveAttendanceByUserID(id int, userID int) (model.Attendance, error) {
	data, err := s.attendanceRepo.RetrieveAttendanceByUserID(id, userID)
	if err != nil {
		return model.Attendance{}, err
	}
	return data, nil
}

func (s attendanceService) RetrieveAttendanceByDate(userID int, scheduleID int, date string) (model.Attendance, error) {
	data, err := s.attendanceRepo.RetrieveAttendanceByDate(userID, scheduleID, date)
	if err != nil {
		return model.Attendance{}, err
	}
	return data, nil
}

func (s attendanceService) UpdateAttendance(id int, attendance model.Attendance) (model.Attendance, error) {
	data, err := s.attendanceRepo.UpdateAttendance(id, attendance)
	if err != nil {
		return model.Attendance{}, err
	}
	return data, nil
}

func (s attendanceService) UpdateAttendanceByUserID(id int, userID int, attendance model.Attendance) (model.Attendance, error) {
	data, err := s.attendanceRepo.UpdateAttendanceByUserID(id, userID, attendance)
	if err != nil {
		return model.Attendance{}, err
	}
	return data, nil
}

func (s attendanceService) DeleteAttendance(id int) error {
	if err := s.attendanceRepo.DeleteAttendance(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s attendanceService) DeleteAttendanceByUserID(id int, userID int) error {
	if err := s.attendanceRepo.DeleteAttendanceByUserID(id, userID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s attendanceService) ListAttendance(attendance model.Attendance, pagination model.Pagination) ([]model.Attendance, error) {
	datas, err := s.attendanceRepo.ListAttendance(attendance, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s attendanceService) ListAttendanceMeta(attendance model.Attendance, pagination model.Pagination) (model.Meta, error) {
	data, err := s.attendanceRepo.ListAttendanceMeta(attendance, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s attendanceService) DropDownAttendance(attendance model.Attendance) ([]model.Attendance, error) {
	datas, err := s.attendanceRepo.DropDownAttendance(attendance)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s attendanceService) CheckIsExist(id int) (isExist bool, err error) {
	return s.attendanceRepo.CheckIsExist(id)
}

func (s attendanceService) CheckIsExistByDate(userID int, scheduleID int, date string) bool {
	return s.attendanceRepo.CheckIsExistByDate(userID, scheduleID, date)
}

func (s attendanceService) CountAttendanceByStatus(userID int, statusAttendance string, startDate string, endDate string) (result int) {
	return s.attendanceRepo.CountAttendanceByStatus(userID, statusAttendance, startDate, endDate)
}
