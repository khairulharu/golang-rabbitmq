package main

import (
	"context"
	"log"
	"time"

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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Khairul Aswad"

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed To Publish A Message")
	log.Printf(" [x] Send %s\n ", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf(" %s: %s", msg, err)
	}
}
