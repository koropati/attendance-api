package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type StudentService interface {
	CreateStudent(student model.Student) (model.Student, error)
	RetrieveStudent(id int) (model.Student, error)
	RetrieveStudentByOwner(id int, ownerID int) (model.Student, error)
	UpdateStudent(id int, student model.Student) (model.Student, error)
	UpdateStudentByOwner(id int, ownerID int, student model.Student) (model.Student, error)
	DeleteStudent(id int) error
	DeleteStudentByOwner(id int, ownerID int) error
	ListStudent(student model.Student, pagination model.Pagination) ([]model.Student, error)
	ListStudentMeta(student model.Student, pagination model.Pagination) (model.Meta, error)
	DropDownStudent(student model.Student) ([]model.Student, error)
	CheckIsExist(id int) (isExist bool)
	CheckIsExistByNIM(nim string, exceptID int) (isExist bool)
}

type studentService struct {
	studentRepo repo.StudentRepo
}

func NewStudentService(studentRepo repo.StudentRepo) StudentService {
	return &studentService{studentRepo: studentRepo}
}

func (s studentService) CreateStudent(student model.Student) (model.Student, error) {
	data, err := s.studentRepo.CreateStudent(student)
	if err != nil {
		return model.Student{}, err
	}
	return data, nil
}

func (s studentService) RetrieveStudent(id int) (model.Student, error) {
	data, err := s.studentRepo.RetrieveStudent(id)
	if err != nil {
		return model.Student{}, err
	}
	return data, nil
}

func (s studentService) RetrieveStudentByOwner(id int, ownerID int) (model.Student, error) {
	data, err := s.studentRepo.RetrieveStudentByOwner(id, ownerID)
	if err != nil {
		return model.Student{}, err
	}
	return data, nil
}

func (s studentService) UpdateStudent(id int, student model.Student) (model.Student, error) {
	data, err := s.studentRepo.UpdateStudent(id, student)
	if err != nil {
		return model.Student{}, err
	}
	return data, nil
}

func (s studentService) UpdateStudentByOwner(id int, ownerID int, student model.Student) (model.Student, error) {
	data, err := s.studentRepo.UpdateStudentByOwner(id, ownerID, student)
	if err != nil {
		return model.Student{}, err
	}
	return data, nil
}

func (s studentService) DeleteStudent(id int) error {
	if err := s.studentRepo.DeleteStudent(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s studentService) DeleteStudentByOwner(id int, ownerID int) error {
	if err := s.studentRepo.DeleteStudentByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s studentService) ListStudent(student model.Student, pagination model.Pagination) ([]model.Student, error) {
	datas, err := s.studentRepo.ListStudent(student, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s studentService) ListStudentMeta(student model.Student, pagination model.Pagination) (model.Meta, error) {
	data, err := s.studentRepo.ListStudentMeta(student, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s studentService) DropDownStudent(student model.Student) ([]model.Student, error) {
	datas, err := s.studentRepo.DropDownStudent(student)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s studentService) CheckIsExist(id int) (isExist bool) {
	return s.studentRepo.CheckIsExist(id)
}

func (s studentService) CheckIsExistByNIM(nim string, exceptID int) (isExist bool) {
	return s.studentRepo.CheckIsExistByNIM(nim, exceptID)
}
