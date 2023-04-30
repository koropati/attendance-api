package main

import (
	"attendance-api/api"
	_ "attendance-api/docs"
	"attendance-api/infra"
	"attendance-api/model"
)

// @title           Attendance API
// @version         1.0
// @description     A Attendance management service API in Go using Gin framework.
// @termsOfService  https://wokdev.com

// @contact.name   Admin CS
// @contact.url    https://attendance.wokdev.com
// @contact.email  wokdev@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @securityDefinitions.apikey BearerTokenAuth
// @in header
// @name Authorization
// @host      localhost:3000
// @BasePath  /v1

func main() {
	i := infra.New("config/config.json")
	i.SetMode()
	i.Migrate(
		&model.Faculty{},
		&model.Major{},
		&model.StudyProgram{},
		&model.User{},
		&model.Student{},
		&model.Teacher{},
		&model.PasswordResetToken{},
		&model.ActivationToken{},
		&model.Subject{},
		&model.Schedule{},
		&model.DailySchedule{},
		&model.UserSchedule{},
		&model.Attendance{},
		&model.AttendanceLog{},
	)

	api.NewServer(i).Run()
}
