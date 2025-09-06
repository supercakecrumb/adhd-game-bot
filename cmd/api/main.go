package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	http_server "github.com/supercakecrumb/adhd-game-bot/internal/infra/http"
	"github.com/supercakecrumb/adhd-game-bot/internal/infra/postgres"
	"github.com/supercakecrumb/adhd-game-bot/internal/usecase"
)

func main() {
	// Load configuration
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	questRepo := postgres.NewQuestRepository(db)
	dungeonRepo := postgres.NewDungeonRepository(db)
	dungeonMemberRepo := postgres.NewDungeonMemberRepository(db)

	// Initialize other components
	uuidGen := postgres.NewUUIDGenerator()
	scheduler := postgres.NewPgScheduler(db)
	txManager := postgres.NewTxManager(db)
	idempotencyRepo := postgres.NewIdempotencyRepository(db)

	// Initialize use case services
	questService := usecase.NewQuestService(questRepo, userRepo, uuidGen, scheduler, idempotencyRepo, txManager)
	dungeonService := usecase.NewDungeonService(dungeonRepo, dungeonMemberRepo, userRepo, uuidGen, txManager)

	// Initialize HTTP server
	server := http_server.NewServer(questService, dungeonService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: server.Router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
