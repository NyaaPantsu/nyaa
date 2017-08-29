# bbcode [![Build Status](https://travis-ci.org/frustra/bbcode.png?branch=master)](http://travis-ci.org/frustra/bbcode)

frustra/bbcode is a fast BBCode compiler for Go. It supports custom tags, safe html output (for user-specified input),
and allows for graceful parsing of syntax errors similar to the output of a regex bbcode compiler.

Visit the godoc here: [http://godoc.org/github.com/frustra/bbcode](http://godoc.org/github.com/frustra/bbcode)

## Usage

To get started compiling some text, create a compiler instance:
```go
compiler := bbcode.NewCompiler(true, true) // autoCloseTags, ignoreUnmatchedClosingTags
fmt.Println(compiler.Compile("[b]Hello World[/b]"))

// Output:
// <b>Hello World</b>
```

## Supported BBCode Syntax
```
[tag]basic tag[/tag]
[tag1][tag2]nested tags[/tag2][/tag1]

[tag=value]tag with value[/tag]
[tag arg=value]tag with named argument[/tag]
[tag="quote value"]tag with quoted value[/tag]

[tag=value foo="hello world" bar=baz]multiple tag arguments[/tag]
```

## Default Tags
 * `[b]text[/b]` --> `<b>text</b>` (b, i, u, and s all map the same)
 * `[url]link[/url]` --> `<a href="link">link</a>`
 * `[url=link]text[/url]` --> `<a href="link">text</a>`
 * `[img]link[/img]` --> `<img src="link">`
 * `[img=link]alt[/img]` --> `<img alt="alt" title="alt" src="link">`
 * `[center]text[/center]` --> `<div style="text-align: center;">text</div>`
 * `[color=red]text[/color]` --> `<span style="color: red;">text</span>`
 * `[size=2]text[/size]` --> `<span class="size2">text</span>`
 * `[quote]text[/quote]` --> `<blockquote><cite>Quote</cite>text</blockquote>`
 * `[quote=Somebody]text[/quote]` --> `<blockquote><cite>Somebody said:</cite>text</blockquote>`
 * `[quote name=Somebody]text[/quote]` --> `<blockquote><cite>Somebody said:</cite>text</blockquote>`
 * `[code][b]anything[/b][/code]` --> `<pre>[b]anything[/b]</pre>`

Lists are not currently implemented as a default tag, but can be added as a custom tag.  
A working implementation of list tags can be found [here](https://gist.github.com/xthexder/44f4b9cec3ed7876780d)

## Adding Custom Tags
Custom tag handlers can be added to a compiler using the `compiler.SetTag(tag, handler)` function:
```go
compiler.SetTag("center", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
	// Create a new div element to output
	out := bbcode.NewHTMLTag("")
	out.Name = "div"

	// Set the style attribute of our output div
	out.Attrs["style"] = "text-align: center;"

	// Returning true here means continue to parse child nodes.
	// This should be false if children are parsed by this tag's handler, like in the [code] tag.
	return out, true
})
```

Tag values can be read from the opening tag like this:  
Main tag value `[tag={value}]`: `node.GetOpeningTag().Value`  
Tag arguments `[tag name={value}]`: `node.GetOpeningTag().Args["name"]`

`bbcode.NewHTMLTag(text)` creates a text node by default. By setting `tag.Name`, the node because an html tag prefixed by the text. The closing html tag is not rendered unless child elements exist. The closing tag can be forced by adding a blank text node:
```go
out := bbcode.NewHTMLTag("")
out.Name = "div"
out.AppendChild(nil) // equivalent to out.AppendChild(bbcode.NewHTMLTag(""))
```

For more examples of tag definitions, look at the default tag implementations in [compiler.go](https://github.com/frustra/bbcode/blob/master/compiler.go)

## Overriding Default Tags
The built-in tags can be overridden simply by redefining the tag with `compiler.SetTag(tag, handler)`

To remove a tag, set the tag handler to nil:
```go
compiler.SetTag("quote", nil)
```

The default tags can also be modified without completely redefining the tag by calling the default handler:
```go
compiler.SetTag("url", func(node *bbcode.BBCodeNode) (*bbcode.HTMLTag, bool) {
	out, appendExpr := bbcode.DefaultTagCompilers["url"](node)
	out.Attrs["class"] = "bbcode-link"
	return out, appendExpr
})
```

## Auto-Close Tags
Input:
```
[center][b]text[/center]
```

Enabled Output:
```html
<div style="text-align: center;"><b>text</b></div>
```
Disabled Output:
```html
<div style="text-align: center;">[b]text</div>
```

## Ignore Unmatched Closing Tags
Input:
```
[center]text[/b][/center]
```

Enabled Output:
```html
<div style="text-align: center;">text</div>
```
Disabled Output:
```html
<div style="text-align: center;">text[/b]</div>
```
