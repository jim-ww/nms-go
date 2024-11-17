package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplate(w http.ResponseWriter, tmpl *template.Template, data any) {
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Println(err)
		http.Error(w, "Could not load page", http.StatusInternalServerError)
	}
}
