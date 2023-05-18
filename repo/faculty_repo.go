package repo

import (
	"attendance-api/model"
	"strings"

	"gorm.io/gorm"
)

type FacultyRepo interface {
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

type facultyRepo struct {
	db *gorm.DB
}

func NewFacultyRepo(db *gorm.DB) FacultyRepo {
	return &facultyRepo{db: db}
}

func (r facultyRepo) CreateFaculty(faculty model.Faculty) (model.Faculty, error) {
	if err := r.db.Table("faculties").Create(&faculty).Error; err != nil {
		return model.Faculty{}, err
	}

	return faculty, nil
}

func (r facultyRepo) RetrieveFaculty(id int) (model.Faculty, error) {
	var faculty model.Faculty
	if err := r.db.First(&faculty, id).Error; err != nil {
		return model.Faculty{}, err
	}
	return faculty, nil
}

func (r facultyRepo) RetrieveFacultyByOwner(id int, ownerID int) (model.Faculty, error) {
	var faculty model.Faculty
	if err := r.db.Model(&model.Faculty{}).Where("id = ? AND owner_id = ?", id, ownerID).First(&faculty).Error; err != nil {
		return model.Faculty{}, err
	}
	return faculty, nil
}

func (r facultyRepo) UpdateFaculty(id int, faculty model.Faculty) (model.Faculty, error) {
	if err := r.db.Model(&model.Faculty{}).Where("id = ?", id).Updates(&faculty).Error; err != nil {
		return model.Faculty{}, err
	}
	return faculty, nil
}

func (r facultyRepo) UpdateFacultyByOwner(id int, ownerID int, faculty model.Faculty) (model.Faculty, error) {
	if err := r.db.Model(&model.Faculty{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(&faculty).Error; err != nil {
		return model.Faculty{}, err
	}
	return faculty, nil
}

func (r facultyRepo) DeleteFaculty(id int) error {
	if err := r.db.Delete(&model.Faculty{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r facultyRepo) DeleteFacultyByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Faculty{}).Error; err != nil {
		return err
	}
	return nil
}

func (r facultyRepo) ListFaculty(faculty model.Faculty, pagination model.Pagination) ([]model.Faculty, error) {
	var faculties []model.Faculty
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("faculties").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterFaculty(query, faculty)
	query = SearchFaculty(query, pagination.Search)
	query = query.Find(&faculties)
	if err := query.Error; err != nil {
		return nil, err
	}

	return faculties, nil
}

func (r facultyRepo) ListFacultyMeta(faculty model.Faculty, pagination model.Pagination) (model.Meta, error) {
	var faculties []model.Faculty
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Faculty{}).Select("count(*)")
	queryTotal = FilterFaculty(queryTotal, faculty)
	queryTotal = SearchFaculty(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("faculties").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterFaculty(query, faculty)
	query = SearchFaculty(query, pagination.Search)
	query = query.Find(&faculties)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(faculties),
	}
	return meta, nil
}

func (r facultyRepo) DropDownFaculty(faculty model.Faculty) ([]model.Faculty, error) {
	var faculties []model.Faculty
	query := r.db.Table("faculties").Order("id desc")
	query = FilterFaculty(query, faculty)
	query = query.Find(&faculties)
	if err := query.Error; err != nil {
		return nil, err
	}
	return faculties, nil
}

func (r facultyRepo) CheckIsExist(id int) (isExist bool) {
	if err := r.db.Table("faculties").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r facultyRepo) CheckIsExistByName(name string, exceptID int) (isExist bool) {
	if err := r.db.Table("faculties").Select("count(*) > 0").Where("LOWER(name) = ? AND id != ?", strings.ToLower(name), exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r facultyRepo) CheckIsExistByCode(code string, exceptID int) (isExist bool) {
	if err := r.db.Table("faculties").Select("count(*) > 0").Where("code = ? AND id != ?", code, exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func FilterFaculty(query *gorm.DB, faculty model.Faculty) *gorm.DB {
	if faculty.Name != "" {
		query = query.Where("name LIKE ?", "%"+faculty.Name+"%")
	}
	if faculty.Code != "" {
		query = query.Where("code LIKE ?", "%"+faculty.Code+"%")
	}
	if faculty.OwnerID > 0 {
		query = query.Where("owner_id = ?", faculty.OwnerID)
	}
	return query
}

func SearchFaculty(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR summary LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}
