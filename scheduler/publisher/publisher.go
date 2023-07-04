package publisher

import (
	"attendance-api/infra"
	"log"
)

type Publisher interface {
	Run()
}

type publisher struct {
	infra infra.Infra
}

func NewPublisher(infra infra.Infra) Publisher {
	return &publisher{
		infra: infra,
	}
}

func (c publisher) Run() {
	config := c.infra.Config().Sub("amqp")
	queueName := config.GetString("queue_name")

	conn := c.infra.AMQP()
	defer conn.Close()

	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	defer amqpChannel.Close()

	stopChan := make(chan bool)
	cronJob := InitCronJob(amqpChannel, queueName)
	cronJob.Start()

	defer cronJob.Stop()
	// Stop for program termination
	<-stopChan

}
func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
