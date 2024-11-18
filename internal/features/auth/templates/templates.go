package templates

import "html/template"

var LoginTmpl *template.Template = template.Must(template.ParseFiles("internal/features/auth/templates/auth.html"))

type login struct {
	Title                 string
	Username              string
	Email                 string
	IsRegisterForm        bool
	AlternativeAction     string
	AlternativeActionText string
	ValidationErrors      map[string][]string
}

func NewLoginFormData(username string, validationErrors map[string][]string) *login {
	return &login{
		Title:                 "Login",
		Username:              username,
		Email:                 "",
		AlternativeAction:     "Register",
		AlternativeActionText: "Don't have an account yet?",
		IsRegisterForm:        false,
		ValidationErrors:      validationErrors,
	}
}

func NewRegisterFormData(username, email string, validationErrors map[string][]string) *login {
	return &login{
		Title:                 "Register",
		Username:              username,
		Email:                 email,
		AlternativeAction:     "Login",
		AlternativeActionText: "Already have an account?",
		IsRegisterForm:        true,
		ValidationErrors:      validationErrors,
	}
}
