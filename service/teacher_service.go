package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type TeacherService interface {
	CreateTeacher(teacher model.Teacher) (model.Teacher, error)
	RetrieveTeacher(id int) (model.Teacher, error)
	RetrieveTeacherByUserID(userID int) (model.Teacher, error)
	RetrieveTeacherByOwner(id int, ownerID int) (model.Teacher, error)
	UpdateTeacher(id int, teacher model.Teacher) (model.Teacher, error)
	UpdateTeacherByOwner(id int, ownerID int, teacher model.Teacher) (model.Teacher, error)
	DeleteTeacher(id int) error
	DeleteTeacherByOwner(id int, ownerID int) error
	ListTeacher(teacher model.Teacher, pagination model.Pagination) ([]model.Teacher, error)
	ListTeacherMeta(teacher model.Teacher, pagination model.Pagination) (model.Meta, error)
	DropDownTeacher(teacher model.Teacher) ([]model.Teacher, error)
	CheckIsExist(id int) (isExist bool)
	CheckIsExistByNip(nip string, exceptID int) (isExist bool)
}

type teacherService struct {
	teacherRepo repo.TeacherRepo
}

func NewTeacherService(teacherRepo repo.TeacherRepo) TeacherService {
	return &teacherService{teacherRepo: teacherRepo}
}

func (s teacherService) CreateTeacher(teacher model.Teacher) (model.Teacher, error) {
	data, err := s.teacherRepo.CreateTeacher(teacher)
	if err != nil {
		return model.Teacher{}, err
	}
	return data, nil
}

func (s teacherService) RetrieveTeacher(id int) (model.Teacher, error) {
	data, err := s.teacherRepo.RetrieveTeacher(id)
	if err != nil {
		return model.Teacher{}, err
	}
	return data, nil
}

func (s teacherService) RetrieveTeacherByUserID(userID int) (model.Teacher, error) {
	data, err := s.teacherRepo.RetrieveTeacherByUserID(userID)
	if err != nil {
		return model.Teacher{}, err
	}
	return data, nil
}

func (s teacherService) RetrieveTeacherByOwner(id int, ownerID int) (model.Teacher, error) {
	data, err := s.teacherRepo.RetrieveTeacherByOwner(id, ownerID)
	if err != nil {
		return model.Teacher{}, err
	}
	return data, nil
}

func (s teacherService) UpdateTeacher(id int, teacher model.Teacher) (model.Teacher, error) {
	data, err := s.teacherRepo.UpdateTeacher(id, teacher)
	if err != nil {
		return model.Teacher{}, err
	}
	return data, nil
}

func (s teacherService) UpdateTeacherByOwner(id int, ownerID int, teacher model.Teacher) (model.Teacher, error) {
	data, err := s.teacherRepo.UpdateTeacherByOwner(id, ownerID, teacher)
	if err != nil {
		return model.Teacher{}, err
	}
	return data, nil
}

func (s teacherService) DeleteTeacher(id int) error {
	if err := s.teacherRepo.DeleteTeacher(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s teacherService) DeleteTeacherByOwner(id int, ownerID int) error {
	if err := s.teacherRepo.DeleteTeacherByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s teacherService) ListTeacher(teacher model.Teacher, pagination model.Pagination) ([]model.Teacher, error) {
	datas, err := s.teacherRepo.ListTeacher(teacher, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s teacherService) ListTeacherMeta(teacher model.Teacher, pagination model.Pagination) (model.Meta, error) {
	data, err := s.teacherRepo.ListTeacherMeta(teacher, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s teacherService) DropDownTeacher(teacher model.Teacher) ([]model.Teacher, error) {
	datas, err := s.teacherRepo.DropDownTeacher(teacher)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s teacherService) CheckIsExist(id int) (isExist bool) {
	return s.teacherRepo.CheckIsExist(id)
}

func (s teacherService) CheckIsExistByNip(nip string, exceptID int) (isExist bool) {
	return s.teacherRepo.CheckIsExistByNip(nip, exceptID)
}
