package handlers

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// RabbitMQ connection variables
var rabbitMQConn *amqp.Connection
var rabbitMQChannel *amqp.Channel

const exchangeName = "user_exchange"
const queueName = "user.registration"
const routingKey = "user.registration"

// Dead Letter Exchange & Queue
const deadLetterExchange = "dead_letter_exchange"
const deadLetterQueue = "dead_letter_queue"

// Initialize RabbitMQ with DLQ setup
func InitRabbitMQ() {
	var err error
	rabbitMQConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	rabbitMQChannel, err = rabbitMQConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	// Declare the Dead Letter Exchange
	err = rabbitMQChannel.ExchangeDeclare(
		deadLetterExchange, // DLX Name
		"fanout",           // Exchange Type (Fanout sends messages to all bound queues)
		true,               // Durable
		false,              // Auto-deleted
		false,              // Internal
		false,              // No-wait
		nil,                // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare Dead Letter Exchange: %v", err)
	}

	// Declare the Dead Letter Queue
	_, err = rabbitMQChannel.QueueDeclare(
		deadLetterQueue, // DLQ Name
		true,            // Durable
		false,           // Auto-delete
		false,           // Exclusive
		false,           // No-wait
		nil,             // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare Dead Letter Queue: %v", err)
	}

	// Bind DLQ to DLX
	err = rabbitMQChannel.QueueBind(
		deadLetterQueue,    // DLQ Name
		"",                 // Routing Key (not needed for fanout)
		deadLetterExchange, // DLX Name
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind Dead Letter Queue: %v", err)
	}

	// Declare the Main Exchange
	err = rabbitMQChannel.ExchangeDeclare(
		exchangeName, // Main Exchange
		"topic",      // Exchange Type
		true,         // Durable
		false,        // Auto-deleted
		false,        // Internal
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	// Declare the Main Queue with Dead Letter Support
	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = deadLetterExchange // Attach DLX
	args["x-message-ttl"] = int32(60000)                // 60s TTL for messages (optional)

	_, err = rabbitMQChannel.QueueDeclare(
		queueName, // Main Queue Name
		true,      // Durable
		false,     // Auto-delete
		false,     // Exclusive
		false,     // No-wait
		args,      // Attach DLX
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Bind Main Queue to Main Exchange
	err = rabbitMQChannel.QueueBind(
		queueName,    // Queue Name
		routingKey,   // Routing Key
		exchangeName, // Exchange Name
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue to exchange: %v", err)
	}

	log.Println("RabbitMQ with Dead Letter Queue initialized successfully!")
}

// PublishToQueue sends user data to RabbitMQ
func PublishToQueue(user User) error {
	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = rabbitMQChannel.Publish(
		exchangeName, // Exchange Name
		routingKey,   // Routing Key
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Printf("Failed to publish message: %v", err)
	} else {
		log.Println("User data published successfully!")
	}

	return err
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
