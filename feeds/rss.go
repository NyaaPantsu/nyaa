package feeds

// rss support
// validation done according to spec here:
//    http://cyber.law.harvard.edu/rss/rss.html

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

// private wrapper around the RssFeed which gives us the <rss>..</rss> xml
type rssFeedXML struct {
	XMLName      xml.Name `xml:"rss"`
	Version      string   `xml:"version,attr"`
	Channel      *RssFeed
	XMLNSTorrent string `xml:"xmlns:torrent,attr"` // modified for Nyaa
}

// RssImage Struct
type RssImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int      `xml:"width,omitempty"`
	Height  int      `xml:"height,omitempty"`
}

// RssTextInput Struct
type RssTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

// RssFeed Struct
type RssFeed struct {
	XMLName        xml.Name `xml:"channel"`
	Title          string   `xml:"title"`       // required
	Link           string   `xml:"link"`        // required
	Description    string   `xml:"description"` // required
	Language       string   `xml:"language,omitempty"`
	Copyright      string   `xml:"copyright,omitempty"`
	ManagingEditor string   `xml:"managingEditor,omitempty"` // Author used
	WebMaster      string   `xml:"webMaster,omitempty"`
	PubDate        string   `xml:"pubDate,omitempty"`       // created or updated
	LastBuildDate  string   `xml:"lastBuildDate,omitempty"` // updated used
	Category       string   `xml:"category,omitempty"`
	Generator      string   `xml:"generator,omitempty"`
	Docs           string   `xml:"docs,omitempty"`
	Cloud          string   `xml:"cloud,omitempty"`
	TTL            int      `xml:"ttl,omitempty"`
	Rating         string   `xml:"rating,omitempty"`
	SkipHours      string   `xml:"skipHours,omitempty"`
	SkipDays       string   `xml:"skipDays,omitempty"`
	Image          *RssImage
	TextInput      *RssTextInput
	Items          []*RssItem
}

// RssItem Struct
type RssItem struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`       // required
	Link        string   `xml:"link"`        // required
	Description string   `xml:"description"` // required
	Author      string   `xml:"author,omitempty"`
	Category    string   `xml:"category,omitempty"`
	Comments    string   `xml:"comments,omitempty"`
	Enclosure   *RssEnclosure
	GUID        string `xml:"guid,omitempty"`    // Id used
	PubDate     string `xml:"pubDate,omitempty"` // created or updated
	Source      string `xml:"source,omitempty"`

	// modified for Nyaa
	FileName      string `xml:"torrent:fileName"`
	Seeds         uint32 `xml:"torrent:seeds"`
	Peers         uint32 `xml:"torrent:peers"`
	InfoHash      string `xml:"torrent:infoHash"`
	ContentLength int64  `xml:"torrent:contentLength"`
	MagnetURI     string `xml:"torrent:magnetURI"`
}

// RssEnclosure Struct
type RssEnclosure struct {
	//RSS 2.0 <enclosure url="http://example.com/file.mp3" length="123456789" type="audio/mpeg" />
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

// Rss Struct
type Rss struct {
	*Feed
}

// create a new RssItem with a generic Item struct's data
func newRssItem(i *Item) *RssItem {
	item := &RssItem{
		Title:       i.Title,
		Link:        i.Link.Href,
		Description: i.Description,
		GUID:        i.ID,
		PubDate:     anyTimeFormat(time.RFC1123Z, i.Created, i.Updated),
		// modified for Nyaa
		FileName:      i.Torrent.FileName,
		Seeds:         i.Torrent.Seeds,
		Peers:         i.Torrent.Peers,
		InfoHash:      i.Torrent.InfoHash,
		ContentLength: i.Torrent.ContentLength,
		MagnetURI:     i.Torrent.MagnetURI,
	}

	intLength, err := strconv.ParseInt(i.Link.Length, 10, 64)

	if err == nil && (intLength > 0 || i.Link.Type != "") {
		item.Enclosure = &RssEnclosure{URL: i.Link.Href, Type: i.Link.Type, Length: i.Link.Length}
	}
	if i.Author != nil {
		item.Author = i.Author.Name
	}
	return item
}

// RssFeed : create a new RssFeed with a generic Feed struct's data
func (r *Rss) RssFeed() *RssFeed {
	pub := anyTimeFormat(time.RFC1123Z, r.Created, r.Updated)
	build := anyTimeFormat(time.RFC1123Z, r.Updated)
	author := ""
	if r.Author != nil {
		author = r.Author.Email
		if len(r.Author.Name) > 0 {
			author = fmt.Sprintf("%s (%s)", r.Author.Email, r.Author.Name)
		}
	}

	channel := &RssFeed{
		Title:          r.Title,
		Link:           r.Link.Href,
		Description:    r.Description,
		ManagingEditor: author,
		PubDate:        pub,
		LastBuildDate:  build,
		Copyright:      r.Copyright,
	}
	for _, i := range r.Items {
		channel.Items = append(channel.Items, newRssItem(i))
	}
	return channel
}

// FeedXML : return an XML-Ready object for an Rss object
func (r *Rss) FeedXML() interface{} {
	// only generate version 2.0 feeds for now
	return r.RssFeed().FeedXML()

}

// FeedXML : return an XML-ready object for an RssFeed object
func (r *RssFeed) FeedXML() interface{} {
	// modified for Nyaa
	return &rssFeedXML{Version: "2.0", Channel: r, XMLNSTorrent: "http://xmlns.nyaa.pantsu.cat/torrent/"}
}
