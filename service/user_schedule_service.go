package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type UserScheduleService interface {
	CreateUserSchedule(userschedule *model.UserSchedule) (*model.UserSchedule, error)
	RetrieveUserSchedule(id int) (*model.UserSchedule, error)
	UpdateUserSchedule(id int, userschedule *model.UserSchedule) (*model.UserSchedule, error)
	DeleteUserSchedule(id int) error
	ListUserSchedule(userschedule *model.UserSchedule, pagination *model.Pagination) (*[]model.UserSchedule, error)
	ListUserScheduleMeta(userschedule *model.UserSchedule, pagination *model.Pagination) (*model.Meta, error)
	DropDownUserSchedule(userschedule *model.UserSchedule) (*[]model.UserSchedule, error)
}

type userScheduleService struct {
	userScheduleRepo repo.UserScheduleRepo
}

func NewUserScheduleService(userScheduleRepo repo.UserScheduleRepo) UserScheduleService {
	return &userScheduleService{userScheduleRepo: userScheduleRepo}
}

func (s *userScheduleService) CreateUserSchedule(userschedule *model.UserSchedule) (*model.UserSchedule, error) {
	data, err := s.userScheduleRepo.CreateUserSchedule(userschedule)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *userScheduleService) RetrieveUserSchedule(id int) (*model.UserSchedule, error) {
	data, err := s.userScheduleRepo.RetrieveUserSchedule(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *userScheduleService) UpdateUserSchedule(id int, userschedule *model.UserSchedule) (*model.UserSchedule, error) {
	data, err := s.userScheduleRepo.UpdateUserSchedule(id, userschedule)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *userScheduleService) DeleteUserSchedule(id int) error {
	if err := s.userScheduleRepo.DeleteUserSchedule(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *userScheduleService) ListUserSchedule(userschedule *model.UserSchedule, pagination *model.Pagination) (*[]model.UserSchedule, error) {
	datas, err := s.userScheduleRepo.ListUserSchedule(userschedule, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *userScheduleService) ListUserScheduleMeta(userschedule *model.UserSchedule, pagination *model.Pagination) (*model.Meta, error) {
	data, err := s.userScheduleRepo.ListUserScheduleMeta(userschedule, pagination)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *userScheduleService) DropDownUserSchedule(userschedule *model.UserSchedule) (*[]model.UserSchedule, error) {
	datas, err := s.userScheduleRepo.DropDownUserSchedule(userschedule)
	if err != nil {
		return nil, err
	}
	return datas, nil
}
