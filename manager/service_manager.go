package manager

import (
	"sync"

	"attendance-api/infra"
	"attendance-api/service"
)

type ServiceManager interface {
	AuthService() service.AuthService
	UserService() service.UserService
	PasswordResetTokenService() service.PasswordResetTokenService
	SubjectService() service.SubjectService
	DailyScheduleService() service.DailyScheduleService
	ScheduleService() service.ScheduleService
	UserScheduleService() service.UserScheduleService
	AttendanceLogService() service.AttendanceLogService
	AttendanceService() service.AttendanceService
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
	authServiceOnce               sync.Once
	userServiceOnce               sync.Once
	passwordResetTokenServiceOnce sync.Once
	subjectServiceOnce            sync.Once
	dailyScheduleServiceOnce      sync.Once
	scheduleServiceOnce           sync.Once
	userScheduleServiceOnce       sync.Once
	attendanceLogServiceOnce      sync.Once
	attendanceServiceOnce         sync.Once
	authService                   service.AuthService
	userService                   service.UserService
	passwordResetTokenService     service.PasswordResetTokenService
	subjectService                service.SubjectService
	dailyScheduleService          service.DailyScheduleService
	scheduleService               service.ScheduleService
	userScheduleService           service.UserScheduleService
	attendanceLogService          service.AttendanceLogService
	attendanceService             service.AttendanceService
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

func (sm *serviceManager) PasswordResetTokenService() service.PasswordResetTokenService {
	passwordResetTokenServiceOnce.Do(func() {
		passwordResetTokenService = sm.repo.PasswordResetTokenRepo()
	})

	return passwordResetTokenService
}

func (sm *serviceManager) SubjectService() service.SubjectService {
	subjectServiceOnce.Do(func() {
		subjectService = sm.repo.SubjectRepo()
	})
	return subjectService
}

func (sm *serviceManager) DailyScheduleService() service.DailyScheduleService {
	dailyScheduleServiceOnce.Do(func() {
		dailyScheduleService = sm.repo.DailyScheduleRepo()
	})
	return dailyScheduleService
}

func (sm *serviceManager) ScheduleService() service.ScheduleService {
	scheduleServiceOnce.Do(func() {
		scheduleService = sm.repo.ScheduleRepo()
	})
	return scheduleService
}

func (sm *serviceManager) UserScheduleService() service.UserScheduleService {
	userScheduleServiceOnce.Do(func() {
		userScheduleService = sm.repo.UserScheduleRepo()
	})
	return userScheduleService
}

func (sm *serviceManager) AttendanceLogService() service.AttendanceLogService {
	attendanceLogServiceOnce.Do(func() {
		attendanceLogService = sm.repo.AttendanceLogRepo()
	})
	return attendanceLogService
}

func (sm *serviceManager) AttendanceService() service.AttendanceService {
	attendanceServiceOnce.Do(func() {
		attendanceService = sm.repo.AttendanceRepo()
	})
	return attendanceService
}
