package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type FacultyService interface {
	CreateFaculty(faculty model.Faculty) (model.Faculty, error)
	RetrieveFaculty(id int) (model.Faculty, error)
	RetrieveFacultyByOwner(id int, ownerID int) (model.Faculty, error)
	UpdateFaculty(id int, faculty model.Faculty) (model.Faculty, error)
	UpdateFacultyByOwner(id int, ownerID int, faculty model.Faculty) (model.Faculty, error)
	DeleteFaculty(id int) error
	DeleteFacultyByOwner(id int, ownerID int) error
	ListFaculty(faculty model.Faculty, pagination model.Pagination) ([]model.Faculty, error)
	ListFacultyMeta(faculty model.Faculty, pagination model.Pagination) (model.Meta, error)
	DropDownFaculty(faculty model.Faculty) ([]model.Faculty, error)
	CheckIsExist(id int) (isExist bool)
	CheckIsExistByName(name string, exceptID int) (isExist bool)
	CheckIsExistByCode(code string, exceptID int) (isExist bool)
}

type facultyService struct {
	facultyRepo repo.FacultyRepo
}

func NewFacultyService(facultyRepo repo.FacultyRepo) FacultyService {
	return &facultyService{facultyRepo: facultyRepo}
}

func (s facultyService) CreateFaculty(faculty model.Faculty) (model.Faculty, error) {
	data, err := s.facultyRepo.CreateFaculty(faculty)
	if err != nil {
		return model.Faculty{}, err
	}
	return data, nil
}

func (s facultyService) RetrieveFaculty(id int) (model.Faculty, error) {
	data, err := s.facultyRepo.RetrieveFaculty(id)
	if err != nil {
		return model.Faculty{}, err
	}
	return data, nil
}

func (s facultyService) RetrieveFacultyByOwner(id int, ownerID int) (model.Faculty, error) {
	data, err := s.facultyRepo.RetrieveFacultyByOwner(id, ownerID)
	if err != nil {
		return model.Faculty{}, err
	}
	return data, nil
}

func (s facultyService) UpdateFaculty(id int, faculty model.Faculty) (model.Faculty, error) {
	data, err := s.facultyRepo.UpdateFaculty(id, faculty)
	if err != nil {
		return model.Faculty{}, err
	}
	return data, nil
}

func (s facultyService) UpdateFacultyByOwner(id int, ownerID int, faculty model.Faculty) (model.Faculty, error) {
	data, err := s.facultyRepo.UpdateFacultyByOwner(id, ownerID, faculty)
	if err != nil {
		return model.Faculty{}, err
	}
	return data, nil
}

func (s facultyService) DeleteFaculty(id int) error {
	if err := s.facultyRepo.DeleteFaculty(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s facultyService) DeleteFacultyByOwner(id int, ownerID int) error {
	if err := s.facultyRepo.DeleteFacultyByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s facultyService) ListFaculty(faculty model.Faculty, pagination model.Pagination) ([]model.Faculty, error) {
	datas, err := s.facultyRepo.ListFaculty(faculty, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s facultyService) ListFacultyMeta(faculty model.Faculty, pagination model.Pagination) (model.Meta, error) {
	data, err := s.facultyRepo.ListFacultyMeta(faculty, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s facultyService) DropDownFaculty(faculty model.Faculty) ([]model.Faculty, error) {
	datas, err := s.facultyRepo.DropDownFaculty(faculty)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s facultyService) CheckIsExist(id int) (isExist bool) {
	return s.facultyRepo.CheckIsExist(id)
}

func (s facultyService) CheckIsExistByName(name string, exceptID int) (isExist bool) {
	return s.facultyRepo.CheckIsExistByName(name, exceptID)
}

func (s facultyService) CheckIsExistByCode(code string, exceptID int) (isExist bool) {
	return s.facultyRepo.CheckIsExistByCode(code, exceptID)
}
