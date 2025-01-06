package note

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jim-ww/nms-go/internal/repository"
	"github.com/jim-ww/nms-go/pkg/page"
	"github.com/jmoiron/sqlx"
)

var validOrderByFields = map[string]bool{
	"title":      true,
	"created_at": true,
	"updated_at": true,
}

type NoteService struct {
	repo *repository.Queries
	db   *sqlx.DB
}

func New(repo *repository.Queries, db *sqlx.DB) *NoteService {
	return &NoteService{
		repo: repo,
		db:   db,
	}
}

func (s NoteService) GetUserNotes(c context.Context, userID uuid.UUID, params page.PaginationParams) ([]*repository.Note, *page.Page, error) {
	if _, fieldValid := validOrderByFields[params.OrderBy]; !fieldValid {
		return nil, nil, errors.New("invalid orderBy field")
	}
	if params.PageSize <= 0 {
		return nil, nil, errors.New("page size must be greater than zero")
	}

	var notes []*repository.Note
	limit := params.PageSize
	offset := params.PageSize * params.PageNumber
	query := fmt.Sprintf(`SELECT * FROM notes WHERE user_id = ? ORDER BY %s %s LIMIT ? OFFSET ?`, params.OrderBy, params.Order)
	if err := s.db.SelectContext(c, &notes, query, userID, limit, offset); err != nil {
		return nil, nil, fmt.Errorf("failed to select notes: %w", err)
	}

	notesCount, err := s.repo.GetUserNotesCount(c, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get notes count: %w", err)
	}

	totalPages := (int(notesCount) + params.PageSize - 1) / params.PageSize // int(notesCount)/params.PageSize
	page := page.NewPage(totalPages, int(notesCount), params.PageNumber)

	return notes, page, nil
}
