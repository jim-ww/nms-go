package note

import "time"

type Note struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Contents  string    `json:"contents"`
	UserID    int64     `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
