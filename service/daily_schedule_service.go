package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type DailyScheduleService interface {
	CreateDailySchedule(dailyschedule *model.DailySchedule) (*model.DailySchedule, error)
	RetrieveDailySchedule(id int) (*model.DailySchedule, error)
	RetrieveDailyScheduleByOwner(id int, ownerID int) (*model.DailySchedule, error)
	UpdateDailySchedule(id int, dailyschedule *model.DailySchedule) (*model.DailySchedule, error)
	UpdateDailyScheduleByOwner(id int, ownerID int, dailyschedule *model.DailySchedule) (*model.DailySchedule, error)
	DeleteDailySchedule(id int) error
	DeleteDailyScheduleByOwner(id int, ownerID int) error
	ListDailySchedule(dailyschedule *model.DailySchedule, pagination *model.Pagination) (*[]model.DailySchedule, error)
	ListDailyScheduleMeta(dailyschedule *model.DailySchedule, pagination *model.Pagination) (*model.Meta, error)
	DropDownDailySchedule(dailyschedule *model.DailySchedule) (*[]model.DailySchedule, error)
	CheckHaveDailySchedule(scheduleID int, day string) (isHaveDailySchedule bool, dailyScheduleID int, err error)
}

type dailyScheduleService struct {
	dailyScheduleRepo repo.DailyScheduleRepo
}

func NewDailyScheduleService(dailyScheduleRepo repo.DailyScheduleRepo) DailyScheduleService {
	return &dailyScheduleService{dailyScheduleRepo: dailyScheduleRepo}
}

func (s *dailyScheduleService) CreateDailySchedule(dailyschedule *model.DailySchedule) (*model.DailySchedule, error) {
	data, err := s.dailyScheduleRepo.CreateDailySchedule(dailyschedule)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *dailyScheduleService) RetrieveDailySchedule(id int) (*model.DailySchedule, error) {
	data, err := s.dailyScheduleRepo.RetrieveDailySchedule(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *dailyScheduleService) RetrieveDailyScheduleByOwner(id int, ownerID int) (*model.DailySchedule, error) {
	data, err := s.dailyScheduleRepo.RetrieveDailyScheduleByOwner(id, ownerID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *dailyScheduleService) UpdateDailySchedule(id int, dailyschedule *model.DailySchedule) (*model.DailySchedule, error) {
	data, err := s.dailyScheduleRepo.UpdateDailySchedule(id, dailyschedule)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *dailyScheduleService) UpdateDailyScheduleByOwner(id int, ownerID int, dailyschedule *model.DailySchedule) (*model.DailySchedule, error) {
	data, err := s.dailyScheduleRepo.UpdateDailyScheduleByOwner(id, ownerID, dailyschedule)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *dailyScheduleService) DeleteDailySchedule(id int) error {
	if err := s.dailyScheduleRepo.DeleteDailySchedule(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *dailyScheduleService) DeleteDailyScheduleByOwner(id int, ownerID int) error {
	if err := s.dailyScheduleRepo.DeleteDailyScheduleByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *dailyScheduleService) ListDailySchedule(dailyschedule *model.DailySchedule, pagination *model.Pagination) (*[]model.DailySchedule, error) {
	datas, err := s.dailyScheduleRepo.ListDailySchedule(dailyschedule, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *dailyScheduleService) ListDailyScheduleMeta(dailyschedule *model.DailySchedule, pagination *model.Pagination) (*model.Meta, error) {
	data, err := s.dailyScheduleRepo.ListDailyScheduleMeta(dailyschedule, pagination)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *dailyScheduleService) DropDownDailySchedule(dailyschedule *model.DailySchedule) (*[]model.DailySchedule, error) {
	datas, err := s.dailyScheduleRepo.DropDownDailySchedule(dailyschedule)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *dailyScheduleService) CheckHaveDailySchedule(scheduleID int, day string) (isHaveDailySchedule bool, dailyScheduleID int, err error) {
	return s.dailyScheduleRepo.CheckHaveDailySchedule(scheduleID, day)
}
