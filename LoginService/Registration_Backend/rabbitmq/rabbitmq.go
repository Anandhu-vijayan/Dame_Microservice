package handlers

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

var rabbitMQConn *amqp.Connection
var rabbitMQChannel *amqp.Channel

// ConnectRabbitMQ initializes the RabbitMQ connection
func ConnectRabbitMQ() {
	var err error
	rabbitMQConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitMQChannel, err = rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	_, err = rabbitMQChannel.QueueDeclare(
		"login_verification",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare login_verification queue: %v", err)
	}

	_, err = rabbitMQChannel.QueueDeclare(
		"login_response",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare login_response queue: %v", err)
	}

	fmt.Println("RabbitMQ connected and queues declared!")
}

// CloseRabbitMQ closes the RabbitMQ connection
func CloseRabbitMQ() {
	if rabbitMQChannel != nil {
		rabbitMQChannel.Close()
	}
	if rabbitMQConn != nil {
		rabbitMQConn.Close()
	}
}
