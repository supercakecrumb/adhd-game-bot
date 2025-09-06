package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
)

func main() {
	// Connect to PostgreSQL database
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "user=postgres password=password dbname=adhd_bot sslmode=disable"
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create schema_migrations table if not exists
	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	// Get list of migration files
	migrations, err := filepath.Glob("internal/infra/postgres/migrations/*.sql")
	if err != nil {
		log.Fatalf("Failed to list migrations: %v", err)
	}

	// Apply migrations in order
	for _, file := range migrations {
		version := filepath.Base(file)[:3] // Get first 3 chars (001, 002 etc.)
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)", version).Scan(&exists)
		if err != nil {
			log.Fatalf("Failed to check migration status: %v", err)
		}

		if !exists {
			fmt.Printf("Applying migration: %s\n", file)
			migrationSQL, err := os.ReadFile(file)
			if err != nil {
				log.Fatalf("Failed to read migration file: %v", err)
			}

			if _, err := db.Exec(string(migrationSQL)); err != nil {
				log.Fatalf("Failed to apply migration %s: %v", file, err)
			}

			if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
				log.Fatalf("Failed to record migration: %v", err)
			}
		}
	}

	fmt.Println("Migrations applied successfully")
}
