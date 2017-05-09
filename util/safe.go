package util

import (
	"html"
	"html/template"
)

func Safe(s string) template.URL {
	return template.URL(s)
}

func SafeText(s string) template.HTML {
	return template.HTML(html.EscapeString(s))
}
