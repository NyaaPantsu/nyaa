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
	ID        int
	Name      string
	Hash      string
	Magnet    string
	Timestamp string
}

type Torrent struct {
	ID          uint      `gorm:"column:torrent_id;primary_key"`
	Name        string    `gorm:"column:torrent_name"`
	Hash        string    `gorm:"column:torrent_hash"`
	Category    int       `gorm:"column:category"`
	SubCategory int       `gorm:"column:sub_category"`
	Status      int       `gorm:"column:status"`
	Date        time.Time `gorm:"column:date"`
	UploaderID  uint      `gorm:"column:uploader"`
	Downloads   int       `gorm:"column:downloads"`
	Stardom     int       `gorm:"column:stardom"`
	Filesize    int64     `gorm:"column:filesize"`
	Description string    `gorm:"column:description"`
	WebsiteLink string    `gorm:"column:website_link"`

	Uploader    *User        `gorm:"ForeignKey:uploader"`
	OldComments []OldComment `gorm:"ForeignKey:torrent_id"`
	Comments    []Comment    `gorm:"ForeignKey:torrent_id"`
}

/* We need a JSON object instead of a Gorm structure because magnet URLs are
   not in the database and have to be generated dynamically */

type ApiResultJSON struct {
	Torrents         []TorrentJSON `json:"torrents"`
	QueryRecordCount int           `json:"queryRecordCount"`
	TotalRecordCount int           `json:"totalRecordCount"`
}

type CommentJSON struct {
	Username string        `json:"username"`
	Content  template.HTML `json:"content"`
	Date     time.Time     `json:"date"`
}

type TorrentJSON struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Status      int           `json:"status"`
	Hash        string        `json:"hash"`
	Date        string        `json:"date"`
	Filesize    string        `json:"filesize"`
	Description template.HTML `json:"description"`
	Comments    []CommentJSON `json:"comments"`
	SubCategory string        `json:"sub_category"`
	Category    string        `json:"category"`
	Magnet      template.URL  `json:"magnet"`
}

// ToJSON converts a model.Torrent to its equivalent JSON structure
func (t *Torrent) ToJSON() TorrentJSON {
	magnet := util.InfoHashToMagnet(strings.TrimSpace(t.Hash), t.Name, config.Trackers...)
	var commentJSON []CommentJSON
	for _, c := range t.OldComments {
		commentJSON = append(commentJSON, CommentJSON{Username: c.Username, Content: template.HTML(c.Content), Date: c.Date})
	}
	for _, c := range t.Comments {
		commentJSON = append(commentJSON, CommentJSON{Username: c.User.Username, Content: template.HTML(c.Content), Date: c.CreatedAt})
	}
	res := TorrentJSON{
		ID:          strconv.FormatUint(uint64(t.ID), 10),
		Name:        html.UnescapeString(t.Name),
		Status:      t.Status,
		Hash:        t.Hash,
		Date:        t.Date.Format(time.RFC3339),
		Filesize:    util.FormatFilesize2(t.Filesize),
		Description: template.HTML(t.Description),
		Comments:    commentJSON,
		SubCategory: strconv.Itoa(t.SubCategory),
		Category:    strconv.Itoa(t.Category),
		Magnet:      util.Safe(magnet)}

	return res
}
