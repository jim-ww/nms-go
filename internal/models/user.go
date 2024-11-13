package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username" validate:"min=3 max=30"`
	Email     string    `json:"email" validate:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
