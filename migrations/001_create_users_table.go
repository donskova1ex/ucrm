package migrations

import (
	"ucrm/internal/user"

	"gorm.io/gorm"
)

func CreateUsersTable() Migration {
	return Migration{
		Version:     "001",
		Description: "Create users table",
		Up: func(db *gorm.DB) error {
			return db.AutoMigrate(&user.User{})
		},
		Down: func(db *gorm.DB) error {
			return db.Migrator().DropTable(&user.User{})
		},
	}
}
