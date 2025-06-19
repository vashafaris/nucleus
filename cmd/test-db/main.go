package main

import (
	"context"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/vashafaris/nucleus/internal/infrastructure"
	"github.com/vashafaris/nucleus/pkg/config"
)

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		godotenv.Load()
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize infrastructure
	infra, err := infrastructure.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize infrastructure: %v", err)
	}
	defer infra.Close()

	fmt.Println("✅ Successfully connected to PostgreSQL and Redis!")

	// Test PostgreSQL query
	var count int
	err = infra.DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if err != nil {
		log.Fatalf("Failed to query products: %v", err)
	}
	fmt.Printf("📊 Found %d products in database\n", count)

	// Test Redis
	ctx := context.Background()
	err = infra.Redis.SetContext(ctx, "test:key", "Hello from Nucleus!", 0)
	if err != nil {
		log.Fatalf("Failed to set Redis key: %v", err)
	}

	value, err := infra.Redis.GetContext(ctx, "test:key")
	if err != nil {
		log.Fatalf("Failed to get Redis key: %v", err)
	}
	fmt.Printf("🔑 Redis test value: %s\n", value)

	// Clean up test key
	infra.Redis.DeleteContext(ctx, "test:key")

	// Check health
	healthErrors := infra.Health()
	if len(healthErrors) == 0 {
		fmt.Println("✅ All infrastructure components are healthy!")
	} else {
		fmt.Println("❌ Health check errors:", healthErrors)
	}
}
