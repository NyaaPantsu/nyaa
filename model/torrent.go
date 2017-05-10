package model

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util"

	"html"
	"html/template"
	"strconv"
	"strings"
	"time"
)

type Feed struct {
	Id        int
	Name      string
	Hash      string
	Magnet    string
	Timestamp string
}

type Torrents struct {
	Category     int          `gorm:"column:category"`
	Status       int          `gorm:"column:status"`
	Sub_Category int          `gorm:"column:sub_category"`
	UploaderId   uint         `gorm:"column:uploader"`
	Downloads    int          `gorm:"column:downloads"`
	Stardom      int          `gorm:"column:stardom"`
	Filesize     int64        `gorm:"column:filesize"`
	Id           uint         `gorm:"column:torrent_id;primary_key"`
	Date         time.Time    `gorm:"column:date"`
	Uploader     *User        `gorm:"ForeignKey:uploader"`
	Name         string       `gorm:"column:torrent_name"`
	Hash         string       `gorm:"column:torrent_hash"`
	Description  string       `gorm:"column:description"`
	WebsiteLink  string       `gorm:"column:website_link"`
	OldComments  []OldComment `gorm:"ForeignKey:torrent_id"`
	Comments     []Comment    `gorm:"ForeignKey:torrent_id"`
}

// Returns the total size of memory recursively allocated for this struct
func (t Torrents) Size() (s int) {
	s += 12 + // numbers and pointers
		4*2 + // string pointer sizes
		// string array sizes
		len(t.Name) + len(t.Hash) + len(t.Description) + len(t.WebsiteLink) +
		2*2 // array pointer length
	s *= 8 // Assume 64 bit OS

	s += t.Uploader.Size()
	for _, c := range t.OldComments {
		s += c.Size()
	}
	for _, c := range t.Comments {
		s += c.Size()
	}

	return
}

/* We need JSON Object instead because of Magnet URL that is not in the database but generated dynamically */

type ApiResultJson struct {
	Torrents         []TorrentsJson `json:"torrents"`
	QueryRecordCount int            `json:"queryRecordCount"`
	TotalRecordCount int            `json:"totalRecordCount"`
}

type CommentsJson struct {
	Username string        `json:"username"`
	Content  template.HTML `json:"content"`
	Date     time.Time     `json:"date"`
}

type TorrentsJson struct {
	Status       int            `json:"status"`
	Downloads    int            `json:"downloads"`
	UploaderId   uint           `json:"uploader_id"`
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Hash         string         `json:"hash"`
	Date         string         `json:"date"`
	Filesize     string         `json:"filesize"`
	Sub_Category string         `json:"sub_category"`
	Category     string         `json:"category"`
	Description  template.HTML  `json:"description"`
	WebsiteLink  template.URL   `json:"website_link"`
	Magnet       template.URL   `json:"magnet"`
	Comments     []CommentsJson `json:"comments"`
}

/* Model Conversion to Json */

func (t *Torrents) ToJson() TorrentsJson {
	magnet := util.InfoHashToMagnet(strings.TrimSpace(t.Hash), t.Name, config.Trackers...)
	commentsJson := make([]CommentsJson, 0, len(t.OldComments)+len(t.Comments))
	for _, c := range t.OldComments {
		commentsJson = append(commentsJson, CommentsJson{Username: c.Username, Content: template.HTML(c.Content), Date: c.Date})
	}
	for _, c := range t.Comments {

		commentsJson = append(commentsJson, CommentsJson{Username: c.User.Username, Content: util.MarkdownToHTML(c.Content), Date: c.CreatedAt})
	}
	res := TorrentsJson{
		Id:           strconv.FormatUint(uint64(t.Id), 10),
		Name:         html.UnescapeString(t.Name),
		Status:       t.Status,
		Hash:         t.Hash,
		Date:         t.Date.Format(time.RFC3339),
		Filesize:     util.FormatFilesize2(t.Filesize),
		Description:  util.MarkdownToHTML(t.Description),
		Comments:     commentsJson,
		Sub_Category: strconv.Itoa(t.Sub_Category),
		Category:     strconv.Itoa(t.Category),
		Downloads:    t.Downloads,
		UploaderId:   t.UploaderId,
		WebsiteLink:  util.Safe(t.WebsiteLink),
		Magnet:       util.Safe(magnet)}

	return res
}

/* Complete the functions when necessary... */

// Map Torrents to TorrentsToJSON without reallocations
func TorrentsToJSON(t []Torrents) []TorrentsJson {
	json := make([]TorrentsJson, len(t))
	for i := range t {
		json[i] = t[i].ToJson()
	}
	return json
}
