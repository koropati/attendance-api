package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type ActivationTokenService interface {
	CreateActivationToken(subject model.ActivationToken) (model.ActivationToken, error)
	RetrieveActivationToken(id int) (model.ActivationToken, error)
	UpdateActivationToken(id int, subject model.ActivationToken) (model.ActivationToken, error)
	DeleteActivationToken(id int) error
	ListActivationToken(subject model.ActivationToken, pagination model.Pagination) ([]model.ActivationToken, error)
	ListActivationTokenMeta(subject model.ActivationToken, pagination model.Pagination) (model.Meta, error)
	DropDownActivationToken(subject model.ActivationToken) ([]model.ActivationToken, error)
	IsValid(token string) (isValid bool, userID uint)
}

type activationTokenService struct {
	activationTokenRepo repo.ActivationTokenRepo
}

func NewActivationTokenService(activationTokenRepo repo.ActivationTokenRepo) ActivationTokenService {
	return &activationTokenService{activationTokenRepo: activationTokenRepo}
}

func (s activationTokenService) CreateActivationToken(subject model.ActivationToken) (model.ActivationToken, error) {
	data, err := s.activationTokenRepo.CreateActivationToken(subject)
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

func (s activationTokenService) UpdateActivationToken(id int, subject model.ActivationToken) (model.ActivationToken, error) {
	data, err := s.activationTokenRepo.UpdateActivationToken(id, subject)
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

func (s activationTokenService) ListActivationToken(subject model.ActivationToken, pagination model.Pagination) ([]model.ActivationToken, error) {
	datas, err := s.activationTokenRepo.ListActivationToken(subject, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s activationTokenService) ListActivationTokenMeta(subject model.ActivationToken, pagination model.Pagination) (model.Meta, error) {
	data, err := s.activationTokenRepo.ListActivationTokenMeta(subject, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s activationTokenService) DropDownActivationToken(subject model.ActivationToken) ([]model.ActivationToken, error) {
	datas, err := s.activationTokenRepo.DropDownActivationToken(subject)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s activationTokenService) IsValid(token string) (isValid bool, userID uint) {
	isValid, userID = s.activationTokenRepo.IsValid(token)
	return
}
