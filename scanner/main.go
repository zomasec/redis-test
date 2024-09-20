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

	requestQueue := "scan_requests"
	responseQueue := "scan_responses"

	// Declare request and response queues
	q, err := ch.QueueDeclare(
		requestQueue, // queue name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a request queue")

	// Declare the response queue
	_, err = ch.QueueDeclare(
		responseQueue, // queue name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a response queue")

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
			log.Printf("Received a scan request: %s", d.Body)
			fmt.Println("Processing file scan...")

			// Simulate a scan by sleeping for a few seconds
			time.Sleep(2 * time.Second)

			// Send response back to the backend every 5 seconds
			scanResponse := fmt.Sprintf("Scan complete for file: %s", d.Body)
			err = ch.Publish(
				"",            // exchange
				responseQueue, // routing key (queue name)
				false,         // mandatory
				false,         // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(scanResponse),
				})
			failOnError(err, "Failed to publish response")
			fmt.Println("Sent scan response: ", scanResponse)
		}
	}()

	log.Printf(" [*] Waiting for scan requests. To exit press CTRL+C")
	<-forever
}
