package handlers

import (
	"fmt"
	"net/http"

	tmpl "github.com/jim-ww/nms-go/internal/features/auth/templates"
)

var loginData = tmpl.NewLoginFormData("", map[string][]string{})
var registerData = tmpl.NewRegisterFormData("", "", map[string][]string{})

func LoginTmpl(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.LoginTmpl.Execute(w, loginData); err != nil {
		fmt.Println(err)
		http.Error(w, "could not load page", http.StatusInternalServerError)
	}
}

func RegisterTmpl(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.LoginTmpl.Execute(w, registerData); err != nil {
		fmt.Println(err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
	}
}
