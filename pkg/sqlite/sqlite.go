package sqlite

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteStorage struct {
	db *sql.DB
}

func NewSqliteStorage(storagePath string) *sql.DB {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		panic(err)
	}
	return db
}

func (ss SqliteStorage) Migrate() {
	ss.db.Exec(`CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY,
		username NOT NULL,
		email NOT NULL
	)
	
	CREATE TABLE IF NOT EXISTS roles(
	
	)

	CREATE TABLE IF NOT EXISTS notes(
	
	)
	`)
}
