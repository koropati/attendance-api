package jobs

import (
	"attendance-api/scheduler"
	"attendance-api/service"
	"fmt"
	"log"
	"time"
)

type AuthJob interface {
	AutoDelete()
}

type authJob struct {
	authService service.AuthService
	task        *scheduler.AddTask
}

func NewAuthJob(
	authService service.AuthService,
	task *scheduler.AddTask,
) AuthJob {
	return &authJob{
		authService: authService,
		task:        task,
	}
}

func (j authJob) AutoDelete() {
	fmt.Println("Execute Task Auth [AUTO DELETE]")
	fmt.Printf("Action: %v\n", j.task.Action)
	fmt.Printf("Body  : %v\n", j.task.Body)
	fmt.Printf("Date  : %v\n", j.task.Date)
	fmt.Printf("TStm  : %v\n", j.task.TimeStamp)

	currentTimeMillis := time.Now().UnixNano() / int64(time.Millisecond)

	err := j.authService.DeleteExpiredAuth(currentTimeMillis)
	if err != nil {
		log.Printf("[Scheduler] [Error] [Auth-AUTO-DELETE] E: %v\n", err)
	} else {
		log.Printf("[Scheduler] [Success] [Auth-AUTO-DELETE] [%v]\n", currentTimeMillis)
	}

}
