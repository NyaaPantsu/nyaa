// Copyright 2015 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"fmt"
	"strconv"
	"strings"
)

type TagCompilerFunc func(*BBCodeNode) (*HTMLTag, bool)

type Compiler struct {
	tagCompilers               map[string]TagCompilerFunc
	defaultCompiler            TagCompilerFunc
	AutoCloseTags              bool
	IgnoreUnmatchedClosingTags bool
}

func NewCompiler(autoCloseTags, ignoreUnmatchedClosingTags bool) Compiler {
	compiler := Compiler{make(map[string]TagCompilerFunc), DefaultTagCompiler, autoCloseTags, ignoreUnmatchedClosingTags}

	for tag, compilerFunc := range DefaultTagCompilers {
		compiler.SetTag(tag, compilerFunc)
	}
	return compiler
}

func (c Compiler) Compile(str string) string {
	tokens := Lex(str)
	tree := Parse(tokens)
	return c.CompileTree(tree).String()
}

func (c Compiler) SetDefault(compiler TagCompilerFunc) {
	if compiler == nil {
		panic("Default tag compiler can't be nil")
	} else {
		c.defaultCompiler = compiler
	}
}

func (c Compiler) SetTag(tag string, compiler TagCompilerFunc) {
	if compiler == nil {
		delete(c.tagCompilers, tag)
	} else {
		c.tagCompilers[tag] = compiler
	}
}

// CompileTree transforms BBCodeNode into an HTML tag.
func (c Compiler) CompileTree(node *BBCodeNode) *HTMLTag {
	var out = NewHTMLTag("")
	if node.ID == TEXT {
		out.Value = node.Value.(string)
		InsertNewlines(out)
		for _, child := range node.Children {
			out.AppendChild(c.CompileTree(child))
		}
	} else if node.ID == CLOSING_TAG {
		if !c.IgnoreUnmatchedClosingTags {
			out.Value = node.Value.(BBClosingTag).Raw
			InsertNewlines(out)
		}
		for _, child := range node.Children {
			out.AppendChild(c.CompileTree(child))
		}
	} else if node.ClosingTag == nil && !c.AutoCloseTags {
		out.Value = node.Value.(BBOpeningTag).Raw
		InsertNewlines(out)
		for _, child := range node.Children {
			out.AppendChild(c.CompileTree(child))
		}
	} else {
		in := node.GetOpeningTag()

		compileFunc, ok := c.tagCompilers[in.Name]
		if !ok {
			compileFunc = c.defaultCompiler
		}
		var appendExpr bool
		node.Compiler = &c
		out, appendExpr = compileFunc(node)
		if appendExpr {
			if len(node.Children) == 0 {
				out.AppendChild(NewHTMLTag(""))
			} else {
				for _, child := range node.Children {
					out.AppendChild(c.CompileTree(child))
				}
			}
		}
	}
	return out
}

func CompileText(in *BBCodeNode) string {
	out := ""
	if in.ID == TEXT {
		out = in.Value.(string)
	}
	for _, child := range in.Children {
		out += CompileText(child)
	}
	return out
}

func CompileRaw(in *BBCodeNode) *HTMLTag {
	out := NewHTMLTag("")
	if in.ID == TEXT {
		out.Value = in.Value.(string)
	} else if in.ID == CLOSING_TAG {
		out.Value = in.Value.(BBClosingTag).Raw
	} else {
		out.Value = in.Value.(BBOpeningTag).Raw
	}
	for _, child := range in.Children {
		out.AppendChild(CompileRaw(child))
	}
	if in.ID == OPENING_TAG && in.ClosingTag != nil {
		tag := NewHTMLTag(in.ClosingTag.Raw)
		out.AppendChild(tag)
	}
	return out
}

var DefaultTagCompilers map[string]TagCompilerFunc
var DefaultTagCompiler TagCompilerFunc

func init() {
	DefaultTagCompiler = func(node *BBCodeNode) (*HTMLTag, bool) {
		out := NewHTMLTag(node.GetOpeningTag().Raw)
		InsertNewlines(out)
		if len(node.Children) == 0 {
			out.AppendChild(NewHTMLTag(""))
		} else {
			for _, child := range node.Children {
				out.AppendChild(node.Compiler.CompileTree(child))
			}
		}
		if node.ClosingTag != nil {
			tag := NewHTMLTag(node.ClosingTag.Raw)
			InsertNewlines(tag)
			out.AppendChild(tag)
		}
		return out, false
	}

	DefaultTagCompilers = make(map[string]TagCompilerFunc)
	DefaultTagCompilers["url"] = func(node *BBCodeNode) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "a"
		value := node.GetOpeningTag().Value
		if value == "" {
			text := CompileText(node)
			if len(text) > 0 {
				out.Attrs["href"] = ValidURL(text)
			}
		} else {
			out.Attrs["href"] = ValidURL(value)
		}
		return out, true
	}

	DefaultTagCompilers["img"] = func(node *BBCodeNode) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "img"
		value := node.GetOpeningTag().Value
		if value == "" {
			out.Attrs["src"] = ValidURL(CompileText(node))
		} else {
			out.Attrs["src"] = ValidURL(value)
			text := CompileText(node)
			if len(text) > 0 {
				out.Attrs["alt"] = text
				out.Attrs["title"] = out.Attrs["alt"]
			}
		}
		return out, false
	}

	DefaultTagCompilers["center"] = func(node *BBCodeNode) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "div"
		out.Attrs["style"] = "text-align: center;"
		return out, true
	}

	DefaultTagCompilers["color"] = func(node *BBCodeNode) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "span"
		sanitize := func(r rune) rune {
			if r == '#' || r == ',' || r == '.' || r == '(' || r == ')' || r == '%' {
				return r
			} else if r >= '0' && r <= '9' {
				return r
			} else if r >= 'a' && r <= 'z' {
				return r
			} else if r >= 'A' && r <= 'Z' {
				return r
			}
			return -1
		}
		color := strings.Map(sanitize, node.GetOpeningTag().Value)
		out.Attrs["style"] = "color: " + color + ";"
		return out, true
	}

	DefaultTagCompilers["size"] = func(node *BBCodeNode) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "span"
		if size, err := strconv.Atoi(node.GetOpeningTag().Value); err == nil {
			out.Attrs["class"] = fmt.Sprintf("size%d", size)
		}
		return out, true
	}

	DefaultTagCompilers["quote"] = func(node *BBCodeNode) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "blockquote"
		who := ""
		in := node.GetOpeningTag()
		if name, ok := in.Args["name"]; ok && name != "" {
			who = name
		} else {
			who = in.Value
		}
		cite := NewHTMLTag("")
		cite.Name = "cite"
		if who != "" {
			cite.AppendChild(NewHTMLTag(who + " said:"))
		} else {
			cite.AppendChild(NewHTMLTag("Quote"))
		}
		return out.AppendChild(cite), true
	}

	DefaultTagCompilers["code"] = func(node *BBCodeNode) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "pre"
		for _, child := range node.Children {
			out.AppendChild(CompileRaw(child))
		}
		return out, false
	}

	for _, tag := range []string{"i", "b", "u", "s"} {
		DefaultTagCompilers[tag] = func(node *BBCodeNode) (*HTMLTag, bool) {
			out := NewHTMLTag("")
			out.Name = node.GetOpeningTag().Name
			return out, true
		}
	}
}
