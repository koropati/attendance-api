package mocks

import (
	"attendance-api/model"
	"time"

	"github.com/stretchr/testify/mock"
)

type AuthRepoMock struct {
	mock.Mock
}

func (m AuthRepoMock) CheckID(id int) bool {
	if err := m.Called(id).Error(0); err != nil {
		return false
	}

	return true
}

func (m AuthRepoMock) CheckUsername(username string) bool {
	if err := m.Called(username).Error(0); err != nil {
		return false
	}

	return true
}

func (m AuthRepoMock) CheckEmail(email string) bool {
	if err := m.Called(email).Error(0); err != nil {
		return false
	}

	return true
}

func (m AuthRepoMock) CheckHandphone(handphone string) bool {
	if err := m.Called(handphone).Error(0); err != nil {
		return false
	}

	return true
}

func (m AuthRepoMock) CheckIsActive(username string) bool {
	if err := m.Called(username).Error(0); err != nil {
		return false
	}

	return true
}

func (m AuthRepoMock) IsSuperAdmin(username string) (bool, error) {
	if err := m.Called(username).Error(0); err != nil {
		return false, err
	}

	return true, nil
}

func (m AuthRepoMock) IsAdmin(username string) (bool, error) {
	if err := m.Called(username).Error(0); err != nil {
		return false, err
	}

	return false, nil
}

func (m AuthRepoMock) IsUser(username string) (bool, error) {
	if err := m.Called(username).Error(0); err != nil {
		return false, err
	}

	return false, nil
}

func (m AuthRepoMock) GetRole(username string) (bool, bool, bool, error) {
	if err := m.Called(username).Error(0); err != nil {
		return false, false, false, err
	}

	return true, false, false, nil
}

func (m AuthRepoMock) GetEmail(username string) (string, error) {
	if err := m.Called(username).Error(0); err != nil {
		return "", err
	}

	return "admin@gmail.com", nil
}

func (m AuthRepoMock) Register(user model.User) error {
	if err := m.Called(user).Error(0); err != nil {
		return err
	}

	return nil
}

func (m AuthRepoMock) Login(username string) (string, error) {
	if err := m.Called(username).Error(0); err != nil {
		return "", err
	}

	return "$2a$10$fk9IPSmo/VYhu5VJm.vPy.5.XVowBHU3otSDAzTBpMR3YpX2cqYwW", nil
}

func (m AuthRepoMock) GetByUsername(username string) (data model.User, err error) {
	if err := m.Called(username).Error(0); err != nil {
		return model.User{}, err
	}

	data.ID = 1
	data.Email = "dewok@gmail.com"
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.FirstName = "Dewok"
	data.LastName = "Satria"
	data.Username = "dewoklucu"
	data.Handphone = "098121342"
	data.Password = "password123"
	data.IsActive = true
	data.IsSuperAdmin = true
	data.IsAdmin = false
	data.IsUser = false

	return data, nil
}

func (m AuthRepoMock) GetByEmail(email string) (data model.User, err error) {
	if err := m.Called(email).Error(0); err != nil {
		return model.User{}, err
	}

	data.ID = 1
	data.Email = "dewok@gmail.com"
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.FirstName = "Dewok"
	data.LastName = "Satria"
	data.Username = "dewoklucu"
	data.Handphone = "098121342"
	data.Password = "password123"
	data.IsActive = true
	data.IsSuperAdmin = true
	data.IsAdmin = false
	data.IsUser = false

	return data, nil
}

func (m AuthRepoMock) GetByID(id uint) (data model.User, err error) {
	if err := m.Called(id).Error(0); err != nil {
		return model.User{}, err
	}

	data.ID = 1
	data.Email = "dewok@gmail.com"
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.FirstName = "Dewok"
	data.LastName = "Satria"
	data.Username = "dewoklucu"
	data.Handphone = "098121342"
	data.Password = "password123"
	data.IsActive = true
	data.IsSuperAdmin = true
	data.IsAdmin = false
	data.IsUser = false

	return data, nil
}

func (m AuthRepoMock) Create(user model.User) error {
	if err := m.Called(user).Error(0); err != nil {
		return err
	}

	return nil
}

func (m AuthRepoMock) Delete(id int) error {
	if err := m.Called(id).Error(0); err != nil {
		return err
	}

	return nil
}

func (m AuthRepoMock) SetActiveUser(id int) (model.User, error) {
	if err := m.Called(id).Error(0); err != nil {
		return model.User{}, err
	}
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")

	gormData := model.GormCustom{
		ID:        1,
		CreatedAt: dateTimeData,
		UpdatedAt: dateTimeData,
	}

	userData := model.User{
		GormCustom:   gormData,
		Username:     "windowsdewa",
		Password:     "Password123",
		FirstName:    "Dewok",
		LastName:     "Satria",
		Handphone:    "081222333440",
		Email:        "windowsdewa.com",
		Intro:        "Hay guysss",
		Profile:      "My Name is Dewok ",
		LastLogin:    dateTimeData,
		IsActive:     true,
		IsSuperAdmin: true,
		IsAdmin:      false,
		IsUser:       false,
	}

	return userData, nil
}

func (m AuthRepoMock) SetDeactiveUser(id int) (model.User, error) {
	if err := m.Called(id).Error(0); err != nil {
		return model.User{}, err
	}

	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")

	gormData := model.GormCustom{
		ID:        1,
		CreatedAt: dateTimeData,
		UpdatedAt: dateTimeData,
	}

	userData := model.User{
		GormCustom:   gormData,
		Username:     "windowsdewa",
		Password:     "Password123",
		FirstName:    "Dewok",
		LastName:     "Satria",
		Handphone:    "081222333440",
		Email:        "windowsdewa.com",
		Intro:        "Hay guysss",
		Profile:      "My Name is Dewok ",
		LastLogin:    dateTimeData,
		IsActive:     true,
		IsSuperAdmin: true,
		IsAdmin:      false,
		IsUser:       false,
	}

	return userData, nil
}

func (m AuthRepoMock) FetchAuth(userID uint, authUUID string) (model.Auth, error) {
	if err := m.Called(userID, authUUID).Error(0); err != nil {
		return model.Auth{}, err
	}

	authData := model.Auth{
		UserID:   1,
		AuthUUID: "qwerty123456",
		Expired:  1234567890,
		TypeAuth: "at",
	}

	return authData, nil
}

func (m AuthRepoMock) DeleteAuth(userID uint, authUUID string) error {
	if err := m.Called(userID, authUUID).Error(0); err != nil {
		return err
	}

	return nil
}

func (m AuthRepoMock) CreateAuth(userID uint, expired int64, typeAuth string) (model.Auth, error) {
	if err := m.Called(userID, expired, typeAuth).Error(0); err != nil {
		return model.Auth{}, err
	}

	authData := model.Auth{
		UserID:   1,
		AuthUUID: "qwerty123456",
		Expired:  1234567890,
		TypeAuth: "at",
	}

	return authData, nil
}
