package repo_test

import (
	"attendance-api/model"
	"attendance-api/repo"
	"strconv"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var gormCustom = model.GormCustom{
	CreatedBy: 1,
	UpdatedBy: 0,
	DeletedBy: 0,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

var scheduleData = model.Schedule{
	GormCustom:   gormCustom,
	Name:         "Robot",
	Code:         "adwd2",
	StartDate:    time.Now(),
	EndDate:      time.Now().AddDate(0, 1, 0),
	LateDuration: 15,
}

func TestCreateSchedule(t *testing.T) {
	t.Run("test normal case repo create schedule", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		query := "INSERT INTO `schedules` (`created_at`,`updated_at`,`deleted_at`,`created_by`,`updated_by`,`deleted_by`,`name`,`code`,`start_date`,`end_date`,`late_duration`) VALUES (?,?,?,?,?,?,?,?,?,?,?)"
		mock.ExpectExec(query).
			WithArgs(AnyTime{}, AnyTime{}, scheduleData.DeletedAt, 1, scheduleData.UpdatedBy, scheduleData.DeletedBy, scheduleData.Name, scheduleData.Code, scheduleData.StartDate, scheduleData.EndDate, scheduleData.LateDuration).
			WillReturnResult(sqlmock.NewResult(1, 1))

		scheduleRepo := repo.NewScheduleRepo(gormDB)
		_, err := scheduleRepo.CreateSchedule(&scheduleData)

		t.Run("test store data with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestRetrieveSchedule(t *testing.T) {
	t.Run("test normal case repo retrieve schedule", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		columns := []string{"id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by", "name", "code", "start_date", "end_date", "late_duration"}
		query := "SELECT * FROM `schedules` WHERE `schedules`.`id` = ? ORDER BY `schedules`.`id` LIMIT 1"
		mock.ExpectQuery(query).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, scheduleData.CreatedAt, scheduleData.UpdatedAt, nil, 1, scheduleData.UpdatedBy, scheduleData.DeletedBy, scheduleData.Name, scheduleData.Code, scheduleData.StartDate, scheduleData.EndDate, scheduleData.LateDuration))

		scheduleRepo := repo.NewScheduleRepo(gormDB)
		_, err := scheduleRepo.RetrieveSchedule(1)

		t.Run("test data retrieve with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestUpdateSchedule(t *testing.T) {
	t.Run("test normal case repo schedule update", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		query := "UPDATE `schedules` SET `created_at`=?,`updated_at`=?,`created_by`=?,`name`=?,`code`=?,`start_date`=?,`end_date`=?,`late_duration`=? WHERE id = ?"
		mock.ExpectExec(query).
			WithArgs(AnyTime{}, AnyTime{}, 1, scheduleData.Name, scheduleData.Code, scheduleData.StartDate, scheduleData.EndDate, scheduleData.LateDuration, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		scheduleRepo := repo.NewScheduleRepo(gormDB)
		_, err := scheduleRepo.UpdateSchedule(1, &scheduleData)

		t.Run("test data updated with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestDeleteSchedule(t *testing.T) {
	t.Run("test normal case repo schedule delete", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		query := "DELETE FROM `schedules` WHERE `schedules`.`id` = ?"
		mock.ExpectExec(query).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		scheduleRepo := repo.NewScheduleRepo(gormDB)
		err := scheduleRepo.DeleteSchedule(1)

		t.Run("test data deleted with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestListSchedule(t *testing.T) {
	t.Run("test normal case repo schedule list", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		query := "SELECT * FROM `schedules` WHERE name LIKE ? AND code LIKE ? AND start_date LIKE ? AND end_date LIKE ? AND late_duration LIKE ? ORDER BY created_at asc LIMIT 2"
		mock.ExpectQuery(query).
			WithArgs("%"+scheduleData.Name+"%", "%"+scheduleData.Code+"%", "%"+scheduleData.StartDate.Format("2006-01-02")+"%", "%"+scheduleData.EndDate.Format("2006-01-02")+"%", "%"+strconv.Itoa(scheduleData.LateDuration)+"%").
			WillReturnRows(sqlmock.NewRows(nil))

		scheduleRepo := repo.NewScheduleRepo(gormDB)
		_, err := scheduleRepo.ListSchedule(&scheduleData, &pagination)

		t.Run("test data list with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestListScheduleMeta(t *testing.T) {
	t.Run("test normal case repo list schedule meta", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		queryListMeta := "SELECT count(*) FROM `schedules` WHERE name LIKE ? AND code LIKE ? AND start_date LIKE ? AND end_date LIKE ? AND late_duration LIKE ?"
		mock.ExpectQuery(queryListMeta).
			WithArgs("%"+scheduleData.Name+"%", "%"+scheduleData.Code+"%", "%"+scheduleData.StartDate.Format("2006-01-02")+"%", "%"+scheduleData.EndDate.Format("2006-01-02")+"%", "%"+strconv.Itoa(scheduleData.LateDuration)+"%").
			WillReturnRows(sqlmock.NewRows(nil))

		scheduleRepo := repo.NewScheduleRepo(gormDB)
		_, err := scheduleRepo.ListScheduleMeta(&scheduleData, &pagination)

		t.Run("test list meta data with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestDropDownSchedule(t *testing.T) {
	t.Run("test normal case repo schedule dropdown", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		query := "SELECT * FROM `schedules` WHERE name LIKE ? AND code LIKE ? AND start_date LIKE ? AND end_date LIKE ? AND late_duration LIKE ?"
		mock.ExpectQuery(query).
			WithArgs("%"+scheduleData.Name+"%", "%"+scheduleData.Code+"%", "%"+scheduleData.StartDate.Format("2006-01-02")+"%", "%"+scheduleData.EndDate.Format("2006-01-02")+"%", "%"+strconv.Itoa(scheduleData.LateDuration)+"%").
			WillReturnRows(sqlmock.NewRows(nil))

		scheduleRepo := repo.NewScheduleRepo(gormDB)
		_, err := scheduleRepo.DropDownSchedule(&scheduleData)

		t.Run("test data dropdown with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}
