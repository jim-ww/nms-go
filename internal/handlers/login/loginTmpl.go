package login

import (
	"html/template"
	"log"
	"net/http"

	"github.com/jim-ww/nms-go/internal/templates"
)

var tmpl = template.Must(template.ParseFiles("web/templates/login.html"))

var loginData = templates.NewLogin("Login", "Register", "Don't have an account yet?", false)
var registerData = templates.NewLogin("Register", "Login", "Already have an account?", true)

func LoginTmpl(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "login.html", templates.NewLogin("Login", "Register", "Don't have an account yet?", false))
	if err != nil {
		log.Fatal("failed to execute template:", err)
	}
}

func RegisterTmpl(w http.ResponseWriter, r *http.Request) {
	// if r.Header.Get("HX-Request") == "true" {
	if err := tmpl.ExecuteTemplate(w, "login.html", registerData); err != nil {
		http.Error(w, "Could not load page", http.StatusInternalServerError)
	}
	// } else {
	// 	// templatePath = "web/templates/base.html"
	// 	var content bytes.Buffer
	// 	if err := tmpl.Execute(&content, registerData); err != nil {
	// 		http.Error(w, "Could not load content", http.StatusInternalServerError)
	// 		return
	// 	}
	// 	data := templates.NewBase("Register", template.HTML(content.String()))
	// 	if err := baseTmpl.Execute(w, data); err != nil {
	// 		http.Error(w, "Could not load page", http.StatusInternalServerError)
	// 	}
	// }
}
