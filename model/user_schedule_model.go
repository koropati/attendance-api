package model

type UserSchedule struct {
	GormCustom
	UserID     int      `json:"user_id" query:"user_id" form:"user_id"`
	ScheduleID uint     `json:"schedule_id" query:"schedule_id" form:"schedule_id"`
	Schedule   Schedule `json:"schedule" gorm:"foreignKey:ScheduleID"`
	User       User     `json:"user" gorm:"foreignKey:UserID"`
	OwnerID    int      `json:"owner_id" gorm:"not null" query:"owner_id" form:"owner_id"`
}

type MySchedule struct {
	ScheduleID   uint    `json:"schedule_id" query:"schedule_id"`
	ScheduleName string  `json:"schedule_name" query:"schedule_name"`
	ScheduleCode string  `json:"schedule_code" query:"schedule_code"`
	StartDate    string  `json:"start_date" query:"start_date"`
	EndDate      string  `json:"end_date" query:"end_date"`
	SubjectID    uint    `json:"subject_id" query:"subject_id"`
	SubjectName  string  `json:"subject_name" query:"subject_name"`
	SubjectCode  string  `json:"subject_code" query:"subject_code"`
	LateDuration int     `json:"late_duration" query:"late_duration"`
	Latitude     float64 `json:"latitude" query:"latitude"`
	Longitude    float64 `json:"longitude" query:"longitude"`
	Radius       int     `json:"radius" query:"radius"`
}

type ListMySchedule struct {
	IndeonesianDate string       `json:"indonesian_date" query:"indonesian_date"`
	Schedules       []MySchedule `json:"schedules" query:"schedules"`
}

type MyScheduleFilter struct {
	Month string `json:"month" query:"month" form:"month"`
	Year  string `json:"year" query:"year" form:"year"`
}

type TodaySchedule struct {
	ScheduleID     uint   `json:"schedule_id"`
	ScheduleName   string `json:"schedule_name"`
	ScheduleCode   string `json:"schedule_code"`
	SubjectID      uint   `json:"subject_id"`
	SubjectName    string `json:"subject_name"`
	StartTime      string `json:"start_time"`
	EndTime        string `json:"end_time"`
	AttendanceID   uint   `json:"attendance_id"`
	ClockInMillis  int64  `json:"clock_in_millis"`
	ClockOutMillis int64  `json:"clock_out_millis"`
	ClockIn        string `json:"clock_in"`
	ClockOut       string `json:"clock_out"`
	TimeZoneIn     int    `json:"time_zone_in"`
	TimeZoneOut    int    `json:"time_zone_out"`
	LocationIn     string `json:"location_in"`
	LocationOut    string `json:"location_out"`
}
