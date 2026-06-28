package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/product-management-server/config"
	"github.com/product-management-server/internal/app"
)

func main() {
	// Load .env file if it exists (optional, won't error if missing)
	_ = godotenv.Load()

	// Load configuration from environment variables
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize and run the application
	a, err := app.New(cfg)
	if err != nil {
		log.Fatal("failed to initialize app: ", err)
	}

	log.Printf("server starting on %s", cfg.Port)
	if err := a.Run(); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
