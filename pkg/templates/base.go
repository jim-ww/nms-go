package templates

import "html/template"

type Base struct {
	Title   string
	Content template.HTML
}

func NewBase(title string, template template.HTML) *Base {
	return &Base{
		Title:   title,
		Content: template,
	}
}
