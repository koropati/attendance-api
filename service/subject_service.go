package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type SubjectService interface {
	CreateSubject(subject *model.Subject) (*model.Subject, error)
	RetrieveSubject(id int) (*model.Subject, error)
	RetrieveSubjectByOwner(id int, ownerID int) (*model.Subject, error)
	UpdateSubject(id int, subject *model.Subject) (*model.Subject, error)
	UpdateSubjectByOwner(id int, ownerID int, subject *model.Subject) (*model.Subject, error)
	DeleteSubject(id int) error
	DeleteSubjectByOwner(id int, ownerID int) error
	ListSubject(subject *model.Subject, pagination *model.Pagination) (*[]model.Subject, error)
	ListSubjectMeta(subject *model.Subject, pagination *model.Pagination) (*model.Meta, error)
	DropDownSubject(subject *model.Subject) (*[]model.Subject, error)
	CheckIsExist(id int) (isExist bool)
}

type subjectService struct {
	subjectRepo repo.SubjectRepo
}

func NewSubjectService(subjectRepo repo.SubjectRepo) SubjectService {
	return &subjectService{subjectRepo: subjectRepo}
}

func (s *subjectService) CreateSubject(subject *model.Subject) (*model.Subject, error) {
	data, err := s.subjectRepo.CreateSubject(subject)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *subjectService) RetrieveSubject(id int) (*model.Subject, error) {
	data, err := s.subjectRepo.RetrieveSubject(id)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *subjectService) RetrieveSubjectByOwner(id int, ownerID int) (*model.Subject, error) {
	data, err := s.subjectRepo.RetrieveSubjectByOwner(id, ownerID)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *subjectService) UpdateSubject(id int, subject *model.Subject) (*model.Subject, error) {
	data, err := s.subjectRepo.UpdateSubject(id, subject)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *subjectService) UpdateSubjectByOwner(id int, ownerID int, subject *model.Subject) (*model.Subject, error) {
	data, err := s.subjectRepo.UpdateSubjectByOwner(id, ownerID, subject)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *subjectService) DeleteSubject(id int) error {
	if err := s.subjectRepo.DeleteSubject(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *subjectService) DeleteSubjectByOwner(id int, ownerID int) error {
	if err := s.subjectRepo.DeleteSubjectByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s *subjectService) ListSubject(subject *model.Subject, pagination *model.Pagination) (*[]model.Subject, error) {
	datas, err := s.subjectRepo.ListSubject(subject, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *subjectService) ListSubjectMeta(subject *model.Subject, pagination *model.Pagination) (*model.Meta, error) {
	data, err := s.subjectRepo.ListSubjectMeta(subject, pagination)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *subjectService) DropDownSubject(subject *model.Subject) (*[]model.Subject, error) {
	datas, err := s.subjectRepo.DropDownSubject(subject)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s *subjectService) CheckIsExist(id int) (isExist bool) {
	return s.subjectRepo.CheckIsExist(id)
}
