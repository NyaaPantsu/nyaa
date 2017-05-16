package feeds

import (
	"bytes"
	"testing"
	"time"
)

var atomOutput = `<?xml version="1.0" encoding="UTF-8"?><feed xmlns="http://www.w3.org/2005/Atom">
  <title>jmoiron.net blog</title>
  <id>http://jmoiron.net/blog</id>
  <updated>2013-01-16T21:52:35-05:00</updated>
  <rights>This work is copyright © Benjamin Button</rights>
  <subtitle>discussion about tech, footie, photos</subtitle>
  <link href="http://jmoiron.net/blog"></link>
  <author>
    <name>Jason Moiron</name>
    <email>jmoiron@jmoiron.net</email>
  </author>
  <entry>
    <title>Limiting Concurrency in Go</title>
    <updated>2013-01-16T21:52:35-05:00</updated>
    <id>tag:jmoiron.net,2013-01-16:/blog/limiting-concurrency-in-go/</id>
    <content type="html">A discussion on controlled parallelism in golang</content>
    <link href="http://jmoiron.net/blog/limiting-concurrency-in-go/"></link>
    <author>
      <name>Jason Moiron</name>
      <email>jmoiron@jmoiron.net</email>
    </author>
  </entry>
  <entry>
    <title>Logic-less Template Redux</title>
    <updated>2013-01-16T21:52:35-05:00</updated>
    <id>tag:jmoiron.net,2013-01-16:/blog/logicless-template-redux/</id>
    <content type="html">More thoughts on logicless templates</content>
    <link href="http://jmoiron.net/blog/logicless-template-redux/"></link>
  </entry>
  <entry>
    <title>Idiomatic Code Reuse in Go</title>
    <updated>2013-01-16T21:52:35-05:00</updated>
    <id>tag:jmoiron.net,2013-01-16:/blog/idiomatic-code-reuse-in-go/</id>
    <content type="html">How to use interfaces &lt;em&gt;effectively&lt;/em&gt;</content>
    <link href="http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"></link>
  </entry>
  <entry>
    <title>Never Gonna Give You Up Mp3</title>
    <updated>2013-01-16T21:52:35-05:00</updated>
    <id>tag:example.com,2013-01-16:/RickRoll.mp3</id>
    <content type="html">Never gonna give you up - Never gonna let you down.</content>
    <link href="http://example.com/RickRoll.mp3" rel="enclosure" type="audio/mpeg" length="123456"></link>
  </entry>
  <entry>
    <title>String formatting in Go</title>
    <updated>2013-01-16T21:52:35-05:00</updated>
    <id>tag:example.com,2013-01-16:/strings</id>
    <content type="html">How to use things like %s, %v, %d, etc.</content>
    <link href="http://example.com/strings"></link>
  </entry>
</feed>`

var rssOutput = `<?xml version="1.0" encoding="UTF-8"?><rss version="2.0">
  <channel>
    <title>jmoiron.net blog</title>
    <link>http://jmoiron.net/blog</link>
    <description>discussion about tech, footie, photos</description>
    <copyright>This work is copyright © Benjamin Button</copyright>
    <managingEditor>jmoiron@jmoiron.net (Jason Moiron)</managingEditor>
    <pubDate>Wed, 16 Jan 2013 21:52:35 -0500</pubDate>
    <item>
      <title>Limiting Concurrency in Go</title>
      <link>http://jmoiron.net/blog/limiting-concurrency-in-go/</link>
      <description>A discussion on controlled parallelism in golang</description>
      <author>Jason Moiron</author>
      <pubDate>Wed, 16 Jan 2013 21:52:35 -0500</pubDate>
    </item>
    <item>
      <title>Logic-less Template Redux</title>
      <link>http://jmoiron.net/blog/logicless-template-redux/</link>
      <description>More thoughts on logicless templates</description>
      <pubDate>Wed, 16 Jan 2013 21:52:35 -0500</pubDate>
    </item>
    <item>
      <title>Idiomatic Code Reuse in Go</title>
      <link>http://jmoiron.net/blog/idiomatic-code-reuse-in-go/</link>
      <description>How to use interfaces &lt;em&gt;effectively&lt;/em&gt;</description>
      <pubDate>Wed, 16 Jan 2013 21:52:35 -0500</pubDate>
    </item>
    <item>
      <title>Never Gonna Give You Up Mp3</title>
      <link>http://example.com/RickRoll.mp3</link>
      <description>Never gonna give you up - Never gonna let you down.</description>
      <enclosure url="http://example.com/RickRoll.mp3" length="123456" type="audio/mpeg"></enclosure>
      <pubDate>Wed, 16 Jan 2013 21:52:35 -0500</pubDate>
    </item>
    <item>
      <title>String formatting in Go</title>
      <link>http://example.com/strings</link>
      <description>How to use things like %s, %v, %d, etc.</description>
      <pubDate>Wed, 16 Jan 2013 21:52:35 -0500</pubDate>
    </item>
  </channel>
</rss>`

