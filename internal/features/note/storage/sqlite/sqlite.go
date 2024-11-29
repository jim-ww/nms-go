package sqlite

import (
	"database/sql"
	"log/slog"

	"github.com/jim-ww/nms-go/pkg/utils/loggers/sl"
)

type NoteRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func New(logger *slog.Logger, db *sql.DB) *NoteRepository {
	return &NoteRepository{
		logger: logger,
		db:     db,
	}
}

func (nr NoteRepository) Migrate() error {
	createNotesTableSQL := `CREATE TABLE IF NOT EXISTS notes(
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT,
		user_id INTEGER NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`
	createIndexNotesSQL := `CREATE INDEX IF NOT EXISTS idx_title ON notes(title)`

	_, err := nr.db.Exec(createNotesTableSQL)
	if err != nil {
		nr.logger.Error("failed to migrate table notes, failed to execute createNotesTableSQL", sl.Err(err))
		return err
	}

	_, err = nr.db.Exec(createIndexNotesSQL)
	if err != nil {
		nr.logger.Error("failed to migrate table notes, failed to execute createIndexNotesSQL", sl.Err(err))
		return err
	}

	return nil
}
