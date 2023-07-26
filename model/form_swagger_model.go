package model

import (
	"time"
)

type ActivationTokenForm struct {
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token" gorm:"type:varchar(64);unique"`
	Valid     time.Time `json:"valid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	DeletedBy int       `json:"deleted_by"`
}

type AttendanceLogForm struct {
	ID           uint    `json:"id" gorm:"primary_key"`
	CreatedBy    int     `json:"created_by"`
	UpdatedBy    int     `json:"updated_by"`
	DeletedBy    int     `json:"deleted_by"`
	AttendanceID uint    `json:"attendance_id"`
	LogType      string  `json:"log_type" gorm:"type:enum('clock_in','clock_out');default:'clock_in'"`
	CheckIn      int64   `json:"check_in"`
	Status       string  `json:"status" gorm:"type:varchar(255)"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	TimeZone     int     `json:"time_zone"`
	Location     string  `json:"location" gorm:"type:varchar(255)"`
}

type AttendanceForm struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	CreatedBy      int       `json:"created_by"`
	UpdatedBy      int       `json:"updated_by"`
	DeletedBy      int       `json:"deleted_by"`
	UserID         int       `json:"user_id"`
	ScheduleID     uint      `json:"schedule_id"`
	Date           time.Time `json:"date" gorm:"type:date;not null"`
	ClockIn        int64     `json:"clock_in"`
	ClockOut       int64     `json:"clock_out"`
	Status         string    `json:"status" gorm:"type:enum('-','late','come_home_early','late_and_home_early');default:'-'"`
	StatusPresence string    `json:"status_presence" gorm:"type:enum('presence','not_presence','sick','leave_attendance');default:'not_presence'"`
	LateIn         string    `json:"late_in" gorm:"type:varchar(8); default:'00:00:00'"`
	EarlyOut       string    `json:"early_out" gorm:"type:varchar(8); default:'00:00:00'"`
	LatitudeIn     float64   `json:"latitude_in"`
	LongitudeIn    float64   `json:"longitude_in"`
	TimeZoneIn     int       `json:"time_zone_in"`
	LocationIn     string    `json:"location_in" gorm:"type:varchar(255)"`
	LatitudeOut    float64   `json:"latitude_out"`
	LongitudeOut   float64   `json:"longitude_out"`
	TimeZoneOut    int       `json:"time_zone_out"`
	LocationOut    string    `json:"location_out" gorm:"type:varchar(255)"`
}

type CheckInDataForm struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	DeletedBy int       `json:"deleted_by"`
	UserID    int       `json:"user_id" query:"user_id"`
	QRCode    string    `json:"qr_code" query:"qr_code"`
	TimeZone  int       `json:"time_zone" query:"time_zone"`
	Latitude  float64   `json:"latitude" query:"latitude"`
	Longitude float64   `json:"longitude" query:"longitude"`
	Location  string    `json:"location" query:"location"`
}

type DailyScheduleForm struct {
	ID         uint      `json:"id" gorm:"primary_key"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	CreatedBy  int       `json:"created_by"`
	UpdatedBy  int       `json:"updated_by"`
	DeletedBy  int       `json:"deleted_by"`
	ScheduleID uint      `json:"schedule_id"`
	Name       string    `json:"name" gorm:"type:enum('sunday','monday','tuesday','wednesday','thursday','friday','saturday');default:'sunday'"`
	StartTime  string    `json:"start_time" gorm:"type:varchar(5)"`
	EndTime    string    `json:"end_time" gorm:"type:varchar(5)"`
	OwnerID    int       `json:"owner_id" gorm:"not null"`
}

type FacultyForm struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	DeletedBy int       `json:"deleted_by"`
	Name      string    `json:"name" gorm:"type:varchar(100)"`
	Code      string    `json:"code" gorm:"unique;type:varchar(25)"`
	Summary   string    `json:"summary" gorm:"type:text"`
	OwnerID   int       `json:"owner_id" gorm:"not null"`
}

type MajorForm struct {
	ID        uint        `json:"id" gorm:"primary_key"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	CreatedBy int         `json:"created_by"`
	UpdatedBy int         `json:"updated_by"`
	DeletedBy int         `json:"deleted_by"`
	Name      string      `json:"name" gorm:"type:varchar(100)"`
	Code      string      `json:"code" gorm:"unique;type:varchar(25)"`
	Summary   string      `json:"summary" gorm:"type:text"`
	FacultyID uint        `json:"faculty_id"`
	Faculty   FacultyForm `json:"faculty"`
	OwnerID   int         `json:"owner_id" gorm:"not null"`
}

type StudyProgramForm struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	DeletedBy int       `json:"deleted_by"`
	Name      string    `json:"name" gorm:"type:varchar(100)"`
	Code      string    `json:"code" gorm:"unique;type:varchar(25)"`
	Summary   string    `json:"summary" gorm:"type:text"`
	MajorID   uint      `json:"major_id"`
	Major     MajorForm `json:"major"`
	OwnerID   int       `json:"owner_id" gorm:"not null"`
}

type PasswordResetTokenForm struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	DeletedBy int       `json:"deleted_by"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token" gorm:"type:varchar(64);unique"`
	Valid     time.Time `json:"valid"`
}

