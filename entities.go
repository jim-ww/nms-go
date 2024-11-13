package main

import (
	"time"
)

// type Role int
//
// const (
// 	ROLE_USER Role = iota
// 	ROLE_ADMIN
// )

type Role struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username" validate:"min=3 max=30"`
	Email     string    `json:"email" validate:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"min=3 max=255"`
	Contents  string    `json:"contents" validate:"max=5000"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
