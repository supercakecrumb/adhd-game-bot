package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
	httpserver "github.com/supercakecrumb/adhd-game-bot/internal/infra/http"
	"github.com/supercakecrumb/adhd-game-bot/internal/infra/postgres"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

func main() {
	// Database connection
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres password=password dbname=adhd_bot sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize PostgreSQL repositories
	userRepo := postgres.NewUserRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	idempotencyRepo := postgres.NewPgIdempotencyRepository(db)
	txManager := postgres.NewTxManager(db)

	// Initialize utility implementations
	uuidGen := postgres.NewUUIDGenerator()
	scheduler := postgres.NewPgScheduler(db)

	// Initialize TaskService
	taskService := usecase.NewTaskService(
		taskRepo,
		userRepo,
		uuidGen,
		scheduler,
		idempotencyRepo,
		txManager,
	)

	// Create and start HTTP server
	server := httpserver.NewServer(taskService)
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting API server on port %s", port)
	if err := http.ListenAndServe(":"+port, server.Router); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
