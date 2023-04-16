package model

import (
	"math"
	"time"
)

// tag
type Schedule struct {
	GormCustom
	Name          string          `json:"name" gorm:"type:varchar(100)"`
	Code          string          `json:"code" gorm:"unique;type:varchar(100)"`
	QRCode        string          `json:"qr_code" gorm:"unique;type:varchar(100)"`
	StartDate     time.Time       `json:"start_date" gorm:"type:date"`
	EndDate       time.Time       `json:"end_date" gorm:"type:date"`
	SubjectID     uint            `json:"subject_id"`
	Subject       Subject         `json:"subject" gorm:"foreignKey:SubjectID"`
	DailySchedule []DailySchedule `json:"daily_schedule" gorm:"foreignKey:ScheduleID"`
	LateDuration  int             `json:"late_duration"` // in minute
	Latitude      float64         `json:"latitude"`
	Longitude     float64         `json:"longitude"`
	Radius        int             `json:"radius"` //in metter
	OwnerID       int             `json:"owner_id" gorm:"not null"`
}

func (data Schedule) InRange(latitudeCheck float64, longitudeCheck float64) (isPassed bool) {
	if data.Radius > 0 && data.Latitude > 0 && data.Longitude > 0 {
		radlat1 := float64(math.Pi * data.Latitude / 180)
		radlat2 := float64(math.Pi * latitudeCheck / 180)

		theta := float64(data.Longitude - longitudeCheck)
		radtheta := float64(math.Pi * theta / 180)

		dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
		if dist > 1 {
			dist = 1
		}

		dist = math.Acos(dist)
		dist = dist * 180 / math.Pi
		dist = dist * 60 * 1.1515
		// Distance in kilometer
		dist = dist * 1.609344
		// Distance in meter
		dist = dist * 1000
		if dist <= float64(data.Radius) {
			return true
		} else {
			return false
		}
	} else {
		return true
	}

}
