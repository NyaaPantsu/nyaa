package nyaafeeds

// rss support
// validation done according to spec here:
//    http://cyber.law.harvard.edu/rss/rss.html

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/feeds"
)

// private wrapper around the RssFeed which gives us the <rss>..</rss> xml
type rssFeedXML struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel *RssFeed
}

type RssImage struct {
	XMLName xml.Name `xml:"image"`
	URL     string   `xml:"url"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
	Width   int      `xml:"width,omitempty"`
	Height  int      `xml:"height,omitempty"`
}

type RssTextInput struct {
	XMLName     xml.Name `xml:"textInput"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Name        string   `xml:"name"`
	Link        string   `xml:"link"`
}

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

type RssItem struct {
	XMLName     xml.Name     `xml:"item"`
	Title       string       `xml:"title"`       // required
	Link        string       `xml:"link"`        // required
	Description string       `xml:"description"` // required
	Author      string       `xml:"author,omitempty"`
	Category    *RssCategory `xml:"category,omitempty"`
	Comments    string       `xml:"comments,omitempty"`
	Enclosure   *RssEnclosure
	GUID        string      `xml:"guid,omitempty"`    // Id used
	PubDate     string      `xml:"pubDate,omitempty"` // created or updated
	Source      string      `xml:"source,omitempty"`
	Torrent     *RssTorrent `xml:"torrent,omitempty"`
	Torznab     *RssTorznab `xml:"torznab,omitempty"`
}

type RssTorrent struct {
	XMLName       xml.Name `xml:"torrent"`
	Xmlns         string   `xml:"xmlns,attr"`
	FileName      string   `xml:"fileName,omitempty"`
	ContentLength string   `xml:"contentLength,omitempty"`
	InfoHash      string   `xml:"infoHash,omitempty"`
	MagnetURI     string   `xml:"magnetUri,omitempty"`
}

type RssTorznab struct {
	XMLName              xml.Name `xml:"torznab"`
	Xmlns                string   `xml:"xmlns,attr"`
	Type                 string   `xml:"type,omitempty"`
	Size                 string   `xml:"size,omitempty"`
	Files                string   `xml:"files,omitempty"`
	Grabs                string   `xml:"grabs,omitempty"`
	Tvdbid               string   `xml:"tvdbid,omitempty"`
	Rageid               string   `xml:"rageid,omitempty"`
	Tvmazeid             string   `xml:"tvmazeid,omitempty"`
	Imdb                 string   `xml:"imdb,omitempty"`
	BannerURL            string   `xml:"bannerurl,omitempty"`
	Infohash             string   `xml:"infohash,omitempty"`
	MagnetURL            string   `xml:"magneturl,omitempty"`
	Seeders              string   `xml:"seeders,omitempty"`
	Leechers             string   `xml:"leechers,omitempty"`
	Peers                string   `xml:"peers,omitempty"`
	SeedType             string   `xml:"seedtype,omitempty"`
	MinimumRatio         string   `xml:"minimumratio,omitempty"`
	MinimumSeedTime      string   `xml:"minimumseedtime,omitempty"`
	DownloadVolumeFactor string   `xml:"downloadvolumefactor,omitempty"`
	UploadVolumeFactor   string   `xml:"uploadvolumefactor,omitempty"`
}

// RssCategory is a category for rss item
type RssCategory struct {
	XMLName xml.Name `xml:"category"`
	Domain  string   `xml:"domain"`
}

type RssEnclosure struct {
	//RSS 2.0 <enclosure url="http://example.com/file.mp3" length="123456789" type="audio/mpeg" />
	XMLName xml.Name `xml:"enclosure"`
	URL     string   `xml:"url,attr"`
	Length  string   `xml:"length,attr"`
	Type    string   `xml:"type,attr"`
}

type Rss struct {
	*feeds.Feed
}

// create a new RssItem with a generic Item struct's data
func newRssItem(i *feeds.Item) *RssItem {
	item := &RssItem{
		Title:       i.Title,
		Link:        i.Link.Href,
		Description: i.Description,
		GUID:        i.Id,
		PubDate:     anyTimeFormat(time.RFC1123Z, i.Created, i.Updated),
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

// returns the first non-zero time formatted as a string or ""
func anyTimeFormat(format string, times ...time.Time) string {
	for _, t := range times {
		if !t.IsZero() {
			return t.Format(format)
		}
	}
	return ""
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

// FeedXml : return an XML-Ready object for an Rss object
func (r *Rss) FeedXml() interface{} {
	// only generate version 2.0 feeds for now
	return r.RssFeed().FeedXml()

}

// FeedXml : return an XML-ready object for an RssFeed object
func (r *RssFeed) FeedXml() interface{} {
	return &rssFeedXML{Version: "2.0", Channel: r}
}
