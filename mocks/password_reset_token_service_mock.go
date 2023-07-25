package mocks

import (
	"attendance-api/model"
	"strconv"
	"time"

	"github.com/stretchr/testify/mock"
)

type PasswordResetTokenServiceMock struct {
	mock.Mock
}

func (m PasswordResetTokenServiceMock) CreatePasswordResetToken(user model.PasswordResetToken) (model.PasswordResetToken, error) {
	if err := m.Called(user).Error(0); err != nil {
		return model.PasswordResetToken{}, err
	}
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")
	gormData := model.GormCustom{
		ID:        1,
		CreatedAt: dateTimeData,
		UpdatedAt: dateTimeData,
	}

	tokenData := model.PasswordResetToken{
		GormCustom: gormData,
		UserID:     1,
		Token:      "v34234234g2564624457",
		Valid:      dateTimeData,
	}

	return tokenData, nil
}

func (m PasswordResetTokenServiceMock) RetrievePasswordResetToken(id int) (model.PasswordResetToken, error) {
	if err := m.Called(id).Error(0); err != nil {
		return model.PasswordResetToken{}, err
	}
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")
	gormData := model.GormCustom{
		ID:        1,
		CreatedAt: dateTimeData,
		UpdatedAt: dateTimeData,
	}

	userData := model.PasswordResetToken{
		GormCustom: gormData,
		UserID:     1,
		Token:      "sfaw452q365",
		Valid:      dateTimeData,
	}

	return userData, nil
}

func (m PasswordResetTokenServiceMock) UpdatePasswordResetToken(id int, user model.PasswordResetToken) (model.PasswordResetToken, error) {
	if err := m.Called(id, user).Error(0); err != nil {
		return model.PasswordResetToken{}, err
	}
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")

	gormData := model.GormCustom{
		ID:        1,
		CreatedAt: dateTimeData,
		UpdatedAt: dateTimeData,
	}

	data := model.PasswordResetToken{
		GormCustom: gormData,
		UserID:     1,
		Token:      "24t23g45262346",
		Valid:      dateTimeData,
	}

	return data, nil
}

func (m PasswordResetTokenServiceMock) DeletePasswordResetToken(id int) error {
	if err := m.Called(id).Error(0); err != nil {
		return err
	}

	return nil
}

func (m PasswordResetTokenServiceMock) ListPasswordResetToken(user model.PasswordResetToken, pagination model.Pagination) ([]model.PasswordResetToken, error) {
	if err := m.Called(user, pagination).Error(0); err != nil {
		return nil, err
	}

	var datas []model.PasswordResetToken
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")
	for i := 1; i <= 3; i++ {
		gormData := model.GormCustom{
			ID:        uint(i),
			CreatedAt: dateTimeData,
			UpdatedAt: dateTimeData,
		}
		data := model.PasswordResetToken{
			GormCustom: gormData,
			UserID:     uint(i),
			Token:      "cwqjh2qr502af" + strconv.Itoa(i),
			Valid:      dateTimeData,
		}
		datas = append(datas, data)
	}

	return datas, nil
}

func (m PasswordResetTokenServiceMock) ListPasswordResetTokenMeta(user model.PasswordResetToken, pagination model.Pagination) (model.Meta, error) {
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

func (m PasswordResetTokenServiceMock) DropDownPasswordResetToken(user model.PasswordResetToken) ([]model.PasswordResetToken, error) {
	if err := m.Called(user).Error(0); err != nil {
		return nil, err
	}

	var datas []model.PasswordResetToken
	dateTimeData, _ := time.Parse("2006-01-02T15:04:05-0700", "2006-01-02T15:04:05-0700")
	for i := 1; i <= 3; i++ {
		gormData := model.GormCustom{
			ID:        uint(i),
			CreatedAt: dateTimeData,
			UpdatedAt: dateTimeData,
		}
		data := model.PasswordResetToken{
			GormCustom: gormData,
			UserID:     uint(i),
			Token:      "cwqjh2qr502af" + strconv.Itoa(i),
			Valid:      dateTimeData,
		}
		datas = append(datas, data)
	}

	return datas, nil
}

func (m PasswordResetTokenServiceMock) IsValid(token string) (isValid bool, userID uint) {
	if err := m.Called(token).Error(0); err != nil {
		return false, 0
	}

	return true, 1
}

func (m PasswordResetTokenServiceMock) DeleteExpiredPasswordResetToken(currentTime time.Time) error {
	if err := m.Called(currentTime).Error(0); err != nil {
		return err
	}
	return nil
}
