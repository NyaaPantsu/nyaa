package sanitize

import (
	"html/template"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownToHTML(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Test   string
		Result template.HTML
	}{
		{"", ""},
		{"> lll", "<blockquote>\n<p>lll</p>\n</blockquote>\n"},
		{"> lll > lol", "<blockquote>\n<p>lll &gt; lol</p>\n</blockquote>\n"}, // Limit number of blockquotes
		{"&gt; lll", "<blockquote>\n<p>lll</p>\n</blockquote>\n"},
		{"\n", ""},
		{"<b>lol</b>", "<p><b>lol</b></p>\n"},                      // keep HTML tags
		{"[b]lol[/b]", "<p>[b]lol[/b]</p>\n"},                      // keep BBCode tags
		{"**[b]lol[/b]**", "<p><strong>[b]lol[/b]</strong></p>\n"}, // Render Markdown
	}
	for _, test := range tests {
		assert.Equal(test.Result, MarkdownToHTML(test.Test), "Should be equal")
	}
}

func TestParseBBCodes(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Test   string
		Result string
	}{
		{"", ""},
		{"&gt;", "&gt;"},                       // keep escaped html
		{"<b>lol</b>", "<b>lol</b>"},           // keep html tags
		{"[b]lol[/b]", "<b>lol</b>"},           // Convert bbcodes
		{"[u][b]lol[/u]", "<u><b>lol</b></u>"}, // Close unclosed tags
	}
	for _, test := range tests {
		assert.Equal(test.Result, ParseBBCodes(test.Test), "Should be equal")
	}
	assert.Contains(ParseBBCodes("[url=http://kk.cc/]lol[/url]"), "rel=\"nofollow\"") // rel="nofollow" for urls
}

func TestRepairHTMLTags(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Test   string
		Result string
	}{
		{"", ""},
		{"&gt;", "&gt;"},                                              // keep escaped html
		{"<b>lol</b>", "<b>lol</b>"},                                  // keep html tags
		{"<b><u>lol</b>", "<b><u>lol</u></b>"},                        // close unclosed tags encapsulated
		{"<b><u>lol", "<b><u>lol</u></b>"},                            // close unclosed tags non encapsulated
		{"<b><u>lol</em>", "<b><u>lol</u></b>"},                       // close unclosed tags non encaptsulated + remove useless end tags
		{"<div><b><u>lol</em></div>", "<div><b><u>lol</u></b></div>"}, // close unclosed tags + remove useless end tags encaptsulated
	}
	for _, test := range tests {
		assert.Equal(test.Result, repairHTMLTags(test.Test), "Should be equal")
	}
}

func TestSanitize(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Test   string
		Result string
	}{
		{"", ""},
		{"[b]lol[/b]", "<b>lol</b>"},                                                                                                                                                                                                        // Should convert bbcodes
		{"&gt;", "&gt;"},                                                                                                                                                                                                                    // keep escaped html
		{"<b>lol</b>", "<b>lol</b>"},                                                                                                                                                                                                        // keep html tags
		{"<b><u>lol</b>", "<b><u>lol</u></b>"},                                                                                                                                                                                              // close unclosed tags encapsulated
		{"<b><u>lol", "<b><u>lol</u></b>"},                                                                                                                                                                                                  // close unclosed tags non encapsulated
		{"<b><u>lol</em>", "<b><u>lol</u></b>"},                                                                                                                                                                                             // close unclosed tags non encaptsulated + remove useless end tags
		{"<div><b><u>lol</em></div>", "<b><u>lol</u></b>"},                                                                                                                                                                                  // close unclosed tags + remove useless end tags encaptsulated and remove div tag
		{"Hello <STYLE>.XSS{background-image:url(\"javascript:alert('XSS')\");}</STYLE><A CLASS=XSS></A>World", "Hello World"},                                                                                                              // Remove css XSS
		{"<a href=\"javascript:alert('XSS1')\" onmouseover=\"alert('XSS2')\">XSS<a>", "XSS"},                                                                                                                                                // Remove javascript xss
		{"<a href=\"http://www.google.com/\"><img src=\"https://ssl.gstatic.com/accounts/ui/logo_2x.png\"/></a>", "<a href=\"http://www.google.com/\" rel=\"nofollow\"><img src=\"https://ssl.gstatic.com/accounts/ui/logo_2x.png\"/></a>"}, // We allow img and linl
		{"<img src=\"data:image/webp;base64,UklGRh4AAABXRUJQVlA4TBEAAAAvAAAAAAfQ//73v/+BiOh/AAA=\">", ""},                                                                                                                                   // But not allow datauri img by default
		{"<objet></object><embed></embed><base><iframe />", ""},                                                                                                                                                                             // Not allowed elements by default
	}
	for _, test := range tests {
		assert.Equal(test.Result, Sanitize(test.Test, "default"), "Should be equal")
	}
}
