package migrations

import (
	"database/sql"
	"embed"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sqlite3
var sqlite3Migrations embed.FS

func MustMigrate(db *sql.DB) {

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}

	src, err := iofs.New(sqlite3Migrations, "sqlite3")
	if err != nil {
		log.Fatalf("Failed to init new Driver from io/fs#FS, err: %v", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "encore", driver)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
}
