package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/vashafaris/nucleus/internal/infrastructure"
	"github.com/vashafaris/nucleus/internal/interfaces/http/router"
	"github.com/vashafaris/nucleus/pkg/config"
)

func main() {
	// Load .env file - try .env.local first for local development
	if err := godotenv.Load(".env.local"); err != nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file not found")
		}
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	log.Printf("Starting %s application...", cfg.App.Name)

	// Initialize infrastructure
	infra, err := infrastructure.NewManager(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize infrastructure: %v", err)
	}
	defer func() {
		if err := infra.Close(); err != nil {
			log.Printf("Error closing infrastructure: %v", err)
		}
	}()

	// Check infrastructure health
	if healthErrors := infra.Health(); len(healthErrors) > 0 {
		log.Printf("Warning: Some infrastructure components are unhealthy: %v", healthErrors)
	} else {
		log.Println("All infrastructure components are healthy")
	}

	// Setup router
	r := router.New(infra)
	r.Setup()

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      r.Engine(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("%s is running on http://localhost:%s", cfg.App.Name, cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
