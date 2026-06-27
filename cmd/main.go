package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/limon/spotsync/internal/config"
	"github.com/limon/spotsync/internal/server"
)

func main() {
	// Load .env file; in production environment variables are set directly
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	cfg := config.Load()

	db := config.ConnectDB(cfg)

	// Start HTTP server
	srv := server.New(cfg, db)
	srv.Start()
}
