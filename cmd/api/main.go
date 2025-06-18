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

	"github.com/vashafaris/nucleus/pkg/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Debug: Print loaded configuration
	log.Printf("Loaded configuration:")
	log.Printf("  App Name: %s", cfg.App.Name)
	log.Printf("  Environment: %s", cfg.App.Env)
	log.Printf("  Port: %s", cfg.App.Port)
	log.Printf("  DB Host: %s", cfg.Database.Host)
	log.Printf("  DB Name: %s", cfg.Database.Name)
	log.Printf("  DB User: %s", cfg.Database.User)

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Setup logger (will be implemented in infrastructure layer)
	log.Printf("Starting %s application on port %s", cfg.App.Name, cfg.App.Port)

	// Create server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      setupRoutes(), // Will be implemented
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("%s is running on http://localhost:%s", cfg.App.Name, cfg.App.Port)

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

	log.Println("Server exited")
}

// setupRoutes will be replaced with actual router setup
func setupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}
