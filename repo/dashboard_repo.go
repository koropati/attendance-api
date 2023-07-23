package repo

import (
	"attendance-api/common/util/converter"
	"attendance-api/model"
	"fmt"
	"log"
	"sync"
	"time"

	"gorm.io/gorm"
)

type DashboardRepo interface {
	RetrieveDashboardAcademic() (result model.DashboardAcademic, err []error)
	RetrieveDashboardUser() (result model.DashboardUser, err error)
	RetrieveDashboardStudent() (result model.DashboardStudent, err error)
	RetrieveDashboardTeacher() (result model.DashboardTeacher, err error)
	RetrieveDashboardAttendance(month, year int) (results []model.DashboardAttendance, err error)
	RetrieveDashboardAttendanceSeries(month, year int) (results []model.AttendanceSeries, err error)
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

func (r dashboardRepo) RetrieveDashboardAttendance(month, year int) (results []model.DashboardAttendance, err error) {

	query := r.db.Table("attendances")
	query = query.Select("STR_TO_DATE(date, '%Y-%m-%d') as date, " +
		"MONTH(STR_TO_DATE(date, '%Y-%m-%d')) as month_period, " +
		"YEAR(STR_TO_DATE(date, '%Y-%m-%d')) as year_period, " +
		"SUM(CASE WHEN status_presence = 'presence' THEN 1 ELSE 0 END) as total_presence, " +
		"SUM(CASE WHEN status_presence = 'not_presence' THEN 1 ELSE 0 END) as total_not_presence, " +
		"SUM(CASE WHEN status_presence = 'sick' THEN 1 ELSE 0 END) as total_sick, " +
		"SUM(CASE WHEN status_presence = 'leave_attendance' THEN 1 ELSE 0 END) as total_leave_attendance, " +
		"SUM(CASE WHEN clock_in = 0 AND clock_out > 0 THEN 1 ELSE 0 END) as total_no_clock_in, " +
		"SUM(CASE WHEN clock_in > 0 AND clock_out = 0 THEN 1 ELSE 0 END) as total_no_clock_out, " +
		"SUM(CASE WHEN status = 'late' THEN 1 ELSE 0 END) as total_late, " +
		"SUM(CASE WHEN status = 'come_home_early' THEN 1 ELSE 0 END) as total_come_home_early, " +
		"SUM(CASE WHEN status = 'late_and_home_early' THEN 1 ELSE 0 END) as total_late_and_home_early")

	if month > 0 && year > 0 {
		query = query.Where("YEAR(STR_TO_DATE(date, '%Y-%m-%d')) = ? AND MONTH(STR_TO_DATE(date, '%Y-%m-%d')) = ?", year, month)
	} else if month > 0 && year <= 0 {
		year = time.Now().Year()
		query = query.Where("YEAR(STR_TO_DATE(date, '%Y-%m-%d')) = ? AND MONTH(STR_TO_DATE(date, '%Y-%m-%d')) = ?", year, month)
	} else if month <= 0 && year > 0 {
		query = query.Where("YEAR(STR_TO_DATE(date, '%Y-%m-%d')) = ?", year)
	}
	query = query.Group("YEAR(STR_TO_DATE(date, '%Y-%m-%d')), MONTH(STR_TO_DATE(date, '%Y-%m-%d'))")

	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}
	return
}

func (r dashboardRepo) RetrieveDashboardAttendanceSeries(month, year int) (results []model.AttendanceSeries, err error) {
	if month <= 0 || year <= 0 {
		month = int(time.Now().Month())
		year = time.Now().Year()
	}

	listDates := converter.GetDatesArray(month, year)
	statusAttendances := []string{"presence", "not_presence", "sick", "leave_attendance"}
	nameAttendances := []string{"Hadir", "Tidak Hadir", "Sakit", "Izin"}

	for i, status := range statusAttendances {
		var dataSeries model.AttendanceSeries
		// var dates []string
		// var datas []int

		dates := make([]string, len(listDates))
		datas := make([]int, len(listDates))

		dataSeries.Name = nameAttendances[i]
		dataSeries.MonthPeriod = month
		dataSeries.YearPeriod = year

		wg := sync.WaitGroup{}
		for j, date := range listDates {
			wg.Add(1)
			go func(j int, date string, status string) {
				count := 0

				query := r.db.Table("attendances").Select("count(*)").Where("status_presence = ? AND DATE(date) = ?", status, date)
				if errGet := query.Find(&count).Error; errGet != nil {
					log.Printf("Error Get Data Status %v Pada Tanggal %v\n", status, date)
					count = 0
				}
				dates[j] = date
				datas[j] = count
				wg.Done()
			}(j, date, status)
		}
		wg.Wait()
		dataSeries.Date = dates
		dataSeries.Data = datas

		results = append(results, dataSeries)
	}

	return

}
