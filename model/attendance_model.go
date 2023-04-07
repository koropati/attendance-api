package model

import "time"

type Attendance struct {
	GormCustom
	UserID         int       `json:"user_id"`
	ScheduleID     uint      `json:"schedule_id"`
	Date           time.Time `json:"date" gorm:"type:date;not null"`
	ClockIn        int64     `json:"clock_in"`
	ClockOut       int64     `json:"clock_out"`
	Status         string    `json:"status" gorm:"type:varchar(255)"`
	StatusPresence string    `json:"status_presence" gorm:"type:varchar(255)"`
	LateIn         string    `json:"late_in" gorm:"type:varchar(5)"`
	EarlyOut       string    `json:"early_out" gorm:"type:varchar(5)"`
	LatitudeIn     float64   `json:"latitude_in"`
	LongitudeIn    float64   `json:"longitude_in"`
	TimeZoneIn     int       `json:"time_zone_in"`
	LocationIn     string    `json:"location_in" gorm:"type:varchar(255)"`
	LatitudeOut    float64   `json:"latitude_out"`
	LongitudeOut   float64   `json:"longitude_out"`
	TimeZoneOut    int       `json:"time_zone_out"`
	LocationOut    string    `json:"location_out"`
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
