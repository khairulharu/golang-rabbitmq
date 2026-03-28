package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed To Open A Channel")
	defer ch.Close()

	queueArgs := amqp.Table{
		amqp.QueueTypeArg: amqp.QueueTypeQuorum,
	}
	q, err := ch.QueueDeclare("hello", true, false, false, false, queueArgs)
	failOnError(err, "Failed To Declare A Queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed To Register A Consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Recieved a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To Exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}
