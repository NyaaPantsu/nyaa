package model

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util"

	"fmt"
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
	DeletedAt *time.Time

	Uploader    *User        `gorm:"ForeignKey:UploaderId"`
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
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Status       int           `json:"status"`
	Hash         string        `json:"hash"`
	Date         string        `json:"date"`
	Filesize     string        `json:"filesize"`
	Description  template.HTML `json:"description"`
	Comments     []CommentJSON `json:"comments"`
	SubCategory  string        `json:"sub_category"`
	Category     string        `json:"category"`
	Downloads    int           `json:"downloads"`
	UploaderID   uint          `json:"uploader_id"`
	UploaderName template.HTML `json:"uploader_name"`
	WebsiteLink  template.URL  `json:"website_link"`
	Magnet       template.URL  `json:"magnet"`
	TorrentLink  template.URL  `json:"torrent"`
}

// ToJSON converts a model.Torrent to its equivalent JSON structure
func (t *Torrent) ToJSON() TorrentJSON {
	magnet := util.InfoHashToMagnet(strings.TrimSpace(t.Hash), t.Name, config.Trackers...)
	commentsJSON := make([]CommentJSON, 0, len(t.OldComments)+len(t.Comments))
	for _, c := range t.OldComments {
		escapedContent := template.HTML(html.EscapeString(c.Content))
		commentsJSON = append(commentsJSON, CommentJSON{Username: c.Username, Content: escapedContent, Date: c.Date})
	}
	for _, c := range t.Comments {
		commentsJSON = append(commentsJSON, CommentJSON{Username: c.User.Username, Content: util.MarkdownToHTML(c.Content), Date: c.CreatedAt})
	}
	uploader := ""
	if t.Uploader != nil {
		uploader = t.Uploader.Username
	}
	torrentlink := ""
	if t.ID <= config.LastOldTorrentID && len(config.TorrentCacheLink) > 0 {
		torrentlink = fmt.Sprintf(config.TorrentCacheLink, t.Hash)
	} else if t.ID > config.LastOldTorrentID && len(config.TorrentStorageLink) > 0 {
		torrentlink = fmt.Sprintf(config.TorrentStorageLink, t.Hash) // TODO: Fix as part of configuration changes
	}
	res := TorrentJSON{
		ID:           strconv.FormatUint(uint64(t.ID), 10),
		Name:         t.Name,
		Status:       t.Status,
		Hash:         t.Hash,
		Date:         t.Date.Format(time.RFC3339),
		Filesize:     util.FormatFilesize2(t.Filesize),
		Description:  util.MarkdownToHTML(t.Description),
		Comments:     commentsJSON,
		SubCategory:  strconv.Itoa(t.SubCategory),
		Category:     strconv.Itoa(t.Category),
		Downloads:    t.Downloads,
		UploaderID:   t.UploaderID,
		UploaderName: util.SafeText(uploader),
		WebsiteLink:  util.Safe(t.WebsiteLink),
		Magnet:       util.Safe(magnet),
		TorrentLink:  util.Safe(torrentlink)}

	return res
}

/* Complete the functions when necessary... */

// Map Torrents to TorrentsToJSON without reallocations
func TorrentsToJSON(t []Torrent) []TorrentJSON { // TODO: Convert to singular version
	json := make([]TorrentJSON, len(t))
	for i := range t {
		json[i] = t[i].ToJSON()
	}
	return json
}
