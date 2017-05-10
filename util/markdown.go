package util

import (
	"github.com/microcosm-cc/bluemonday"
	md "github.com/russross/blackfriday"

	"html/template"
	"strings"
)
//Some default rules, plus and minus some.
var mdOptions = 0 |
	md.EXTENSION_AUTOLINK |
	md.EXTENSION_HARD_LINE_BREAK |
	md.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK |
	md.EXTENSION_NO_INTRA_EMPHASIS |
	md.EXTENSION_SPACE_HEADERS |
	md.EXTENSION_STRIKETHROUGH

var htmlFlags = 0 |
	md.HTML_USE_XHTML |
	md.HTML_SMARTYPANTS_FRACTIONS |
	md.HTML_SAFELINK |
	md.HTML_NOREFERRER_LINKS |
	md.HTML_HREF_TARGET_BLANK

func init() {
	HtmlMdRenderer = md.HtmlRenderer(htmlFlags, "", "")
}
var HtmlMdRenderer md.Renderer

// TODO: restrict certain types of markdown
func MarkdownToHTML(markdown string) template.HTML {
	if len(markdown) >= 3) && markdown[:3] == "&gt;" {
		markdown = ">" + markdown[3:]
	}
	markdown = strings.Replace(markdown,"\n&gt;","\n>")
	unsafe := md.MarkdownOptions([]byte(markdown), HtmlMdRenderer, md.Options{Extensions: mdOptions})
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return template.HTML(html)
}
