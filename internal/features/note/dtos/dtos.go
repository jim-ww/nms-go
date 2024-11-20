package dtos

import "time"

type NoteSummaryDTO struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type NoteDetailDTO struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Contents  string    `json:"contents"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
