package main

import (
	"flag"
	"log"

	"gotask/internal/config"
	"gotask/internal/database"
)

func main() {
	direction := flag.String("direction", "up", "migration direction: up, down, reset, or version")
	flag.Parse()

	cfg := config.Load()
	switch *direction {
	case "up":
		if err := database.Migrate(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
			log.Fatalf("migrate database: %v", err)
		}
		log.Println("database migrations completed")
	case "down":
		if err := database.Rollback(cfg.DatabaseURL, cfg.MigrationsPath); err != nil {
			log.Fatalf("rollback migration: %v", err)
		}
		log.Println("database migration rolled back")
	case "reset":
		if err := database.Reset(cfg.DatabaseURL); err != nil {
			log.Fatalf("reset database migrations: %v", err)
		}
		log.Println("database schema cleared")
	case "version":
		version, dirty, err := database.Version(cfg.DatabaseURL, cfg.MigrationsPath)
		if err != nil {
			log.Fatalf("read migration version: %v", err)
		}
		log.Printf("migration version: %d, dirty: %t", version, dirty)
	default:
		log.Fatalf("invalid migration direction %q; use up, down, reset, or version", *direction)
	}
}
