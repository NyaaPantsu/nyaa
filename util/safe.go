package util

import (
	"html"
	"html/template"
)

func Safe(s string) template.URL {
	return template.URL(html.EscapeString(s))
}

func SafeText(s string) template.HTML {
	return template.HTML(html.EscapeString(s))
}
