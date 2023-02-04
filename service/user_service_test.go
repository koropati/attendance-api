package service_test

import (
	"attendance-api/mocks"
	"attendance-api/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListUser(t *testing.T) {
	t.Run("test normal case service list user", func(t *testing.T) {
		UserRepoMock := new(mocks.UserRepoMock)
		UserRepoMock.On("ListUser", mock.AnythingOfType("*model.User"), mock.AnythingOfType("*model.Pagination")).Return(nil)
		// UserRepoMock.On("ListUser", mock.AnythingOfType("*model.User, *model.Pagination")).Return(nil)

		userService := service.NewUserService(UserRepoMock)
		_, err := userService.ListUser(&u, &p)

		t.Run("test list user", func(t *testing.T) {
			assert.Equal(t, true, err == nil)
		})
	})
}

func TestListUserMeta(t *testing.T) {
	t.Run("test normal case service list user meta", func(t *testing.T) {
		UserRepoMock := new(mocks.UserRepoMock)
		UserRepoMock.On("ListUserMeta", mock.AnythingOfType("*model.User"), mock.AnythingOfType("*model.Pagination")).Return(nil)

		userService := service.NewUserService(UserRepoMock)
		_, err := userService.ListUserMeta(&u, &p)

		t.Run("test list user meta", func(t *testing.T) {
			assert.Equal(t, true, err == nil)
		})
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("test normal case service create user", func(t *testing.T) {
		UserRepoMock := new(mocks.UserRepoMock)
		UserRepoMock.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)

		userService := service.NewUserService(UserRepoMock)
		_, err := userService.CreateUser(&u)

		t.Run("test create user", func(t *testing.T) {
			assert.Equal(t, true, err == nil)
		})
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("test normal case service update user", func(t *testing.T) {
		UserRepoMock := new(mocks.UserRepoMock)
		UserRepoMock.On("UpdateUser", mock.AnythingOfType("int, *model.User")).Return(nil)

		userService := service.NewUserService(UserRepoMock)
		_, err := userService.UpdateUser(1, &u)

		t.Run("test update user", func(t *testing.T) {
			assert.Equal(t, true, err == nil)
		})
	})
}

func TestHardDeleteUser(t *testing.T) {
	t.Run("test normal case service hard Delete user", func(t *testing.T) {
		UserRepoMock := new(mocks.UserRepoMock)
		UserRepoMock.On("HardDeleteUser", mock.AnythingOfType("int")).Return(nil)

		userService := service.NewUserService(UserRepoMock)
		err := userService.HardDeleteUser(1)

		t.Run("test hard delete user", func(t *testing.T) {
			assert.Equal(t, true, err == nil)
		})
	})
}

func TestSetActiveUser(t *testing.T) {
	t.Run("test normal case service set active user", func(t *testing.T) {
		UserRepoMock := new(mocks.UserRepoMock)
		UserRepoMock.On("SetActiveUser", mock.AnythingOfType("int")).Return(nil)

		userService := service.NewUserService(UserRepoMock)
		_, err := userService.SetActiveUser(1)

		t.Run("test set active user", func(t *testing.T) {
			assert.Equal(t, true, err == nil)
		})
	})
}

func TestSetDeactiveUser(t *testing.T) {
	t.Run("test normal case service set deactive user", func(t *testing.T) {
		UserRepoMock := new(mocks.UserRepoMock)
		UserRepoMock.On("SetDeactiveUser", mock.AnythingOfType("int")).Return(nil)

		userService := service.NewUserService(UserRepoMock)
		_, err := userService.SetDeactiveUser(1)

		t.Run("test set deactive user", func(t *testing.T) {
			assert.Equal(t, true, err == nil)
		})
	})
}
