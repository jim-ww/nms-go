package middleware

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/jim-ww/nms-go/internal/templates"
)

var baseTmpl = template.Must(template.ParseFiles("web/templates/base.html"))

type CustomResponseWriter struct {
	http.ResponseWriter
	buffer *bytes.Buffer
}

func (crw *CustomResponseWriter) Write(b []byte) (int, error) {
	crw.buffer.Write(b)
	return len(b), nil
}

func WrapHTMXWithBaseTemplate(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header.Get("HX-Request") == "true" {
			next.ServeHTTP(w, r)
		} else {

			var content bytes.Buffer
			rw := CustomResponseWriter{buffer: &content}

			next.ServeHTTP(&rw, r)

			data := templates.NewBase("NMS", template.HTML(content.String()))
			if err := baseTmpl.Execute(w, data); err != nil {
				http.Error(w, "Could not load page", http.StatusInternalServerError)
			}
		}
	})
}
