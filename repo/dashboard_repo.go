package repo

import (
	"attendance-api/model"
	"fmt"
	"log"

	"gorm.io/gorm"
)

type DashboardRepo interface {
	RetrieveDashboardAcademic() (result model.DashboardAcademic, err []error)
	RetrieveDashboardUser() (result model.DashboardUser, err error)
	RetrieveDashboardStudent() (result model.DashboardStudent, err error)
	RetrieveDashboardTeacher() (result model.DashboardTeacher, err error)
}

type dashboardRepo struct {
	db *gorm.DB
}

func NewDashboardRepo(db *gorm.DB) DashboardRepo {
	return &dashboardRepo{db: db}
}

func (r dashboardRepo) RetrieveDashboardAcademic() (result model.DashboardAcademic, err []error) {

	if errFaculty := r.db.Table("faculties").Select("count(*)").Find(&result.TotalFaculty).Error; errFaculty != nil {
		log.Printf("Error RetrieveDashboardAcademic() [Faculty] E: %v\n", errFaculty)
		err = append(err, errFaculty)
	}

	if errMajor := r.db.Table("majors").Select("count(*)").Find(&result.TotalMajor).Error; errMajor != nil {
		log.Printf("Error RetrieveDashboardAcademic() [Major] E: %v\n", errMajor)
		err = append(err, errMajor)
	}

	if errStudyProgram := r.db.Table("study_programs").Select("count(*)").Find(&result.TotalStudyProgram).Error; errStudyProgram != nil {
		log.Printf("Error RetrieveDashboardAcademic() [Study Program] E: %v\n", errStudyProgram)
		err = append(err, errStudyProgram)
	}

	if errSubject := r.db.Table("subjects").Select("count(*)").Find(&result.TotalSubject).Error; errSubject != nil {
		log.Printf("Error RetrieveDashboardAcademic() [Subject] E: %v\n", errSubject)
		err = append(err, errSubject)
	}

	if errSchedule := r.db.Table("schedules").Select("count(*)").Find(&result.TotalSchedule).Error; errSchedule != nil {
		log.Printf("Error RetrieveDashboardAcademic() [Schedule] E: %v\n", errSchedule)
		err = append(err, errSchedule)
	}

	return result, err
}

func (r dashboardRepo) RetrieveDashboardUser() (result model.DashboardUser, err error) {
	rawQuery := fmt.Sprintf(`SELECT COUNT(*) AS total_user,
									SUM(CASE WHEN is_active = 1 THEN 1 ELSE 0 END) AS total_user_active,
									SUM(CASE WHEN is_active = 0 THEN 1 ELSE 0 END) AS total_user_non_active,
									SUM(CASE WHEN is_super_admin = 1 THEN 1 ELSE 0 END) AS total_super_admin 
							FROM %s`, "users")
	if err := r.db.Raw(rawQuery).Scan(&result).Error; err != nil {
		return model.DashboardUser{}, err
	}
	return
}

func (r dashboardRepo) RetrieveDashboardStudent() (result model.DashboardStudent, err error) {
	rawQuery := fmt.Sprintf(`SELECT COUNT(s.id) AS total_student,
									SUM(CASE WHEN u.is_active = 1 THEN 1 ELSE 0 END) AS total_student_active,
									SUM(CASE WHEN u.is_active = 0 THEN 1 ELSE 0 END) AS total_student_non_active 
							FROM %s s LEFT JOIN users u ON s.user_id= u.id`, "students")
	if err := r.db.Raw(rawQuery).Scan(&result).Error; err != nil {
		return model.DashboardStudent{}, err
	}
	return
}

func (r dashboardRepo) RetrieveDashboardTeacher() (result model.DashboardTeacher, err error) {
	rawQuery := fmt.Sprintf(`SELECT COUNT(t.id) AS total_teacher,
									SUM(CASE WHEN u.is_active = 1 THEN 1 ELSE 0 END) AS total_teacher_active,
									SUM(CASE WHEN u.is_active = 0 THEN 1 ELSE 0 END) AS total_teacher_non_active 
							FROM %s t LEFT JOIN users u ON t.user_id = u.id`, "teachers")
	if err := r.db.Raw(rawQuery).Scan(&result).Error; err != nil {
		return model.DashboardTeacher{}, err
	}
	return
}
