package sqlite

import "database/sql"

type SqliteRepository struct {
	*sql.DB
}

func (ur SqliteRepository) IsUsernameTaken(username string) (exists bool, err error) {
	stmt, err := ur.Prepare(`SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// IsEmailTaken(email string) bool
// CreateUser(username, email, hashedPassword string, role models.Role) (createdUserID int64, err error)
