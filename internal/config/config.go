package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr        string
	DatabaseURL     string
	JWTSecret       string
	ResendAPIKey    string
	ResendFromEmail string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		HTTPAddr:        env("HTTP_ADDR", ":8080"),
		DatabaseURL:     env("DATABASE_URL", "postgres://gotask:gotask@localhost:5432/gotask?sslmode=disable"),
		JWTSecret:       env("JWT_SECRET", "change-me-in-production"),
		ResendAPIKey:    env("RESEND_API_KEY", ""),
		ResendFromEmail: env("RESEND_FROM_EMAIL", ""),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
