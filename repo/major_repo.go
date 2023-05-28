package repo

import (
	"attendance-api/model"
	"strings"

	"gorm.io/gorm"
)

type MajorRepo interface {
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
	DropDownByFaculty(facultyID int) ([]model.Major, error)
	CheckIsExist(id int) (isExist bool)
	CheckIsExistByName(name string, exceptID int) (isExist bool)
	CheckIsExistByCode(code string, exceptID int) (isExist bool)
}

type majorRepo struct {
	db *gorm.DB
}

func NewMajorRepo(db *gorm.DB) MajorRepo {
	return &majorRepo{db: db}
}

func (r majorRepo) CreateMajor(major model.Major) (model.Major, error) {
	query := r.db.Table("majors")
	query = PreloadMajor(query)
	if err := query.Create(&major).Error; err != nil {
		return model.Major{}, err
	}

	query2 := r.db.Table("majors")
	query2 = PreloadMajor(query2)
	if err := query2.Where("id = ?", major.ID).First(&major).Error; err != nil {
		return model.Major{}, err
	}

	return major, nil
}

func (r majorRepo) RetrieveMajor(id int) (model.Major, error) {
	var major model.Major
	query := r.db.Table("majors")
	query = PreloadMajor(query)
	if err := query.Where("id = ?", id).First(&major).Error; err != nil {
		return model.Major{}, err
	}
	return major, nil
}

func (r majorRepo) RetrieveMajorByOwner(id int, ownerID int) (model.Major, error) {
	var major model.Major
	query := r.db.Table("majors")
	query = PreloadMajor(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).First(&major).Error; err != nil {
		return model.Major{}, err
	}
	return major, nil
}

func (r majorRepo) UpdateMajor(id int, major model.Major) (model.Major, error) {
	query := r.db.Table("majors")
	query = PreloadMajor(query)
	if err := query.Where("id = ?", id).Updates(&major).Error; err != nil {
		return model.Major{}, err
	}

	query2 := r.db.Table("majors")
	query2 = PreloadMajor(query2)
	if err := query2.Where("id = ?", id).First(&major).Error; err != nil {
		return model.Major{}, err
	}
	return major, nil
}

func (r majorRepo) UpdateMajorByOwner(id int, ownerID int, major model.Major) (model.Major, error) {
	query := r.db.Table("majors")
	query = PreloadMajor(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).Updates(&major).Error; err != nil {
		return model.Major{}, err
	}

	query2 := r.db.Table("majors")
	query2 = PreloadMajor(query2)
	if err := query2.Where("id = ? AND owner_id = ?", id, ownerID).First(&major).Error; err != nil {
		return model.Major{}, err
	}
	return major, nil
}

func (r majorRepo) DeleteMajor(id int) error {
	if err := r.db.Delete(&model.Major{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r majorRepo) DeleteMajorByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Major{}).Error; err != nil {
		return err
	}
	return nil
}

func (r majorRepo) ListMajor(major model.Major, pagination model.Pagination) ([]model.Major, error) {
	var majors []model.Major
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("majors").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = PreloadMajor(query)
	query = FilterMajor(query, major)
	query = SearchMajor(query, pagination.Search)
	query = query.Find(&majors)
	if err := query.Error; err != nil {
		return nil, err
	}

	return majors, nil
}

func (r majorRepo) ListMajorMeta(major model.Major, pagination model.Pagination) (model.Meta, error) {
	var majors []model.Major
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Major{}).Select("count(*)")
	queryTotal = FilterMajor(queryTotal, major)
	queryTotal = SearchMajor(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("majors").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterMajor(query, major)
	query = SearchMajor(query, pagination.Search)
	query = query.Find(&majors)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(majors),
	}
	return meta, nil
}

func (r majorRepo) DropDownMajor(major model.Major) ([]model.Major, error) {
	var majors []model.Major
	query := r.db.Table("majors").Order("id desc")
	query = PreloadMajor(query)
	query = FilterMajor(query, major)
	query = query.Find(&majors)
	if err := query.Error; err != nil {
		return nil, err
	}
	return majors, nil
}

func (r majorRepo) DropDownByFaculty(facultyID int) ([]model.Major, error) {
	var majors []model.Major
	query := r.db.Table("majors").Where("faculty_id = ?", facultyID).Order("id desc")
	query = PreloadMajor(query)
	query = query.Find(&majors)
	if err := query.Error; err != nil {
		return nil, err
	}
	return majors, nil
}

func (r majorRepo) CheckIsExist(id int) (isExist bool) {
	if err := r.db.Table("majors").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r majorRepo) CheckIsExistByName(name string, exceptID int) (isExist bool) {
	if err := r.db.Table("majors").Select("count(*) > 0").Where("LOWER(name) = ? AND id != ?", strings.ToLower(name), exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r majorRepo) CheckIsExistByCode(code string, exceptID int) (isExist bool) {
	if err := r.db.Table("majors").Select("count(*) > 0").Where("code = ? AND id != ?", code, exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func FilterMajor(query *gorm.DB, major model.Major) *gorm.DB {
	if major.Name != "" {
		query = query.Where("name LIKE ?", "%"+major.Name+"%")
	}
	if major.Code != "" {
		query = query.Where("code LIKE ?", "%"+major.Code+"%")
	}
	if major.OwnerID > 0 {
		query = query.Where("owner_id = ?", major.OwnerID)
	}
	return query
}

func SearchMajor(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR summary LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

func PreloadMajor(query *gorm.DB) *gorm.DB {
	query = query.Preload("Faculty")
	return query
}
