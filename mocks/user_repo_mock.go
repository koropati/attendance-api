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

func (m UserRepoMock) CheckID(id int) bool {
	if err := m.Called(id).Error(0); err != nil {
		return false
	}

	return true
}

func (m UserRepoMock) CheckUsername(username string) bool {
	if err := m.Called(username).Error(0); err != nil {
		return false
	}

	return true
}

func (m UserRepoMock) CheckEmail(email string) bool {
	if err := m.Called(email).Error(0); err != nil {
		return false
	}

	return true
}

func (m UserRepoMock) CheckHandphone(handphone string) bool {
	if err := m.Called(handphone).Error(0); err != nil {
		return false
	}

	return true
}

func (m UserRepoMock) CheckIsActive(username string) bool {
	if err := m.Called(username).Error(0); err != nil {
		return false
	}

	return true
}

func (m UserRepoMock) CheckUpdateUsername(id int, username string) bool {
	if err := m.Called(id, username).Error(0); err != nil {
		return false
	}

	return true
}

func (m UserRepoMock) CheckUpdateEmail(id int, email string) bool {
	if err := m.Called(id, email).Error(0); err != nil {
		return false
	}

	return true
}

func (m UserRepoMock) CheckUpdateHandphone(id int, handphone string) bool {
	if err := m.Called(id, handphone).Error(0); err != nil {
		return false
	}

	return true
}

func (m UserRepoMock) ListUser(user model.User, pagination model.Pagination) ([]model.User, error) {
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

	return users, nil
}

func (m UserRepoMock) ListUserMeta(user model.User, pagination model.Pagination) (model.Meta, error) {
	if err := m.Called(user, pagination).Error(0); err != nil {
		return model.Meta{}, err
	}

	metaData := model.Meta{
		TotalPage:     1,
		CurrentPage:   1,
		TotalRecord:   3,
		CurrentRecord: 3,
	}

	return metaData, nil
}

func (m UserRepoMock) CreateUser(user model.User) (model.User, error) {
	if err := m.Called(user).Error(0); err != nil {
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

func (m UserRepoMock) UpdateUser(id int, user model.User) (model.User, error) {
	if err := m.Called(id, user).Error(0); err != nil {
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

func (m UserRepoMock) UpdateProfile(id int, user model.User) (model.User, error) {
	if err := m.Called(id, user).Error(0); err != nil {
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

func (m UserRepoMock) DeleteUser(id int) error {
	if err := m.Called(id).Error(0); err != nil {
		return err
	}

	return nil
}

func (m UserRepoMock) RetrieveUser(id int) (result model.User, err error) {
	if err := m.Called(id).Error(0); err != nil {
		return model.User{}, err
	}
	return
}

func (m UserRepoMock) RetrieveUserByUsername(username string) (data model.User, err error) {
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

func (m UserRepoMock) RetrieveUserByEmail(email string) (data model.User, err error) {
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

func (m UserRepoMock) SetActiveUser(id int) (model.User, error) {
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

func (m UserRepoMock) SetDeactiveUser(id int) (model.User, error) {
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

func (m UserRepoMock) DropDownUser(user model.User) ([]model.UserDropDown, error) {
	if err := m.Called(user).Error(0); err != nil {
		return nil, err
	}

	var users []model.UserDropDown
	for i := 1; i <= 3; i++ {

		userData := model.UserDropDown{
			ID:        uint(i),
			Username:  "windowsdewa" + strconv.Itoa(i),
			FirstName: "Dewok",
			LastName:  "Satria " + strconv.Itoa(i),
			Handphone: "08122233344" + strconv.Itoa(i),
			Email:     "windowsdewa" + strconv.Itoa(i) + ".com",
		}
		users = append(users, userData)
	}

	return users, nil
}

func (m UserRepoMock) UpdatePassword(userPasswordForm model.UserUpdatePasswordForm) error {
	if err := m.Called(userPasswordForm).Error(0); err != nil {
		return err
	}

	return nil
}

func (m UserRepoMock) GetPassword(id int) (hashPassword string, err error) {
	if err := m.Called(id).Error(0); err != nil {
		return "", err
	}

	return "hashPasswordNya", nil
}

func (m UserRepoMock) GetAbility(user model.User) []model.Ability {
	if err := m.Called(user).Error(0); err != nil {
		return nil
	}

	var ability []model.Ability
	for i := 1; i <= 1; i++ {

		ability = append(ability, model.Ability{
			Action:  "read",
			Subject: "Auth",
		})
	}
	return ability
}
