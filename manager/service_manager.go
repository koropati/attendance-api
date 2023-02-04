package manager

import (
	"sync"

	"attendance-api/infra"
	"attendance-api/service"
)

type ServiceManager interface {
	AuthService() service.AuthService
	UserService() service.UserService
}

type serviceManager struct {
	infra infra.Infra
	repo  RepoManager
}

func NewServiceManager(infra infra.Infra) ServiceManager {
	return &serviceManager{
		infra: infra,
		repo:  NewRepoManager(infra),
	}
}

var (
	authServiceOnce sync.Once
	authService     service.AuthService
	userServiceOnce sync.Once
	userService     service.UserService
)

func (sm *serviceManager) AuthService() service.AuthService {
	authServiceOnce.Do(func() {
		authService = sm.repo.AuthRepo()
	})

	return authService
}

func (sm *serviceManager) UserService() service.UserService {
	userServiceOnce.Do(func() {
		userService = sm.repo.UserRepo()
	})

	return userService
}
