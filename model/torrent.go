package model

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util"

//	"encoding/json"
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
	Id           uint         `gorm:"column:torrent_id;primary_key"`
	Name         string       `gorm:"column:torrent_name"`
	Hash         string       `gorm:"column:torrent_hash"`
	Category     int          `gorm:"column:category"`
	Sub_Category int          `gorm:"column:sub_category"`
	Status       int          `gorm:"column:status"`
	Date         time.Time    `gorm:"column:date"`
	UploaderId   uint         `gorm:"column:uploader"`
	Downloads    int          `gorm:"column:downloads"`
	Stardom      int          `gorm:"column:stardom"`
	Filesize     int64        `gorm:"column:filesize"`
	Description  string       `gorm:"column:description"`
	WebsiteLink  string       `gorm:"column:website_link"`

	Uploader     *User        `gorm:"ForeignKey:UploaderId"`
	OldComments  []OldComment `gorm:"ForeignKey:Id"`
	Comments     []Comment    `gorm:"ForeignKey:Id"`
}

/* We need JSON Object instead because of Magnet URL that is not in the database but generated dynamically
--------------------------------------------------------------------------------------------------------------
JSON Models Oject
--------------------------------------------------------------------------------------------------------------
*/

type ApiResultJson struct {
	Torrents         []TorrentsJson `json:"torrents"`
	QueryRecordCount int            `json:"queryRecordCount"`
	TotalRecordCount int            `json:"totalRecordCount"`
}

type OldCommentsJson struct {
	C  template.HTML `json:"c"`
	Us string        `json:"us"`
	Un string        `json:"un"`
	UI int           `json:"ui"`
	T  int           `json:"t"`
	Av string        `json:"av"`
	ID string        `json:"id"`
}

type CommentsJson struct {
	Content  template.HTML `json:"content"`
	Username string        `json:"username"`
}

type TorrentsJson struct {
	Id           string         `json:"id"`
	Name         string         `json:"name"`
	Status       int            `json:"status"`
	Hash         string         `json:"hash"`
	Date         string         `json:"date"`
	Filesize     string         `json:"filesize"`
	Description  template.HTML  `json:"description"`
	Comments     []CommentsJson `json:"comments"`
	Sub_Category string         `json:"sub_category"`
	Category     string         `json:"category"`
	Magnet       template.URL   `json:"magnet"`
}

/* Model Conversion to Json */

func (t *Torrents) ToJson() TorrentsJson {
	magnet := util.InfoHashToMagnet(strings.TrimSpace(t.Hash), t.Name, config.Trackers...)
	//offset := 0
	var commentsJson []CommentsJson
	/*if len(t.OldComments) != 0 {
		b := []OldCommentsJson{}
		err := json.Unmarshal([]byte(t.OldComments), &b)
		if err == nil {
			commentsJson = make([]CommentsJson, len(t.Comments)+len(b))
			offset = len(b)
			for i, commentJson := range b {
				commentsJson[i] = CommentsJson{Content: commentJson.C,
					Username: commentJson.Un}
			}
		} else {
			commentsJson = make([]CommentsJson, len(t.Comments))
		}
	} else {
		commentsJson = make([]CommentsJson, len(t.Comments))
	}
	for i, comment := range t.Comments {
		commentsJson[i+offset] = CommentsJson{Content: template.HTML(comment.Content), Username: comment.Username}
	}*/
	res := TorrentsJson{
		Id:           strconv.FormatUint(uint64(t.Id), 10),
		Name:         html.UnescapeString(t.Name),
		Status:       t.Status,
		Hash:         t.Hash,
		Date:         t.Date.Format(time.RFC3339),
		Filesize:     util.FormatFilesize2(t.Filesize),
		Description:  template.HTML(t.Description),
		Comments:     commentsJson,
		Sub_Category: strconv.Itoa(t.Sub_Category),
		Category:     strconv.Itoa(t.Category),
		Magnet:       util.Safe(magnet)}

	return res
}

/* Complete the functions when necessary... */
