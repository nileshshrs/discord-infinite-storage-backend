package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/nileshshrs/infinite-storage/application"
	"github.com/nileshshrs/infinite-storage/config"
	"github.com/nileshshrs/infinite-storage/db"
)

func main() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	// Connect to MongoDB
	collection, err := db.Connect(cfg.URI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Start app
	app := application.New(collection)
	if err := app.Start(context.Background()); err != nil {
		fmt.Println("Error starting application:", err)
	}
}
