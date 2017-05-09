package util

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"

	"html/template"
)

// TODO restrict certain types of markdown
func MarkdownToHTML(markdown string) template.HTML {
	unsafe := blackfriday.MarkdownCommon([]byte(markdown))
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return template.HTML(html)
}
