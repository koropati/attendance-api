package model

type UserSchedule struct {
	GormCustom
	UserID     int      `json:"user_id"`
	ScheduleID uint     `json:"schedule_id"`
	Schedule   Schedule `json:"schedule" gorm:"foreignKey:ScheduleID"`
	User       User     `json:"user" gorm:"foreignKey:UserID"`
	OwnerID    int      `json:"owner_id" gorm:"not null"`
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

type TodaySchedule struct {
	ScheduleID   uint   `json:"schedule_id"`
	ScheduleName string `json:"schedule_name"`
	ScheduleCode string `json:"schedule_code"`
	SubjectID    uint   `json:"subject_id"`
	SubjectName  string `json:"subject_name"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}
