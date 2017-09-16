package sanitize

import (
	"bytes"
	"encoding/xml"
	"html/template"
	"io"
	"regexp"
	"strings"

	"github.com/NyaaPantsu/nyaa/utils/log"

	"github.com/frustra/bbcode"
	"github.com/microcosm-cc/bluemonday"
	md "github.com/russross/blackfriday"
	"golang.org/x/net/html"
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

type htmlTag struct {
	XMLName xml.Name `xml:"html"`
	Body    body     `xml:"body"`
}

type body struct {
	Content string `xml:",innerxml"`
}

func init() {
	HTMLMdRenderer = md.HtmlRenderer(htmlFlags, "", "")
	BBCodesRenderer = bbcode.NewCompiler(true, true) // autoCloseTags, ignoreUnmatchedClosingTags
	BBCodesRenderer.SetTag("url", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
		out, appendExpr := bbcode.DefaultTagCompilers["url"](node)
		out.Attrs["rel"] = "nofollow"
		return out, appendExpr
	})
}

// HTMLMdRenderer render for markdown to html
var HTMLMdRenderer md.Renderer

// BBCodesRenderer render bbcodes to html
var BBCodesRenderer bbcode.Compiler

// MarkdownToHTML : convert markdown to html
// TODO: restrict certain types of markdown
func MarkdownToHTML(markdown string) template.HTML {
	if len(markdown) >= 4 && markdown[:4] == "&gt;" {
		markdown = ">" + markdown[4:]
	}
	markdown = strings.Replace(markdown, "\n&gt;", "\n>", -1)
	unsafe := md.MarkdownOptions([]byte(markdown), HTMLMdRenderer, md.Options{Extensions: mdOptions})
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	return template.HTML(html)
}

// Sanitize :
/* Sanitize a message passed as a string according to a setted model or allowing a set of html tags and output a string
 */
