package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/maynguyen24/sever/configs"
	"github.com/maynguyen24/sever/internal/database"
	"github.com/maynguyen24/sever/internal/server"
	"github.com/maynguyen24/sever/pkg/logger"
	"github.com/maynguyen24/sever/pkg/snowflake"
)

func main() {
	// 1. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, continuing with system environments")
	}

	// Load Configuration centralized
	cfg := configs.LoadConfig()

	// Initialize logger
	logger.InitLogger()
	defer logger.Log.Sync()

	// Initialize ID generator node (machine ID 1)
	snowflake.Init(1)

	// 3. Connect to Database (with error handling)
	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	log.Println("Database connected successfully!")
	defer database.Close()

	// Initialize the server and all routes using Functional Options pattern
	srv := server.NewServer(
		server.WithPort(cfg.Port),
		server.WithConfig(cfg),
	)
	log.Fatalf("Server exited: %v", srv.Start())
}
