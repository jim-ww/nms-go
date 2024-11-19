package repository

import (
	"errors"

	"github.com/jim-ww/nms-go/internal/features/user"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
)

type UserRepository interface {
	IsUsernameTaken(username string) (taken bool, err error)
	IsEmailTaken(email string) (taken bool, err error)
	CreateUser(username, email, hashedPassword string, role user.Role) (createdID int64, err error)
}
