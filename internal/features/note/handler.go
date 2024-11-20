package note

import (
	"log/slog"
	"net/http"
)

type Handler struct {
	logger *slog.Logger
	// noteService *NoteService
	// userService *UserService
}

func NewHandler(logger *slog.Logger) *Handler { //noteService *NoteService, userService *UserService) *Handler {
	return &Handler{
		logger: logger,
		// noteService: noteService,
		// userService: userService,
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from dashboard"))
}
