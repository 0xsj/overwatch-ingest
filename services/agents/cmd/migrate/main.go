// services/agents/cmd/migrate/main.go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/0xsj/scout/services/agents/config"
	"github.com/0xsj/scout/services/agents/internal/infrastructure/persistence/postgres"
	_ "github.com/lib/pq"
)

func main() {
	// Parse flags
	up := flag.Bool("up", false, "Run migrations up")
	down := flag.Bool("down", false, "Run migrations down")
	create := flag.Bool("create", false, "Create database if it doesn't exist")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load(false)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	pgConfig := cfg.Postgres()

	// If create flag is set, create the database first
	if *create {
		if err := createDatabase(pgConfig); err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		log.Printf("✓ Database '%s' created or already exists", pgConfig.Database())
	}

	// Connect to the agents database
	db, err := sql.Open("postgres", pgConfig.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping to verify connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Run migrations
	if *up {
		log.Println("Running migrations up...")
		if err := postgres.MigrateUp(db); err != nil {
			log.Fatalf("Failed to run migrations up: %v", err)
		}
		log.Println("✓ Migrations up completed")
	} else if *down {
		log.Println("Running migrations down...")
		if err := postgres.MigrateDown(db); err != nil {
			log.Fatalf("Failed to run migrations down: %v", err)
		}
		log.Println("✓ Migrations down completed")
	} else {
		log.Println("Usage:")
		log.Println("  -create    Create database if it doesn't exist")
		log.Println("  -up        Run migrations up")
		log.Println("  -down      Run migrations down")
		os.Exit(1)
	}
}

// createDatabase creates the database if it doesn't exist
func createDatabase(cfg config.PostgresConfig) error {
	// Connect to default postgres database
	defaultDSN := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=postgres sslmode=disable",
		cfg.Host(), cfg.Port(), cfg.User(), cfg.Password(),
	)

	db, err := sql.Open("postgres", defaultDSN)
	if err != nil {
		return fmt.Errorf("failed to connect to postgres: %w", err)
	}
	defer db.Close()

	// Check if database exists
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(query, cfg.Database()).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if exists {
		log.Printf("Database '%s' already exists", cfg.Database())
		return nil
	}

	// Create database
	createQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.Database())
	if _, err := db.Exec(createQuery); err != nil {
		return fmt.Errorf("failed to create database: %w", err)
	}

	return nil
}