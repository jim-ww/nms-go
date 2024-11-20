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
	"github.com/jim-ww/nms-go/pkg/utils/loggers/sl"
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

	lh.logger.Debug("got register dto, executing authService.RegisterUser()", sl.RegisterDTO(dto))

	token, validationErrors, err := lh.authService.RegisterUser(dto)
	if err != nil {
		handlers.RenderError(w, r, "Registration error", http.StatusInternalServerError)
		return
	}

	// if no errors, set token cookie and redirect to home page
	if !validationErrors.HasErrors() {
		lh.logger.Debug("Setting token cookie")
		http.SetCookie(w, auth.NewTokenCookie(token))
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// if has errors, return validation errors to form
	data := templates.NewRegisterFormData(dto.Username, dto.Email, validationErrors.TranslateValidationErrors())
	lh.tmplHandler.RenderTemplate(w, r, lh.tmpl, data)
}
