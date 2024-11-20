package sqlite

import (
	"database/sql"
	"log/slog"

	"github.com/jim-ww/nms-go/internal/features/user"
	"github.com/jim-ww/nms-go/pkg/utils/loggers/sl"
)

type UserRepository struct {
	logger *slog.Logger
	db     *sql.DB
}

func NewUserRepository(logger *slog.Logger, db *sql.DB) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (repo UserRepository) Migrate() {
	_, err := repo.db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
    id INTEGER,
    username TEXT(30) NOT NULL UNIQUE,
    email TEXT(255) NOT NULL UNIQUE,
    password TEXT(255) NOT NULL,
    PRIMARY KEY(id AUTOINCREMENT)
	)`)
	if err != nil {
		repo.logger.Error("failed to migrate table users, failed to execute query", sl.Err(err))
		return
	}
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

func (repo UserRepository) CreateUser(username, email, hashedPassword string, role user.Role) (createdID int64, err error) {
	stmt, err := repo.db.Prepare(`INSERT INTO users (username, email, password) VALUES (?,?,?)`)
	if err != nil {
		repo.logger.Error("Failed to prepare stmt CreateUser()", sl.Err(err))
		return 0, err
	}

	result, err := stmt.Exec(username, email, hashedPassword)
	if err != nil {
		repo.logger.Error("Failed to execute stmt CreateUser()", sl.Err(err))
		// TODO handle UNIQUE CONSTRAINT err, return repository.ErrUserAlreadyExists
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		repo.logger.Error("Failed to get lastInsertId from CreateUser() stmt", sl.Err(err))
		return 0, err
	}

	return userID, nil
}
