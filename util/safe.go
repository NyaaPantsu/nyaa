package util

import "html/template"

func Safe(s string) template.URL { // TODO: Inline function, or expand since it's unsafe (like, really?)
	return template.URL(s)
}
