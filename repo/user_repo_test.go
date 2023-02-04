package repo_test

import (
	"attendance-api/repo"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

type AnyTime struct{}

func (a AnyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestCreateUser(testCreate *testing.T) {
	testCreate.Run("test normal case repo create user", func(testCreate *testing.T) {
		gormDB, mockCreate := MockGormDB()

		queryCreate := "INSERT INTO `users` (`created_at`,`updated_at`,`deleted_at`,`created_by`,`updated_by`,`deleted_by`,`username`,`password`,`first_name`,`first_name`,`handphone`,`email`,`intro`,`profile`,`last_login`,`is_super_admin`,`is_active`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
		mockCreate.ExpectExec(queryCreate).
			WithArgs(AnyTime{}, AnyTime{}, u.DeletedAt, 1, u.UpdatedBy, u.DeletedBy, u.Username, u.Password, u.FirstName, u.LastName, u.Handphone, u.Email, u.Intro, u.Profile, sqlmock.AnyArg(), u.IsSuperAdmin, false).
			WillReturnResult(sqlmock.NewResult(1, 1))

		userRepo := repo.NewUserRepo(gormDB)
		_, errCreate := userRepo.CreateUser(&u)

		testCreate.Run("test store data with no error", func(testCreate *testing.T) {
			assert.Equal(testCreate, nil, errCreate)
		})
	})
}

func TestUpdateUser(testUpdate *testing.T) {
	testUpdate.Run("test normal case repo update", func(testUpdate *testing.T) {
		gormDB, mockUpdate := MockGormDB()
		queryUpdate := "UPDATE `users` SET `created_at`=?,`updated_at`=?,`created_by`=?,`username`=?,`password`=?,`first_name`=?,`last_name`=?,`handphone`=?,`email`=?,`intro`=?,`profile`=?,`is_super_admin`=? WHERE id = ?"
		mockUpdate.ExpectExec(queryUpdate).
			WithArgs(AnyTime{}, AnyTime{}, 1, u.Username, u.Password, u.FirstName, u.LastName, u.Handphone, u.Email, u.Intro, u.Profile, u.IsSuperAdmin, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		userRepo := repo.NewUserRepo(gormDB)
		_, errUpdate := userRepo.UpdateUser(1, &u)

		testUpdate.Run("test data updated with no error", func(testUpdate *testing.T) {
			assert.Equal(testUpdate, nil, errUpdate)
		})
	})
}

func TestListUser(t *testing.T) {
	t.Run("test normal case repo list user", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		queryList := "SELECT * FROM `users` WHERE username LIKE ? AND first_name LIKE ? AND last_name LIKE ? AND handphone LIKE ? AND email LIKE ? AND is_super_admin LIKE ? ORDER BY created_at asc LIMIT 2"
		mock.ExpectQuery(queryList).
			WithArgs("%"+u.Username+"%", "%"+u.FirstName+"%", "%"+u.LastName+"%", "%"+u.Handphone+"%", "%"+u.Email+"%", u.IsSuperAdmin).
			WillReturnRows(sqlmock.NewRows(nil))

		userRepo := repo.NewUserRepo(gormDB)
		_, err := userRepo.ListUser(&u, &pagination)

		t.Run("test list data with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestListUserMeta(t *testing.T) {
	t.Run("test normal case repo list user meta", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		queryListMeta := "SELECT count(*) FROM `users` WHERE username LIKE ? AND first_name LIKE ? AND last_name LIKE ? AND handphone LIKE ? AND email LIKE ? AND is_super_admin LIKE ?"
		mock.ExpectQuery(queryListMeta).
			WithArgs("%"+u.Username+"%", "%"+u.FirstName+"%", "%"+u.LastName+"%", "%"+u.Handphone+"%", "%"+u.Email+"%", u.IsSuperAdmin).
			WillReturnRows(sqlmock.NewRows(nil))

		userRepo := repo.NewUserRepo(gormDB)
		_, err := userRepo.ListUserMeta(&u, &pagination)

		t.Run("test list meta data with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestSetActiveUser(t *testing.T) {
	t.Run("test normal case repo set active user", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		querySetActive := "UPDATE `users` SET `is_active`=?,`updated_at`=? WHERE id = ?"

		mock.ExpectExec(querySetActive).
			WithArgs(true, AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		userRepo := repo.NewUserRepo(gormDB)
		_, err := userRepo.SetActiveUser(1)

		t.Run("test data set active user with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestSetDeactiveUser(t *testing.T) {
	t.Run("test normal case repo set deactive user", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		querySetNonActive := "UPDATE `users` SET `is_active`=?,`updated_at`=? WHERE id = ?"

		mock.ExpectExec(querySetNonActive).
			WithArgs(false, AnyTime{}, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		userRepo := repo.NewUserRepo(gormDB)
		_, err := userRepo.SetDeactiveUser(1)

		t.Run("test data set deactive user with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestHardDeleteUser(t *testing.T) {
	t.Run("test normal case repo delete", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		queryDelete := "DELETE FROM `users` WHERE `users`.`id` = ?"

		mock.ExpectExec(queryDelete).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		userRepo := repo.NewUserRepo(gormDB)
		err := userRepo.HardDeleteUser(1)

		t.Run("test data deleted with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}
