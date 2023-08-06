package model

type Attendance struct {
	GormCustom
	UserID         int             `json:"user_id" query:"user_id" form:"user_id"`
	ScheduleID     uint            `json:"schedule_id" query:"schedule_id" form:"schedule_id"`
	User           User            `json:"user" gorm:"foreignKey:UserID" query:"user" form:"user"`
	Schedule       Schedule        `json:"schedule" gorm:"foreignKey:ScheduleID" query:"schedule" form:"schedule"`
	Date           string          `json:"date" gorm:"type:date;not null" query:"date" form:"date"`
	ClockIn        int64           `json:"clock_in" query:"clock_in" form:"clock_in"`
	ClockOut       int64           `json:"clock_out" query:"clock_out" form:"clock_out"`
	Status         string          `json:"status" gorm:"type:enum('-','late','come_home_early','late_and_home_early');default:'-'" query:"status" form:"status"`
	StatusPresence string          `json:"status_presence" gorm:"type:enum('presence','not_presence','sick','leave_attendance');default:'not_presence'" query:"status_presence" form:"status_presence"`
	LateIn         string          `json:"late_in" gorm:"type:varchar(8); default:'00:00:00'" query:"late_in" form:"late_in"`
	EarlyOut       string          `json:"early_out" gorm:"type:varchar(8); default:'00:00:00'" query:"early_out" form:"early_out"`
	LatitudeIn     float64         `json:"latitude_in" query:"latitude_in" form:"latitude_in"`
	LongitudeIn    float64         `json:"longitude_in" query:"longitude_in" form:"longitude_in"`
	TimeZoneIn     int             `json:"time_zone_in" query:"time_zone_in" form:"time_zone_in"`
	LocationIn     string          `json:"location_in" gorm:"type:varchar(255)" query:"location_in" form:"location_in"`
	LatitudeOut    float64         `json:"latitude_out" query:"latitude_out" form:"latitude_out"`
	LongitudeOut   float64         `json:"longitude_out" query:"longitude_out" form:"longitude_out"`
	TimeZoneOut    int             `json:"time_zone_out" query:"time_zone_out" form:"time_zone_out"`
	LocationOut    string          `json:"location_out" gorm:"type:varchar(255)" query:"location_out" form:"location_out"`
	AttendanceLog  []AttendanceLog `json:"attendance_log" gorm:"foreignKey:AttendanceID" query:"attendance_log" form:"attendance_log"`
}

type QuickUpdateAttendance struct {
	StatusPresence string `json:"status_presence" query:"status_presence" form:"status_presence"`
}

type CheckInData struct {
	UserID    int     `json:"user_id" query:"user_id" form:"user_id"`
	QRCode    string  `json:"qr_code" query:"qr_code" form:"qr_code"`
	TimeZone  int     `json:"time_zone" query:"time_zone" form:"time_zone"`
	Latitude  float64 `json:"latitude" query:"latitude" form:"latitude"`
	Longitude float64 `json:"longitude" query:"longitude" form:"longitude"`
	Location  string  `json:"location" query:"location" form:"location"`
}

// 'presence','not_presence','sick','leave_attendance'
type AttendanceSummary struct {
	Presence        int `json:"presence" query:"presence" form:"presence"`
	NotPresence     int `json:"not_presence" query:"not_presence" form:"not_presence"`
	Sick            int `json:"sick" query:"sick" form:"sick"`
	LeaveAttendance int `json:"leave_attendance" query:"leave_attendance" form:"leave_attendance"`
}

func (data Attendance) GenerateStatusPresence() (statusPresence string) {
	if data.StatusPresence == "" || data.StatusPresence == "not_presence" {
		if data.ClockIn > 0 || data.ClockOut > 0 {
			return "presence"
		} else {
			return "not_presence"
		}
	} else {
		return data.StatusPresence
	}
}

func (data Attendance) GenerateStatus() (status string) {
	if data.StatusPresence == "presence" {
		if (data.LateIn != "00:00:00" && data.LateIn != "") && (data.EarlyOut != "00:00:00" && data.EarlyOut != "") {
			return "late_and_home_early"
		} else if data.LateIn != "00:00:00" && data.LateIn != "" && (data.EarlyOut == "00:00:00" || data.EarlyOut == "") {
			return "late"
		} else if (data.LateIn == "00:00:00" || data.LateIn == "") && data.EarlyOut != "00:00:00" && data.EarlyOut != "" {
			return "come_home_early"
		} else {
			return "-"
		}
	} else {
		return "-"
	}
}
