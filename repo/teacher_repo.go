package repo

import (
	"attendance-api/model"

	"gorm.io/gorm"
)

type TeacherRepo interface {
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
	CheckIsExistByNIP(nip string, exceptID int) (isExist bool)
}

type teacherRepo struct {
	db *gorm.DB
}

func NewTeacherRepo(db *gorm.DB) TeacherRepo {
	return &teacherRepo{db: db}
}

func (r teacherRepo) CreateTeacher(teacher model.Teacher) (model.Teacher, error) {
	query := r.db.Table("teachers")
	query = PreloadTeacher(query)
	if err := query.Create(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}

	query2 := r.db.Table("teachers")
	query2 = PreloadTeacher(query2)
	if err := query2.Where("id = ?", teacher.ID).First(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}

	return teacher, nil
}

func (r teacherRepo) RetrieveTeacher(id int) (model.Teacher, error) {
	var teacher model.Teacher
	query := r.db.Table("teachers")
	query = PreloadTeacher(query)
	if err := query.Where("id = ?", id).First(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}
	return teacher, nil
}

func (r teacherRepo) RetrieveTeacherByUserID(userID int) (model.Teacher, error) {
	var teacher model.Teacher
	query := r.db.Table("teachers")
	query = PreloadTeacher(query)
	if err := query.Where("user_id = ?", userID).First(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}
	return teacher, nil
}

func (r teacherRepo) RetrieveTeacherByOwner(id int, ownerID int) (model.Teacher, error) {
	var teacher model.Teacher
	query := r.db.Table("teachers")
	query = PreloadTeacher(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).First(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}
	return teacher, nil
}

func (r teacherRepo) UpdateTeacher(id int, teacher model.Teacher) (model.Teacher, error) {
	query := r.db.Table("teachers")
	query = PreloadTeacher(query)
	if err := query.Where("id = ?", id).Updates(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}

	query2 := r.db.Table("teachers")
	query2 = PreloadTeacher(query2)
	if err := query2.Where("id = ?", id).First(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}
	return teacher, nil
}

func (r teacherRepo) UpdateTeacherByOwner(id int, ownerID int, teacher model.Teacher) (model.Teacher, error) {
	query := r.db.Table("teachers")
	query = PreloadTeacher(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).Updates(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}

	query2 := r.db.Table("teachers")
	query2 = PreloadTeacher(query2)
	if err := query2.Where("id = ? AND owner_id = ?", id, ownerID).First(&teacher).Error; err != nil {
		return model.Teacher{}, err
	}
	return teacher, nil
}

func (r teacherRepo) DeleteTeacher(id int) error {
	if err := r.db.Delete(&model.Teacher{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r teacherRepo) DeleteTeacherByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Teacher{}).Error; err != nil {
		return err
	}
	return nil
}

func (r teacherRepo) ListTeacher(teacher model.Teacher, pagination model.Pagination) ([]model.Teacher, error) {
	var teachers []model.Teacher
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("teachers").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = PreloadTeacher(query)
	query = FilterTeacher(query, teacher)
	query = SearchTeacher(query, pagination.Search)
	query = query.Find(&teachers)
	if err := query.Error; err != nil {
		return nil, err
	}

	return teachers, nil
}

func (r teacherRepo) ListTeacherMeta(teacher model.Teacher, pagination model.Pagination) (model.Meta, error) {
	var teachers []model.Teacher
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Teacher{}).Select("count(*)")
	queryTotal = FilterTeacher(queryTotal, teacher)
	queryTotal = SearchTeacher(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("teachers").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterTeacher(query, teacher)
	query = SearchTeacher(query, pagination.Search)
	query = query.Find(&teachers)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(teachers),
	}
	return meta, nil
}

func (r teacherRepo) DropDownTeacher(teacher model.Teacher) ([]model.Teacher, error) {
	var teachers []model.Teacher
	query := r.db.Table("teachers").Order("id desc")
	query = PreloadTeacher(query)
	query = FilterTeacher(query, teacher)
	query = query.Find(&teachers)
	if err := query.Error; err != nil {
		return nil, err
	}
	return teachers, nil
}

func (r teacherRepo) CheckIsExist(id int) (isExist bool) {
	if err := r.db.Table("teachers").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r teacherRepo) CheckIsExistByNIP(nip string, exceptID int) (isExist bool) {
	if err := r.db.Table("teachers").Select("count(*) > 0").Where("nip = ? AND id != ?", nip, exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func FilterTeacher(query *gorm.DB, teacher model.Teacher) *gorm.DB {
	if teacher.NIP != "" {
		query = query.Where("nip LIKE ?", "%"+teacher.NIP+"%")
	}
	if teacher.Gender != "" {
		query = query.Where("gender LIKE ?", "%"+teacher.Gender+"%")
	}
	if teacher.DOB != "" {
		query = query.Where("DATE(dob) = ?", teacher.DOB)
	}
	if teacher.FacultyID > 0 {
		query = query.Where("faculty_id = ?", teacher.FacultyID)
	}
	if teacher.MajorID > 0 {
		query = query.Where("major_id = ?", teacher.MajorID)
	}
	if teacher.StudyProgramID > 0 {
		query = query.Where("study_program_id = ?", teacher.StudyProgramID)
	}
	return query
}

func SearchTeacher(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("nim LIKE ? OR gender LIKE ? OR DATE(dob) LIKE ? OR address LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

func PreloadTeacher(query *gorm.DB) *gorm.DB {
	query = query.Preload("User")
	query = query.Preload("Faculty")
	query = query.Preload("Major")
	query = query.Preload("Major.Faculty")
	query = query.Preload("StudyProgram")
	query = query.Preload("StudyProgram.Major")
	query = query.Preload("StudyProgram.Major.Faculty")
	return query
}
