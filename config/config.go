package config

import (
	"log"
	"os"
)

type Config struct {
	URI  string
	Port string
}

func Load() *Config {
	cfg := &Config{
		URI:  os.Getenv("URI"),            // use URI as env variable
		Port: getEnv("PORT", "6278"),
	}

	if cfg.URI == "" {
		log.Fatal("URI environment variable is not set")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
