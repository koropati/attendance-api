package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type AttendanceService interface {
	CreateAttendance(attendance *model.Attendance) (*model.Attendance, error)
	RetrieveAttendance(id int) (*model.Attendance, error)
	UpdateAttendance(id int, attendance *model.Attendance) (*model.Attendance, error)
	DeleteAttendance(id int) error
	ListAttendance(attendance *model.Attendance, pagination *model.Pagination) (*[]model.Attendance, error)
	ListAttendanceMeta(attendance *model.Attendance, pagination *model.Pagination) (*model.Meta, error)
	DropDownAttendance(attendance *model.Attendance) (*[]model.Attendance, error)
}

type attendanceService struct {
	attendanceRepo repo.AttendanceRepo
}

func NewAttendanceService(attendanceRepo repo.AttendanceRepo) AttendanceService {
	return &attendanceService{attendanceRepo: attendanceRepo}
}

func (s *attendanceService) CreateAttendance(attendance *model.Attendance) (*model.Attendance, error) {
	data, err := s.attendanceRepo.CreateAttendance(attendance)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *attendanceService) RetrieveAttendance(id int) (*model.Attendance, error) {
	data, err := s.attendanceRepo.RetrieveAttendance(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *attendanceService) UpdateAttendance(id int, attendance *model.Attendance) (*model.Attendance, error) {
	data, err := s.attendanceRepo.UpdateAttendance(id, attendance)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *attendanceService) DeleteAttendance(id int) error {
	if err := s.attendanceRepo.DeleteAttendance(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *attendanceService) ListAttendance(attendance *model.Attendance, pagination *model.Pagination) (*[]model.Attendance, error) {
	datas, err := s.attendanceRepo.ListAttendance(attendance, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *attendanceService) ListAttendanceMeta(attendance *model.Attendance, pagination *model.Pagination) (*model.Meta, error) {
	data, err := s.attendanceRepo.ListAttendanceMeta(attendance, pagination)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *attendanceService) DropDownAttendance(attendance *model.Attendance) (*[]model.Attendance, error) {
	datas, err := s.attendanceRepo.DropDownAttendance(attendance)
	if err != nil {
		return nil, err
	}
	return datas, nil
}
