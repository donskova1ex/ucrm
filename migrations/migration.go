package migrations

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	Version     string
	Description string
	Up          func(*gorm.DB) error
	Down        func(*gorm.DB) error
}

type MigrationRecord struct {
	ID          uint      `gorm:"primaryKey"`
	Version     string    `gorm:"uniqueIndex;not null"`
	Description string    `gorm:"not null"`
	AppliedAt   time.Time `gorm:"not null"`
}

type Migrator struct {
	db         *gorm.DB
	migrations []Migration
	tableName  string
}

func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
		tableName:  "schema_migrations",
	}
}

func (m *Migrator) AddMigration(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

func (m *Migrator) addMigrationsTable() error {
	return m.db.AutoMigrate(&MigrationRecord{})
}

func (m *Migrator) getAppliedMigrations() (map[string]bool, error) {
	var records []MigrationRecord
	if err := m.db.Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}

	applied := make(map[string]bool)
	for _, record := range records {
		applied[record.Version] = true
	}
	return applied, nil
}

func (m *Migrator) Up(ctx context.Context) error {
	if err := m.addMigrationsTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	for _, migration := range m.migrations {
		if applied[migration.Version] {
			continue
		}

		fmt.Printf("Applying migration %s: %s\n", migration.Version, migration.Description)

		tx := m.db.Begin()
		if tx.Error != nil {
			return fmt.Errorf("failed to begin transaction: %w", tx.Error)
		}

		if err := migration.Up(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
		}

		record := MigrationRecord{
			Version:     migration.Version,
			Description: migration.Description,
			AppliedAt:   time.Now(),
		}
		if err := tx.Create(&record).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}

		if err := tx.Commit().Error; err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Version, err)
		}

		fmt.Printf("Successfully applied migration %s\n", migration.Version)
	}

	return nil
}

func (m *Migrator) Down(ctx context.Context) error {
	if err := m.addMigrationsTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	var lastRecord MigrationRecord
	if err := m.db.Order("applied_at DESC").First(&lastRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to get last migration: %w", err)
	}

	var migration *Migration
	for i := range m.migrations {
		if m.migrations[i].Version == lastRecord.Version {
			migration = &m.migrations[i]
			break
		}
	}

	if migration == nil {
		return fmt.Errorf("migration %s not found in migration list", lastRecord.Version)
	}

	if migration.Down == nil {
		return fmt.Errorf("migration %s does not support rollback", lastRecord.Version)
	}

	fmt.Printf("Rolling back migration %s: %s\n", migration.Version, migration.Description)

	tx := m.db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	if err := migration.Down(tx); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to rollback migration %s: %w", migration.Version, err)
	}

	if err := tx.Where("version = ?", lastRecord.Version).Delete(&MigrationRecord{}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to remove migration record %s: %w", migration.Version, err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit rollback for migration %s: %w", migration.Version, err)
	}

	fmt.Printf("Successfully rolled back migration %s\n", migration.Version)
	return nil
}

func (m *Migrator) Status(ctx context.Context) error {
	if err := m.addMigrationsTable(); err != nil {
		return fmt.Errorf("failed to ensure migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return err
	}

	fmt.Println("Migration Status:")
	fmt.Println("=================")

	for _, migration := range m.migrations {
		status := "PENDING"
		if applied[migration.Version] {
			status = "APPLIED"
		}
		fmt.Printf("%s - %s: %s\n", migration.Version, status, migration.Description)
	}

	return nil
}
