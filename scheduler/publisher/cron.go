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

	c.AddFunc("* * * * *", TaskAttendance(amqpChannel, queueName))

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
