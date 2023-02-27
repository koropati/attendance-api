package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type ScheduleService interface {
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

type scheduleService struct {
	scheduleRepo repo.ScheduleRepo
}

func NewScheduleService(scheduleRepo repo.ScheduleRepo) ScheduleService {
	return &scheduleService{scheduleRepo: scheduleRepo}
}

func (s *scheduleService) CreateSchedule(schedule *model.Schedule) (*model.Schedule, error) {
	data, err := s.scheduleRepo.CreateSchedule(schedule)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *scheduleService) RetrieveSchedule(id int) (*model.Schedule, error) {
	data, err := s.scheduleRepo.RetrieveSchedule(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *scheduleService) RetrieveScheduleByOwner(id int, ownerID int) (*model.Schedule, error) {
	data, err := s.scheduleRepo.RetrieveScheduleByOwner(id, ownerID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *scheduleService) UpdateSchedule(id int, schedule *model.Schedule) (*model.Schedule, error) {
	data, err := s.scheduleRepo.UpdateSchedule(id, schedule)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *scheduleService) UpdateScheduleByOwner(id int, ownerID int, schedule *model.Schedule) (*model.Schedule, error) {
	data, err := s.scheduleRepo.UpdateScheduleByOwner(id, ownerID, schedule)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *scheduleService) DeleteSchedule(id int) error {
	if err := s.scheduleRepo.DeleteSchedule(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *scheduleService) DeleteScheduleByOwner(id int, ownerID int) error {
	if err := s.scheduleRepo.DeleteScheduleByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *scheduleService) ListSchedule(schedule *model.Schedule, pagination *model.Pagination) (*[]model.Schedule, error) {
	datas, err := s.scheduleRepo.ListSchedule(schedule, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *scheduleService) ListScheduleMeta(schedule *model.Schedule, pagination *model.Pagination) (*model.Meta, error) {
	data, err := s.scheduleRepo.ListScheduleMeta(schedule, pagination)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *scheduleService) DropDownSchedule(schedule *model.Schedule) (*[]model.Schedule, error) {
	datas, err := s.scheduleRepo.DropDownSchedule(schedule)
	if err != nil {
		return nil, err
	}
	return datas, nil
}
