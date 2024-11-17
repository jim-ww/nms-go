package auth

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jim-ww/nms-go/internal/auth/templates"
	"github.com/jim-ww/nms-go/pkg/utils/handlers"
)

// TODO split on 2 handler files?

type AuthHandler struct {
	*AuthService
	*slog.Logger
}

func New(userService *AuthService, log *slog.Logger) *AuthHandler {
	return &AuthHandler{AuthService: userService, Logger: log}
}

var tmpl = template.Must(template.ParseFiles("web/templates/login.html"))

var loginData = templates.NewLogin("Login", "", "", "Register", "Don't have an account yet?", false, nil)
var registerData = templates.NewLogin("Register", "", "", "Login", "Already have an account?", true, nil)

func LoginTmpl(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, loginData); err != nil {
		fmt.Println(err)
		http.Error(w, "could not load page", http.StatusInternalServerError)
	}
}

func RegisterTmpl(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.Execute(w, registerData); err != nil {
		fmt.Println(err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
	}
}

func (auth *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	dto := LoginDTO{
		Username: username,
		Password: password,
	}
	fmt.Println(dto)

	// TODO validate

	// TODO add authentication logic

}

func (lh *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Println(err)
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	dto := &RegisterDTO{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}

	lh.Logger.Info("got dto:", slog.Any("RegisterDTO", dto))

	token, validationErrors, err := lh.AuthService.RegisterUser(dto)
	if err != nil {
		http.Error(w, "Registration error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, NewTokenCookie(token))

	data := templates.NewLogin(registerData.Title, dto.Username, dto.Email, registerData.AlternativeAction, registerData.AlternativeActionText, registerData.IsRegisterForm, validationErrors)

	handlers.RenderTemplate(w, tmpl, data)
}