func TestFeed(t *testing.T) {
	now, err := time.Parse(time.RFC3339, "2013-01-16T21:52:35-05:00")
	if err != nil {
		t.Error(err)
	}
	tz := time.FixedZone("EST", -5*60*60)
	now = now.In(tz)

	feed := &Feed{
		Title:       "jmoiron.net blog",
		Link:        &Link{Href: "http://jmoiron.net/blog"},
		Description: "discussion about tech, footie, photos",
		Author:      &Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
		Created:     now,
		Copyright:   "This work is copyright © Benjamin Button",
	}

	feed.Items = []*Item{
		{
			Title:       "Limiting Concurrency in Go",
			Link:        &Link{Href: "http://jmoiron.net/blog/limiting-concurrency-in-go/"},
			Description: "A discussion on controlled parallelism in golang",
			Author:      &Author{Name: "Jason Moiron", Email: "jmoiron@jmoiron.net"},
			Created:     now,
		},
		{
			Title:       "Logic-less Template Redux",
			Link:        &Link{Href: "http://jmoiron.net/blog/logicless-template-redux/"},
			Description: "More thoughts on logicless templates",
			Created:     now,
		},
		{
			Title:       "Idiomatic Code Reuse in Go",
			Link:        &Link{Href: "http://jmoiron.net/blog/idiomatic-code-reuse-in-go/"},
			Description: "How to use interfaces <em>effectively</em>",
			Created:     now,
		},
		{
			Title:       "Never Gonna Give You Up Mp3",
			Link:        &Link{Href: "http://example.com/RickRoll.mp3", Length: "123456", Type: "audio/mpeg"},
			Description: "Never gonna give you up - Never gonna let you down.",
			Created:     now,
		},
		{
			Title:       "String formatting in Go",
			Link:        &Link{Href: "http://example.com/strings"},
			Description: "How to use things like %s, %v, %d, etc.",
			Created:     now,
		}}

	atom, err := feed.ToAtom()
	if err != nil {
		t.Errorf("unexpected error encoding Atom: %v", err)
	}
	if atom != atomOutput {
		t.Errorf("Atom not what was expected.  Got:\n%s\n\nExpected:\n%s\n", atom, atomOutput)
	}
	var buf bytes.Buffer
	if err := feed.WriteAtom(&buf); err != nil {
		t.Errorf("unexpected error writing Atom: %v", err)
	}
	if got := buf.String(); got != atomOutput {
		t.Errorf("Atom not what was expected.  Got:\n%s\n\nExpected:\n%s\n", got, atomOutput)
	}

	rss, err := feed.ToRss()
	if err != nil {
		t.Errorf("unexpected error encoding RSS: %v", err)
	}
	if rss != rssOutput {
		t.Errorf("Rss not what was expected.  Got:\n%s\n\nExpected:\n%s\n", rss, rssOutput)
	}
	buf.Reset()
	if err := feed.WriteRss(&buf); err != nil {
		t.Errorf("unexpected error writing RSS: %v", err)
	}
	if got := buf.String(); got != rssOutput {
		t.Errorf("Rss not what was expected.  Got:\n%s\n\nExpected:\n%s\n", got, rssOutput)
	}
}
