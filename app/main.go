package main

import (
	"attendance-api/api"
	"attendance-api/infra"
	"attendance-api/model"
)

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
