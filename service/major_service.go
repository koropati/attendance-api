package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type MajorService interface {
	CreateMajor(major model.Major) (model.Major, error)
	RetrieveMajor(id int) (model.Major, error)
	RetrieveMajorByOwner(id int, ownerID int) (model.Major, error)
	UpdateMajor(id int, major model.Major) (model.Major, error)
	UpdateMajorByOwner(id int, ownerID int, major model.Major) (model.Major, error)
	DeleteMajor(id int) error
	DeleteMajorByOwner(id int, ownerID int) error
	ListMajor(major model.Major, pagination model.Pagination) ([]model.Major, error)
	ListMajorMeta(major model.Major, pagination model.Pagination) (model.Meta, error)
	DropDownMajor(major model.Major) ([]model.Major, error)
	CheckIsExist(id int) (isExist bool)
}

type majorService struct {
	majorRepo repo.MajorRepo
}

func NewMajorService(majorRepo repo.MajorRepo) MajorService {
	return &majorService{majorRepo: majorRepo}
}

func (s majorService) CreateMajor(major model.Major) (model.Major, error) {
	data, err := s.majorRepo.CreateMajor(major)
	if err != nil {
		return model.Major{}, err
	}
	return data, nil
}

func (s majorService) RetrieveMajor(id int) (model.Major, error) {
	data, err := s.majorRepo.RetrieveMajor(id)
	if err != nil {
		return model.Major{}, err
	}
	return data, nil
}

func (s majorService) RetrieveMajorByOwner(id int, ownerID int) (model.Major, error) {
	data, err := s.majorRepo.RetrieveMajorByOwner(id, ownerID)
	if err != nil {
		return model.Major{}, err
	}
	return data, nil
}

func (s majorService) UpdateMajor(id int, major model.Major) (model.Major, error) {
	data, err := s.majorRepo.UpdateMajor(id, major)
	if err != nil {
		return model.Major{}, err
	}
	return data, nil
}

func (s majorService) UpdateMajorByOwner(id int, ownerID int, major model.Major) (model.Major, error) {
	data, err := s.majorRepo.UpdateMajorByOwner(id, ownerID, major)
	if err != nil {
		return model.Major{}, err
	}
	return data, nil
}

func (s majorService) DeleteMajor(id int) error {
	if err := s.majorRepo.DeleteMajor(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s majorService) DeleteMajorByOwner(id int, ownerID int) error {
	if err := s.majorRepo.DeleteMajorByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s majorService) ListMajor(major model.Major, pagination model.Pagination) ([]model.Major, error) {
	datas, err := s.majorRepo.ListMajor(major, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s majorService) ListMajorMeta(major model.Major, pagination model.Pagination) (model.Meta, error) {
	data, err := s.majorRepo.ListMajorMeta(major, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s majorService) DropDownMajor(major model.Major) ([]model.Major, error) {
	datas, err := s.majorRepo.DropDownMajor(major)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s majorService) CheckIsExist(id int) (isExist bool) {
	return s.majorRepo.CheckIsExist(id)
}
