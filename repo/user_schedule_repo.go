package repo

import (
	"attendance-api/common/util/converter"
	"attendance-api/model"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"
)

type UserScheduleRepo interface {
	CreateUserSchedule(userschedule model.UserSchedule) (model.UserSchedule, error)
	RetrieveUserSchedule(id int) (model.UserSchedule, error)
	RetrieveUserScheduleByOwner(id int, ownerID int) (model.UserSchedule, error)
	ListMySchedule(userID int, filter model.MyScheduleFilter) (results []model.ListMySchedule, err error)
	ListTodaySchedule(userID int, dayName string) ([]model.TodaySchedule, error)
	ListUserInRule(scheduleID int, user model.Student, pagination model.Pagination) ([]model.Student, error)
	ListUserInRuleMeta(scheduleID int, user model.Student, pagination model.Pagination) (model.Meta, error)
	ListUserNotInRule(scheduleID int, user model.Student, pagination model.Pagination) ([]model.Student, error)
	ListUserNotInRuleMeta(scheduleID int, user model.Student, pagination model.Pagination) (model.Meta, error)
	UpdateUserSchedule(id int, userschedule model.UserSchedule) (model.UserSchedule, error)
	UpdateUserScheduleByOwner(id int, ownerID int, userschedule model.UserSchedule) (model.UserSchedule, error)
	DeleteUserSchedule(id int) error
	DeleteUserScheduleByOwner(id int, ownerID int) error
	RemoveUserFromSchedule(scheduleID int, userID int) error
	RemoveUserFromScheduleByOwner(scheduleID int, userID int, ownerID int) error
	ListUserSchedule(userschedule model.UserSchedule, pagination model.Pagination) ([]model.UserSchedule, error)
	ListUserScheduleMeta(userschedule model.UserSchedule, pagination model.Pagination) (model.Meta, error)
	DropDownUserSchedule(userschedule model.UserSchedule) ([]model.UserSchedule, error)
	CheckHaveSchedule(userID int, date time.Time) (isHaveSchedule bool, scheduleID int, err error)
	CheckUserInSchedule(scheduleID int, userID int) bool
	CountByScheduleID(scheduleID int) (total int)
	GetAll() (results []model.UserSchedule, err error)
	GetAllByTodayRange() (results []model.UserSchedule, err error)
}

type userScheduleRepo struct {
	db *gorm.DB
}

func NewUserScheduleRepo(db *gorm.DB) UserScheduleRepo {
	return userScheduleRepo{db: db}
}

func (r userScheduleRepo) CreateUserSchedule(userschedule model.UserSchedule) (result model.UserSchedule, err error) {
	if err := r.db.Table("user_schedules").Create(&userschedule).Error; err != nil {
		return model.UserSchedule{}, err
	}
	if err := PreloadUserSchedule(r.db.Table("user_schedules")).Where("id = ?", userschedule.ID).First(&result).Error; err != nil {
		return model.UserSchedule{}, err
	}
	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility()
	return
}

func (r userScheduleRepo) RetrieveUserSchedule(id int) (result model.UserSchedule, err error) {
	if err := PreloadUserSchedule(r.db.Table("user_schedules")).Where("id = ?", id).First(&result).Error; err != nil {
		return model.UserSchedule{}, err
	}
	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility()
	return
}

func (r userScheduleRepo) RetrieveUserScheduleByOwner(id int, ownerID int) (result model.UserSchedule, err error) {
	if err := PreloadUserSchedule(r.db.Table("user_schedules")).Where("id = ? AND owner_id = ?", id, ownerID).First(&result).Error; err != nil {
		return model.UserSchedule{}, err
	}
	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility()
	return
}

func (r userScheduleRepo) UpdateUserSchedule(id int, userschedule model.UserSchedule) (result model.UserSchedule, err error) {
	if err := r.db.Model(&model.UserSchedule{}).Where("id = ?", id).Updates(&userschedule).Error; err != nil {
		return model.UserSchedule{}, err
	}
	if err := PreloadUserSchedule(r.db.Table("user_schedules")).Where("id = ?", id).First(&result).Error; err != nil {
		return model.UserSchedule{}, err
	}
	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility()
	return
}

