package service

import (
	"attendance-api/model"
	"attendance-api/repo"
	"time"
)

type ActivationTokenService interface {
	CreateActivationToken(activationToken model.ActivationToken) (model.ActivationToken, error)
	RetrieveActivationToken(id int) (model.ActivationToken, error)
	UpdateActivationToken(id int, activationToken model.ActivationToken) (model.ActivationToken, error)
	DeleteActivationToken(id int) error
	ListActivationToken(activationToken model.ActivationToken, pagination model.Pagination) ([]model.ActivationToken, error)
	ListActivationTokenMeta(activationToken model.ActivationToken, pagination model.Pagination) (model.Meta, error)
	DropDownActivationToken(activationToken model.ActivationToken) ([]model.ActivationToken, error)
	IsValid(token string) (isValid bool, userID uint)
	DeleteExpiredActivationToken(currentTime time.Time) error
}

type activationTokenService struct {
	activationTokenRepo repo.ActivationTokenRepo
}

func NewActivationTokenService(activationTokenRepo repo.ActivationTokenRepo) ActivationTokenService {
	return &activationTokenService{activationTokenRepo: activationTokenRepo}
}

func (s activationTokenService) CreateActivationToken(activationToken model.ActivationToken) (model.ActivationToken, error) {
	data, err := s.activationTokenRepo.CreateActivationToken(activationToken)
	if err != nil {
		return model.ActivationToken{}, err
	}
	return data, nil
}

func (s activationTokenService) RetrieveActivationToken(id int) (model.ActivationToken, error) {
	data, err := s.activationTokenRepo.RetrieveActivationToken(id)
	if err != nil {
		return model.ActivationToken{}, err
	}
	return data, nil
}

func (s activationTokenService) UpdateActivationToken(id int, activationToken model.ActivationToken) (model.ActivationToken, error) {
	data, err := s.activationTokenRepo.UpdateActivationToken(id, activationToken)
	if err != nil {
		return model.ActivationToken{}, err
	}
	return data, nil
}

func (s activationTokenService) DeleteActivationToken(id int) error {
	if err := s.activationTokenRepo.DeleteActivationToken(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s activationTokenService) ListActivationToken(activationToken model.ActivationToken, pagination model.Pagination) ([]model.ActivationToken, error) {
	datas, err := s.activationTokenRepo.ListActivationToken(activationToken, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s activationTokenService) ListActivationTokenMeta(activationToken model.ActivationToken, pagination model.Pagination) (model.Meta, error) {
	data, err := s.activationTokenRepo.ListActivationTokenMeta(activationToken, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s activationTokenService) DropDownActivationToken(activationToken model.ActivationToken) ([]model.ActivationToken, error) {
	datas, err := s.activationTokenRepo.DropDownActivationToken(activationToken)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s activationTokenService) IsValid(token string) (isValid bool, userID uint) {
	isValid, userID = s.activationTokenRepo.IsValid(token)
	return
}

func (s activationTokenService) DeleteExpiredActivationToken(currentTime time.Time) error {
	err := s.activationTokenRepo.DeleteExpiredActivationToken(currentTime)
	if err != nil {
		return err
	}
	return nil
}
