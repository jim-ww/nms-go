package templates

type Login struct {
	Title                 string
	IsRegisterForm        bool
	AlternativeAction     string
	AlternativeActionText string
}

func NewLogin(title, alternativeAction, alternativeActionText string, isRegisterForm bool) *Login {
	return &Login{
		Title:                 title,
		AlternativeAction:     alternativeAction,
		AlternativeActionText: alternativeActionText,
		IsRegisterForm:        isRegisterForm,
	}
}
