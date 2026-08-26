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
	MigrationsPath  string
	SwaggerUsername string
	SwaggerPassword string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		HTTPAddr:        env("HTTP_ADDR", ":8080"),
		DatabaseURL:     env("DATABASE_URL", "postgres://gotask:gotask@localhost:5432/gotask?sslmode=disable"),
		JWTSecret:       env("JWT_SECRET", "change-me-in-production"),
		ResendAPIKey:    env("RESEND_API_KEY", ""),
		ResendFromEmail: env("RESEND_FROM_EMAIL", ""),
		MigrationsPath:  env("MIGRATIONS_PATH", "migrations"),
		SwaggerUsername: env("SWAGGER_USERNAME", ""),
		SwaggerPassword: env("SWAGGER_PASSWORD", ""),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
