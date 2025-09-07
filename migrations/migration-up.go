package main

import (
	"fmt"
	"log"
	"os"
	"ucrm/internal/application"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// TODO: Добавить модель User{email, password и т.п.} и добавить миграцию на эту таблицу
func main() {
	err := godotenv.Load(".env.local")
	if err != nil {
		panic(err)
	}
	log.Print("Successfully loaded .env")

	dbPath := os.Getenv("DB_DSN")

	db, err := gorm.Open(postgres.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {

	}

	err = initTable(db, &application.Application{})
	log.Print("Successfully migrated")
}

func initTable(db *gorm.DB, tables ...any) error {
	for _, table := range tables {
		err := db.AutoMigrate(table)
		if err != nil {
			err := db.Migrator().DropTable(table)
			if err != nil {
				return fmt.Errorf("failed dropping table organizations: %w", err)
			}
			return fmt.Errorf("failed creating table organizations: %w", err)
		}
	}
	return nil
}
