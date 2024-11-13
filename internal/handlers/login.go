package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jim-ww/nms-go/internal/dtos"
	"github.com/jim-ww/nms-go/internal/templates"
	"github.com/jim-ww/nms-go/internal/validators"
)

var tmplPath = "web/templates/login.html"
var tmpl = template.Must(template.New(tmplPath).ParseFiles(tmplPath))

func LoginTmpl(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "login.html", templates.NewLoginData("Login", "Register", "Don't have an account yet?", false))
	if err != nil {
		log.Fatal("failed to execute template:", err)
	}
}

func RegisterTmpl(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "login.html", templates.NewLoginData("Register", "Login", "Already have an account?", true))
	if err != nil {
		log.Fatal("failed to execute template:", err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")

	dto := dtos.LoginDTO{
		Username: username,
		Password: password,
	}
	fmt.Println("loginData:", dto)

	// TODO validate

	// TODO add authentication logic

}

func Register(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	dto := &dtos.RegisterDTO{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
		Email:    r.FormValue("email"),
	}
	validators.ValidateRegisterDTO(dto)
	// TODO validate

	// TODO check if username or email already exists

	// TODO add authentication logic
}
