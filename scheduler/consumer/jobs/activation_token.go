package jobs

import (
	"attendance-api/scheduler"
	"attendance-api/service"
	"fmt"
	"log"
	"time"
)

type ActivationTokenJob interface {
	AutoDelete()
}

type activationTokenJob struct {
	activationTokenService service.ActivationTokenService
	task                   *scheduler.AddTask
}

func NewActivationTokenJob(
	activationTokenService service.ActivationTokenService,
	task *scheduler.AddTask,
) ActivationTokenJob {
	return &activationTokenJob{
		activationTokenService: activationTokenService,
		task:                   task,
	}
}

func (j activationTokenJob) AutoDelete() {
	fmt.Println("Execute Task ActivationToken [AUTO DELETE]")
	fmt.Printf("Action: %v\n", j.task.Action)
	fmt.Printf("Body  : %v\n", j.task.Body)
	fmt.Printf("Date  : %v\n", j.task.Date)
	fmt.Printf("TStm  : %v\n", j.task.TimeStamp)

	currentTime := time.Now()
	expiredTime := currentTime.Add(-15 * time.Minute)

	err := j.activationTokenService.DeleteExpiredActivationToken(expiredTime)
	if err != nil {
		log.Printf("[Scheduler] [Error] [ActivationToken-AUTO-DELETE] E: %v\n", err)
	} else {
		log.Printf("[Scheduler] [Success] [ActivationToken-AUTO-DELETE] [%v]\n", expiredTime)
	}

}
