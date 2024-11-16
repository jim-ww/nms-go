package templates

type Login struct {
	Title                 string
	Username              string
	Email                 string
	IsRegisterForm        bool
	AlternativeAction     string
	AlternativeActionText string
	ValidationErrors      map[string][]string
}

func NewLogin(title, username, email, alternativeAction, alternativeActionText string, isRegisterForm bool, validationErrors map[string][]string) *Login {
	return &Login{
		Title:                 title,
		Username:              username,
		Email:                 email,
		AlternativeAction:     alternativeAction,
		AlternativeActionText: alternativeActionText,
		IsRegisterForm:        isRegisterForm,
		ValidationErrors:      validationErrors,
	}
}
