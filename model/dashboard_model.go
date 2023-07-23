package model

type DashboardAcademic struct {
	TotalFaculty      int `json:"total_faculty" query:"total_faculty" form:"total_faculty"`
	TotalMajor        int `json:"total_major" query:"total_major" form:"total_major"`
	TotalStudyProgram int `json:"total_study_program" query:"total_study_program" form:"total_study_program"`
	TotalSubject      int `json:"total_subject" query:"total_subject" form:"total_subject"`
	TotalSchedule     int `json:"total_schedule" query:"total_schedule" form:"total_schedule"`
}

type DashboardUser struct {
	TotalUser          int `json:"total_user" query:"total_user" form:"total_user"`
	TotalUserActive    int `json:"total_user_active" query:"total_user_active" form:"total_user_active"`
	TotalUserNonActive int `json:"total_user_non_active" query:"total_user_non_active" form:"total_user_non_active"`
	TotalSuperAdmin    int `json:"total_super_admin" query:"total_super_admin" form:"total_super_admin"`
}

type DashboardStudent struct {
	TotalStudent          int `json:"total_student" query:"total_student" form:"total_student"`
	TotalStudentActive    int `json:"total_student_active" query:"total_student_active" form:"total_student_active"`
	TotalStudentNonActive int `json:"total_student_non_active" query:"total_student_non_active" form:"total_student_non_active"`
}

type DashboardTeacher struct {
	TotalTeacher          int `json:"total_teacher" query:"total_teacher" form:"total_teacher"`
	TotalTeacherActive    int `json:"total_teacher_active" query:"total_teacher_active" form:"total_teacher_active"`
	TotalTeacherNonActive int `json:"total_teacher_non_active" query:"total_teacher_non_active" form:"total_teacher_non_active"`
}

type DashboardAttendance struct {
	Date                  string `json:"date" query:"date" form:"date"`
	MonthPeriod           int    `json:"month_period" query:"month_period" form:"month_period"`
	YearPeriod            int    `json:"year_period" query:"year_period" form:"year_period"`
	TotalPresence         int    `json:"total_presence" query:"total_presence" form:"total_presence"`
	TotalNotPresence      int    `json:"total_not_presence" query:"total_not_presence" form:"total_not_presence"`
	TotalSick             int    `json:"total_sick" query:"total_sick" form:"total_sick"`
	TotalLeaveAttendance  int    `json:"total_leave_attendance" query:"total_leave_attendance" form:"total_leave_attendance"`
	TotalNoClockIn        int    `json:"total_no_clock_in" query:"total_no_clock_in" form:"total_no_clock_in"`
	TotalNoClockOut       int    `json:"total_no_clock_out" query:"total_no_clock_out" form:"total_no_clock_out"`
	TotalLate             int    `json:"total_late" query:"total_late" form:"total_late"`
	TotalComeHomeEarly    int    `json:"total_come_home_early" query:"total_come_home_early" form:"total_come_home_early"`
	TotalLateAndHomeEarly int    `json:"total_late_and_home_early" query:"total_late_and_home_early" form:"total_late_and_home_early"`
}

type AttendanceSeries struct {
	Name        string   `json:"name" query:"name" form:"name"`
	Data        []int    `json:"data" query:"data" form:"data"`
	Date        []string `json:"date" query:"date" form:"date"`
	MonthPeriod int      `json:"month_period" query:"month_period" form:"month_period"`
	YearPeriod  int      `json:"year_period" query:"year_period" form:"year_period"`
}
