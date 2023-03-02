package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type UserService interface {
	ListUser(user *model.User, pagination *model.Pagination) (*[]model.User, error)
	ListUserMeta(user *model.User, pagination *model.Pagination) (*model.Meta, error)
	CreateUser(user *model.User) (*model.User, error)
	RetrieveUser(id int) (*model.User, error)
	UpdateUser(id int, user *model.User) (*model.User, error)
	DeleteUser(id int) error
	SetActiveUser(id int) (*model.User, error)
	SetDeactiveUser(id int) (*model.User, error)
}

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) ListUser(user *model.User, pagination *model.Pagination) (*[]model.User, error) {
	users, err := s.userRepo.ListUser(user, pagination)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (s *userService) ListUserMeta(user *model.User, pagination *model.Pagination) (*model.Meta, error) {
	meta, err := s.userRepo.ListUserMeta(user, pagination)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (s *userService) CreateUser(user *model.User) (*model.User, error) {
	newUser, err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}

func (s *userService) RetrieveUser(id int) (*model.User, error) {
	data, err := s.userRepo.RetrieveUser(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *userService) UpdateUser(id int, user *model.User) (*model.User, error) {
	userUpdate, err := s.userRepo.UpdateUser(id, user)
	if err != nil {
		return nil, err
	}
	return userUpdate, nil
}

func (s *userService) DeleteUser(id int) error {
	err := s.userRepo.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) SetActiveUser(id int) (*model.User, error) {
	userUpdate, err := s.userRepo.SetActiveUser(id)
	if err != nil {
		return nil, err
	}
	return userUpdate, nil
}

func (s *userService) SetDeactiveUser(id int) (*model.User, error) {
	userUpdate, err := s.userRepo.SetDeactiveUser(id)
	if err != nil {
		return nil, err
	}
	return userUpdate, nil
}
