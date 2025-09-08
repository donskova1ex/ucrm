package migrations

import (
	"ucrm/internal/client_application"

	"gorm.io/gorm"
)

func CreateClientApplicationsTable() Migration {
	return Migration{
		Version:     "002",
		Description: "Create client applications table",
		Up: func(db *gorm.DB) error {
			if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
				return err
			}
			return db.AutoMigrate(&client_application.ClientApplication{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&client_application.ClientApplication{})
		},
	}
}
