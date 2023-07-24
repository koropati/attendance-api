package jobs

import (
	"attendance-api/scheduler"
	"attendance-api/service"
	"fmt"
	"log"
	"time"
)

type PasswordResetTokenJob interface {
	AutoDelete()
}

type passwordResetTokenJob struct {
	passwordResetTokenService service.PasswordResetTokenService
	task                      *scheduler.AddTask
}

func NewPasswordResetTokenJob(
	passwordResetTokenService service.PasswordResetTokenService,
	task *scheduler.AddTask,
) PasswordResetTokenJob {
	return &passwordResetTokenJob{
		passwordResetTokenService: passwordResetTokenService,
		task:                      task,
	}
}

func (j passwordResetTokenJob) AutoDelete() {
	fmt.Println("Execute Task PasswordResetToken [AUTO DELETE]")
	fmt.Printf("Action: %v\n", j.task.Action)
	fmt.Printf("Body  : %v\n", j.task.Body)
	fmt.Printf("Date  : %v\n", j.task.Date)
	fmt.Printf("TStm  : %v\n", j.task.TimeStamp)

	currentTime := time.Now()
	expiredTime := currentTime.Add(-15 * time.Minute)

	err := j.passwordResetTokenService.DeleteExpiredPasswordResetToken(expiredTime)
	if err != nil {
		log.Printf("[Scheduler] [Error] [PasswordResetToken-AUTO-DELETE] E: %v\n", err)
	} else {
		log.Printf("[Scheduler] [Success] [PasswordResetToken-AUTO-DELETE] [%v]\n", expiredTime)
	}

}
