package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vashafaris/nucleus/pkg/config"
)

func main() {
	// Show current working directory
	wd, _ := os.Getwd()
	fmt.Printf("Current directory: %s\n", wd)

	// Try loading .env files
	if err := godotenv.Load(".env.local"); err != nil {
		fmt.Printf("Could not load .env.local: %v\n", err)
		if err := godotenv.Load(); err != nil {
			fmt.Printf("Could not load .env: %v\n", err)
		}
	}

	// Show environment variables
	fmt.Println("\nEnvironment variables:")
	fmt.Printf("DB_HOST=%s\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT=%s\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_USER=%s\n", os.Getenv("DB_USER"))
	fmt.Printf("DB_PASSWORD=%s\n", os.Getenv("DB_PASSWORD"))
	fmt.Printf("DB_NAME=%s\n", os.Getenv("DB_NAME"))

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Show loaded config
	fmt.Println("\nLoaded configuration:")
	fmt.Printf("DB Host: %s\n", cfg.Database.Host)
	fmt.Printf("DB Port: %s\n", cfg.Database.Port)
	fmt.Printf("DB User: %s\n", cfg.Database.User)
	fmt.Printf("DB Name: %s\n", cfg.Database.Name)
	fmt.Printf("Redis Host: %s\n", cfg.Redis.Host)
	fmt.Printf("Redis Port: %s\n", cfg.Redis.Port)
}
