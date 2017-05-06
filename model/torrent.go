package model

import (
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/util"

	"encoding/json"
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
	Id              int            `gorm:"column:torrent_id"`
	Name            string         `gorm:"column:torrent_name"`
	Category        int            `gorm:"column:category_id"`
	Sub_Category    int            `gorm:"column:sub_category_id"`
	Status          int            `gorm:"column:status_id"`
	Hash            string         `gorm:"column:torrent_hash"`
	Date            int64          `gorm:"column:date"`
	Downloads       int            `gorm:"column:downloads"`
	Filesize        int64          `gorm:"column:filesize"`
	Description     []byte         `gorm:"column:description"`
	Comments        []byte         `gorm:"column:comments"`
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

type CommentsJson struct {
	C  template.HTML `json:"c"`
	Us string        `json:"us"`
	Un string        `json:"un"`
	UI int           `json:"ui"`
	T  int           `json:"t"`
	Av string        `json:"av"`
	ID string        `json:"id"`
}

type TorrentsJson struct {
	Id           string          `json:"id"`
	Name         string          `json:"name"`
	Status       int             `json:"status"`
	Hash         string          `json:"hash"`
	Date         string          `json:"date"`
	Filesize     string          `json:"filesize"`
	Description  template.HTML   `json:"description"`
	Comments     []CommentsJson  `json:"comments"`
	Sub_Category string          `json:"sub_category"`
	Category     string          `json:"category"`
	Magnet       template.URL    `json:"magnet"`
}

/* Model Conversion to Json */

func (t *Torrents) ToJson() TorrentsJson {
	magnet := "magnet:?xt=urn:btih:" + strings.TrimSpace(t.Hash) + "&dn=" + t.Name + config.Trackers
	b := []CommentsJson{}
	_ = json.Unmarshal([]byte(t.Comments), &b)
	res := TorrentsJson{
		Id:           strconv.Itoa(t.Id),
		Name:         html.UnescapeString(t.Name),
		Status:       t.Status,
		Hash:         t.Hash,
		Date:         time.Unix(t.Date, 0).Format(time.RFC3339),
		Filesize:     util.FormatFilesize(t.Filesize),
		Description:  template.HTML(t.Description),
		Comments:     b,
		Sub_Category: strconv.Itoa(t.Sub_Category),
		Category:     strconv.Itoa(t.Category),
		Magnet:       util.Safe(magnet)}

	return res
}

/* Complete the functions when necessary... */
