package model

type ActivationTokenResponseData struct {
	Code    int                 `json:"code"`
	Data    ActivationTokenForm `json:"data"`
	Message string              `json:"message"`
}

type ActivationTokenResponseList struct {
	Code    int                   `json:"code"`
	Data    []ActivationTokenForm `json:"data"`
	Meta    Meta                  `json:"meta"`
	Message string                `json:"message"`
}

type AttendanceLogResponseData struct {
	Code    int               `json:"code"`
	Data    AttendanceLogForm `json:"data"`
	Message string            `json:"message"`
}

type AttendanceLogResponseList struct {
	Code    int                 `json:"code"`
	Data    []AttendanceLogForm `json:"data"`
	Meta    Meta                `json:"meta"`
	Message string              `json:"message"`
}

type AttendanceResponseData struct {
	Code    int            `json:"code"`
	Data    AttendanceForm `json:"data"`
	Message string         `json:"message"`
}

type AttendanceResponseList struct {
	Code    int              `json:"code"`
	Data    []AttendanceForm `json:"data"`
	Meta    Meta             `json:"meta"`
	Message string           `json:"message"`
}

type CheckInDataResponseData struct {
	Code    int             `json:"code"`
	Data    CheckInDataForm `json:"data"`
	Message string          `json:"message"`
}

type CheckInDataResponseList struct {
	Code    int               `json:"code"`
	Data    []CheckInDataForm `json:"data"`
	Meta    Meta              `json:"meta"`
	Message string            `json:"message"`
}

type DailyScheduleResponseData struct {
	Code    int               `json:"code"`
	Data    DailyScheduleForm `json:"data"`
	Message string            `json:"message"`
}

type DailyScheduleResponseList struct {
	Code    int                 `json:"code"`
	Data    []DailyScheduleForm `json:"data"`
	Meta    Meta                `json:"meta"`
	Message string              `json:"message"`
}

type FacultyResponseData struct {
	Code    int         `json:"code"`
	Data    FacultyForm `json:"data"`
	Message string      `json:"message"`
}

type FacultyResponseList struct {
	Code    int           `json:"code"`
	Data    []FacultyForm `json:"data"`
	Meta    Meta          `json:"meta"`
	Message string        `json:"message"`
}

type MajorResponseData struct {
	Code    int       `json:"code"`
	Data    MajorForm `json:"data"`
	Message string    `json:"message"`
}

type MajorResponseList struct {
	Code    int         `json:"code"`
	Data    []MajorForm `json:"data"`
	Meta    Meta        `json:"meta"`
	Message string      `json:"message"`
}

type PasswordResetTokenResponseData struct {
	Code    int                    `json:"code"`
	Data    PasswordResetTokenForm `json:"data"`
	Message string                 `json:"message"`
}

type PasswordResetTokenResponseList struct {
	Code    int                      `json:"code"`
	Data    []PasswordResetTokenForm `json:"data"`
	Meta    Meta                     `json:"meta"`
	Message string                   `json:"message"`
}

type StudyProgramResponseData struct {
	Code    int              `json:"code"`
	Data    StudyProgramForm `json:"data"`
	Message string           `json:"message"`
}

type StudyProgramResponseList struct {
	Code    int                `json:"code"`
	Data    []StudyProgramForm `json:"data"`
	Meta    Meta               `json:"meta"`
	Message string             `json:"message"`
}

type SubjectResponseData struct {
	Code    int         `json:"code"`
	Data    SubjectForm `json:"data"`
	Message string      `json:"message"`
}

type SubjectResponseList struct {
	Code    int           `json:"code"`
	Data    []SubjectForm `json:"data"`
	Meta    Meta          `json:"meta"`
	Message string        `json:"message"`
}

type ScheduleResponseData struct {
	Code    int          `json:"code"`
	Data    ScheduleForm `json:"data"`
	Message string       `json:"message"`
}

type ScheduleResponseList struct {
	Code    int            `json:"code"`
	Data    []ScheduleForm `json:"data"`
	Meta    Meta           `json:"meta"`
	Message string         `json:"message"`
}

type UserResponseData struct {
	Code    int      `json:"code"`
	Data    UserForm `json:"data"`
	Message string   `json:"message"`
}

type UserResponseList struct {
	Code    int        `json:"code"`
	Data    []UserForm `json:"data"`
	Meta    Meta       `json:"meta"`
	Message string     `json:"message"`
}

type StudentResponseData struct {
	Code    int         `json:"code"`
	Data    StudentForm `json:"data"`
	Message string      `json:"message"`
}

type StudentResponseList struct {
	Code    int           `json:"code"`
	Data    []StudentForm `json:"data"`
	Meta    Meta          `json:"meta"`
	Message string        `json:"message"`
}

type TeacherResponseData struct {
	Code    int         `json:"code"`
	Data    TeacherForm `json:"data"`
	Message string      `json:"message"`
}

type TeacherResponseList struct {
	Code    int           `json:"code"`
	Data    []TeacherForm `json:"data"`
	Meta    Meta          `json:"meta"`
	Message string        `json:"message"`
}

type UserForgotPasswordResponseData struct {
	Code    int                    `json:"code"`
	Data    UserForgotPasswordForm `json:"data"`
	Message string                 `json:"message"`
}

type UserForgotPasswordResponseList struct {
	Code    int                      `json:"code"`
	Data    []UserForgotPasswordForm `json:"data"`
	Meta    Meta                     `json:"meta"`
	Message string                   `json:"message"`
}

type UserUpdatePasswordResponseData struct {
	Code    int                    `json:"code"`
	Data    UserUpdatePasswordForm `json:"data"`
	Message string                 `json:"message"`
}

type UserUpdatePasswordResponseList struct {
	Code    int                      `json:"code"`
	Data    []UserUpdatePasswordForm `json:"data"`
	Meta    Meta                     `json:"meta"`
	Message string                   `json:"message"`
}

type UserScheduleResponseData struct {
	Code    int              `json:"code"`
	Data    UserScheduleForm `json:"data"`
	Message string           `json:"message"`
}

type UserScheduleResponseList struct {
	Code    int                `json:"code"`
	Data    []UserScheduleForm `json:"data"`
	Message string             `json:"message"`
}

type MyScheduleResponseData struct {
	Code    int        `json:"code"`
	Data    MySchedule `json:"data"`
	Message string     `json:"message"`
}

type MyScheduleResponseList struct {
	Code    int          `json:"code"`
	Data    []MySchedule `json:"data"`
	Message string       `json:"message"`
}

type TodayScheduleResponseData struct {
	Code    int           `json:"code"`
	Data    TodaySchedule `json:"data"`
	Message string        `json:"message"`
}

type TodayScheduleResponseList struct {
	Code    int             `json:"code"`
	Data    []TodaySchedule `json:"data"`
	Message string          `json:"message"`
}

type AuthDataResponseData struct {
	UserData  UserForm  `json:"user_data"`
	TokenData TokenData `json:"token_data"`
}