type SubjectForm struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	DeletedBy int       `json:"deleted_by"`
	Name      string    `json:"name" gorm:"type:varchar(100)"`
	Code      string    `json:"code" gorm:"unique;type:varchar(25)"`
	Summary   string    `json:"summary" gorm:"type:text"`
	OwnerID   int       `json:"owner_id" gorm:"not null"`
}

type ScheduleForm struct {
	ID            uint                `json:"id" gorm:"primary_key"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	CreatedBy     int                 `json:"created_by"`
	UpdatedBy     int                 `json:"updated_by"`
	DeletedBy     int                 `json:"deleted_by"`
	Name          string              `json:"name" gorm:"type:varchar(100)"`
	Code          string              `json:"code" gorm:"unique;type:varchar(100)"`
	QRCode        string              `json:"qr_code" gorm:"unique;type:varchar(100)"`
	StartDate     string              `json:"start_date" gorm:"type:date"`
	EndDate       string              `json:"end_date" gorm:"type:date"`
	SubjectID     uint                `json:"subject_id"`
	Subject       SubjectForm         `json:"subject" gorm:"foreignKey:SubjectID"`
	DailySchedule []DailyScheduleForm `json:"daily_schedule" gorm:"foreignKey:ScheduleID"`
	LateDuration  int                 `json:"late_duration"` // in minute
	Latitude      float64             `json:"latitude"`
	Longitude     float64             `json:"longitude"`
	Radius        int                 `json:"radius"` //in metter
	UserInRule    int                 `json:"user_in_rule" gorm:"-"`
	OwnerID       int                 `json:"owner_id" gorm:"not null"`
	Owner         UserForm            `json:"owner"`
}

type UserForm struct {
	Username     string    `json:"username"`
	Password     string    `json:"password"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	Handphone    string    `json:"handphone"`
	Email        string    `json:"email"`
	Intro        string    `json:"intro"`
	Profile      string    `json:"profile"`
	IsActive     bool      `json:"is_active"`
	IsUser       bool      `json:"is_user"`
	IsAdmin      bool      `json:"is_admin"`
	IsSuperAdmin bool      `json:"is_super_admin"`
	LastLogin    string    `json:"last_login"`
	Role         string    `json:"role"`
	UserAbility  []Ability `json:"user_abilities"`
	Avatar       string    `json:"avatar"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CreatedBy    int       `json:"created_by"`
	UpdatedBy    int       `json:"updated_by"`
	DeletedBy    int       `json:"deleted_by"`
}

type UserForgotPasswordForm struct {
	Email           string `json:"email"`
	Activation      string `json:"activation"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserUpdatePasswordForm struct {
	ID              uint   `json:"id" query:"id"`
	CurrentPassword string `json:"current_password" query:"current_password"`
	NewPassword     string `json:"new_password" query:"new_password"`
	ConfirmPassword string `json:"confirm_password" query:"confirm_password"`
}

type UserScheduleForm struct {
	ID         uint         `json:"id" gorm:"primary_key"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	CreatedBy  int          `json:"created_by"`
	UpdatedBy  int          `json:"updated_by"`
	DeletedBy  int          `json:"deleted_by"`
	UserID     int          `json:"user_id"`
	ScheduleID uint         `json:"schedule_id"`
	Schedule   ScheduleForm `json:"schedule" gorm:"foreignKey:ScheduleID"`
	User       UserForm     `json:"user" gorm:"foreignKey:UserID"`
	OwnerID    int          `json:"owner_id" gorm:"not null"`
}

type StudentForm struct {
	UserID         uint             `json:"user_id"`
	User           UserForm         `json:"user"`
	NIM            string           `json:"nim" gorm:"type:varchar(20);unique"`
	DOB            string           `json:"dob" gorm:"type:date"`
	FacultyID      uint             `json:"faculty_id"`
	Faculty        FacultyForm      `json:"faculty"`
	MajorID        uint             `json:"major_id"`
	Major          MajorForm        `json:"major"`
	StudyProgramID uint             `json:"study_program_id"`
	StudyProgram   StudyProgramForm `json:"study_program"`
	Address        string           `json:"address" gorm:"type:varchar(255)"`
	Gender         string           `json:"gender" gorm:"type:enum('laki-laki','perempuan');default:'laki-laki'"`
}

type TeacherForm struct {
	UserID         uint             `json:"user_id"`
	User           UserForm         `json:"user"`
	Nip            string           `json:"nip" gorm:"type:varchar(20);unique"`
	DOB            string           `json:"dob" gorm:"type:date"`
	FacultyID      uint             `json:"faculty_id"`
	Faculty        FacultyForm      `json:"faculty"`
	MajorID        uint             `json:"major_id"`
	Major          MajorForm        `json:"major"`
	StudyProgramID uint             `json:"study_program_id"`
	StudyProgram   StudyProgramForm `json:"study_program"`
	Address        string           `json:"address" gorm:"type:varchar(255)"`
	Gender         string           `json:"gender" gorm:"type:enum('laki-laki','perempuan');default:'laki-laki'"`
}
