package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type AuthService interface {
	CheckUsername(username string) bool
	CheckEmail(email string) bool
	CheckHandphone(handphone string) bool
	CheckIsActive(username string) bool
	IsSuperAdmin(username string) (bool, error)
	IsAdmin(username string) (bool, error)
	IsUser(username string) (bool, error)
	GetRole(username string) (bool, bool, bool, error)
	GetEmail(username string) (string, error)
	Register(user *model.User) error
	Login(username string) (string, error)
	CheckID(id int) bool
	GetByUsername(username string) (user *model.User, err error)
	GetByEmail(email string) (user *model.User, err error)
	Create(user *model.User) error
	Delete(id int) error
}

type authService struct {
	authRepo repo.AuthRepo
}

func NewAuthService(authRepo repo.AuthRepo) AuthService {
	return &authService{authRepo: authRepo}
}

func (s *authService) CheckUsername(username string) bool {
	return s.authRepo.CheckUsername(username)
}

func (s *authService) CheckEmail(email string) bool {
	return s.authRepo.CheckEmail(email)
}

func (s *authService) CheckHandphone(handphone string) bool {
	return s.authRepo.CheckHandphone(handphone)
}

func (s *authService) CheckIsActive(username string) bool {
	return s.authRepo.CheckIsActive(username)
}

func (s *authService) IsSuperAdmin(username string) (bool, error) {
	isSuperAdmin, err := s.authRepo.IsSuperAdmin(username)
	if err != nil {
		return false, err
	}
	return isSuperAdmin, nil
}

func (s *authService) IsAdmin(username string) (bool, error) {
	isAdmin, err := s.authRepo.IsAdmin(username)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func (s *authService) IsUser(username string) (bool, error) {
	isUser, err := s.authRepo.IsUser(username)
	if err != nil {
		return false, err
	}
	return isUser, nil
}

func (s *authService) GetRole(username string) (bool, bool, bool, error) {
	isSuper, isAdmin, isUser, err := s.authRepo.GetRole(username)
	if err != nil {
		return false, false, false, err
	}
	return isSuper, isAdmin, isUser, nil
}

func (s *authService) GetEmail(username string) (string, error) {
	role, err := s.authRepo.GetEmail(username)
	if err != nil {
		return "", err
	}
	return role, nil
}

func (s *authService) Register(user *model.User) error {
	if err := s.authRepo.Register(user); err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(username string) (string, error) {
	password, err := s.authRepo.Login(username)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (s *authService) CheckID(id int) bool {
	return s.authRepo.CheckID(id)
}

func (s *authService) GetByUsername(username string) (user *model.User, err error) {
	userData, err := s.authRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (s *authService) GetByEmail(email string) (user *model.User, err error) {
	userData, err := s.authRepo.GetByUsername(email)
	if err != nil {
		return nil, err
	}
	return userData, nil
}

func (s *authService) Create(user *model.User) error {
	if err := s.authRepo.Create(user); err != nil {
		return err
	}

	return nil
}

func (s *authService) Delete(id int) error {
	if err := s.authRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
