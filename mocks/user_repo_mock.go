package mocks

import (
	"attendance-api/model"
	"strconv"
	"time"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (m *UserRepoMock) ListUser(user *model.User, pagination *model.Pagination) (*[]model.User, error) {
	if err := m.Called(user, pagination).Error(0); err != nil {
		return nil, err
	}

	var users []model.User
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")
	for i := 1; i <= 3; i++ {
		gormData := model.GormCustom{
			ID:        uint(i),
			CreatedAt: dateTimeData,
			UpdatedAt: dateTimeData,
		}
		userData := model.User{
			GormCustom:   gormData,
			Username:     "windowsdewa" + strconv.Itoa(i),
			Password:     "Password123",
			FirstName:    "Dewok",
			LastName:     "Satria " + strconv.Itoa(i),
			Handphone:    "08122233344" + strconv.Itoa(i),
			Email:        "windowsdewa" + strconv.Itoa(i) + ".com",
			Intro:        "Hay guysss",
			Profile:      "My Name is Dewok " + strconv.Itoa(i),
			LastLogin:    dateTimeData,
			IsActive:     true,
			IsSuperAdmin: true,
			IsAdmin:      false,
			IsUser:       false,
		}
		users = append(users, userData)
	}

	return &users, nil
}

func (m *UserRepoMock) ListUserMeta(user *model.User, pagination *model.Pagination) (*model.Meta, error) {
	if err := m.Called(user, pagination).Error(0); err != nil {
		return nil, err
	}

	metaData := model.Meta{
		TotalPage:     1,
		CurrentPage:   1,
		TotalRecord:   3,
		CurrentRecord: 3,
	}

	return &metaData, nil
}

func (m *UserRepoMock) CreateUser(user *model.User) (*model.User, error) {
	if err := m.Called(user).Error(0); err != nil {
		return nil, err
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

	return &userData, nil
}

func (m *UserRepoMock) UpdateUser(id int, user *model.User) (*model.User, error) {
	if err := m.Called(id, user).Error(0); err != nil {
		return nil, err
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

	return &userData, nil
}

func (m *UserRepoMock) HardDeleteUser(id int) error {
	if err := m.Called(id).Error(0); err != nil {
		return err
	}

	return nil
}

func (m *UserRepoMock) SetActiveUser(id int) (*model.User, error) {
	if err := m.Called(id).Error(0); err != nil {
		return nil, err
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

	return &userData, nil
}

func (m *UserRepoMock) SetDeactiveUser(id int) (*model.User, error) {
	if err := m.Called(id).Error(0); err != nil {
		return nil, err
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

	return &userData, nil
}
