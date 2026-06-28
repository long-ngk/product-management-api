package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/product-management-server/internal/handler"
	"github.com/product-management-server/internal/infrastructure"
	"github.com/product-management-server/internal/repository"
	"github.com/product-management-server/internal/service"
)

func main() {
	// Load .env file if it exists (optional, won't error if missing)
	_ = godotenv.Load()

	// Read configuration from environment variables
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	// Initialize database connection
	db, err := infrastructure.NewDatabaseConnection(dsn)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close()

	// Wire dependencies: Repository → Service → Handler
	repo := repository.NewProductRepository(db)
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)

	// Setup Gin router with default middleware (logger + recovery)
	router := gin.Default()

	// Register routes
	h.RegisterRoutes(router)

	// Start HTTP server
	log.Printf("server starting on %s", port)
	if err := router.Run(port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
