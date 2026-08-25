package main

import (
	"log"

	"gotask/internal/auth"
	"gotask/internal/config"
	"gotask/internal/database"
	gotaskemail "gotask/internal/email"
	"gotask/internal/me"
	apiresponse "gotask/internal/response"
	"gotask/internal/task"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate database: %v", err)
	}

	taskHandler := task.NewHandler(task.NewService(task.NewRepository(db)))
	emailService := gotaskemail.NewResendService(cfg.ResendAPIKey, cfg.ResendFromEmail)
	authRepository := auth.NewRepository(db)
	authHandler := auth.NewHandler(auth.NewService(authRepository, cfg.JWTSecret, emailService))
	meHandler := me.NewHandler(me.NewService(me.NewRepository(db)))
	protectedMiddleware := auth.JWTMiddlewareWithUserStore(cfg.JWTSecret, authRepository)

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		apiresponse.JSON(c, 200, gin.H{"status": "ok"})
	})
	taskHandler.RegisterRoutes(router, protectedMiddleware)
	authHandler.RegisterRoutes(router)
	meHandler.RegisterRoutes(router, protectedMiddleware)

	log.Printf("API listening on %s", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
