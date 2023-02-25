package model

import "time"

// tag
type Schedule struct {
	GormCustom
	Name          string          `json:"name" gorm:"type:varchar(100)"`
	Code          string          `json:"code" gorm:"unique;type:varchar(100)"`
	StartDate     time.Time       `json:"start_date" gorm:"type:date"`
	EndDate       time.Time       `json:"end_date" gorm:"type:date"`
	SubjectID     uint            `json:"subject_id"`
	Subject       Subject         `json:"subject" gorm:"foreignKey:SubjectID"`
	DailySchedule []DailySchedule `json:"daily_schedule" gorm:"foreignKey:ScheduleID"`
	LateDuration  int             `json:"late_duration"`
	OwnerID       int             `json:"owner_id" gorm:"not null"`
}
