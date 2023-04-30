package converter

import (
	"fmt"
	"time"

	"github.com/zsefvlol/timezonemapper"
)

func GetTimeZone(latitude, longitude float64) (timeZoneCode int) {
	timezone := timezonemapper.LatLngToTimezoneString(latitude, longitude)
	loc, _ := time.LoadLocation(timezone)
	// Parse time string with location
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), loc)
	_, offset := t.Zone()
	hoursGMT := int(offset / 3600)
	return hoursGMT
}

func MillisToTimeString(timeMillis int64, timeZone int) (timeString string) {
	if timeMillis <= 0 {
		return "--:--"
	}
	if timeZone == 0 {
		dateTime := time.Unix(0, timeMillis*int64(time.Millisecond))
		return dateTime.Format("15:04")
	} else {
		millisTimeZone := int64(timeZone) * 3600000
		millisFinal := timeMillis + millisTimeZone
		dateTime := time.Unix(0, millisFinal*int64(time.Millisecond))
		return dateTime.UTC().Format("15:04")
	}
}

func FormatDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func GetDayName(myTime time.Time) (dayName string) {
	days := []string{"sunday", "monday", "tuesday", "wednesday", "thursday", "friday", "saturday"}

	return days[int(myTime.Weekday())]
}

func MonthInterval(y int, m time.Month) (firstDay, lastDay time.Time) {
	firstDay = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	lastDay = time.Date(y, m+1, 1, 0, 0, 0, -1, time.UTC)
	return firstDay, lastDay
}
