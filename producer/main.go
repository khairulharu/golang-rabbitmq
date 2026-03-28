package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type EmailMessage struct {
	Name  string
	Email string
}

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
	q, err := ch.QueueDeclare(
		"email_queue",
		true,
		false,
		false,
		false,
		queueArgs,
	)
	failOnError(err, "Failed To Declare A Queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	msg := EmailMessage{
		Name:  "Khairul",
		Email: "khairul@gmail.com",
	}

	body, err := json.Marshal(msg)
	failOnError(err, "Err Marshalling Message")

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		})
	failOnError(err, "Failed To Publish A Message")
	log.Printf(" [x] Send %s\n ", body)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf(" %s: %s", msg, err)
	}
}
