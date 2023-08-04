package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type RoleAbilityService interface {
	CreateRoleAbility(roleAbility model.RoleAbility) (model.RoleAbility, error)
	RetrieveRoleAbility(id int) (model.RoleAbility, error)
	RetrieveRoleAbilityByRole(isSuperAdmin bool, isAdmin bool, isUser bool) ([]model.Ability, error)
	UpdateRoleAbility(id int, roleAbility model.RoleAbility) (model.RoleAbility, error)
	DeleteRoleAbility(id int) error
	ListRoleAbility(roleAbility model.RoleAbility, pagination model.Pagination) ([]model.RoleAbility, error)
	ListRoleAbilityMeta(roleAbility model.RoleAbility, pagination model.Pagination) (model.Meta, error)
	DropDownRoleAbility(roleAbility model.RoleAbility) ([]model.RoleAbility, error)
	CheckIsExist(id int) (isExist bool)
	CheckIsExistByActionAndSubject(action string, subject string, exceptID int) (isExist bool)
}

type roleAbilityService struct {
	roleAbilityRepo repo.RoleAbilityRepo
}

func NewRoleAbilityService(roleAbilityRepo repo.RoleAbilityRepo) RoleAbilityService {
	return &roleAbilityService{roleAbilityRepo: roleAbilityRepo}
}

func (s roleAbilityService) CreateRoleAbility(roleAbility model.RoleAbility) (model.RoleAbility, error) {
	data, err := s.roleAbilityRepo.CreateRoleAbility(roleAbility)
	if err != nil {
		return model.RoleAbility{}, err
	}
	return data, nil
}

func (s roleAbilityService) RetrieveRoleAbility(id int) (model.RoleAbility, error) {
	data, err := s.roleAbilityRepo.RetrieveRoleAbility(id)
	if err != nil {
		return model.RoleAbility{}, err
	}
	return data, nil
}

func (s roleAbilityService) RetrieveRoleAbilityByRole(isSuperAdmin bool, isAdmin bool, isUser bool) ([]model.Ability, error) {
	data, err := s.roleAbilityRepo.RetrieveRoleAbilityByRole(isSuperAdmin, isAdmin, isUser)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s roleAbilityService) UpdateRoleAbility(id int, roleAbility model.RoleAbility) (model.RoleAbility, error) {
	data, err := s.roleAbilityRepo.UpdateRoleAbility(id, roleAbility)
	if err != nil {
		return model.RoleAbility{}, err
	}
	return data, nil
}

func (s roleAbilityService) DeleteRoleAbility(id int) error {
	if err := s.roleAbilityRepo.DeleteRoleAbility(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s roleAbilityService) ListRoleAbility(roleAbility model.RoleAbility, pagination model.Pagination) ([]model.RoleAbility, error) {
	datas, err := s.roleAbilityRepo.ListRoleAbility(roleAbility, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s roleAbilityService) ListRoleAbilityMeta(roleAbility model.RoleAbility, pagination model.Pagination) (model.Meta, error) {
	data, err := s.roleAbilityRepo.ListRoleAbilityMeta(roleAbility, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s roleAbilityService) DropDownRoleAbility(roleAbility model.RoleAbility) ([]model.RoleAbility, error) {
	datas, err := s.roleAbilityRepo.DropDownRoleAbility(roleAbility)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s roleAbilityService) CheckIsExist(id int) (isExist bool) {
	return s.roleAbilityRepo.CheckIsExist(id)
}

func (s roleAbilityService) CheckIsExistByActionAndSubject(action string, subject string, exceptID int) (isExist bool) {
	return s.roleAbilityRepo.CheckIsExistByActionAndSubject(action, subject, exceptID)
}
