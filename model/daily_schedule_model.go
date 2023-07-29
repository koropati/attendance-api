package model

import (
	"attendance-api/common/util/converter"
	"time"
)

// tag
type DailySchedule struct {
	GormCustom
	ScheduleID uint   `json:"schedule_id" query:"schedule_id" form:"schedule_id"`
	Name       string `json:"name" gorm:"type:enum('sunday','monday','tuesday','wednesday','thursday','friday','saturday');default:'sunday'" query:"name" form:"name"`
	StartTime  string `json:"start_time" gorm:"type:varchar(5)" query:"start_time" form:"start_time"`
	EndTime    string `json:"end_time" gorm:"type:varchar(5)" query:"end_time" form:"end_time"`
	OwnerID    int    `json:"owner_id" gorm:"not null" query:"owner_id" form:"owner_id"`
}

func (dailySchedule DailySchedule) IsToday() (isToday bool) {
	now := time.Now()

	if dailySchedule.Name == converter.GetDayName(now) {
		return true
	} else {
		return false
	}
}
