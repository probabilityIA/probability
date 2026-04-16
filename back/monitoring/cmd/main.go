package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/secamc93/probability/back/monitoring/internal/app"
	"github.com/secamc93/probability/back/monitoring/internal/infra/primary/handlers"
	"github.com/secamc93/probability/back/monitoring/internal/infra/secondary/docker"
	"github.com/secamc93/probability/back/monitoring/internal/infra/secondary/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	_ = godotenv.Load()

	port := getEnv("HTTP_PORT", "3070")
	jwtSecret := getEnv("JWT_SECRET", "monitoring-secret-change-in-prod")
	ginMode := getEnv("GIN_MODE", "debug")

	// Database
	db, err := connectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("Database connected")

	// Docker client
	composeProject := getEnv("COMPOSE_PROJECT", "probability")
	dockerClient, err := docker.New(composeProject)
	if err != nil {
		log.Fatalf("Failed to connect to Docker: %v", err)
	}

	// Repository
	userRepo := repository.New(db)

	// Use case
	useCase := app.New(dockerClient, userRepo, jwtSecret)

	// HTTP server
	gin.SetMode(ginMode)
	router := gin.Default()

	// CORS
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Register routes
	h := handlers.New(useCase, jwtSecret)
	h.RegisterRoutes(router)

	fmt.Printf("Monitoring API running on :%s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectDB() (*gorm.DB, error) {
	sslmode := getEnv("PGSSLMODE", "disable")
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASS", "postgres"),
		getEnv("DB_NAME", "postgres"),
		getEnv("DB_PORT", "5433"),
		sslmode,
	)

	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
