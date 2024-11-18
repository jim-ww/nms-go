package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/jim-ww/nms-go/internal/features/auth"
	"github.com/jim-ww/nms-go/internal/features/auth/dtos"
	"github.com/jim-ww/nms-go/internal/features/auth/templates"
	"github.com/jim-ww/nms-go/pkg/utils/handlers"
)

type AuthHandler struct {
	*auth.AuthService
	*slog.Logger
}

func New(userService *auth.AuthService, log *slog.Logger) *AuthHandler {
	return &AuthHandler{AuthService: userService, Logger: log}
}

func (auth *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	dto := dtos.NewLoginDTO(r.FormValue("username"), r.FormValue("password"))

	fmt.Println(dto)
}

func (lh *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	dto := dtos.NewRegisterDTO(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))

	lh.Logger.Info("got dto:", slog.Any("RegisterDTO", dto))

	token, validationErrors, err := lh.AuthService.RegisterUser(dto)
	if err != nil {
		http.Error(w, "Registration error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, auth.NewTokenCookie(token))

	data := templates.NewRegisterFormData(dto.Username, dto.Email, validationErrors)

	handlers.RenderTemplate(w, templates.LoginTmpl, data)
}
