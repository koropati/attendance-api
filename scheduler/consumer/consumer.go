package consumer

import (
	"attendance-api/infra"
	"attendance-api/manager"
	"attendance-api/scheduler"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Consumer interface {
	Run()
}

type consumer struct {
	infra   infra.Infra
	service manager.ServiceManager
}

func NewConsumer(infra infra.Infra) Consumer {
	return &consumer{
		infra:   infra,
		service: manager.NewServiceManager(infra),
	}
}

func (c consumer) Run() {
	config := c.infra.Config().Sub("amqp")
	queueName := config.GetString("queue_name")

	conn := c.infra.AMQP()
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	queue, err := amqpChannel.QueueDeclare(queueName, true, false, false, false, nil)
	handleError(err, fmt.Sprintf(`Could not declare "%s" queue`, queueName))

	err = amqpChannel.Qos(1, 0, false)
	handleError(err, "Could not configure QoS")

	messageChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", d.Body)

			addTask := &scheduler.AddTask{}

			err := json.Unmarshal(d.Body, addTask)

			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			InitTask(addTask)

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}

		}
	}()

	// Stop for program termination
	<-stopChan

}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
