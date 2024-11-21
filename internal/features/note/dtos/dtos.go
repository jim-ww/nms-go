package dtos

import "time"

type NoteSummaryDTO struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewSummaryDTO(id int64, title string, createdAt time.Time, updatedAt time.Time) *NoteSummaryDTO {
	return &NoteSummaryDTO{
		ID:        id,
		Title:     title,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}

type NoteDetailDTO struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Contents  string    `json:"contents"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewDetailDTO(id int64, title string, contents string, createdAt time.Time, updatedAt time.Time) *NoteDetailDTO {
	return &NoteDetailDTO{
		ID:        id,
		Title:     title,
		Contents:  contents,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
