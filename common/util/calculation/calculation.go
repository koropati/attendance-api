package calculation

import (
	"attendance-api/common/util/converter"
	"log"
	"time"
)

func CalculateLateDuration(scheduleIn string, clockIn int64, timeZoneIn int, lateTolerance int) (reuslt string) {
	timeStringIn := converter.MillisToTimeString(clockIn, timeZoneIn)
	timeIn, err := time.Parse("2006-01-02 15:04", "2006-01-02"+" "+timeStringIn)
	if err != nil {
		log.Printf("Err time parse (in): %v\n", err)
		return "00:00:00"
	}
	timeSchedule, err := time.Parse("2006-01-02 15:04", "2006-01-02"+" "+scheduleIn)
	if err != nil {
		log.Printf("Err time parse (schedule) : %v\n", err)
		return "00:00:00"
	}
	if lateTolerance > 0 {
		timeSchedule.Add(-time.Minute * time.Duration(lateTolerance))
	}
	diff := timeIn.Sub(timeSchedule)

	if diff <= 0 {
		return "00:00:00"
	} else {
		return converter.FormatDuration(diff)
	}
}

func CalculateEarlyDuration(scheduleOut string, clockOut int64, timeZoneOut int) (result string) {
	timeStringOut := converter.MillisToTimeString(clockOut, timeZoneOut)
	timeOut, err := time.Parse("2006-01-02 15:04", "2006-01-02"+" "+timeStringOut)
	if err != nil {
		log.Printf("Err time parse (in): %v\n", err)
		return "00:00:00"
	}
	timeSchedule, err := time.Parse("2006-01-02 15:04", "2006-01-02"+" "+scheduleOut)
	if err != nil {
		log.Printf("Err time parse (schedule) : %v\n", err)
		return "00:00:00"
	}
	diff := timeSchedule.Sub(timeOut)
	if diff <= 0 {
		return "00:00:00"
	} else {
		return converter.FormatDuration(diff)
	}
}
