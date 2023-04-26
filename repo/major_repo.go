package repo

import (
	"attendance-api/model"

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
	CheckIsExist(id int) (isExist bool)
}

type majorRepo struct {
	db *gorm.DB
}

func NewMajorRepo(db *gorm.DB) MajorRepo {
	return &majorRepo{db: db}
}

func (r majorRepo) CreateMajor(major model.Major) (model.Major, error) {
	if err := r.db.Table("majors").Create(&major).Error; err != nil {
		return model.Major{}, err
	}

	return major, nil
}

func (r majorRepo) RetrieveMajor(id int) (model.Major, error) {
	var major model.Major
	if err := r.db.First(&major, id).Error; err != nil {
		return model.Major{}, err
	}
	return major, nil
}

func (r majorRepo) RetrieveMajorByOwner(id int, ownerID int) (model.Major, error) {
	var major model.Major
	if err := r.db.Model(&model.Major{}).Where("id = ? AND owner_id = ?", id, ownerID).First(&major).Error; err != nil {
		return model.Major{}, err
	}
	return major, nil
}

func (r majorRepo) UpdateMajor(id int, major model.Major) (model.Major, error) {
	if err := r.db.Model(&model.Major{}).Where("id = ?", id).Updates(&major).Error; err != nil {
		return model.Major{}, err
	}
	return major, nil
}

func (r majorRepo) UpdateMajorByOwner(id int, ownerID int, major model.Major) (model.Major, error) {
	if err := r.db.Model(&model.Major{}).Where("id = ? AND owner_id = ?", id, ownerID).Updates(&major).Error; err != nil {
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
	query := r.db.Table("majors")
	query = FilterMajor(query, major)
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
