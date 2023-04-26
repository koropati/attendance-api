package repo

import (
	"attendance-api/model"
	"strings"

	"gorm.io/gorm"
)

type StudyProgramRepo interface {
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

type studyProgramRepo struct {
	db *gorm.DB
}

func NewStudyProgramRepo(db *gorm.DB) StudyProgramRepo {
	return &studyProgramRepo{db: db}
}

func (r studyProgramRepo) CreateStudyProgram(studyProgram model.StudyProgram) (model.StudyProgram, error) {
	query := r.db.Table("study_programs")
	query = PreloadStudyProgram(query)
	if err := query.Create(&studyProgram).Error; err != nil {
		return model.StudyProgram{}, err
	}

	query2 := r.db.Table("study_programs")
	query2 = PreloadStudyProgram(query2)
	if err := query2.Where("id = ?", studyProgram.ID).First(&studyProgram).Error; err != nil {
		return model.StudyProgram{}, err
	}

	return studyProgram, nil
}

func (r studyProgramRepo) RetrieveStudyProgram(id int) (model.StudyProgram, error) {
	var studyProgram model.StudyProgram
	query := r.db.Table("study_programs")
	query = PreloadStudyProgram(query)
	if err := query.Where("id = ?", id).First(&studyProgram).Error; err != nil {
		return model.StudyProgram{}, err
	}
	return studyProgram, nil
}

func (r studyProgramRepo) RetrieveStudyProgramByOwner(id int, ownerID int) (model.StudyProgram, error) {
	var studyProgram model.StudyProgram
	query := r.db.Table("study_programs")
	query = PreloadStudyProgram(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).First(&studyProgram).Error; err != nil {
		return model.StudyProgram{}, err
	}
	return studyProgram, nil
}

func (r studyProgramRepo) UpdateStudyProgram(id int, studyProgram model.StudyProgram) (model.StudyProgram, error) {
	query := r.db.Table("study_programs")
	query = PreloadStudyProgram(query)
	if err := query.Where("id = ?", id).Updates(&studyProgram).Error; err != nil {
		return model.StudyProgram{}, err
	}

	query2 := r.db.Table("study_programs")
	query2 = PreloadStudyProgram(query2)
	if err := query2.Where("id = ?", id).First(&studyProgram).Error; err != nil {
		return model.StudyProgram{}, err
	}
	return studyProgram, nil
}

func (r studyProgramRepo) UpdateStudyProgramByOwner(id int, ownerID int, studyProgram model.StudyProgram) (model.StudyProgram, error) {
	query := r.db.Table("study_programs")
	query = PreloadStudyProgram(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).Updates(&studyProgram).Error; err != nil {
		return model.StudyProgram{}, err
	}
	return studyProgram, nil
}

func (r studyProgramRepo) DeleteStudyProgram(id int) error {
	if err := r.db.Delete(&model.StudyProgram{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r studyProgramRepo) DeleteStudyProgramByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.StudyProgram{}).Error; err != nil {
		return err
	}
	return nil
}

func (r studyProgramRepo) ListStudyProgram(studyProgram model.StudyProgram, pagination model.Pagination) ([]model.StudyProgram, error) {
	var studyPrograms []model.StudyProgram
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("study_programs").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = PreloadStudyProgram(query)
	query = FilterStudyProgram(query, studyProgram)
	query = SearchStudyProgram(query, pagination.Search)
	query = query.Find(&studyPrograms)
	if err := query.Error; err != nil {
		return nil, err
	}

	return studyPrograms, nil
}

func (r studyProgramRepo) ListStudyProgramMeta(studyProgram model.StudyProgram, pagination model.Pagination) (model.Meta, error) {
	var studyPrograms []model.StudyProgram
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.StudyProgram{}).Select("count(*)")
	queryTotal = FilterStudyProgram(queryTotal, studyProgram)
	queryTotal = SearchStudyProgram(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("study_programs").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterStudyProgram(query, studyProgram)
	query = SearchStudyProgram(query, pagination.Search)
	query = query.Find(&studyPrograms)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(studyPrograms),
	}
	return meta, nil
}

func (r studyProgramRepo) DropDownStudyProgram(studyProgram model.StudyProgram) ([]model.StudyProgram, error) {
	var studyPrograms []model.StudyProgram
	query := r.db.Table("study_programs")
	query = PreloadStudyProgram(query)
	query = FilterStudyProgram(query, studyProgram)
	query = query.Find(&studyPrograms)
	if err := query.Error; err != nil {
		return nil, err
	}
	return studyPrograms, nil
}

func (r studyProgramRepo) CheckIsExist(id int) (isExist bool) {
	if err := r.db.Table("study_programs").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r studyProgramRepo) CheckIsExistByName(name string, majorID int, exceptID int) (isExist bool) {
	if err := r.db.Table("study_programs").Select("count(*) > 0").Where("LOWER(name) = ? AND major_id = ? AND id != ?", strings.ToLower(name), majorID, exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r studyProgramRepo) CheckIsExistByCode(code string, exceptID int) (isExist bool) {
	if err := r.db.Table("study_programs").Select("count(*) > 0").Where("code = ? AND id != ?", code, exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func FilterStudyProgram(query *gorm.DB, studyProgram model.StudyProgram) *gorm.DB {
	if studyProgram.Name != "" {
		query = query.Where("name LIKE ?", "%"+studyProgram.Name+"%")
	}
	if studyProgram.Code != "" {
		query = query.Where("code LIKE ?", "%"+studyProgram.Code+"%")
	}
	if studyProgram.OwnerID > 0 {
		query = query.Where("owner_id = ?", studyProgram.OwnerID)
	}
	return query
}

func SearchStudyProgram(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR summary LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

func PreloadStudyProgram(query *gorm.DB) *gorm.DB {
	query = query.Preload("Major")
	return query
}
