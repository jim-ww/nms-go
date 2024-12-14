package dtos

import (
	"time"

	"github.com/google/uuid"
)

type NoteSummaryDTO struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewSummaryDTO(id uuid.UUID, title string, createdAt time.Time, updatedAt time.Time) *NoteSummaryDTO {
	return &NoteSummaryDTO{
		ID:        id,
		Title:     title,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type NoteDetailDTO struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Contents  string    `json:"contents"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewDetailDTO(id uuid.UUID, title string, contents string, createdAt time.Time, updatedAt time.Time) *NoteDetailDTO {
	return &NoteDetailDTO{
		ID:        id,
		Title:     title,
		Contents:  contents,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
