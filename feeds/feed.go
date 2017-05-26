package feeds

import (
	"encoding/xml"
	"io"
	"time"
)

// Link Struct
type Link struct {
	Href, Rel, Type, Length string
}

// Author Struct
type Author struct {
	Name, Email string
}

// Torrent Struct modified for Nyaa
type Torrent struct {
	FileName      string
	Seeds         uint32
	Peers         uint32
	InfoHash      string
	ContentLength int64
	MagnetURI     string
}

// Item Struct
type Item struct {
	Title       string
	Link        *Link
	Author      *Author
	Description string // used as description in rss, summary in atom
	ID          string // used as guid in rss, id in atom
	Updated     time.Time
	Created     time.Time

	Torrent *Torrent // modified for Nyaa
}

// Feed Struct
type Feed struct {
	Title       string
	Link        *Link
	Description string
	Author      *Author
	Updated     time.Time
	Created     time.Time
	ID          string
	Subtitle    string
	Items       []*Item
	Copyright   string
}

// Add a new Item to a Feed
func (f *Feed) Add(item *Item) {
	f.Items = append(f.Items, item)
}

// returns the first non-zero time formatted as a string or ""
func anyTimeFormat(format string, times ...time.Time) string {
	for _, t := range times {
		if !t.IsZero() {
			return t.Format(format)
		}
	}
	return ""
}

// XMLFeed : interface used by ToXML to get a object suitable for exporting XML.
type XMLFeed interface {
	FeedXML() interface{}
}

// ToXML : turn a feed object (either a Feed, AtomFeed, or RssFeed) into xml
// returns an error if xml marshaling fails
func ToXML(feed XMLFeed) (string, error) {
	x := feed.FeedXML()
	data, err := xml.MarshalIndent(x, "", "  ")
	if err != nil {
		return "", err
	}
	// strip empty line from default xml header
	s := xml.Header[:len(xml.Header)-1] + string(data)
	return s, nil
}

// WriteXML : Write a feed object (either a Feed, AtomFeed, or RssFeed) as XML into
// the writer. Returns an error if XML marshaling fails.
func WriteXML(feed XMLFeed, w io.Writer) error {
	x := feed.FeedXML()
	// write default xml header, without the newline
	if _, err := w.Write([]byte(xml.Header[:len(xml.Header)-1])); err != nil {
		return err
	}
	e := xml.NewEncoder(w)
	e.Indent("", "  ")
	return e.Encode(x)
}

// ToAtom : creates an Atom representation of this feed
func (f *Feed) ToAtom() (string, error) {
	a := &Atom{f}
	return ToXML(a)
}

// WriteAtom : Writes an Atom representation of this feed to the writer.
func (f *Feed) WriteAtom(w io.Writer) error {
	return WriteXML(&Atom{f}, w)
}

// ToRss : creates an Rss representation of this feed
func (f *Feed) ToRss() (string, error) {
	r := &Rss{f}
	return ToXML(r)
}

// WriteRss : Writes an RSS representation of this feed to the writer.
func (f *Feed) WriteRss(w io.Writer) error {
	return WriteXML(&Rss{f}, w)
}
