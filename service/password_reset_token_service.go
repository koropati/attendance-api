package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type PasswordResetTokenService interface {
	CreatePasswordResetToken(subject model.PasswordResetToken) (model.PasswordResetToken, error)
	RetrievePasswordResetToken(id int) (model.PasswordResetToken, error)
	UpdatePasswordResetToken(id int, subject model.PasswordResetToken) (model.PasswordResetToken, error)
	DeletePasswordResetToken(id int) error
	ListPasswordResetToken(subject model.PasswordResetToken, pagination model.Pagination) ([]model.PasswordResetToken, error)
	ListPasswordResetTokenMeta(subject model.PasswordResetToken, pagination model.Pagination) (model.Meta, error)
	DropDownPasswordResetToken(subject model.PasswordResetToken) ([]model.PasswordResetToken, error)
}

type passwordResetTokenService struct {
	passwordResetTokenRepo repo.PasswordResetTokenRepo
}

func NewPasswordResetTokenService(passwordResetTokenRepo repo.PasswordResetTokenRepo) PasswordResetTokenService {
	return &passwordResetTokenService{passwordResetTokenRepo: passwordResetTokenRepo}
}

func (s passwordResetTokenService) CreatePasswordResetToken(subject model.PasswordResetToken) (model.PasswordResetToken, error) {
	data, err := s.passwordResetTokenRepo.CreatePasswordResetToken(subject)
	if err != nil {
		return model.PasswordResetToken{}, err
	}
	return data, nil
}

func (s passwordResetTokenService) RetrievePasswordResetToken(id int) (model.PasswordResetToken, error) {
	data, err := s.passwordResetTokenRepo.RetrievePasswordResetToken(id)
	if err != nil {
		return model.PasswordResetToken{}, err
	}
	return data, nil
}

func (s passwordResetTokenService) UpdatePasswordResetToken(id int, subject model.PasswordResetToken) (model.PasswordResetToken, error) {
	data, err := s.passwordResetTokenRepo.UpdatePasswordResetToken(id, subject)
	if err != nil {
		return model.PasswordResetToken{}, err
	}
	return data, nil
}

func (s passwordResetTokenService) DeletePasswordResetToken(id int) error {
	if err := s.passwordResetTokenRepo.DeletePasswordResetToken(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s passwordResetTokenService) ListPasswordResetToken(subject model.PasswordResetToken, pagination model.Pagination) ([]model.PasswordResetToken, error) {
	datas, err := s.passwordResetTokenRepo.ListPasswordResetToken(subject, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s passwordResetTokenService) ListPasswordResetTokenMeta(subject model.PasswordResetToken, pagination model.Pagination) (model.Meta, error) {
	data, err := s.passwordResetTokenRepo.ListPasswordResetTokenMeta(subject, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s passwordResetTokenService) DropDownPasswordResetToken(subject model.PasswordResetToken) ([]model.PasswordResetToken, error) {
	datas, err := s.passwordResetTokenRepo.DropDownPasswordResetToken(subject)
	if err != nil {
		return nil, err
	}
	return datas, nil
}
