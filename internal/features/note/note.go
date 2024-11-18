package note

import "time"

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title" validate:"min=3 max=255"`
	Contents  string    `json:"contents" validate:"max=5000"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
