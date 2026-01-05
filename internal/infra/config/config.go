package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr          string
	ReadHeaderTimeout time.Duration
	ShutdownTimeout   time.Duration
}

func Load() (*Config, error) {
	// Load .env file, but don't fail if it's missing (e.g. in production or docker)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := &Config{
		HTTPAddr:          ":8080",
		ReadHeaderTimeout: 10 * time.Second,
		ShutdownTimeout:   10 * time.Second,
	}

	if addr := os.Getenv("HTTP_ADDR"); addr != "" {
		cfg.HTTPAddr = addr
	}

	if val := os.Getenv("HTTP_READ_HEADER_TIMEOUT"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			cfg.ReadHeaderTimeout = d
		} else {
			log.Printf("Invalid HTTP_READ_HEADER_TIMEOUT %q, using default", val)
		}
	}

	if val := os.Getenv("HTTP_SHUTDOWN_TIMEOUT"); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			cfg.ShutdownTimeout = d
		} else {
			log.Printf("Invalid HTTP_SHUTDOWN_TIMEOUT %q, using default", val)
		}
	}

	return cfg, nil
}
