package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPAddr             string
	DatabaseURL          string
	JWTSecret            string
	ResendAPIKey         string
	ResendFromEmail      string
	MigrationsPath       string
	SwaggerUsername      string
	SwaggerPassword      string
	CORSAllowedOrigins   string
	CORSAllowedMethods   string
	CORSAllowedHeaders   string
	CORSExposedHeaders   string
	CORSAllowCredentials string
	CORSMaxAge           string
}

func Load() Config {
	_ = godotenv.Load()

	return Config{
		HTTPAddr:             env("HTTP_ADDR", ":8080"),
		DatabaseURL:          env("DATABASE_URL", "postgres://gotask:gotask@localhost:5432/gotask?sslmode=disable"),
		JWTSecret:            env("JWT_SECRET", "change-me-in-production"),
		ResendAPIKey:         env("RESEND_API_KEY", ""),
		ResendFromEmail:      env("RESEND_FROM_EMAIL", ""),
		MigrationsPath:       env("MIGRATIONS_PATH", "migrations"),
		SwaggerUsername:      env("SWAGGER_USERNAME", ""),
		SwaggerPassword:      env("SWAGGER_PASSWORD", ""),
		CORSAllowedOrigins:   env("CORS_ALLOWED_ORIGINS", ""),
		CORSAllowedMethods:   env("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS"),
		CORSAllowedHeaders:   env("CORS_ALLOWED_HEADERS", "Authorization,Content-Type"),
		CORSExposedHeaders:   env("CORS_EXPOSED_HEADERS", ""),
		CORSAllowCredentials: env("CORS_ALLOW_CREDENTIALS", "false"),
		CORSMaxAge:           env("CORS_MAX_AGE", "12h"),
	}
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
