package main

import (
	"log"

	"gotask/internal/api/customer/auth"
	"gotask/internal/api/customer/me"
	"gotask/internal/api/customer/task"
	"gotask/internal/config"
	"gotask/internal/cors"
	"gotask/internal/database"
	gotaskemail "gotask/internal/email"
	"gotask/internal/ratelimit"
	response "gotask/internal/utils/response"

	"github.com/gin-gonic/gin"
	"time"
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
	statisticsHandler := task.NewStatisticsHandler(task.NewStatisticsService(task.NewStatisticsRepository(db)))
	emailService := gotaskemail.NewResendService(cfg.ResendAPIKey, cfg.ResendFromEmail)
	authRepository := auth.NewRepository(db)
	authHandler := auth.NewHandler(auth.NewService(authRepository, cfg.JWTSecret, emailService))
	meHandler := me.NewHandler(me.NewService(me.NewRepository(db), emailService))
	protectedMiddleware := auth.JWTMiddlewareWithUserStore(cfg.JWTSecret, authRepository)
	apiRateLimit := ratelimit.New(10, time.Minute).Middleware()
	emailRateLimit := ratelimit.New(3, time.Minute).Middleware()
	corsMiddleware, err := cors.New(cors.Config{
		AllowedOrigins:   cfg.CORSAllowedOrigins,
		AllowedMethods:   cfg.CORSAllowedMethods,
		AllowedHeaders:   cfg.CORSAllowedHeaders,
		ExposedHeaders:   cfg.CORSExposedHeaders,
		AllowCredentials: cfg.CORSAllowCredentials,
		MaxAge:           cfg.CORSMaxAge,
	})
	if err != nil {
		log.Fatalf("configure CORS: %v", err)
	}

	router := gin.Default()
	if err := router.SetTrustedProxies(nil); err != nil {
		log.Fatalf("configure trusted proxies: %v", err)
	}
	router.Use(corsMiddleware)
	registerSwaggerRoutes(router, cfg.SwaggerUsername, cfg.SwaggerPassword)
	router.Use(apiRateLimit)
	router.GET("/health", health)
	statisticsHandler.RegisterRoutes(router, protectedMiddleware)
	taskHandler.RegisterRoutes(router, protectedMiddleware)
	authHandler.RegisterRoutes(router, emailRateLimit)
	meHandler.RegisterRoutes(router, protectedMiddleware, emailRateLimit)

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
// @Failure 429 {object} response.ErrorEnvelope
// @Router /health [get]
func health(c *gin.Context) {
	response.JSON(c, 200, gin.H{"status": "ok"})
}
