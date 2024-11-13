package templates

type LoginData struct {
	Title                 string
	IsRegisterForm        bool
	AlternativeAction     string
	AlternativeActionText string
}

func NewLoginData(title, alternativeAction, alternativeActionText string, isRegisterForm bool) LoginData {
	return LoginData{
		Title:                 title,
		AlternativeAction:     alternativeAction,
		AlternativeActionText: alternativeActionText,
		IsRegisterForm:        isRegisterForm,
	}
}
