package user

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	ROLE_USER  Role = "user"
	ROLE_ADMIN Role = "admin"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
