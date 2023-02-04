package mocks

import (
	"attendance-api/model"
	"time"

	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) CheckUsername(username string) bool {
	if err := m.Called(username).Error(0); err != nil {
		return false
	}

	return true
}

func (m *AuthServiceMock) CheckEmail(email string) bool {
	if err := m.Called(email).Error(0); err != nil {
		return false
	}

	return true
}

func (m *AuthServiceMock) CheckHandphone(handphone string) bool {
	if err := m.Called(handphone).Error(0); err != nil {
		return false
	}

	return true
}

func (m *AuthServiceMock) CheckIsActive(username string) bool {
	if err := m.Called(username).Error(0); err != nil {
		return false
	}

	return true
}

func (m *AuthServiceMock) GetRole(username string) (string, error) {
	if err := m.Called(username).Error(0); err != nil {
		return "", err
	}

	return "user", nil
}

func (m *AuthServiceMock) GetEmail(username string) (string, error) {
	if err := m.Called(username).Error(0); err != nil {
		return "", err
	}

	return "admin@gmail.com", nil
}

func (m *AuthServiceMock) Register(user *model.User) error {
	if err := m.Called(user).Error(0); err != nil {
		return err
	}

	return nil
}

func (m *AuthServiceMock) Login(username string) (string, error) {
	if err := m.Called(username).Error(0); err != nil {
		return "", err
	}

	return "$2a$10$fk9IPSmo/VYhu5VJm.vPy.5.XVowBHU3otSDAzTBpMR3YpX2cqYwW", nil
}

func (m *AuthServiceMock) CheckID(id int) bool {
	if err := m.Called(id).Error(0); err != nil {
		return false
	}

	return true
}

func (m *AuthServiceMock) GetByUsername(username string) (data *model.User, err error) {
	if err := m.Called(username).Error(0); err != nil {
		return nil, err
	}

	data.ID = 1
	data.Email = "dewok@gmail.com"
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.Name = "Dewok Satria"
	data.IsActive = true
	data.Username = "dewoklucu"
	data.Handphone = "098121342"
	data.Role = "user"
	data.Password = "password123"

	return data, nil
}

func (m *AuthServiceMock) GetByEmail(email string) (data *model.User, err error) {
	if err := m.Called(email).Error(0); err != nil {
		return nil, err
	}

	data.ID = 1
	data.Email = "dewok@gmail.com"
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()
	data.Name = "Dewok Satria"
	data.IsActive = true
	data.Username = "dewoklucu"
	data.Handphone = "098121342"
	data.Role = "user"
	data.Password = "password123"

	return data, nil
}

func (m *AuthServiceMock) Create(user *model.User) error {
	if err := m.Called(user).Error(0); err != nil {
		return err
	}

	return nil
}

func (m *AuthServiceMock) Delete(id int) error {
	if err := m.Called(id).Error(0); err != nil {
		return err
	}

	return nil
}
