package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type AuthService interface {
	CheckID(id int) bool
	CheckUsername(username string) bool
	CheckEmail(email string) bool
	CheckHandphone(handphone string) bool
	CheckIsActive(username string) bool
	IsSuperAdmin(username string) (bool, error)
	IsAdmin(username string) (bool, error)
	IsUser(username string) (bool, error)
	GetRole(username string) (bool, bool, bool, error)
	GetEmail(username string) (string, error)
	Register(user model.User) error
	Login(username string) (string, error)
	GetByUsername(username string) (user model.User, err error)
	GetByEmail(email string) (user model.User, err error)
	GetByID(id uint) (user model.User, err error)
	Create(user model.User) error
	Delete(id int) error
	SetActiveUser(id int) (model.User, error)
	SetDeactiveUser(id int) (model.User, error)
	FetchAuth(userID uint, authUUID string) (model.Auth, error)
	DeleteAuth(userID uint, authUUID string) error
	CreateAuth(userID uint, expired int64, typeAuth string) (model.Auth, error)
	DeleteExpiredAuth(currentMillis int64) error
	SetNewPassword(userID int, password string) error
}

type authService struct {
	authRepo repo.AuthRepo
}

func NewAuthService(authRepo repo.AuthRepo) AuthService {
	return &authService{authRepo: authRepo}
}

func (s authService) CheckID(id int) bool {
	return s.authRepo.CheckID(id)
}

func (s authService) CheckUsername(username string) bool {
	return s.authRepo.CheckUsername(username)
}

func (s authService) CheckEmail(email string) bool {
	return s.authRepo.CheckEmail(email)
}

func (s authService) CheckHandphone(handphone string) bool {
	return s.authRepo.CheckHandphone(handphone)
}

func (s authService) CheckIsActive(username string) bool {
	return s.authRepo.CheckIsActive(username)
}

func (s authService) IsSuperAdmin(username string) (bool, error) {
	isSuperAdmin, err := s.authRepo.IsSuperAdmin(username)
	if err != nil {
		return false, err
	}
	return isSuperAdmin, nil
}

func (s authService) IsAdmin(username string) (bool, error) {
	isAdmin, err := s.authRepo.IsAdmin(username)
	if err != nil {
		return false, err
	}
	return isAdmin, nil
}

func (s authService) IsUser(username string) (bool, error) {
	isUser, err := s.authRepo.IsUser(username)
	if err != nil {
		return false, err
	}
	return isUser, nil
}

func (s authService) GetRole(username string) (bool, bool, bool, error) {
	isSuper, isAdmin, isUser, err := s.authRepo.GetRole(username)
	if err != nil {
		return false, false, false, err
	}
	return isSuper, isAdmin, isUser, nil
}

func (s authService) GetEmail(username string) (string, error) {
	role, err := s.authRepo.GetEmail(username)
	if err != nil {
		return "", err
	}
	return role, nil
}

func (s authService) Register(user model.User) error {
	if err := s.authRepo.Register(user); err != nil {
		return err
	}

	return nil
}

func (s authService) Login(username string) (string, error) {
	password, err := s.authRepo.Login(username)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (s authService) GetByUsername(username string) (user model.User, err error) {
	userData, err := s.authRepo.GetByUsername(username)
	if err != nil {
		return model.User{}, err
	}
	return userData, nil
}

func (s authService) GetByEmail(email string) (user model.User, err error) {
	userData, err := s.authRepo.GetByUsername(email)
	if err != nil {
		return model.User{}, err
	}
	return userData, nil
}

func (s authService) GetByID(id uint) (user model.User, err error) {
	userData, err := s.authRepo.GetByID(id)
	if err != nil {
		return model.User{}, err
	}
	return userData, nil
}

func (s authService) Create(user model.User) error {
	if err := s.authRepo.Create(user); err != nil {
		return err
	}

	return nil
}

func (s authService) Delete(id int) error {
	if err := s.authRepo.Delete(id); err != nil {
		return err
	}

	return nil
}

func (s authService) SetActiveUser(id int) (model.User, error) {
	userUpdate, err := s.authRepo.SetActiveUser(id)
	if err != nil {
		return model.User{}, err
	}
	return userUpdate, nil
}

func (s authService) SetDeactiveUser(id int) (model.User, error) {
	userUpdate, err := s.authRepo.SetDeactiveUser(id)
	if err != nil {
		return model.User{}, err
	}
	return userUpdate, nil
}

func (s authService) FetchAuth(userID uint, authUUID string) (model.Auth, error) {
	auth, err := s.authRepo.FetchAuth(userID, authUUID)
	if err != nil {
		return model.Auth{}, err
	}
	return auth, nil
}

func (s authService) DeleteAuth(userID uint, authUUID string) error {
	err := s.authRepo.DeleteAuth(userID, authUUID)
	if err != nil {
		return err
	}
	return nil
}

func (s authService) CreateAuth(userID uint, expired int64, typeAuth string) (model.Auth, error) {
	auth, err := s.authRepo.CreateAuth(userID, expired, typeAuth)
	if err != nil {
		return model.Auth{}, err
	}
	return auth, nil
}

func (s authService) DeleteExpiredAuth(currentMillis int64) error {
	err := s.authRepo.DeleteExpiredAuth(currentMillis)
	if err != nil {
		return err
	}
	return nil
}

func (s authService) SetNewPassword(userID int, hashPassword string) error {
	err := s.authRepo.SetNewPassword(userID, hashPassword)
	if err != nil {
		return err
	}
	return nil
}
