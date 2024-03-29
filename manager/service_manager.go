package manager

import (
	"sync"

	"attendance-api/infra"
	"attendance-api/service"
)

type ServiceManager interface {
	FacultyService() service.FacultyService
	MajorService() service.MajorService
	StudyProgramService() service.StudyProgramService
	AuthService() service.AuthService
	UserService() service.UserService
	StudentService() service.StudentService
	TeacherService() service.TeacherService
	PasswordResetTokenService() service.PasswordResetTokenService
	ActivationTokenService() service.ActivationTokenService
	SubjectService() service.SubjectService
	DailyScheduleService() service.DailyScheduleService
	ScheduleService() service.ScheduleService
	UserScheduleService() service.UserScheduleService
	AttendanceLogService() service.AttendanceLogService
	AttendanceService() service.AttendanceService
	DashboardService() service.DashboardService
	RoleAbilityService() service.RoleAbilityService
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
	facultyServiceOnce            sync.Once
	majorServiceOnce              sync.Once
	studyProgramServiceOnce       sync.Once
	authServiceOnce               sync.Once
	userServiceOnce               sync.Once
	studentServiceOnce            sync.Once
	teacherServiceOnce            sync.Once
	passwordResetTokenServiceOnce sync.Once
	activationTokenServiceOnce    sync.Once
	subjectServiceOnce            sync.Once
	dailyScheduleServiceOnce      sync.Once
	scheduleServiceOnce           sync.Once
	userScheduleServiceOnce       sync.Once
	attendanceLogServiceOnce      sync.Once
	attendanceServiceOnce         sync.Once
	dashboardServiceOnce          sync.Once
	roleAbilityServiceOnce        sync.Once
	facultyService                service.FacultyService
	majorService                  service.MajorService
	studyProgramService           service.StudyProgramService
	authService                   service.AuthService
	userService                   service.UserService
	studentService                service.StudentService
	teacherService                service.TeacherService
	passwordResetTokenService     service.PasswordResetTokenService
	activationTokenService        service.ActivationTokenService
	subjectService                service.SubjectService
	dailyScheduleService          service.DailyScheduleService
	scheduleService               service.ScheduleService
	userScheduleService           service.UserScheduleService
	attendanceLogService          service.AttendanceLogService
	attendanceService             service.AttendanceService
	dashboardService              service.DashboardService
	roleAbilityService            service.RoleAbilityService
)

func (sm *serviceManager) FacultyService() service.FacultyService {
	facultyServiceOnce.Do(func() {
		facultyService = sm.repo.FacultyRepo()
	})
	return facultyService
}

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

func (sm *serviceManager) StudentService() service.StudentService {
	studentServiceOnce.Do(func() {
		studentService = sm.repo.StudentRepo()
	})
	return studentService
}

func (sm *serviceManager) TeacherService() service.TeacherService {
	teacherServiceOnce.Do(func() {
		teacherService = sm.repo.TeacherRepo()
	})
	return teacherService
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

func (sm *serviceManager) DashboardService() service.DashboardService {
	dashboardServiceOnce.Do(func() {
		dashboardService = sm.repo.DashboardRepo()
	})
	return dashboardService
}

func (sm *serviceManager) RoleAbilityService() service.RoleAbilityService {
	roleAbilityServiceOnce.Do(func() {
		roleAbilityService = sm.repo.RoleAbilityRepo()
	})
	return roleAbilityService
}
