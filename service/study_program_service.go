package service

import (
	"attendance-api/model"
	"attendance-api/repo"
)

type StudyProgramService interface {
	CreateStudyProgram(studyProgram model.StudyProgram) (model.StudyProgram, error)
	RetrieveStudyProgram(id int) (model.StudyProgram, error)
	RetrieveStudyProgramByOwner(id int, ownerID int) (model.StudyProgram, error)
	UpdateStudyProgram(id int, studyProgram model.StudyProgram) (model.StudyProgram, error)
	UpdateStudyProgramByOwner(id int, ownerID int, studyProgram model.StudyProgram) (model.StudyProgram, error)
	DeleteStudyProgram(id int) error
	DeleteStudyProgramByOwner(id int, ownerID int) error
	ListStudyProgram(studyProgram model.StudyProgram, pagination model.Pagination) ([]model.StudyProgram, error)
	ListStudyProgramMeta(studyProgram model.StudyProgram, pagination model.Pagination) (model.Meta, error)
	DropDownStudyProgram(studyProgram model.StudyProgram) ([]model.StudyProgram, error)
	CheckIsExist(id int) (isExist bool)
	CheckIsExistByName(name string, majorID int, exceptID int) (isExist bool)
	CheckIsExistByCode(code string, exceptID int) (isExist bool)
}

type studyProgramService struct {
	studyProgramRepo repo.StudyProgramRepo
}

func NewStudyProgramService(studyProgramRepo repo.StudyProgramRepo) StudyProgramService {
	return &studyProgramService{studyProgramRepo: studyProgramRepo}
}

func (s studyProgramService) CreateStudyProgram(studyProgram model.StudyProgram) (model.StudyProgram, error) {
	data, err := s.studyProgramRepo.CreateStudyProgram(studyProgram)
	if err != nil {
		return model.StudyProgram{}, err
	}
	return data, nil
}

func (s studyProgramService) RetrieveStudyProgram(id int) (model.StudyProgram, error) {
	data, err := s.studyProgramRepo.RetrieveStudyProgram(id)
	if err != nil {
		return model.StudyProgram{}, err
	}
	return data, nil
}

func (s studyProgramService) RetrieveStudyProgramByOwner(id int, ownerID int) (model.StudyProgram, error) {
	data, err := s.studyProgramRepo.RetrieveStudyProgramByOwner(id, ownerID)
	if err != nil {
		return model.StudyProgram{}, err
	}
	return data, nil
}

func (s studyProgramService) UpdateStudyProgram(id int, studyProgram model.StudyProgram) (model.StudyProgram, error) {
	data, err := s.studyProgramRepo.UpdateStudyProgram(id, studyProgram)
	if err != nil {
		return model.StudyProgram{}, err
	}
	return data, nil
}

func (s studyProgramService) UpdateStudyProgramByOwner(id int, ownerID int, studyProgram model.StudyProgram) (model.StudyProgram, error) {
	data, err := s.studyProgramRepo.UpdateStudyProgramByOwner(id, ownerID, studyProgram)
	if err != nil {
		return model.StudyProgram{}, err
	}
	return data, nil
}

func (s studyProgramService) DeleteStudyProgram(id int) error {
	if err := s.studyProgramRepo.DeleteStudyProgram(id); err != nil {
		return err
	} else {
		return nil
	}
}

func (s studyProgramService) DeleteStudyProgramByOwner(id int, ownerID int) error {
	if err := s.studyProgramRepo.DeleteStudyProgramByOwner(id, ownerID); err != nil {
		return err
	} else {
		return nil
	}
}

func (s studyProgramService) ListStudyProgram(studyProgram model.StudyProgram, pagination model.Pagination) ([]model.StudyProgram, error) {
	datas, err := s.studyProgramRepo.ListStudyProgram(studyProgram, pagination)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s studyProgramService) ListStudyProgramMeta(studyProgram model.StudyProgram, pagination model.Pagination) (model.Meta, error) {
	data, err := s.studyProgramRepo.ListStudyProgramMeta(studyProgram, pagination)
	if err != nil {
		return model.Meta{}, err
	}
	return data, nil
}

func (s studyProgramService) DropDownStudyProgram(studyProgram model.StudyProgram) ([]model.StudyProgram, error) {
	datas, err := s.studyProgramRepo.DropDownStudyProgram(studyProgram)
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (s studyProgramService) CheckIsExist(id int) (isExist bool) {
	return s.studyProgramRepo.CheckIsExist(id)
}

func (s studyProgramService) CheckIsExistByName(name string, majorID int, exceptID int) (isExist bool) {
	return s.studyProgramRepo.CheckIsExistByName(name, majorID, exceptID)
}

func (s studyProgramService) CheckIsExistByCode(code string, exceptID int) (isExist bool) {
	return s.studyProgramRepo.CheckIsExistByCode(code, exceptID)
}
