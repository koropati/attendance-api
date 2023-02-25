package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type AttendanceLogService interface {
	CreateAttendanceLog(attendancelog *model.AttendanceLog) (*model.AttendanceLog, error)
	RetrieveAttendanceLog(id int) (*model.AttendanceLog, error)
	UpdateAttendanceLog(id int, attendancelog *model.AttendanceLog) (*model.AttendanceLog, error)
	DeleteAttendanceLog(id int) error
	ListAttendanceLog(attendancelog *model.AttendanceLog, pagination *model.Pagination) (*[]model.AttendanceLog, error)
	ListAttendanceLogMeta(attendancelog *model.AttendanceLog, pagination *model.Pagination) (*model.Meta, error)
	DropDownAttendanceLog(attendancelog *model.AttendanceLog) (*[]model.AttendanceLog, error)
}

type attendanceLogService struct {
	attendanceLogRepo repo.AttendanceLogRepo
}

func NewAttendanceLogService(attendanceLogRepo repo.AttendanceLogRepo) AttendanceLogService {
	return &attendanceLogService{attendanceLogRepo: attendanceLogRepo}
}

func (s *attendanceLogService) CreateAttendanceLog(attendancelog *model.AttendanceLog) (*model.AttendanceLog, error) {
	data, err := s.attendanceLogRepo.CreateAttendanceLog(attendancelog)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *attendanceLogService) RetrieveAttendanceLog(id int) (*model.AttendanceLog, error) {
	data, err := s.attendanceLogRepo.RetrieveAttendanceLog(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *attendanceLogService) UpdateAttendanceLog(id int, attendancelog *model.AttendanceLog) (*model.AttendanceLog, error) {
	data, err := s.attendanceLogRepo.UpdateAttendanceLog(id, attendancelog)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *attendanceLogService) DeleteAttendanceLog(id int) error {
	if err := s.attendanceLogRepo.DeleteAttendanceLog(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *attendanceLogService) ListAttendanceLog(attendancelog *model.AttendanceLog, pagination *model.Pagination) (*[]model.AttendanceLog, error) {
	datas, err := s.attendanceLogRepo.ListAttendanceLog(attendancelog, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *attendanceLogService) ListAttendanceLogMeta(attendancelog *model.AttendanceLog, pagination *model.Pagination) (*model.Meta, error) {
	data, err := s.attendanceLogRepo.ListAttendanceLogMeta(attendancelog, pagination)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *attendanceLogService) DropDownAttendanceLog(attendancelog *model.AttendanceLog) (*[]model.AttendanceLog, error) {
	datas, err := s.attendanceLogRepo.DropDownAttendanceLog(attendancelog)
	if err != nil {
		return nil, err
	}
	return datas, nil
}
