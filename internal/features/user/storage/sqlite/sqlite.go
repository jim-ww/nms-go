package sqlite

import (
	"database/sql"
	"errors"
	"log/slog"
	"time"

	"github.com/jim-ww/nms-go/internal/features/user"
	"github.com/jim-ww/nms-go/internal/features/user/storage"
	"github.com/jim-ww/nms-go/pkg/utils/loggers/sl"
	"github.com/mattn/go-sqlite3"
)

type UserRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func New(logger *slog.Logger, db *sql.DB) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (repo UserRepository) Migrate() error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT(30) NOT NULL UNIQUE,
			email TEXT(255) NOT NULL UNIQUE,
			password TEXT(255) NOT NULL,
			role TEXT(10) NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
	);`
	createIndexUsernameSQL := `CREATE INDEX IF NOT EXISTS idx_username ON users(username);`
	createIndexEmailSQL := `CREATE INDEX IF NOT EXISTS idx_email ON users(email);`

	if _, err := repo.db.Exec(createTableSQL); err != nil {
		repo.logger.Error("failed to migrate table users, failed to execute query 'createTableSQL'", sl.Err(err))
		return err
	}

	if _, err := repo.db.Exec(createIndexUsernameSQL); err != nil {
		repo.logger.Error("failed to migrate table users, failed to execute query 'createIndexUsernameSQL'", sl.Err(err))
		return err
	}
	if _, err := repo.db.Exec(createIndexEmailSQL); err != nil {
		repo.logger.Error("failed to migrate table users, failed to execute query 'createIndexEmailSQL'", sl.Err(err))
		return err
	}
	return nil
}

func (repo UserRepository) IsUsernameTaken(username string) (taken bool, err error) {
	repo.logger.Debug("executing IsUsernameTaken() sqlite query...")

	stmt, err := repo.db.Prepare(`SELECT EXISTS(SELECT id FROM users WHERE username = ?)`)
	if err != nil {
		repo.logger.Error("failed to initialize stmt IsUsernameTaken()", sl.Err(err))
		return false, err
	}

	err = stmt.QueryRow(username).Scan(&taken)
	if err != nil {
		repo.logger.Error("failed to execute stmt IsUsernameTaken()", sl.Err(err))
		return false, err
	}

	repo.logger.Debug("executed IsUsernameTaken() sqlite query succesfully", slog.String("username", username), slog.Bool("taken", taken))

	return taken, nil
}

func (repo UserRepository) IsEmailTaken(email string) (taken bool, err error) {
	repo.logger.Debug("executing IsEmailTaken() sqlite query...")

	stmt, err := repo.db.Prepare(`SELECT EXISTS(SELECT id FROM users WHERE email = ?)`)
	if err != nil {
		repo.logger.Error("failed to initialize stmt IsEmailTaken()", sl.Err(err))
		return false, err
	}

	err = stmt.QueryRow(email).Scan(&taken)
	if err != nil {
		repo.logger.Error("failed to execute stmt IsEmailTaken()", sl.Err(err))
		return false, err
	}

	repo.logger.Debug("executed IsEmailTaken() sqlite query succesfully", slog.String("email", email), slog.Bool("taken", taken))

	return taken, nil
}

func (repo UserRepository) Create(username, email, hashedPassword string, role user.Role) (createdID int64, err error) {
	stmt, err := repo.db.Prepare(`INSERT INTO users (username, email, password, role, created_at, updated_at) VALUES (?,?,?,?,?,?)`)
	if err != nil {
		repo.logger.Error("Failed to prepare stmt CreateUser()", sl.Err(err))
		return 0, err
	}

	now := time.Now()
	result, err := stmt.Exec(username, email, hashedPassword, role, now, now)
	if err != nil {

		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			repo.logger.Debug("Failed to execute stmt CreateUser(), user already exists", sl.Err(err))
			return 0, storage.ErrUserAlreadyExists
		}

		repo.logger.Error("Failed to execute stmt CreateUser()", sl.Err(err))
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		repo.logger.Error("Failed to get lastInsertId from CreateUser() stmt", sl.Err(err))
		return 0, err
	}

	return userID, nil
}

func (repo UserRepository) GetByUsername(username string) (user user.User, err error) {
	repo.logger.Debug("executing GetUserByUsername() sqlite query...")
	stmt, err := repo.db.Prepare(`SELECT id, username, email, password, role, created_at, updated_at FROM users WHERE username = ?`)
	if err != nil {
		repo.logger.Error("failed to prepare stmt GetUserByUsername()", sl.Err(err))
		return user, err
	}

	if err = stmt.QueryRow(username).Scan(&user.ID, &user.Username, &user.Email, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {

		// TODO test
		if errors.Is(err, sql.ErrNoRows) {
			repo.logger.Debug("username does not exist", sl.Err(err))
			return user, storage.ErrUsernameDoesNotExist
		}

		repo.logger.Error("failed to execute stmt GetUserByUsername()", sl.Err(err))
		return user, err
	}
	repo.logger.Debug("executed GetUserByUsername() sqlite query succesfully", slog.String("username", username))

	return user, nil
}