func (r userScheduleRepo) UpdateUserScheduleByOwner(id int, ownerID int, userschedule model.UserSchedule) (result model.UserSchedule, err error) {
	if err := PreloadUserSchedule(r.db.Table("user_schedules")).Where("id = ? AND owner_id = ?", id, ownerID).Updates(&userschedule).Error; err != nil {
		return model.UserSchedule{}, err
	}
	if err := PreloadUserSchedule(r.db.Table("user_schedules")).Where("id = ?", id).First(&result).Error; err != nil {
		return model.UserSchedule{}, err
	}
	result.User.Role = result.User.GetRole()
	result.User.Avatar = result.User.GetAvatar()
	result.User.UserAbilities = result.User.GetAbility()
	return
}

func (r userScheduleRepo) DeleteUserSchedule(id int) error {
	if err := r.db.Delete(&model.UserSchedule{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (r userScheduleRepo) RemoveUserFromSchedule(scheduleID int, userID int) error {
	var dataToDelete model.UserSchedule

	if err := r.db.Model(&model.UserSchedule{}).Where("schedule_id = ? AND user_id = ?", scheduleID, userID).First(&dataToDelete).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&model.UserSchedule{}, dataToDelete.ID).Error; err != nil {
		return err
	}
	return nil
}

func (r userScheduleRepo) DeleteUserScheduleByOwner(id int, ownerID int) error {
	if err := r.db.Where("id = ? AND owner_id = ?", id, ownerID).Delete(&model.UserSchedule{}).Error; err != nil {
		return err
	}
	return nil
}

func (r userScheduleRepo) RemoveUserFromScheduleByOwner(scheduleID int, userID int, ownerID int) error {
	var dataToDelete model.UserSchedule

	if err := r.db.Model(&model.UserSchedule{}).Where("schedule_id = ? AND user_id = ? AND owner_id = ?", scheduleID, userID, ownerID).First(&dataToDelete).Error; err != nil {
		return err
	}

	if err := r.db.Delete(&model.UserSchedule{}, dataToDelete.ID).Error; err != nil {
		return err
	}
	return nil
}

func (r userScheduleRepo) ListMySchedule(userID int, filter model.MyScheduleFilter) (results []model.ListMySchedule, err error) {

	var month int
	var year int

	if filter.Month != "" && filter.Year != "" {
		month, _ = strconv.Atoi(filter.Month)
		year, _ = strconv.Atoi(filter.Year)
	} else {
		month = int(time.Now().Month())
		year = time.Now().Year()
	}

	listDates := converter.GetDatesArray(month, year)

	finalResult := make([]model.ListMySchedule, len(listDates))

	wg := sync.WaitGroup{}

	for i, date := range listDates {
		wg.Add(1)
		go func(i int, date string, userID int) {
			dayName, errorDayName := converter.GetEnglishDayName(date)
			if errorDayName != nil {
				log.Printf("Error Get Day Name E: %v\n", errorDayName)
			}
			indonesianDate, errIndoDate := converter.FormatTanggalIndonesia(date)
			if errIndoDate != nil {
				log.Printf("Error Get Indonesian Date E: %v\n", errIndoDate)
			}

			var listMySchedule model.ListMySchedule

			var mySchedule []model.MySchedule

			rawQuery := fmt.Sprintf(`
					SELECT 
					us.schedule_id as schedule_id, 
					s.name as schedule_name, 
					s.code as schedule_code, 
					DATE(s.start_date) as start_date, 
					DATE(s.end_date) as end_date, 
					s.subject_id as subject_id, 
					sbj.name as subject_name, 
					sbj.code as subject_code, 
					s.late_duration as late_duration, 
					s.latitude as latitude, 
					s.longitude as longitude, 
					s.radius as radius 
					FROM user_schedules us 
					LEFT JOIN schedules s ON us.schedule_id = s.id 
					LEFT JOIN subjects sbj ON s.subject_id = sbj.id
					LEFT JOIN daily_schedules ds ON s.id = ds.schedule_id 
					WHERE us.user_id = %d AND (DATE('%s') BETWEEN DATE(s.start_date) AND DATE(s.end_date)) AND ds.name = '%s' AND us.deleted_at IS NULL`, userID, date, dayName)
			if err := r.db.Raw(rawQuery).Scan(&mySchedule).Error; err != nil {
				log.Printf("Error Get Day Name E: %v\n", errorDayName)
			}
			if len(mySchedule) > 0 {
				listMySchedule.IndeonesianDate = indonesianDate
				listMySchedule.Schedules = mySchedule
			} else {
				mySchedule = append(mySchedule, model.MySchedule{
					ScheduleID:   0,
					ScheduleName: "Tidak Memiliki Jadwal",
					ScheduleCode: "-",
					StartDate:    date,
					EndDate:      date,
					SubjectID:    0,
					SubjectName:  "-",
					SubjectCode:  "-",
					LateDuration: 0,
					Latitude:     0,
					Longitude:    0,
					Radius:       0,
				})
				listMySchedule.IndeonesianDate = indonesianDate
				listMySchedule.Schedules = mySchedule
			}
			finalResult[i] = listMySchedule
			wg.Done()
		}(i, date, userID)
	}
	wg.Wait()

	return finalResult, nil
}

func (r userScheduleRepo) ListTodaySchedule(userID int, dayName string) (results []model.TodaySchedule, err error) {
	today := time.Now().Format("2006-01-02")
	query := fmt.Sprintf(`
	SELECT 
	s.id as schedule_id, 
	s.name as schedule_name, 
	s.code as schedule_code, 
	sbj.id as subject_id, 
	sbj.name as subject_name, 
	ds.start_time as start_time, 
	ds.end_time as end_time 
	FROM user_schedules us 
	LEFT JOIN schedules s ON us.schedule_id = s.id 
	LEFT JOIN subjects sbj ON s.subject_id = sbj.id 
	LEFT JOIN daily_schedules ds ON us.schedule_id = ds.schedule_id 
	WHERE us.user_id = %d AND ds.name = '%s' AND '%s' BETWEEN DATE(s.start_date) AND DATE(s.end_date) AND us.deleted_at IS NULL`, userID, dayName, today)
	if err := r.db.Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	return
}

func (r userScheduleRepo) ListUserSchedule(userschedule model.UserSchedule, pagination model.Pagination) ([]model.UserSchedule, error) {
	var userschedules []model.UserSchedule
	offset := (pagination.Page - 1) * pagination.Limit
	query := PreloadUserSchedule(r.db.Table("user_schedules")).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterUserSchedule(query, userschedule)
	query = SearchUserSchedule(query, pagination.Search)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for i, data := range userschedules {
		wg.Add(1)
		go func(i int, data model.UserSchedule) {
			userschedules[i].User.Role = data.User.GetRole()
			userschedules[i].User.Avatar = data.User.GetAvatar()
			userschedules[i].User.UserAbilities = data.User.GetAbility()
			wg.Done()
		}(i, data)
	}
	wg.Wait()
	return userschedules, nil
}

func (r userScheduleRepo) ListUserInRule(scheduleID int, student model.Student, pagination model.Pagination) ([]model.Student, error) {
	var userID []int
	if student.OwnerID > 0 {
		if err := r.db.Model(&[]model.UserSchedule{}).Select("user_id").Where("schedule_id = ? AND owner_id = ?", scheduleID, student.OwnerID).Find(&userID).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Model(&[]model.UserSchedule{}).Select("user_id").Where("schedule_id = ?", scheduleID).Find(&userID).Error; err != nil {
			return nil, err
		}
	}

	var students []model.Student
	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Model(&[]model.Student{})
	query = PreloadStudent(query)
	query = query.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = query.Where("user_id IN (?)", userID)
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
			students[i].User.Avatar = student.User.GetAvatar()
			students[i].User.Role = student.User.GetRole()
			students[i].User.UserAbilities = student.User.GetAbility()
			wg.Done()
		}(i, student)
	}
	wg.Wait()

	return students, nil
}

