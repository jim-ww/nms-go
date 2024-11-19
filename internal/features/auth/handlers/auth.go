package handlers

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jim-ww/nms-go/internal/features/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	"github.com/jim-ww/nms-go/internal/features/auth/templates"
	"github.com/jim-ww/nms-go/pkg/utils/handlers"
)

type AuthHandler struct {
	authService *auth.AuthService
	logger      *slog.Logger
	tmpl        *template.Template
	tmplHandler *handlers.TmplHandler
}

func NewAuthHandler(userService *auth.AuthService, log *slog.Logger, tmplHandler *handlers.TmplHandler) *AuthHandler {
	templ := template.Must(template.ParseFiles("web/templates/auth.html"))
	return &AuthHandler{
		authService: userService,
		logger:      log,
		tmpl:        templ,
		tmplHandler: tmplHandler,
	}
}

func (auth *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handlers.RenderError(w, r, "Unable to parse form", http.StatusBadRequest)
		return
	}
	dto := dtos.NewLoginDTO(r.FormValue("username"), r.FormValue("password"))

	fmt.Println(dto)
}

func (lh *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		handlers.RenderError(w, r, "Unable to parse form", http.StatusBadRequest)
		return
	}
	dto := dtos.NewRegisterDTO(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))

	lh.logger.Info("got register dto", slog.String("dto-username", dto.Username), slog.String("dto-email", dto.Email))

	token, validationErrors, err := lh.authService.RegisterUser(dto)
	if err != nil {
		handlers.RenderError(w, r, "Registration error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, auth.NewTokenCookie(token))

	data := templates.NewRegisterFormData(dto.Username, dto.Email, validationErrors)
	lh.tmplHandler.RenderTemplate(w, r, lh.tmpl, data)
}
