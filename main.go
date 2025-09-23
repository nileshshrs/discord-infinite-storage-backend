package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/nileshshrs/infinite-storage/application"
	"github.com/nileshshrs/infinite-storage/bot"
	"github.com/nileshshrs/infinite-storage/config"
	"github.com/nileshshrs/infinite-storage/db"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Load config
	cfg := config.Load()

	// Connect to MongoDB
	collection, err := db.Connect(cfg.URI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Start HTTP server in background
	app := application.New(collection, cfg)
	go func() {
		if err := app.Start(context.Background()); err != nil {
			fmt.Println("Error starting application:", err)
		}
	}()

	// Start Discord bot
	dg := bot.Run(cfg)

	// Graceful shutdown for both HTTP + Bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	dg.Close()
	log.Println("Application stopped")
}