func (r userScheduleRepo) ListUserNotInRule(scheduleID int, student model.Student, pagination model.Pagination) ([]model.Student, error) {
	var userID []int

	if student.OwnerID > 0 {
		if err := r.db.Model(&[]model.UserSchedule{}).Select("user_id").Where("schedule_id = ? AND owner_id = ?", scheduleID, student.OwnerID).Find(&userID).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.db.Model(&[]model.UserSchedule{}).Select("user_id").Where("schedule_id = ?", scheduleID).Find(&userID).Error; err != nil {
			return nil, err
		}
	}
	if len(userID) <= 0 {
		userID = append(userID, 0)
	}

	var students []model.Student
	offset := (pagination.Page - 1) * pagination.Limit
	query := r.db.Model(&[]model.Student{})
	query = PreloadStudent(query)
	query = query.Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = query.Not("user_id", userID)
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
			students[i].User.Avatar = student.User.GetAvatar()
			students[i].User.Role = student.User.GetRole()
			students[i].User.UserAbilities = student.User.GetAbility()
			wg.Done()
		}(i, student)
	}
	wg.Wait()

	return students, nil
}

func (r userScheduleRepo) ListUserScheduleMeta(userschedule model.UserSchedule, pagination model.Pagination) (model.Meta, error) {
	var userschedules []model.UserSchedule
	var totalRecord int
	var totalPage int

	queryTotal := r.db.Model(&model.UserSchedule{}).Select("count(*)")
	queryTotal = FilterUserSchedule(queryTotal, userschedule)
	queryTotal = SearchUserSchedule(queryTotal, pagination.Search)
	queryTotal = queryTotal.Scan(&totalRecord)
	if err := queryTotal.Error; err != nil {
		return model.Meta{}, err
	}

	totalPage = int(totalRecord / pagination.Limit)
	if totalRecord%pagination.Limit > 0 {
		totalPage += 1
	}

	offset := (pagination.Page - 1) * pagination.Limit

	query := r.db.Table("user_schedules").Limit(pagination.Limit).Offset(offset).Order(pagination.Sort)
	query = FilterUserSchedule(query, userschedule)
	query = SearchUserSchedule(query, pagination.Search)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return model.Meta{}, err
	}

	meta := model.Meta{
		CurrentPage:   pagination.Page,
		TotalPage:     totalPage,
		TotalRecord:   totalRecord,
		CurrentRecord: len(userschedules),
	}
	return meta, nil
}

