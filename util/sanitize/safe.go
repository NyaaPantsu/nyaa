package sanitize

import (
	"html"
	"html/template"
)

// Safe : make un url safe
func Safe(s string) template.URL {
	return template.URL(html.EscapeString(s))
}

// SafeText : make a string safe
func SafeText(s string) template.HTML {
	return template.HTML(html.EscapeString(s))
}
