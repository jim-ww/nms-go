package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

type TmplHandler struct {
	logger *slog.Logger
}

func NewTmplHandler(logger *slog.Logger) *TmplHandler {
	return &TmplHandler{
		logger: logger,
	}
}

func (th TmplHandler) RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl *template.Template, data any) {
	if err := tmpl.Execute(w, data); err != nil {
		th.logger.Error("failed to execute template", slog.Any("template", tmpl), slog.Any("data", data))
		RenderError(w, r, "Could not load page", http.StatusInternalServerError)
		return
	}
}

func RenderError(w http.ResponseWriter, r *http.Request, errorMsg string, errorCode int) {
	http.Error(w, errorMsg, errorCode)
}
