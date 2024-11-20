package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jim-ww/nms-go/internal/features/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	authService "github.com/jim-ww/nms-go/internal/features/auth/services/auth"
	jwtService "github.com/jim-ww/nms-go/internal/features/auth/services/jwt"
	"github.com/jim-ww/nms-go/pkg/utils/handlers"
	"github.com/jim-ww/nms-go/pkg/utils/loggers/sl"
)

type AuthHandler struct {
	authService *authService.AuthService
	jwtService  *jwtService.JWTService
	logger      *slog.Logger
	tmpl        *template.Template
	tmplHandler *handlers.TmplHandler
}

func NewAuthHandler(userService *authService.AuthService, log *slog.Logger, tmplHandler *handlers.TmplHandler) *AuthHandler {
	templ := template.Must(template.ParseFiles("web/templates/auth.html"))
	return &AuthHandler{
		authService: userService,
		logger:      log,
		tmpl:        templ,
		tmplHandler: tmplHandler,
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		ah.logger.Error("unable to parse form", sl.Err(err))
		handlers.RenderError(w, r, "unable to parse form", http.StatusBadRequest)
		return
	}
	dto := dtos.NewLoginDTO(r.FormValue("username"), r.FormValue("password"))
	ah.logger.Debug("got login dto, executing authService.LoginUser()", sl.LoginDTO(dto))

	token, validationErrors, err := ah.authService.LoginUser(dto)
	if err != nil {
		ah.logger.Error("Failed to execute authService.LoginUser()", sl.Err(err))
		handlers.RenderError(w, r, "Registration error", http.StatusInternalServerError)
		return
	}

	// if has errors, return validation errors to form
	if validationErrors.HasErrors() {
		ah.logger.Debug("dto has validation errors, returning them to login form", slog.Any("validationErrors", validationErrors))
		data := auth.NewLoginFormData(dto.Username, validationErrors.TranslateValidationErrors())
		ah.tmplHandler.RenderTemplate(w, r, ah.tmpl, data)
		return
	}

	// if no errors, set token cookie and redirect to home page
	ah.logger.Debug("Setting token cookie")
	http.SetCookie(w, ah.jwtService.NewTokenCookie(token))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (lh *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		lh.logger.Error("Unable to parse form", sl.Err(err))
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

	// if has errors, return validation errors to form
	if validationErrors.HasErrors() {
		data := auth.NewRegisterFormData(dto.Username, dto.Email, validationErrors.TranslateValidationErrors())
		lh.tmplHandler.RenderTemplate(w, r, lh.tmpl, data)
		return
	}

	// if no errors, set token cookie and redirect to home page
	lh.logger.Debug("Setting token cookie")
	http.SetCookie(w, lh.jwtService.NewTokenCookie(token))
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
