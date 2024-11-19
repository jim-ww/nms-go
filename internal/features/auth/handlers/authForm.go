package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jim-ww/nms-go/internal/features/auth/templates"
	"github.com/jim-ww/nms-go/pkg/utils/handlers"
)

type AuthFormHandler struct {
	tmpl        *template.Template
	logger      *slog.Logger
	tmplHandler *handlers.TmplHandler
}

var loginData = templates.NewLoginFormData("", map[string][]string{})
var registerData = templates.NewRegisterFormData("", "", map[string][]string{})

func NewAuthFormHandler(logger *slog.Logger, tmplHandler *handlers.TmplHandler) *AuthFormHandler {
	templatePath := "web/templates/auth.html"
	templ, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Error("failed to initialize authFormHandler, failed to parse template", slog.String("template-path", templatePath))
		panic(err)
	}

	return &AuthFormHandler{
		tmpl:        templ,
		logger:      logger,
		tmplHandler: tmplHandler,
	}
}

func (afh AuthFormHandler) LoginTmpl(w http.ResponseWriter, r *http.Request) {
	afh.logger.Debug("rendering LoginTmpl")
	afh.tmplHandler.RenderTemplate(w, r, afh.tmpl, loginData)
}

func (afh AuthFormHandler) RegisterTmpl(w http.ResponseWriter, r *http.Request) {
	afh.logger.Debug("rendering RegisterTmpl")
	afh.tmplHandler.RenderTemplate(w, r, afh.tmpl, registerData)
}
