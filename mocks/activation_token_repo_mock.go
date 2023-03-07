package mocks

import (
	"attendance-api/model"
	"strconv"
	"time"

	"github.com/stretchr/testify/mock"
)

type ActivationTokenRepoMock struct {
	mock.Mock
}

func (m *ActivationTokenRepoMock) CreateActivationToken(token *model.ActivationToken) (*model.ActivationToken, error) {
	if err := m.Called(token).Error(0); err != nil {
		return nil, err
	}
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")
	gormData := model.GormCustom{
		ID:        1,
		CreatedAt: dateTimeData,
		UpdatedAt: dateTimeData,
	}

	tokenData := model.ActivationToken{
		GormCustom: gormData,
		UserID:     1,
		Token:      "v34234234g2564624457",
		Valid:      dateTimeData,
	}

	return &tokenData, nil
}

func (m *ActivationTokenRepoMock) RetrieveActivationToken(id int) (result *model.ActivationToken, err error) {
	if err := m.Called(id).Error(0); err != nil {
		return nil, err
	}
	return
}

func (m *ActivationTokenRepoMock) UpdateActivationToken(id int, user *model.ActivationToken) (*model.ActivationToken, error) {
	if err := m.Called(id, user).Error(0); err != nil {
		return nil, err
	}

	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")

	gormData := model.GormCustom{
		ID:        1,
		CreatedAt: dateTimeData,
		UpdatedAt: dateTimeData,
	}

	data := model.ActivationToken{
		GormCustom: gormData,
		UserID:     1,
		Token:      "24t23g45262346",
		Valid:      dateTimeData,
	}

	return &data, nil
}

func (m *ActivationTokenRepoMock) DeleteActivationToken(id int) error {
	if err := m.Called(id).Error(0); err != nil {
		return err
	}

	return nil
}

func (m *ActivationTokenRepoMock) ListActivationToken(user *model.ActivationToken, pagination *model.Pagination) (*[]model.ActivationToken, error) {
	if err := m.Called(user, pagination).Error(0); err != nil {
		return nil, err
	}

	var datas []model.ActivationToken
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")
	for i := 1; i <= 3; i++ {
		gormData := model.GormCustom{
			ID:        uint(i),
			CreatedAt: dateTimeData,
			UpdatedAt: dateTimeData,
		}
		data := model.ActivationToken{
			GormCustom: gormData,
			UserID:     uint(i),
			Token:      "cwqjh2qr502af" + strconv.Itoa(i),
			Valid:      dateTimeData,
		}
		datas = append(datas, data)
	}

	return &datas, nil
}

func (m *ActivationTokenRepoMock) ListActivationTokenMeta(user *model.ActivationToken, pagination *model.Pagination) (*model.Meta, error) {
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

func (m *ActivationTokenRepoMock) DropDownActivationToken(user *model.ActivationToken) (*[]model.ActivationToken, error) {
	if err := m.Called(user).Error(0); err != nil {
		return nil, err
	}

	var datas []model.ActivationToken
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")
	for i := 1; i <= 3; i++ {
		gormData := model.GormCustom{
			ID:        uint(i),
			CreatedAt: dateTimeData,
			UpdatedAt: dateTimeData,
		}
		data := model.ActivationToken{
			GormCustom: gormData,
			UserID:     uint(i),
			Token:      "cwqjh2qr502af" + strconv.Itoa(i),
			Valid:      dateTimeData,
		}
		datas = append(datas, data)
	}

	return &datas, nil
}

func (m *ActivationTokenRepoMock) IsValid(token string) (isValid bool, userID uint) {
	if err := m.Called(token).Error(0); err != nil {
		return false, 0
	}

	return true, 1
}
