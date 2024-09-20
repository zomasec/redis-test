package main

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// Helper function to check errors
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func connectRabbitMQ() (*amqp.Connection, error) {
	for {
		conn, err := amqp.Dial("amqp://rabbitmq")
		if err != nil {
			fmt.Println("Failed to connect to RabbitMQ, retrying in 5 seconds...")
			time.Sleep(5 * time.Second) // Wait before retrying
		} else {
			return conn, nil
		}
	}
}

func main() {
	conn, err := connectRabbitMQ()
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"scan_requests", // queue name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			fmt.Println("Processing file scan...")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
