package main

import (
	"encoding/json"
	"fmt"
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

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed To Register A Consumer")

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			var emailMsg EmailMessage
			err := json.Unmarshal(msg.Body, &emailMsg)
			if err != nil {
				log.Println("Error Unmarsahl Body:", err)
				msg.Nack(false, false)
				continue
			}

			fmt.Printf("📧 Sending email to %s (%s)...\n", emailMsg.Email, emailMsg.Name)
			time.Sleep(1 * time.Second) // simulasi proses
			fmt.Printf("✅ Email terkirim ke %s!\n", emailMsg.Email)

			// Ack → kasih tau RabbitMQ pesan udah diproses
			msg.Ack(false)
		}
	}()

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s : %s", msg, err)
	}
}
