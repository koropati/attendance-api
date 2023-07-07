package service

import (
	"attendance-api/model"
	"attendance-api/repo"
	"time"
)

type UserScheduleService interface {
	CreateUserSchedule(userschedule model.UserSchedule) (model.UserSchedule, error)
	RetrieveUserSchedule(id int) (model.UserSchedule, error)
	RetrieveUserScheduleByOwner(id int, ownerID int) (model.UserSchedule, error)
	UpdateUserSchedule(id int, userschedule model.UserSchedule) (model.UserSchedule, error)
	UpdateUserScheduleByOwner(id int, ownerID int, userschedule model.UserSchedule) (model.UserSchedule, error)
	DeleteUserSchedule(id int) error
	DeleteUserScheduleByOwner(id int, ownerID int) error
	RemoveUserFromSchedule(scheduleID int, userID int) error
	RemoveUserFromScheduleByOwner(scheduleID int, userID int, ownerID int) error
	ListMySchedule(userID int) ([]model.MySchedule, error)
	ListTodaySchedule(userID int, dayName string) ([]model.TodaySchedule, error)
	ListUserSchedule(userschedule model.UserSchedule, pagination model.Pagination) ([]model.UserSchedule, error)
	ListUserScheduleMeta(userschedule model.UserSchedule, pagination model.Pagination) (model.Meta, error)
	ListUserInRule(scheduleID int, student model.Student, pagination model.Pagination) ([]model.Student, error)
	ListUserInRuleMeta(scheduleID int, student model.Student, pagination model.Pagination) (model.Meta, error)
	ListUserNotInRule(scheduleID int, student model.Student, pagination model.Pagination) ([]model.Student, error)
	ListUserNotInRuleMeta(scheduleID int, student model.Student, pagination model.Pagination) (model.Meta, error)
	DropDownUserSchedule(userschedule model.UserSchedule) ([]model.UserSchedule, error)
	CheckHaveSchedule(userID int, date time.Time) (isHaveSchedule bool, scheduleID int, err error)
	CheckUserInSchedule(scheduleID int, userID int) bool
	CountByScheduleID(scheduleID int) (total int)
	GetAll() (results []model.UserSchedule, err error)
	GetAllByTodayRange() (results []model.UserSchedule, err error)
}

type userScheduleService struct {
	userScheduleRepo repo.UserScheduleRepo
}

func NewUserScheduleService(userScheduleRepo repo.UserScheduleRepo) UserScheduleService {
	return userScheduleService{userScheduleRepo: userScheduleRepo}
}

func (s userScheduleService) CreateUserSchedule(userschedule model.UserSchedule) (model.UserSchedule, error) {
	data, err := s.userScheduleRepo.CreateUserSchedule(userschedule)
	if err != nil {
		return model.UserSchedule{}, err
	}
	return data, nil
}

func (s userScheduleService) RetrieveUserSchedule(id int) (model.UserSchedule, error) {
	data, err := s.userScheduleRepo.RetrieveUserSchedule(id)
	if err != nil {
		return model.UserSchedule{}, err
	}
	return data, nil
}

func (s userScheduleService) RetrieveUserScheduleByOwner(id int, ownerID int) (model.UserSchedule, error) {
	data, err := s.userScheduleRepo.RetrieveUserScheduleByOwner(id, ownerID)
	if err != nil {
		return model.UserSchedule{}, err
	}
	return data, nil
}

func (s userScheduleService) UpdateUserSchedule(id int, userschedule model.UserSchedule) (model.UserSchedule, error) {
	data, err := s.userScheduleRepo.UpdateUserSchedule(id, userschedule)
	if err != nil {
		return model.UserSchedule{}, err
	}
	return data, nil
}

func (s userScheduleService) UpdateUserScheduleByOwner(id int, ownerID int, userschedule model.UserSchedule) (model.UserSchedule, error) {
	data, err := s.userScheduleRepo.UpdateUserScheduleByOwner(id, ownerID, userschedule)
	if err != nil {
		return model.UserSchedule{}, err
	}
	return data, nil
}

func (s userScheduleService) DeleteUserSchedule(id int) error {
	if err := s.userScheduleRepo.DeleteUserSchedule(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s userScheduleService) RemoveUserFromSchedule(scheduleID int, userID int) error {
	if err := s.userScheduleRepo.RemoveUserFromSchedule(scheduleID, userID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s userScheduleService) DeleteUserScheduleByOwner(id int, ownerID int) error {
	if err := s.userScheduleRepo.DeleteUserScheduleByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s userScheduleService) RemoveUserFromScheduleByOwner(scheduleID int, userID int, ownerID int) error {
	if err := s.userScheduleRepo.RemoveUserFromScheduleByOwner(scheduleID, userID, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s userScheduleService) ListMySchedule(userID int) (results []model.MySchedule, err error) {
	return s.userScheduleRepo.ListMySchedule(userID)
}

func (s userScheduleService) ListTodaySchedule(userID int, dayName string) (results []model.TodaySchedule, err error) {
	return s.userScheduleRepo.ListTodaySchedule(userID, dayName)
}

func (s userScheduleService) ListUserSchedule(userschedule model.UserSchedule, pagination model.Pagination) ([]model.UserSchedule, error) {
	datas, err := s.userScheduleRepo.ListUserSchedule(userschedule, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s userScheduleService) ListUserScheduleMeta(userschedule model.UserSchedule, pagination model.Pagination) (model.Meta, error) {
	data, err := s.userScheduleRepo.ListUserScheduleMeta(userschedule, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s userScheduleService) ListUserInRule(scheduleID int, student model.Student, pagination model.Pagination) ([]model.Student, error) {
	datas, err := s.userScheduleRepo.ListUserInRule(scheduleID, student, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s userScheduleService) ListUserInRuleMeta(scheduleID int, student model.Student, pagination model.Pagination) (model.Meta, error) {
	data, err := s.userScheduleRepo.ListUserInRuleMeta(scheduleID, student, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s userScheduleService) ListUserNotInRule(scheduleID int, student model.Student, pagination model.Pagination) ([]model.Student, error) {
	datas, err := s.userScheduleRepo.ListUserNotInRule(scheduleID, student, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s userScheduleService) ListUserNotInRuleMeta(scheduleID int, student model.Student, pagination model.Pagination) (model.Meta, error) {
	data, err := s.userScheduleRepo.ListUserNotInRuleMeta(scheduleID, student, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s userScheduleService) DropDownUserSchedule(userschedule model.UserSchedule) ([]model.UserSchedule, error) {
	datas, err := s.userScheduleRepo.DropDownUserSchedule(userschedule)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s userScheduleService) CheckHaveSchedule(userID int, date time.Time) (isHaveSchedule bool, scheduleID int, err error) {
	return s.userScheduleRepo.CheckHaveSchedule(userID, date)
}

func (s userScheduleService) CheckUserInSchedule(scheduleID int, userID int) bool {
	return s.userScheduleRepo.CheckUserInSchedule(scheduleID, userID)
}

func (s userScheduleService) CountByScheduleID(scheduleID int) (total int) {
	return s.userScheduleRepo.CountByScheduleID(scheduleID)
}

func (s userScheduleService) GetAll() (resutls []model.UserSchedule, err error) {
	return s.userScheduleRepo.GetAll()
}

func (s userScheduleService) GetAllByTodayRange() (resutls []model.UserSchedule, err error) {
	return s.userScheduleRepo.GetAllByTodayRange()
}
