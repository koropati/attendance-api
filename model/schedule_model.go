package model

import (
	"attendance-api/common/util/converter"
	"fmt"
	"math"
	"time"
)

// tag
type Schedule struct {
	GormCustom
	Name          string          `json:"name" gorm:"type:varchar(100)" query:"name" form:"name"`
	Code          string          `json:"code" gorm:"unique;type:varchar(100)" query:"code" form:"code"`
	QRCode        string          `json:"qr_code" gorm:"unique;type:varchar(100)" query:"qr_code" form:"qr_code"`
	StartDate     string          `json:"start_date" gorm:"type:date" query:"start_date" form:"start_date"`
	EndDate       string          `json:"end_date" gorm:"type:date" query:"end_date" form:"end_date"`
	SubjectID     uint            `json:"subject_id" query:"subject_id" form:"subject_id"`
	Subject       Subject         `json:"subject" gorm:"foreignKey:SubjectID" query:"subject" form:"subject"`
	DailySchedule []DailySchedule `json:"daily_schedule" gorm:"foreignKey:ScheduleID" query:"daily_schedule" form:"daily_schedule"`
	LateDuration  int             `json:"late_duration" query:"late_duration" form:"late_duration"` // in minute
	Latitude      float64         `json:"latitude" query:"latitude" form:"latitude"`
	Longitude     float64         `json:"longitude" query:"longitude" form:"longitude"`
	Radius        int             `json:"radius" query:"radius" form:"radius"` //in metter
	UserInRule    int             `json:"user_in_rule" gorm:"-" query:"user_in_rule" form:"user_in_rule"`
	OwnerID       int             `json:"owner_id" gorm:"not null" query:"owner_id" form:"owner_id"`
	Owner         User            `json:"owner" gorm:"foreignKey:OwnerID" query:"owner" form:"owner"`
}

func (data Schedule) IsTodaySchedule() (isTodaySchedule bool) {
	inRange := false
	isToday := false
	today := time.Now()
	startDate, err := time.Parse("2006-01-02", converter.GetOnlyDateString(data.StartDate))
	if err != nil {
		fmt.Println("Format tanggal mulai tidak valid")
		return false
	}
	endDate, err := time.Parse("2006-01-02", converter.GetOnlyDateString(data.EndDate))
	if err != nil {
		fmt.Println("Format tanggal selesai tidak valid")
		return false
	}

	if today.After(startDate) && today.Before(endDate) {
		inRange = true
		for _, daily := range data.DailySchedule {
			if daily.IsToday() {
				isToday = true
				break
			}
		}
		if inRange && isToday {
			return true
		} else {
			return false
		}
	} else {
		return false
	}

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
