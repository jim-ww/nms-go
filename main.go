package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/login", LoggingMiddleware(LoginPage))
	log.Fatal(http.ListenAndServe(":3434", nil))
}

func LoginPage(w http.ResponseWriter, r *http.Request) {

	tmplPath := "views/login.html"
	t := template.Must(template.New(tmplPath).ParseFiles(tmplPath))

	err := t.ExecuteTemplate(w, "login.html", AuthTmpl{"Login", false, "Register", "Don't have an account yet?"})
	if err != nil {
		fmt.Printf("error executing template: %v\n", err)
		return
	}

	fmt.Fprint(w, err)
}

func AuthLogin(w http.ResponseWriter, r *http.Request) {

}

type AuthTmpl struct {
	Title                 string
	IsRegisterForm        bool
	AlternativeAction     string
	AlternativeActionText string
}

func LoggingMiddleware(f func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s to %s from %s \n", r.Method, r.RequestURI, r.RemoteAddr)
		f(w, r)
	}
}

type JWTService interface {
	NewSession(userID int64, roleName string) (jwt string, err error)
	DecodeSession(jwt string) (Session, error)
}
