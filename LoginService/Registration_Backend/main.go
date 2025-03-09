package main

import (
	"fmt"
	"log"
	"net/http"
	"registration_backend/database"
	"registration_backend/handlers"
	"registration_backend/storage" // Import storage package

	"github.com/rs/cors" // Import CORS package
)

func main() {
	// Connect to Database
	database.ConnectDB()
	defer database.CloseDB()

	// Connect to RabbitMQ
	handlers.InitRabbitMQ()
	defer handlers.CloseRabbitMQ()

	// Initialize MinIO
	storage.InitMinio() // 🛠️ Add this line to initialize MinIO

	// Setup HTTP Server
	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.RegisterUserHandler)

	// Enable CORS
	handler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow requests from frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}).Handler(mux)

	fmt.Println("🚀 Starting Registration Service on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
