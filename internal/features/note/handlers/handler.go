package note

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"github.com/jim-ww/nms-go/internal/features/note/dtos"
	noteTempl "github.com/jim-ww/nms-go/internal/features/note/template"
	userDTO "github.com/jim-ww/nms-go/internal/features/user/dtos"
	"github.com/jim-ww/nms-go/pkg/utils/handlers"
)

type Handler struct {
	logger       *slog.Logger
	templ        *template.Template
	templHandler *handlers.TmplHandler
	// noteService *NoteService
	// userService *UserService
}

func New(logger *slog.Logger, templHandler *handlers.TmplHandler) *Handler { //noteService *NoteService, userService *UserService) *Handler {
	templ := template.Must(template.ParseFiles("web/templates/dashboard.html"))
	return &Handler{
		logger:       logger,
		templ:        templ,
		templHandler: templHandler,
		// noteService: noteService,
		// userService: userService,
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {

	// TODO remove this mock data
	notes := []*dtos.NoteSummaryDTO{
		dtos.NewSummaryDTO(1, "Some Title", time.Now(), time.Now()),
		dtos.NewSummaryDTO(2, "Second", time.Now(), time.Now()),
		dtos.NewSummaryDTO(3, "secret", time.Now(), time.Now()),
	}
	selectedNote := dtos.NewDetailDTO(notes[0].ID, notes[0].Title, "Hello World!", notes[0].CreatedAt, notes[0].UpdatedAt)
	user := userDTO.New("bob", "bob@bob.com", len(notes))

	// TODO
	h.logger.Debug("rendering dashboard...")
	_ = noteTempl.New(notes, selectedNote, user)
	h.templHandler.RenderTemplate(w, r, h.templ, nil)
}
