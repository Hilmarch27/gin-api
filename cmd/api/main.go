package main

import (
	"log"
	"time"

	"github.com/Hilmarch27/gin-api/internal/delivery/http/handler"
	"github.com/Hilmarch27/gin-api/internal/delivery/http/router"
	"github.com/Hilmarch27/gin-api/internal/domain"
	"github.com/Hilmarch27/gin-api/internal/repository"
	"github.com/Hilmarch27/gin-api/internal/usecase"
	"github.com/Hilmarch27/gin-api/pkg/config"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate database
	err = cfg.DB.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatal(err)
	}

	// Initialize repositories
	userRepo := repository.NewUserRepository(cfg.DB)

	// Initialize usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, cfg.JWTSecret, time.Hour*1)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUsecase)

	// Initialize Gin engine
	engine := gin.Default()

	// Initialize routers
	authRouter := router.NewAuthRouter(authHandler, cfg.JWTSecret)
	apiRouter := router.NewApiRouter(cfg.JWTSecret)

	// Setup main router
	mainRouter := router.NewRouter(engine, authRouter, apiRouter, []byte(cfg.JWTSecret))
	mainRouter.SetupRoutes()

	// Start server
	if err := engine.Run(":3027"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}