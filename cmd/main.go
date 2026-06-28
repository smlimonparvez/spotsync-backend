package main

import (
	"log"

	"spotsync/internal/config"
	"spotsync/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file; in production environment variables are set directly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	// Build typed config from env vars
	cfg := config.Load()

	// Open database connection using config
	db := config.ConnectDB(cfg)

	// Start HTTP server (registers routes + auto-migrates)
	srv := server.New(cfg, db)
	srv.Start()
}
