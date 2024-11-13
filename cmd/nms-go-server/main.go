package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/jim-ww/nms-go/internal/middleware"
	"github.com/jim-ww/nms-go/internal/templates"
)

func main() {

	http.Handle("/web/static/", http.StripPrefix("/web/static/", http.FileServer(http.Dir("./web/static"))))

	http.HandleFunc("/login", middleware.LogMiddleware(LoginPage))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func LoginPage(w http.ResponseWriter, r *http.Request) {

	tmplPath := "web/templates/login.html"
	t := template.Must(template.New(tmplPath).ParseFiles(tmplPath))

	err := t.ExecuteTemplate(w, "login.html", templates.NewLoginData("Login", "Register", "Don't have an account yet?", false))
	if err != nil {
		fmt.Printf("error executing template: %v\n", err)
		return
	}

	fmt.Fprint(w, err)
}
