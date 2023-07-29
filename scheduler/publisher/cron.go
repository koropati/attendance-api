package publisher

import (
	"attendance-api/scheduler"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
	"gopkg.in/robfig/cron.v2"
)

func InitCronJob(amqpChannel *amqp.Channel, queueName string) (cronJob *cron.Cron) {
	c := cron.New()

	c.AddFunc("* * * * *", TaskAttendance(amqpChannel, queueName))         //tiap jam 00:00 dini hari
	c.AddFunc("@hourly", TaskAuth(amqpChannel, queueName))                 //tiap 1 Jam
	c.AddFunc("0 3 * * *", TaskActivationToken(amqpChannel, queueName))    //tiap jam 03:00 dini hari
	c.AddFunc("0 3 * * *", TaskPasswordResetToken(amqpChannel, queueName)) //tiap jam 03:00 dini hari

	return c
}

func TaskAttendance(amqpChannel *amqp.Channel, queueName string) func() {
	return func() {
		fmt.Println("Task Attendance")

		addTask := scheduler.AddTask{
			Action:    "attendance",
			Body:      "update",
			Date:      time.Now().Format("2006-01-02"),
			TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
		}

		PushMessage(amqpChannel, addTask, queueName)

	}
}

func TaskAuth(amqpChannel *amqp.Channel, queueName string) func() {
	return func() {
		fmt.Println("Task Auth")

		addTask := scheduler.AddTask{
			Action:    "auth",
			Body:      "auto_delete",
			Date:      time.Now().Format("2006-01-02"),
			TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
		}

		PushMessage(amqpChannel, addTask, queueName)

	}
}

func TaskActivationToken(amqpChannel *amqp.Channel, queueName string) func() {
	return func() {
		fmt.Println("Task Activation Token")

		addTask := scheduler.AddTask{
			Action:    "activation_token",
			Body:      "auto_delete",
			Date:      time.Now().Format("2006-01-02"),
			TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
		}

		PushMessage(amqpChannel, addTask, queueName)

	}
}

func TaskPasswordResetToken(amqpChannel *amqp.Channel, queueName string) func() {
	return func() {
		fmt.Println("Task Password Reset Token")

		addTask := scheduler.AddTask{
			Action:    "password_reset_token",
			Body:      "auto_delete",
			Date:      time.Now().Format("2006-01-02"),
			TimeStamp: time.Now().Format("2006-01-02 15:04:05"),
		}

		PushMessage(amqpChannel, addTask, queueName)

	}
}

func PushMessage(amqpChannel *amqp.Channel, addTask scheduler.AddTask, queueName string) {
	queue, err := amqpChannel.QueueDeclare(queueName, true, false, false, false, nil)
	handleError(err, fmt.Sprintf(`Could not declare "%s" queue`, queueName))

	body, err := json.Marshal(addTask)
	if err != nil {
		handleError(err, "Error encoding JSON")
	}

	err = amqpChannel.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})

	if err != nil {
		log.Printf("Error publishing message: %s", err)
	}
}
