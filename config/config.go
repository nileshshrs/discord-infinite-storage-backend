package config

import (
	"log"
	"os"
)

type Config struct {
	URI              string
	Port             string
	DiscordToken     string
	DiscordClientID  string
	DiscordChannelID string
}

func Load() *Config {
	cfg := &Config{
		URI:              os.Getenv("URI"),
		Port:             getEnv("PORT", "6278"),
		DiscordToken:     os.Getenv("DISCORD_TOKEN"),
		DiscordClientID:  os.Getenv("DISCORD_CLIENT_ID"),
		DiscordChannelID: os.Getenv("DISCORD_CHANNEL_ID"),
	}

	if cfg.URI == "" {
		log.Fatal("URI environment variable is not set")
	}
	if cfg.DiscordToken == "" {
		log.Println("warning: DISCORD_TOKEN is not set")
	}
	if cfg.DiscordClientID == "" {
		log.Println("warning: DISCORD_CLIENT_ID is not set")
	}
	if cfg.DiscordChannelID == "" {
		log.Println("warning: DISCORD_CHANNEL_ID is not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
