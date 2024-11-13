package handlers

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jim-ww/nms-go/internal/templates"
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
