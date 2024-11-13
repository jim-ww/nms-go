package repositories

import (
	"github.com/jim-ww/nms-go/internal/dtos"
)

type UserRepository interface {
	IsUsernameTaken(username string) bool
	IsEmailTaken(email string) bool
	CreateUser(dto *dtos.RegisterDTO) error
}

type NoteRepository interface {
}

type RoleRepository interface {
}
