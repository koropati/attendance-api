package model

import (
	"strings"
	"time"
)

// tag
type DailySchedule struct {
	GormCustom
	ScheduleID uint   `json:"schedule_id"`
	Name       string `json:"name" gorm:"type:enum('sunday','monday','tuesday','wednesday','thursday','friday','saturday');default:'sunday'"`
	StartTime  string `json:"start_time" gorm:"type:varchar(5)"`
	EndTime    string `json:"end_time" gorm:"type:varchar(5)"`
	OwnerID    int    `json:"owner_id" gorm:"not null"`
}

func (dailySchedule DailySchedule) IsToday() (isToday bool) {
	now := time.Now()
	day := now.Weekday()

	dayName := strings.ToLower(day.String())
	if dailySchedule.Name == dayName {
		return true
	} else {
		return false
	}
}
