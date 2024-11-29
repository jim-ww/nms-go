package storage

import (
	"errors"

	"github.com/jim-ww/nms-go/internal/features/user"
)

var (
	ErrUserAlreadyExists    = errors.New("username or email already exists")
	ErrUsernameDoesNotExist = errors.New("username does not exist")
)

type UserRepository interface {
	IsUsernameTaken(username string) (taken bool, err error)
	IsEmailTaken(email string) (taken bool, err error)
	Create(username, email, hashedPassword string, role user.Role) (createdID int64, err error)
	GetByUsername(username string) (user user.User, err error)
	Migrate() error
	// GetByID(id int64) (user user.User, err error)
	// GetByEmail(email string) (user user.User, err error)
}