func Sanitize(msg string, elements ...string) string {
	// Convert BBCodes to HTML
	//msg = ParseBBCodes(msg)
	// Repair HTML
	msg = repairHTMLTags(msg) // We repair possible broken html tags
	// HTML Sanitize
	p := bluemonday.NewPolicy()
	if len(elements) > 0 {
		if elements[0] == "default" { // default model same as UGC without div
			///////////////////////
			// Global attributes //
			///////////////////////

			// "class" is not permitted as we are not allowing users to style their own
			// content

			p.AllowStandardAttributes()

			//////////////////////////////
			// Global URL format policy //
			//////////////////////////////

			p.AllowStandardURLs()

			////////////////////////////////
			// Declarations and structure //
			////////////////////////////////

			// "xml" "xslt" "DOCTYPE" "html" "head" are not permitted as we are
			// expecting user generated content to be a fragment of HTML and not a full
			// document.

			//////////////////////////
			// Sectioning root tags //
			//////////////////////////

			// "article" and "aside" are permitted and takes no attributes
			p.AllowElements("article", "aside")

			// "body" is not permitted as we are expecting user generated content to be a fragment
			// of HTML and not a full document.

			// "details" is permitted, including the "open" attribute which can either
			// be blank or the value "open".
			p.AllowAttrs(
				"open",
			).Matching(regexp.MustCompile(`(?i)^(|open)$`)).OnElements("details")

			// "fieldset" is not permitted as we are not allowing forms to be created.

			// "figure" is permitted and takes no attributes
			p.AllowElements("figure")

			// "nav" is not permitted as it is assumed that the site (and not the user)
			// has defined navigation elements

			// "section" is permitted and takes no attributes
			p.AllowElements("section")

			// "summary" is permitted and takes no attributes
			p.AllowElements("summary")

			//////////////////////////
			// Headings and footers //
			//////////////////////////

			// "footer" is not permitted as we expect user content to be a fragment and
			// not structural to this extent

			// "h1" through "h6" are permitted and take no attributes
			p.AllowElements("h1", "h2", "h3", "h4", "h5", "h6")

			// "header" is not permitted as we expect user content to be a fragment and
			// not structural to this extent

			// "hgroup" is permitted and takes no attributes
			p.AllowElements("hgroup")

			/////////////////////////////////////
			// Content grouping and separating //
			/////////////////////////////////////

			// "blockquote" is permitted, including the "cite" attribute which must be
			// a standard URL.
			p.AllowAttrs("cite").OnElements("blockquote")

			// "br" "div" "hr" "p" "span" "wbr" are permitted and take no attributes
			p.AllowElements("br", "hr", "p", "span", "wbr")

			///////////
			// Links //
			///////////

			// "a" is permitted
			p.AllowAttrs("href").OnElements("a")

			// "area" is permitted along with the attributes that map image maps work
			p.AllowAttrs("name").Matching(
				regexp.MustCompile(`^([\p{L}\p{N}_-]+)$`),
			).OnElements("map")
			p.AllowAttrs("alt").Matching(bluemonday.Paragraph).OnElements("area")
			p.AllowAttrs("coords").Matching(
				regexp.MustCompile(`^([0-9]+,)+[0-9]+$`),
			).OnElements("area")
			p.AllowAttrs("href").OnElements("area")
			p.AllowAttrs("rel").Matching(bluemonday.SpaceSeparatedTokens).OnElements("area")
			p.AllowAttrs("shape").Matching(
				regexp.MustCompile(`(?i)^(default|circle|rect|poly)$`),
			).OnElements("area")
			p.AllowAttrs("usemap").Matching(
				regexp.MustCompile(`(?i)^#[\p{L}\p{N}_-]+$`),
			).OnElements("img")

			// "link" is not permitted

			/////////////////////
			// Phrase elements //
			/////////////////////

			// The following are all inline phrasing elements
			p.AllowElements("abbr", "acronym", "cite", "code", "dfn", "em",
				"figcaption", "mark", "s", "samp", "strong", "sub", "sup", "var")

			// "q" is permitted and "cite" is a URL and handled by URL policies
			p.AllowAttrs("cite").OnElements("q")

			// "time" is permitted
			p.AllowAttrs("datetime").Matching(bluemonday.ISO8601).OnElements("time")

			////////////////////
			// Style elements //
			////////////////////

			// block and inline elements that impart no semantic meaning but style the
			// document
			p.AllowElements("b", "i", "pre", "small", "strike", "tt", "u")

			// "style" is not permitted as we are not yet sanitising CSS and it is an
			// XSS attack vector

			//////////////////////
			// HTML5 Formatting //
			//////////////////////

			// "bdi" "bdo" are permitted
			p.AllowAttrs("dir").Matching(bluemonday.Direction).OnElements("bdi", "bdo")

			// "rp" "rt" "ruby" are permitted
			p.AllowElements("rp", "rt", "ruby")

			///////////////////////////
			// HTML5 Change tracking //
			///////////////////////////

			// "del" "ins" are permitted
			p.AllowAttrs("cite").Matching(bluemonday.Paragraph).OnElements("del", "ins")
			p.AllowAttrs("datetime").Matching(bluemonday.ISO8601).OnElements("del", "ins")

			///////////
			// Lists //
			///////////

			p.AllowLists()

			////////////
			// Tables //
			////////////

			p.AllowTables()

			///////////
			// Forms //
			///////////

			// By and large, forms are not permitted. However there are some form
			// elements that can be used to present data, and we do permit those
			//
			// "button" "fieldset" "input" "keygen" "label" "output" "select" "datalist"
			// "textarea" "optgroup" "option" are all not permitted

			// "meter" is permitted
			p.AllowAttrs(
				"value",
				"min",
				"max",
				"low",
				"high",
				"optimum",
			).Matching(bluemonday.Number).OnElements("meter")

			// "progress" is permitted
			p.AllowAttrs("value", "max").Matching(bluemonday.Number).OnElements("progress")

			//////////////////////
			// Embedded content //
			//////////////////////

			// Vast majority not permitted
			// "audio" "canvas" "embed" "iframe" "object" "param" "source" "svg" "track"
			// "video" are all not permitted

			p.AllowImages()
		} else if elements[0] == "comment" {
			p.AllowElements("b", "strong", "em", "i", "u", "blockquote", "q")
			p.AllowImages()
			p.AllowStandardURLs()
			p.AllowAttrs("cite").OnElements("blockquote", "q")
			p.AllowAttrs("href").OnElements("a")
			p.AddTargetBlankToFullyQualifiedLinks(true)
		} else { // allowing set of html tags
			p.AllowElements(elements...)
		}
	}
	return p.Sanitize(msg)
}

// repairHTMLTags Should close any opened tags and strip any empty end tags
func repairHTMLTags(brokenHTML string) string {
	reader := strings.NewReader(brokenHTML)
	root, err := html.Parse(reader)
	if !log.CheckError(err) {
		return ""
	}
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, root)

	fixedHTML := htmlTag{}
	err = xml.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&fixedHTML)
	if !log.CheckError(err) {
		return ""
	}
	return fixedHTML.Body.Content
}

// ParseBBCodes returns the bbcode compiler with the bbcode tags to parse
func ParseBBCodes(msg string) string {
	msg = BBCodesRenderer.Compile(msg)
	// For some reason, BBCodes compiler return escaped html
	// We need to unescape it
	return html.UnescapeString(msg)
}
