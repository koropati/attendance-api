package model

import (
	"time"
)

type Attendance struct {
	GormCustom
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

type CheckInData struct {
	UserID    int     `json:"user_id" query:"user_id"`
	QRCode    string  `json:"qr_code" query:"qr_code"`
	TimeZone  int     `json:"time_zone" query:"time_zone"`
	Latitude  float64 `json:"latitude" query:"latitude"`
	Longitude float64 `json:"longitude" query:"longitude"`
	Location  string  `json:"location" query:"location"`
}

func (data Attendance) GenerateStatusPresence() (statusPresence string) {
	if data.StatusPresence == "" {
		if data.ClockIn > 0 || data.ClockOut > 0 {
			return "presence"
		} else {
			return "not_presence"
		}
	} else {
		return data.StatusPresence
	}
}

func (data Attendance) GenerateStatus() (status string) {
	if data.StatusPresence == "presence" {
		if (data.LateIn != "00:00:00" && data.LateIn != "") && (data.EarlyOut != "00:00:00" && data.EarlyOut != "") {
			return "late_and_home_early"
		} else if data.LateIn != "00:00:00" && data.LateIn != "" && (data.EarlyOut == "00:00:00" || data.EarlyOut == "") {
			return "late"
		} else if (data.LateIn == "00:00:00" || data.LateIn == "") && data.EarlyOut != "00:00:00" && data.EarlyOut != "" {
			return "come_home_early"
		} else {
			return "-"
		}
	} else {
		return "-"
	}
}
