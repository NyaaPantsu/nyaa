package util

import 	"html/template"

func Safe(s string) template.URL {
	return template.URL(s)
}