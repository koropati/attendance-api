package model

type UserSchedule struct {
	GormCustom
	UserID     int      `json:"user_id"`
	ScheduleID uint     `json:"schedule_id"`
	Schedule   Schedule `json:"schedule" gorm:"foreignKey:ScheduleID"`
	User       User     `json:"user" gorm:"foreignKey:UserID"`
	OwnerID    int      `json:"owner_id" gorm:"not null"`
}
