package middleware

import (
	"bytes"
	"html/template"
	"log/slog"
	"net/http"

	"github.com/jim-ww/nms-go/pkg/templates"
	"github.com/jim-ww/nms-go/pkg/utils/handlers"
)

type HTMXMiddlewareBaseWrapper struct {
	logger *slog.Logger
	tmpl   *template.Template
}

func NewHTMXMiddlewareBaseWrapper(logger *slog.Logger) *HTMXMiddlewareBaseWrapper {
	logger.Debug("Parsing base.html template...")

	templatePath := "web/templates/base.html"

	templ, err := template.ParseFiles(templatePath)
	if err != nil {
		logger.Error("Failed to parse base.html template", slog.String("template-path", templatePath))
		panic(err)
	}

	return &HTMXMiddlewareBaseWrapper{
		logger: logger,
		tmpl:   templ,
	}
}

func (wr HTMXMiddlewareBaseWrapper) WrapHTMXWithBaseTemplate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("HX-Request") == "true" {
			wr.logger.Debug("request is HX-request, serving directly", slog.Any("request-uri", r.RequestURI))
			next.ServeHTTP(w, r)
		} else {
			wr.logger.Debug("request is not HX-request, serving with base template", slog.Any("request-uri", r.RequestURI))

			wr.logger.Debug("Serving data to CustomResponseWriter...", slog.Any("request-uri", r.RequestURI))
			var content bytes.Buffer
			rw := NewCustomResponseWriter(&content, wr.logger)

			next.ServeHTTP(rw, r)
			wr.logger.Debug("Served data to CustomResponseWriter", slog.Any("request-uri", r.RequestURI))

			data := templates.NewBase("NMS", template.HTML(content.String()))
			if err := wr.tmpl.Execute(w, data); err != nil {
				wr.logger.Error("Failed to execute template", slog.String("template-name", wr.tmpl.Name()), slog.Any("data", data))
				handlers.RenderError(w, r, "Could not load page", http.StatusInternalServerError)
			}
		}
	})
}

type CustomResponseWriter struct {
	*DiscardResponseWriter
	buffer *bytes.Buffer
}

func NewCustomResponseWriter(buffer *bytes.Buffer, logger *slog.Logger) *CustomResponseWriter {
	drw := NewDiscardResponseWriter(logger)

	return &CustomResponseWriter{
		DiscardResponseWriter: drw,
		buffer:                buffer,
	}
}

type DiscardResponseWriter struct {
	logger *slog.Logger
}

func (d *DiscardResponseWriter) Header() http.Header {
	d.logger.Debug("discardResponseWriter Header() called")
	return map[string][]string{}
}

func (d *DiscardResponseWriter) Write([]byte) (int, error) {
	d.logger.Debug("discardResponseWriter Write() called")
	return 0, nil
}

func (d *DiscardResponseWriter) WriteHeader(int) {
	d.logger.Debug("discardResponseWriter WriteHeader() called")
}

func NewDiscardResponseWriter(logger *slog.Logger) *DiscardResponseWriter {
	return &DiscardResponseWriter{
		logger: logger,
	}
}

func (crw *CustomResponseWriter) Write(b []byte) (int, error) {
	crw.buffer.Write(b)
	return len(b), nil
}
