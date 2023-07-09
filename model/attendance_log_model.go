package model

type AttendanceLog struct {
	GormCustom
	AttendanceID uint    `json:"attendance_id" query:"attendance_id" form:"attendance_id"`
	LogType      string  `json:"log_type" gorm:"type:enum('clock_in','clock_out');default:'clock_in'" query:"log_type" form:"log_type"`
	CheckIn      int64   `json:"check_in" query:"check_in" form:"check_in"`
	Status       string  `json:"status" gorm:"type:varchar(255)" query:"status" form:"status"`
	Latitude     float64 `json:"latitude" query:"latitude" form:"latitude"`
	Longitude    float64 `json:"longitude" query:"longitude" form:"longitude"`
	TimeZone     int     `json:"time_zone" query:"time_zone" form:"time_zone"`
	Location     string  `json:"location" gorm:"type:varchar(255)" query:"location" form:"location"`
}
