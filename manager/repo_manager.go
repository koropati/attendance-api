package manager

import (
	"sync"

	"attendance-api/infra"
	"attendance-api/repo"
)

type RepoManager interface {
	AuthRepo() repo.AuthRepo
	UserRepo() repo.UserRepo
	PasswordResetTokenRepo() repo.PasswordResetTokenRepo
	ActivationTokenRepo() repo.ActivationTokenRepo
	SubjectRepo() repo.SubjectRepo
	DailyScheduleRepo() repo.DailyScheduleRepo
	ScheduleRepo() repo.ScheduleRepo
	UserScheduleRepo() repo.UserScheduleRepo
	AttendanceLogRepo() repo.AttendanceLogRepo
	AttendanceRepo() repo.AttendanceRepo
	MajorRepo() repo.MajorRepo
	StudyProgramRepo() repo.StudyProgramRepo
}

type repoManager struct {
	infra infra.Infra
}

func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra: infra}
}

var (
	majorRepoOnce              sync.Once
	studyProgramRepoOnce       sync.Once
	authRepoOnce               sync.Once
	userRepoOnce               sync.Once
	passwordResetTokenRepoOnce sync.Once
	activationTokenRepoOnce    sync.Once
	subjectRepoOnce            sync.Once
	dailyScheduleRepoOnce      sync.Once
	scheduleRepoOnce           sync.Once
	userScheduleRepoOnce       sync.Once
	attendanceLogRepoOnce      sync.Once
	attendanceRepoOnce         sync.Once
	majorRepo                  repo.MajorRepo
	studyProgramRepo           repo.StudyProgramRepo
	authRepo                   repo.AuthRepo
	userRepo                   repo.UserRepo
	passwordResetTokenRepo     repo.PasswordResetTokenRepo
	activationTokenRepo        repo.ActivationTokenRepo
	subjectRepo                repo.SubjectRepo
	dailyScheduleRepo          repo.DailyScheduleRepo
	scheduleRepo               repo.ScheduleRepo
	userScheduleRepo           repo.UserScheduleRepo
	attendanceLogRepo          repo.AttendanceLogRepo
	attendanceRepo             repo.AttendanceRepo
)

func (rm *repoManager) MajorRepo() repo.MajorRepo {
	majorRepoOnce.Do(func() {
		majorRepo = repo.NewMajorRepo(rm.infra.GormDB())
	})
	return majorRepo
}

func (rm *repoManager) StudyProgramRepo() repo.StudyProgramRepo {
	studyProgramRepoOnce.Do(func() {
		studyProgramRepo = repo.NewStudyProgramRepo(rm.infra.GormDB())
	})
	return studyProgramRepo
}

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

func (rm *repoManager) PasswordResetTokenRepo() repo.PasswordResetTokenRepo {
	passwordResetTokenRepoOnce.Do(func() {
		passwordResetTokenRepo = repo.NewPasswordResetTokenRepo(rm.infra.GormDB())
	})
	return passwordResetTokenRepo
}

func (rm *repoManager) ActivationTokenRepo() repo.ActivationTokenRepo {
	activationTokenRepoOnce.Do(func() {
		activationTokenRepo = repo.NewActivationTokenRepo(rm.infra.GormDB())
	})
	return activationTokenRepo
}

func (rm *repoManager) SubjectRepo() repo.SubjectRepo {
	subjectRepoOnce.Do(func() {
		subjectRepo = repo.NewSubjectRepo(rm.infra.GormDB())
	})
	return subjectRepo
}

func (rm *repoManager) DailyScheduleRepo() repo.DailyScheduleRepo {
	dailyScheduleRepoOnce.Do(func() {
		dailyScheduleRepo = repo.NewDailyScheduleRepo(rm.infra.GormDB())
	})
	return dailyScheduleRepo
}

func (rm *repoManager) ScheduleRepo() repo.ScheduleRepo {
	scheduleRepoOnce.Do(func() {
		scheduleRepo = repo.NewScheduleRepo(rm.infra.GormDB())
	})
	return scheduleRepo
}

func (rm *repoManager) UserScheduleRepo() repo.UserScheduleRepo {
	userScheduleRepoOnce.Do(func() {
		userScheduleRepo = repo.NewUserScheduleRepo(rm.infra.GormDB())
	})
	return userScheduleRepo
}

func (rm *repoManager) AttendanceLogRepo() repo.AttendanceLogRepo {
	attendanceLogRepoOnce.Do(func() {
		attendanceLogRepo = repo.NewAttendanceLogRepo(rm.infra.GormDB())
	})
	return attendanceLogRepo
}

func (rm *repoManager) AttendanceRepo() repo.AttendanceRepo {
	attendanceRepoOnce.Do(func() {
		attendanceRepo = repo.NewAttendanceRepo(rm.infra.GormDB())
	})
	return attendanceRepo
}
