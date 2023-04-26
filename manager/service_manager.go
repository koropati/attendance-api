package manager

import (
	"sync"

	"attendance-api/infra"
	"attendance-api/service"
)

type ServiceManager interface {
	MajorService() service.MajorService
	StudyProgramService() service.StudyProgramService
	AuthService() service.AuthService
	UserService() service.UserService
	PasswordResetTokenService() service.PasswordResetTokenService
	ActivationTokenService() service.ActivationTokenService
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
	majorServiceOnce              sync.Once
	studyProgramServiceOnce       sync.Once
	authServiceOnce               sync.Once
	userServiceOnce               sync.Once
	passwordResetTokenServiceOnce sync.Once
	activationTokenServiceOnce    sync.Once
	subjectServiceOnce            sync.Once
	dailyScheduleServiceOnce      sync.Once
	scheduleServiceOnce           sync.Once
	userScheduleServiceOnce       sync.Once
	attendanceLogServiceOnce      sync.Once
	attendanceServiceOnce         sync.Once
	majorService                  service.MajorService
	studyProgramService           service.StudyProgramService
	authService                   service.AuthService
	userService                   service.UserService
	passwordResetTokenService     service.PasswordResetTokenService
	activationTokenService        service.ActivationTokenService
	subjectService                service.SubjectService
	dailyScheduleService          service.DailyScheduleService
	scheduleService               service.ScheduleService
	userScheduleService           service.UserScheduleService
	attendanceLogService          service.AttendanceLogService
	attendanceService             service.AttendanceService
)

func (sm *serviceManager) MajorService() service.MajorService {
	majorServiceOnce.Do(func() {
		majorService = sm.repo.MajorRepo()
	})
	return majorService
}

func (sm *serviceManager) StudyProgramService() service.StudyProgramService {
	studyProgramServiceOnce.Do(func() {
		studyProgramService = sm.repo.StudyProgramRepo()
	})
	return studyProgramService
}

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

func (sm *serviceManager) ActivationTokenService() service.ActivationTokenService {
	activationTokenServiceOnce.Do(func() {
		activationTokenService = sm.repo.ActivationTokenRepo()
	})

	return activationTokenService
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
