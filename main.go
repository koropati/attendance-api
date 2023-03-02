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

// @contact.name   Dewok satria
// @contact.url    https://twitter.com/dewok_satria
// @contact.email  dewa.ketut.satriawan@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3000
// @BasePath  /v1

func main() {
	i := infra.New("config/config.json")
	i.SetMode()
	i.Migrate(
		&model.User{},
		&model.Subject{},
		&model.Schedule{},
	)

	api.NewServer(i).Run()
}
