package main

import (
	"context"
	"flag"
	"log"
	"os"
	"ucrm/migrations"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var (
		command = flag.String("command", "up", "Migration command: up, down, status")
	)
	flag.Parse()

	if err := godotenv.Load(".env.local"); err != nil {
		log.Printf("Warning: failed to load .env.local: %v", err)
	}

	dbPath := os.Getenv("DB_DSN")
	if dbPath == "" {
		log.Fatal("DB_DSN environment variable is required")
	}

	db, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	migrator := migrations.NewMigrator(db)

	migrator.AddMigration(migrations.CreateUsersTable())
	migrator.AddMigration(migrations.CreateClientApplicationsTable())

	ctx := context.Background()
	
	switch *command {
	case "up":
		if err := migrator.Up(ctx); err != nil {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		log.Println("All migrations applied successfully")
	case "down":
		if err := migrator.Down(ctx); err != nil {
			log.Fatalf("Failed to rollback migration: %v", err)
		}
		log.Println("Migration rolled back successfully")
	case "status":
		if err := migrator.Status(ctx); err != nil {
			log.Fatalf("Failed to get migration status: %v", err)
		}
	default:
		log.Fatalf("Unknown command: %s. Use: up, down, or status", *command)
	}
}
