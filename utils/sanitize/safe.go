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

// ClearEmpty removes empty string entries from a map
func ClearEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	if len(r) == 0 {
		r = append(r, "")
	}
	return r
}
