package main

import (
	"log"

	"gotask/internal/api/customer/auth"
	"gotask/internal/api/customer/me"
	"gotask/internal/api/customer/task"
	"gotask/internal/config"
	"gotask/internal/database"
	gotaskemail "gotask/internal/email"
	response "gotask/internal/utils/response"

	"github.com/gin-gonic/gin"
)

// @title GoTask API
// @version 1.0
// @description REST API for GoTask customer accounts and tasks.
// @BasePath /
// @schemes http https
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Supply a customer JWT as `Bearer <token>`.
func main() {
	cfg := config.Load()
	if cfg.SwaggerUsername == "" || cfg.SwaggerPassword == "" {
		log.Fatal("SWAGGER_USERNAME and SWAGGER_PASSWORD must both be set")
	}

	db, err := database.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}

	taskHandler := task.NewHandler(task.NewService(task.NewRepository(db)))
	emailService := gotaskemail.NewResendService(cfg.ResendAPIKey, cfg.ResendFromEmail)
	authRepository := auth.NewRepository(db)
	authHandler := auth.NewHandler(auth.NewService(authRepository, cfg.JWTSecret, emailService))
	meHandler := me.NewHandler(me.NewService(me.NewRepository(db), emailService))
	protectedMiddleware := auth.JWTMiddlewareWithUserStore(cfg.JWTSecret, authRepository)

	router := gin.Default()
	router.GET("/health", health)
	registerSwaggerRoutes(router, cfg.SwaggerUsername, cfg.SwaggerPassword)
	taskHandler.RegisterRoutes(router, protectedMiddleware)
	authHandler.RegisterRoutes(router)
	meHandler.RegisterRoutes(router, protectedMiddleware)

	log.Printf("API listening on %s", cfg.HTTPAddr)
	if err := router.Run(cfg.HTTPAddr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

// health godoc
// @Summary Check API health
// @Description Reports whether the API server is available.
// @Tags System
// @Produce json
// @Success 200 {object} response.SuccessEnvelope
// @Router /health [get]
func health(c *gin.Context) {
	response.JSON(c, 200, gin.H{"status": "ok"})
}
