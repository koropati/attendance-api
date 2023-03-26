package converter

import (
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
