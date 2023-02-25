package model

type AttendanceLog struct {
	GormCustom
	AttendanceID uint    `json:"attendance_id"`
	LogType      string  `json:"log_type" gorm:"type:enum('clock_in','clock_out');default:'clock_in'"`
	CheckIn      int64   `json:"check_in"`
	Status       string  `json:"status" gorm:"type:varchar(255)"`
	Latitude     float64 `json:"latitude"`
	Longitude    float64 `json:"longitude"`
	TimeZone     int     `json:"time_zone"`
	Location     string  `json:"location" gorm:"type:varchar(255)"`
}
