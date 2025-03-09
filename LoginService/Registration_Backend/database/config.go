package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func ConnectDB() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file. Make sure it exists in the root directory")
	}

	// Read database configuration from environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Validate environment variables
	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Database environment variables are not set properly")
	}

	// Construct the PostgreSQL connection string
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	// Create context
	ctx := context.Background()

	// Connect to PostgreSQL using the context
	DB, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL database successfully!")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		fmt.Println("Database connection closed.")
	}
}
