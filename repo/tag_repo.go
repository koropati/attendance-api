package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type TagRepo interface {
	CreateTag(tag *model.Tag) (*model.Tag, error)
	RetrieveTag(id int) (*model.Tag, error)
	UpdateTag(id int, tag *model.Tag) (*model.Tag, error)
	DeleteTag(id int) error
	ListTag(tag *model.Tag, pagination *model.Pagination) (*[]model.Tag, error)
	ListTagMeta(tag *model.Tag, pagination *model.Pagination) (*model.Meta, error)
	DropDownTag(tag *model.Tag) (*[]model.Tag, error)
}

type tagRepo struct {
	db *gorm.DB
}

func NewTagRepo(db *gorm.DB) TagRepo {
	return &tagRepo{db: db}
}

func (r *tagRepo) CreateTag(tag *model.Tag) (*model.Tag, error) {
	if err := r.db.Table("tags").Create(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *tagRepo) RetrieveTag(id int) (*model.Tag, error) {
	var tag model.Tag
	if err := r.db.First(&tag, id).Error; err != nil {
		return nil, err
	}
	return &tag, nil
}

func (r *tagRepo) UpdateTag(id int, tag *model.Tag) (*model.Tag, error) {
	if err := r.db.Model(&model.Tag{}).Where("id = ?", id).Updates(&tag).Error; err != nil {
		return nil, err
	}
	return tag, nil
}

func (r *tagRepo) DeleteTag(id int) error {
	if err := r.db.Delete(&model.Tag{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r *tagRepo) ListTag(tag *model.Tag, pagination *model.Pagination) (*[]model.Tag, error) {
	var tags []model.Tag
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("tags").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterTag(query, tag)
	query = query.Find(&tags)
	if err := query.Error; err != nil {
		return nil, err
	}

	return &tags, nil
}

func (r *tagRepo) ListTagMeta(tag *model.Tag, pagination *model.Pagination) (*model.Meta, error) {
	var tags []model.Tag
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Tag{}).Select("count(*)")
	queryTotal = FilterTag(queryTotal, tag)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return nil, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(tags),
	}
	return &meta, nil
}

func (r *tagRepo) DropDownTag(tag *model.Tag) (*[]model.Tag, error) {
	var tags []model.Tag
	query := r.db.Table("tags")
	query = FilterTag(query, tag)
	query = query.Find(&tags)
	if err := query.Error; err != nil {
		return nil, err
	}
	return &tags, nil
}

func FilterTag(query *gorm.DB, tag *model.Tag) *gorm.DB {
	if tag.Title != "" {
		query = query.Where("title LIKE ?", "%"+tag.Title+"%")
	}
	if tag.Slug != "" {
		query = query.Where("slug LIKE ?", "%"+tag.Slug+"%")
	}
	if tag.MetaTitle != "" {
		query = query.Where("meta_title LIKE ?", "%"+tag.MetaTitle+"%")
	}
	if tag.Content != "" {
		query = query.Where("content LIKE ?", "%"+tag.Content+"%")
	}
	return query
}
