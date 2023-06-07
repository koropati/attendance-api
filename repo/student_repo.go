package repo

import (
	"attendance-api/model"
	"sync"

	"gorm.io/gorm"
)

type StudentRepo interface {
	CreateStudent(student model.Student) (model.Student, error)
	RetrieveStudent(id int) (model.Student, error)
	RetrieveStudentByUserID(userID int) (model.Student, error)
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

type studentRepo struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) StudentRepo {
	return &studentRepo{db: db}
}

func (r studentRepo) CreateStudent(student model.Student) (model.Student, error) {
	query := r.db.Table("students")
	query = PreloadStudent(query)
	if err := query.Create(&student).Error; err != nil {
		return model.Student{}, err
	}

	query2 := r.db.Table("students")
	query2 = PreloadStudent(query2)
	if err := query2.Where("id = ?", student.ID).First(&student).Error; err != nil {
		return model.Student{}, err
	}
	student.Avatar = student.GetAvatar()
	return student, nil
}

func (r studentRepo) RetrieveStudent(id int) (model.Student, error) {
	var student model.Student
	query := r.db.Table("students")
	query = PreloadStudent(query)
	if err := query.Where("id = ?", id).First(&student).Error; err != nil {
		return model.Student{}, err
	}
	student.Avatar = student.GetAvatar()
	return student, nil
}

func (r studentRepo) RetrieveStudentByUserID(userID int) (model.Student, error) {
	var student model.Student
	query := r.db.Table("students")
	query = PreloadStudent(query)
	if err := query.Where("user_id = ?", userID).First(&student).Error; err != nil {
		return model.Student{}, err
	}
	student.Avatar = student.GetAvatar()
	return student, nil
}

func (r studentRepo) RetrieveStudentByOwner(id int, ownerID int) (model.Student, error) {
	var student model.Student
	query := r.db.Table("students")
	query = PreloadStudent(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).First(&student).Error; err != nil {
		return model.Student{}, err
	}
	student.Avatar = student.GetAvatar()
	return student, nil
}

func (r studentRepo) UpdateStudent(id int, student model.Student) (model.Student, error) {
	query := r.db.Table("students")
	query = PreloadStudent(query)
	if err := query.Where("id = ?", id).Updates(&student).Error; err != nil {
		return model.Student{}, err
	}

	query2 := r.db.Table("students")
	query2 = PreloadStudent(query2)
	if err := query2.Where("id = ?", id).First(&student).Error; err != nil {
		return model.Student{}, err
	}
	student.Avatar = student.GetAvatar()
	return student, nil
}

func (r studentRepo) UpdateStudentByOwner(id int, ownerID int, student model.Student) (model.Student, error) {
	query := r.db.Table("students")
	query = PreloadStudent(query)
	if err := query.Where("id = ? AND owner_id = ?", id, ownerID).Updates(&student).Error; err != nil {
		return model.Student{}, err
	}

	query2 := r.db.Table("students")
	query2 = PreloadStudent(query2)
	if err := query2.Where("id = ? AND owner_id = ?", id, ownerID).First(&student).Error; err != nil {
		return model.Student{}, err
	}
	student.Avatar = student.GetAvatar()
	return student, nil
}

func (r studentRepo) DeleteStudent(id int) error {
	if err := r.db.Unscoped().Delete(&model.Student{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r studentRepo) DeleteStudentByOwner(id int, ownerID int) error {
	if err := r.db.Unscoped().Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.Student{}).Error; err != nil {
		return err
	}
	return nil
}

func (r studentRepo) ListStudent(student model.Student, pagination model.Pagination) ([]model.Student, error) {
	var students []model.Student
	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("students").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = PreloadStudent(query)
	query = FilterStudent(query, student)
	query = SearchStudent(query, pagination.Search)
	query = query.Find(&students)
	if err := query.Error; err != nil {
		return nil, err
	}

	wg := sync.WaitGroup{}
	for i, student := range students {
		wg.Add(1)
		go func(i int, student model.Student) {
			students[i].Avatar = student.GetAvatar()
			wg.Done()
		}(i, student)
	}
	wg.Wait()

	return students, nil
}

func (r studentRepo) ListStudentMeta(student model.Student, pagination model.Pagination) (model.Meta, error) {
	var students []model.Student
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.Student{}).Select("count(*)")
	queryTotal = FilterStudent(queryTotal, student)
	queryTotal = SearchStudent(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Table("students").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterStudent(query, student)
	query = SearchStudent(query, pagination.Search)
	query = query.Find(&students)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(students),
	}
	return meta, nil
}

func (r studentRepo) DropDownStudent(student model.Student) ([]model.Student, error) {
	var students []model.Student
	query := r.db.Table("students").Order("id desc")
	query = PreloadStudent(query)
	query = FilterStudent(query, student)
	query = query.Find(&students)
	if err := query.Error; err != nil {
		return nil, err
	}
	return students, nil
}

func (r studentRepo) CheckIsExist(id int) (isExist bool) {
	if err := r.db.Table("students").Select("count(*) > 0").Where("id = ?", id).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func (r studentRepo) CheckIsExistByNIM(nim string, exceptID int) (isExist bool) {
	if err := r.db.Table("students").Select("count(*) > 0").Where("nim = ? AND id != ?", nim, exceptID).Find(&isExist).Error; err != nil {
		return false
	}
	return
}

func FilterStudent(query *gorm.DB, student model.Student) *gorm.DB {
	if student.NIM != "" {
		query = query.Where("nim LIKE ?", "%"+student.NIM+"%")
	}
	if student.Gender != "" {
		query = query.Where("gender LIKE ?", "%"+student.Gender+"%")
	}
	if student.DOB != "" {
		query = query.Where("DATE(dob) = ?", student.DOB)
	}
	if student.FacultyID > 0 {
		query = query.Where("faculty_id = ?", student.FacultyID)
	}
	if student.MajorID > 0 {
		query = query.Where("major_id = ?", student.MajorID)
	}
	if student.StudyProgramID > 0 {
		query = query.Where("study_program_id = ?", student.StudyProgramID)
	}
	return query
}

func SearchStudent(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("nim LIKE ? OR gender LIKE ? OR DATE(dob) LIKE ? OR address LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

func PreloadStudent(query *gorm.DB) *gorm.DB {
	query = query.Preload("User")
	query = query.Preload("Faculty")
	query = query.Preload("Major")
	query = query.Preload("Major.Faculty")
	query = query.Preload("StudyProgram")
	query = query.Preload("StudyProgram.Major")
	query = query.Preload("StudyProgram.Major.Faculty")
	return query
}
