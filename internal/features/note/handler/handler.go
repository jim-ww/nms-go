package handler

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jim-ww/nms-go/internal/features/note/services/note"
	"github.com/jim-ww/nms-go/internal/templates"
	"github.com/jim-ww/nms-go/pkg/page"
	"github.com/labstack/echo/v4"
)

type NoteHandler struct {
	srv *note.NoteService
}

func New(service *note.NoteService) *NoteHandler {
	return &NoteHandler{
		srv: service,
	}
}

// TODO implement pagination via url params
// TODO fetch notes
// TODO fetch user profile
func (h *NoteHandler) Dashboard(c echo.Context) error {
	ctxUserID := c.Get("user_id")
	userID, ok := ctxUserID.(uuid.UUID)
	if !ok {
		return fmt.Errorf("%w: ctxUserID: %s", errors.New("failed to get/convert userID from echo context"), ctxUserID)
	}
	c.Logger().Debug("userID ", userID)

	paginationParams := page.NewPaginationParams(1, 10, page.DESC, "title")

	// TODO handle page
	notes, _, err := h.srv.GetUserNotes(c.Request().Context(), userID, *paginationParams)
	if err != nil {
		return err
	}

	tmpl := templates.Layout("Dashboard", templates.Dashboard(notes))
	return tmpl.Render(c.Request().Context(), c.Response().Writer)
}
