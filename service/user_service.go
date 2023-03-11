package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type UserService interface {
	CheckID(id int) bool
	CheckUsername(username string) bool
	CheckEmail(email string) bool
	CheckHandphone(handphone string) bool
	CheckIsActive(username string) bool
	CheckUpdateUsername(id int, username string) bool
	CheckUpdateEmail(id int, email string) bool
	CheckUpdateHandphone(id int, handphone string) bool
	ListUser(user *model.User, pagination *model.Pagination) (*[]model.User, error)
	ListUserMeta(user *model.User, pagination *model.Pagination) (*model.Meta, error)
	CreateUser(user *model.User) (*model.User, error)
	RetrieveUser(id int) (*model.User, error)
	RetrieveUserByUsername(username string) (user *model.User, err error)
	RetrieveUserByEmail(email string) (user *model.User, err error)
	UpdateUser(id int, user *model.User) (*model.User, error)
	DeleteUser(id int) error
	SetActiveUser(id int) (*model.User, error)
	SetDeactiveUser(id int) (*model.User, error)
	DropDownUser(user *model.User) (*[]model.UserDropDown, error)
	GetPassword(id int) (hashPassword string, err error)
	UpdatePassword(userPasswordData *model.UserUpdatePasswordForm) error
}

type userService struct {
	userRepo repo.UserRepo
}

func NewUserService(userRepo repo.UserRepo) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CheckID(id int) bool {
	return s.userRepo.CheckID(id)
}

func (s *userService) CheckUsername(username string) bool {
	return s.userRepo.CheckUsername(username)
}

func (s *userService) CheckEmail(email string) bool {
	return s.userRepo.CheckEmail(email)
}

func (s *userService) CheckHandphone(handphone string) bool {
	return s.userRepo.CheckHandphone(handphone)
}

func (s *userService) CheckIsActive(username string) bool {
	return s.userRepo.CheckIsActive(username)
}

func (s *userService) CheckUpdateUsername(id int, username string) bool {
	return s.userRepo.CheckUpdateUsername(id, username)
}

func (s *userService) CheckUpdateEmail(id int, email string) bool {
	return s.userRepo.CheckUpdateEmail(id, email)
}

func (s *userService) CheckUpdateHandphone(id int, handphone string) bool {
	return s.userRepo.CheckUpdateHandphone(id, handphone)
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

func (s *userService) RetrieveUserByUsername(username string) (user *model.User, err error) {
	userData, err := s.userRepo.RetrieveUserByUsername(username)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (s *userService) RetrieveUserByEmail(email string) (user *model.User, err error) {
	userData, err := s.userRepo.RetrieveUserByUsername(email)
	if err != nil {
		return nil, err
	}
	return userData, nil
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

func (s *userService) DropDownUser(user *model.User) (*[]model.UserDropDown, error) {
	datas, err := s.userRepo.DropDownUser(user)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *userService) GetPassword(id int) (hashPassword string, err error) {
	hashPassword, err = s.userRepo.GetPassword(id)
	if err != nil {
		return "", err
	}
	return
}

func (s *userService) UpdatePassword(userPasswordData *model.UserUpdatePasswordForm) error {
	err := s.userRepo.UpdatePassword(userPasswordData)
	if err != nil {
		return err
	}
	return nil
}
