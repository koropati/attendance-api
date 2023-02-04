package manager

import (
	"sync"

	"attendance-api/infra"
	"attendance-api/repo"
)

type RepoManager interface {
	AuthRepo() repo.AuthRepo
	UserRepo() repo.UserRepo
}

type repoManager struct {
	infra infra.Infra
}

func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra: infra}
}

var (
	authRepoOnce sync.Once
	authRepo     repo.AuthRepo
	userRepoOnce sync.Once
	userRepo     repo.UserRepo
)

func (rm *repoManager) AuthRepo() repo.AuthRepo {
	authRepoOnce.Do(func() {
		authRepo = repo.NewAuthRepo(rm.infra.GormDB())
	})

	return authRepo
}

func (rm *repoManager) UserRepo() repo.UserRepo {
	userRepoOnce.Do(func() {
		userRepo = repo.NewUserRepo(rm.infra.GormDB())
	})

	return userRepo
}
