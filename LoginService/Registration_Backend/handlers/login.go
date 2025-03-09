package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"registration_backend/database"

	"github.com/streadway/amqp"
	"golang.org/x/crypto/bcrypt"
)

// Structs for login handling
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Listen for login requests
func ListenForLoginRequests() {
	msgs, err := rabbitMQChannel.Consume(
		"login_verification", "", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("Failed to consume login_verification queue: %v", err)
	}

	fmt.Println("Listening for login requests...")

	for msg := range msgs {
		var loginReq LoginRequest
		if err := json.Unmarshal(msg.Body, &loginReq); err != nil {
			log.Printf("Error decoding login request: %v", err)
			continue
		}

		fmt.Printf("Received login request for email: %s\n", loginReq.Email)

		// Validate credentials
		response := validateUser(loginReq)

		// Send response to login_response queue
		sendLoginResponse(msg.CorrelationId, response)
	}
}

// Validate login credentials
func validateUser(loginReq LoginRequest) LoginResponse {
	var storedPassword string
	err := database.DB.QueryRow(context.Background(),
		"SELECT password FROM user_registration WHERE email=$1", loginReq.Email).Scan(&storedPassword)

	if err != nil {
		return LoginResponse{Status: "error", Message: "User not found"}
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginReq.Password)); err != nil {
		return LoginResponse{Status: "error", Message: "Invalid credentials"}
	}

	return LoginResponse{
		Status:  "success",
		Message: "Login successful",
	}
}

// Send login response to RabbitMQ
func sendLoginResponse(correlationId string, response LoginResponse) {
	body, _ := json.Marshal(response)
	err := rabbitMQChannel.Publish(
		"", "login_response", false, false,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          body,
			CorrelationId: correlationId,
		},
	)
	if err != nil {
		log.Printf("Failed to send login response: %v", err)
	} else {
		fmt.Println("Login response sent for:", correlationId)
	}
}
