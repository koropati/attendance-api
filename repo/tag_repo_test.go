package repo_test

import (
	"attendance-api/model"
	"attendance-api/repo"
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

var tagData = model.Tag{
	GormCustom: gormCustom,
	Title:      "Robot",
	MetaTitle:  "robot robotika",
	Slug:       "robot",
	Content:    "Robotika Arduino",
}

func TestCreateTag(t *testing.T) {
	t.Run("test normal case repo create tag", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		query := "INSERT INTO `tags` (`created_at`,`updated_at`,`deleted_at`,`created_by`,`updated_by`,`deleted_by`,`title`,`meta_title`,`slug`,`content`) VALUES (?,?,?,?,?,?,?,?,?,?)"
		mock.ExpectExec(query).
			WithArgs(AnyTime{}, AnyTime{}, tagData.DeletedAt, 1, tagData.UpdatedBy, tagData.DeletedBy, tagData.Title, tagData.MetaTitle, tagData.Slug, tagData.Content).
			WillReturnResult(sqlmock.NewResult(1, 1))

		tagRepo := repo.NewTagRepo(gormDB)
		_, err := tagRepo.CreateTag(&tagData)

		t.Run("test store data with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestRetrieveTag(t *testing.T) {
	t.Run("test normal case repo retrieve tag", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		columns := []string{"id", "created_at", "updated_at", "deleted_at", "created_by", "updated_by", "deleted_by", "title", "meta_title", "slug", "content"}
		query := "SELECT * FROM `tags` WHERE `tags`.`id` = ? ORDER BY `tags`.`id` LIMIT 1"
		mock.ExpectQuery(query).WithArgs(1).WillReturnRows(sqlmock.NewRows(columns).AddRow(1, tagData.CreatedAt, tagData.UpdatedAt, nil, 1, tagData.UpdatedBy, tagData.DeletedBy, tagData.Title, tagData.MetaTitle, tagData.Slug, tagData.Content))

		tagRepo := repo.NewTagRepo(gormDB)
		_, err := tagRepo.RetrieveTag(1)

		t.Run("test data retrieve with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestUpdateTag(t *testing.T) {
	t.Run("test normal case repo tag update", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		query := "UPDATE `tags` SET `created_at`=?,`updated_at`=?,`created_by`=?,`title`=?,`meta_title`=?,`slug`=?,`content`=? WHERE id = ?"
		mock.ExpectExec(query).
			WithArgs(AnyTime{}, AnyTime{}, 1, tagData.Title, tagData.MetaTitle, tagData.Slug, tagData.Content, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		tagRepo := repo.NewTagRepo(gormDB)
		_, err := tagRepo.UpdateTag(1, &tagData)

		t.Run("test data updated with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestDeleteTag(t *testing.T) {
	t.Run("test normal case repo tag delete", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		query := "DELETE FROM `tags` WHERE `tags`.`id` = ?"
		mock.ExpectExec(query).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		tagRepo := repo.NewTagRepo(gormDB)
		err := tagRepo.DeleteTag(1)

		t.Run("test data deleted with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestListTag(t *testing.T) {
	t.Run("test normal case repo tag list", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		query := "SELECT * FROM `tags` WHERE title LIKE ? AND slug LIKE ? AND meta_title LIKE ? AND content LIKE ? ORDER BY created_at asc LIMIT 2"
		mock.ExpectQuery(query).
			WithArgs("%"+tagData.Title+"%", "%"+tagData.Slug+"%", "%"+tagData.MetaTitle+"%", "%"+tagData.Content+"%").
			WillReturnRows(sqlmock.NewRows(nil))

		tagRepo := repo.NewTagRepo(gormDB)
		_, err := tagRepo.ListTag(&tagData, &pagination)

		t.Run("test data list with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestListTagMeta(t *testing.T) {
	t.Run("test normal case repo list tag meta", func(t *testing.T) {
		gormDB, mock := MockGormDB()

		queryListMeta := "SELECT count(*) FROM `tags` WHERE title LIKE ? AND slug LIKE ? AND meta_title LIKE ? AND content LIKE ?"
		mock.ExpectQuery(queryListMeta).
			WithArgs("%"+tagData.Title+"%", "%"+tagData.Slug+"%", "%"+tagData.MetaTitle+"%", "%"+tagData.Content+"%").
			WillReturnRows(sqlmock.NewRows(nil))

		tagRepo := repo.NewTagRepo(gormDB)
		_, err := tagRepo.ListTagMeta(&tagData, &pagination)

		t.Run("test list meta data with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}

func TestDropDownTag(t *testing.T) {
	t.Run("test normal case repo tag dropdown", func(t *testing.T) {
		gormDB, mock := MockGormDB()
		query := "SELECT * FROM `tags` WHERE title LIKE ? AND slug LIKE ? AND meta_title LIKE ? AND content LIKE ?"
		mock.ExpectQuery(query).
			WithArgs("%"+tagData.Title+"%", "%"+tagData.Slug+"%", "%"+tagData.MetaTitle+"%", "%"+tagData.Content+"%").
			WillReturnRows(sqlmock.NewRows(nil))

		tagRepo := repo.NewTagRepo(gormDB)
		_, err := tagRepo.DropDownTag(&tagData)

		t.Run("test data dropdown with no error", func(t *testing.T) {
			assert.Equal(t, nil, err)
		})
	})
}
