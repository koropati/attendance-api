package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type RoleAbilityRepo interface {
	CreateRoleAbility(roleAbility model.RoleAbility) (model.RoleAbility, error)
	RetrieveRoleAbility(id int) (model.RoleAbility, error)
	RetrieveRoleAbilityByRole(isSuperAdmin bool, isAdmin bool, isUser bool) ([]model.Ability, error)
	UpdateRoleAbility(id int, roleAbility model.RoleAbility) (model.RoleAbility, error)
	DeleteRoleAbility(id int) error
	ListRoleAbility(roleAbility model.RoleAbility, pagination model.Pagination) ([]model.RoleAbility, error)
	ListRoleAbilityMeta(roleAbility model.RoleAbility, pagination model.Pagination) (model.Meta, error)
	DropDownRoleAbility(roleAbility model.RoleAbility) ([]model.RoleAbility, error)
	CheckIsExist(id int) (isExist bool)
	CheckIsExistByActionAndSubject(action string, subject string, exceptID int) (isExist bool)
}

type roleAbilityRepo struct {
	db *gorm.DB
}

func NewRoleAbilityRepo(db *gorm.DB) RoleAbilityRepo {
	return &roleAbilityRepo{db: db}
}

func (r roleAbilityRepo) CreateRoleAbility(roleAbility model.RoleAbility) (model.RoleAbility, error) {
	if err := r.db.Table("role_abilities").Create(&roleAbility).Error; err != nil {
		return model.RoleAbility{}, err
	}

	return roleAbility, nil
}

func (r roleAbilityRepo) RetrieveRoleAbility(id int) (model.RoleAbility, error) {
	var roleAbility model.RoleAbility
	if err := r.db.First(&roleAbility, id).Error; err != nil {
		return model.RoleAbility{}, err
	}
	return roleAbility, nil
}

func (r roleAbilityRepo) RetrieveRoleAbilityByRole(isSuperAdmin bool, isAdmin bool, isUser bool) ([]model.Ability, error) {
	var results []model.Ability

	if isSuperAdmin {
		if err := r.db.Select("action, subject").Table("role_abilities").Where("is_super_admin = ?", isSuperAdmin).Find(&results).Error; err != nil {
			return nil, err
		}
	} else if isAdmin {
		if err := r.db.Select("action, subject").Table("role_abilities").Where("is_admin = ?", isAdmin).Find(&results).Error; err != nil {
			return nil, err
		}
	} else if isUser {
		if err := r.db.Select("action, subject").Table("role_abilities").Where("is_user = ?", isUser).Find(&results).Error; err != nil {
			return nil, err
		}
	} else {
		return nil, nil
	}

	return results, nil
}

func (r roleAbilityRepo) UpdateRoleAbility(id int, roleAbility model.RoleAbility) (model.RoleAbility, error) {
	if err := r.db.Model(&model.RoleAbility{}).Where("id = ?", id).Updates(&roleAbility).Error; err != nil {
		return model.RoleAbility{}, err
	}
	return roleAbility, nil
}

func (r roleAbilityRepo) DeleteRoleAbility(id int) error {
	if err := r.db.Delete(&model.RoleAbility{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r roleAbilityRepo) ListRoleAbility(roleAbility model.RoleAbility, pagination model.Pagination) ([]model.RoleAbility, error) {
	var role_abilities []model.RoleAbility
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("role_abilities").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterRoleAbility(query, roleAbility)
	query = SearchRoleAbility(query, pagination.Search)
	query = query.Find(&role_abilities)
	if err := query.Error; err != nil {
		return nil, err
	}

	return role_abilities, nil
}

func (r roleAbilityRepo) ListRoleAbilityMeta(roleAbility model.RoleAbility, pagination model.Pagination) (model.Meta, error) {
	var role_abilities []model.RoleAbility
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.RoleAbility{}).Select("count(*)")
	queryTotal = FilterRoleAbility(queryTotal, roleAbility)
	queryTotal = SearchRoleAbility(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("role_abilities").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterRoleAbility(query, roleAbility)
	query = SearchRoleAbility(query, pagination.Search)
	query = query.Find(&role_abilities)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(role_abilities),
	}
	return meta, nil
}

func (r roleAbilityRepo) DropDownRoleAbility(roleAbility model.RoleAbility) ([]model.RoleAbility, error) {
	var role_abilities []model.RoleAbility
	query := r.db.Table("role_abilities").Order("id desc")
	query = FilterRoleAbility(query, roleAbility)
	query = query.Find(&role_abilities)
	if err := query.Error; err != nil {
		return nil, err
	}
	return role_abilities, nil
}

func (r roleAbilityRepo) CheckIsExist(id int) (isExist bool) {
	if err := r.db.Table("role_abilities").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r roleAbilityRepo) CheckIsExistByActionAndSubject(action string, subject string, exceptID int) (isExist bool) {
	if err := r.db.Table("role_abilities").Select("count(*) > 0").Where("action = ? AND subject = ? AND id != ?", action, subject, exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func FilterRoleAbility(query *gorm.DB, roleAbility model.RoleAbility) *gorm.DB {
	if roleAbility.Action != "" {
		query = query.Where("action LIKE ?", "%"+roleAbility.Action+"%")
	}
	if roleAbility.Subject != "" {
		query = query.Where("subject LIKE ?", "%"+roleAbility.Subject+"%")
	}
	return query
}

func SearchRoleAbility(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR summary LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}
