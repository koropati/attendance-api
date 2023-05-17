package main

import (
	"attendance-api/api"
	_ "attendance-api/docs"
	"attendance-api/infra"
	"attendance-api/model"
	"log"
	"os"
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
	handleCommand()
}

func handleCommand() {
	if len(os.Args) >= 2 {
		switch command := os.Args[1]; command {
		case "scheduler":
			log.Printf("Hello gan \n")
		case "server":
			i := infra.New("config/config.json")
			i.SetMode()
			api.NewServer(i).Run()
		case "migrate":
			log.Printf("Running Auto Migration...\n")
			i := infra.New("config/config.json")
			i.SetMode()
			i.Migrate(
				&model.Faculty{},
				&model.Major{},
				&model.StudyProgram{},
				&model.User{},
				&model.Auth{},
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
			log.Printf("Berhasil Melakukan Migrasi Database!\n")
			os.Exit(0)
		case "help":
			log.Printf("Available List Command:\n")
			log.Printf("- go run main.go server    (to start server process)\n")
			log.Printf("- go run main.go scheduler (to start scheduler process)\n")
			log.Printf("- go run main.go migrate   (to start migration process)\n")
		default:
			log.Printf("It's Working!\n")
		}
	} else {
		log.Printf("Program It's Working!, you must select opration to start a session.\n")
		log.Printf("List Command:\n")
		log.Printf("- go run main.go server    (to start server process)\n")
		log.Printf("- go run main.go scheduler (to start scheduler process)\n")
		log.Printf("- go run main.go migrate   (to start migration process)\n")
		log.Printf("- go run main.go help      (to see list of command)\n")
	}

}