func (r userScheduleRepo) ListUserInRuleMeta(scheduleID int, student model.Student, pagination model.Pagination) (model.Meta, error) {
	var students []model.Student
	var totalRecord int
	var totalPage int

	var userID []int
	if student.OwnerID > 0 {
		if err := r.db.Table("user_schedules").Select("user_id").Where("schedule_id = ? AND owner_id = ?", scheduleID, student.OwnerID).Find(&userID).Error; err != nil {
			return model.Meta{}, err
		}
	} else {
		if err := r.db.Table("user_schedules").Select("user_id").Where("schedule_id = ?", scheduleID).Find(&userID).Error; err != nil {
			return model.Meta{}, err
		}
	}

	queryTotal := r.db.Model(&model.Student{}).Select("count(*)")
	queryTotal = queryTotal.Where("user_id IN (?)", userID)
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
	query = query.Where("user_id IN (?)", userID)
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

func (r userScheduleRepo) ListUserNotInRuleMeta(scheduleID int, student model.Student, pagination model.Pagination) (model.Meta, error) {
	var students []model.User
	var totalRecord int
	var totalPage int

	var userID []int
	if student.OwnerID > 0 {
		if err := r.db.Table("user_schedules").Select("user_id").Where("schedule_id = ? AND owner_id = ?", scheduleID, student.OwnerID).Find(&userID).Error; err != nil {
			return model.Meta{}, err
		}
	} else {
		if err := r.db.Table("user_schedules").Select("user_id").Where("schedule_id = ?", scheduleID).Find(&userID).Error; err != nil {
			return model.Meta{}, err
		}
	}

	if len(userID) <= 0 {
		userID = append(userID, 0)
	}

	queryTotal := r.db.Model(&model.Student{}).Select("count(*)")
	queryTotal = queryTotal.Not("user_id", userID)
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
	query = query.Not("user_id", userID)
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

func (r userScheduleRepo) DropDownUserSchedule(userschedule model.UserSchedule) ([]model.UserSchedule, error) {
	var userschedules []model.UserSchedule
	query := PreloadUserSchedule(r.db.Table("user_schedules")).Order("id desc")
	query = FilterUserSchedule(query, userschedule)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for i, data := range userschedules {
		wg.Add(1)
		go func(i int, data model.UserSchedule) {
			userschedules[i].User.Role = data.User.GetRole()
			userschedules[i].User.Avatar = data.User.GetAvatar()
			userschedules[i].User.UserAbilities = data.User.GetAbility()
			wg.Done()
		}(i, data)
	}
	wg.Wait()
	return userschedules, nil
}

func (r userScheduleRepo) GetAll() (results []model.UserSchedule, err error) {
	var userschedules []model.UserSchedule
	query := PreloadUserSchedule(r.db.Table("user_schedules")).Order("id desc")
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for i, data := range userschedules {
		wg.Add(1)
		go func(i int, data model.UserSchedule) {
			userschedules[i].User.Role = data.User.GetRole()
			userschedules[i].User.Avatar = data.User.GetAvatar()
			userschedules[i].User.UserAbilities = data.User.GetAbility()
			wg.Done()
		}(i, data)
	}
	wg.Wait()
	return userschedules, nil
}

func (r userScheduleRepo) GetAllByTodayRange() (resutls []model.UserSchedule, err error) {
	today := time.Now().UTC().Truncate(24 * time.Hour)

	var idSchedule []int

	queryID := r.db.Model(&[]model.Schedule{}).Select("id")
	queryID = queryID.Where("start_date <= ? AND end_date >= ?", today, today)
	if err := queryID.Find(&idSchedule).Error; err != nil {
		return nil, err
	}

	if len(idSchedule) <= 0 {
		return nil, nil
	}

	var userschedules []model.UserSchedule
	query := PreloadUserSchedule(r.db.Table("user_schedules")).Order("id desc")
	query = query.Where("schedule_id IN ?", idSchedule)
	query = query.Find(&userschedules)
	if err := query.Error; err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	for i, data := range userschedules {
		wg.Add(1)
		go func(i int, data model.UserSchedule) {
			userschedules[i].User.Role = data.User.GetRole()
			userschedules[i].User.Avatar = data.User.GetAvatar()
			userschedules[i].User.UserAbilities = data.User.GetAbility()
			wg.Done()
		}(i, data)
	}
	wg.Wait()
	return userschedules, nil
}

func (r userScheduleRepo) CheckHaveSchedule(userID int, date time.Time) (isHaveSchedule bool, scheduleID int, err error) {
	type DataSchedule struct {
		IsHaveSchedule bool `json:"is_have_schedule" query:"is_have_schedule"`
		ScheduleID     int  `json:"schedule_id" query:"schedule_id"`
	}
	var data DataSchedule
	rawQuery := fmt.Sprintf(`SELECT COUNT(*) > 0 as is_have_schedule, us.schedule_id as schedule_id 
	FROM user_schedules us 
	LEFT JOIN schedules s ON us.schedule_id = s.id 
	WHERE us.user_id = %d AND '%v' BETWEEN s.start_date AND s.end_date`, userID, date)

	if err := r.db.Raw(rawQuery).Scan(&data).Error; err != nil {
		return data.IsHaveSchedule, data.ScheduleID, err
	}
	return data.IsHaveSchedule, data.ScheduleID, nil
}

func (r userScheduleRepo) CheckUserInSchedule(scheduleID int, userID int) (isHave bool) {
	rawQuery := fmt.Sprintf(`SELECT COUNT(*) > 0 FROM user_schedules us WHERE us.user_id = %d AND us.schedule_id = %d`, userID, scheduleID)
	if err := r.db.Raw(rawQuery).Scan(&isHave).Error; err != nil {
		return false
	}
	return
}

func (r userScheduleRepo) CountByScheduleID(scheduleID int) (total int) {
	if err := r.db.Table("user_schedules").Select("count(*)").Where("schedule_id = ? AND user_id != ?", scheduleID, 0).Find(&total).Error; err != nil {
		total = 0
	}
	return
}

func FilterUserSchedule(query *gorm.DB, userschedule model.UserSchedule) *gorm.DB {
	if userschedule.UserID > 0 {
		query = query.Where("user_id = ?", userschedule.UserID)
	}
	if userschedule.ScheduleID > 0 {
		query = query.Where("schedule_id = ?", userschedule.ScheduleID)
	}
	if userschedule.OwnerID > 0 {
		query = query.Where("owner_id = ?", userschedule.OwnerID)
	}
	return query
}

func SearchUserSchedule(query *gorm.DB, search string) *gorm.DB {
	if search != "" {
		query = query.Where("user_id LIKE ? OR schedule_id LIKE ? OR owner_id LIKE ? ", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	return query
}

func PreloadUserSchedule(query *gorm.DB) *gorm.DB {
	query = query.Preload("User")
	query = query.Preload("Schedule")
	return query
}
