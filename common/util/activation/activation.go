package activation

import (
	"attendance-api/model"
	"log"
	"strconv"
	"strings"
	"time"
)

type Activation interface {
	Generate(durationInHours int) (output string)
	Valid(inputString string) (userID int, valid bool)
}

type activation struct {
	a model.User
}

func New(a model.User) Activation {
	return &activation{a: a}
}

func (a *activation) Generate(durationInHours int) (output string) {
	// ID User _ Valid Until
	timeValidUntil := time.Now().Local().Add(time.Hour * time.Duration(durationInHours)).Format("2006-01-02T15:04:05")
	output = strconv.Itoa(int(a.a.ID)) + "_" + timeValidUntil
	return
}

func (a *activation) Valid(inputString string) (userID int, valid bool) {
	data := strings.Split(inputString, "_")
	userID, err := strconv.Atoi(data[0])
	if err != nil {
		log.Printf("[Error] [Convert UserID to Int] E: %v", err)
	}
	dateTime, err := time.Parse("2006-01-02T15:04:05", data[1])
	if err != nil {
		log.Printf("[Error] [Parse DateTime String to Time] E: %v", err)
	}

	valid = time.Now().Before(dateTime)
	return
}
