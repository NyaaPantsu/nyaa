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

type Categories struct {
	Id             int              `gorm:"column:category_id"`
	Name           string           `gorm:"column:category_name"`
	Torrents       []Torrents       `gorm:"ForeignKey:category_id;AssociationForeignKey:category_id"`
	Sub_Categories []Sub_Categories `gorm:"ForeignKey:parent_id;AssociationForeignKey:category_id"`
}

type Sub_Categories struct {
	Id        int        `gorm:"column:sub_category_id"`
	Name      string     `gorm:"column:sub_category_name"`
	Parent_id int        `gorm:"column:parent_id"`
	Torrents  []Torrents `gorm:"ForeignKey:sub_category_id;AssociationForeignKey:sub_category_id"`
}

type Statuses struct {
	Status_id   int
	Status_name string
	Torrents    []Torrents `gorm:"ForeignKey:status_id;AssociationForeignKey:status_id"`
}

type Torrents struct {
	Id              int            `gorm:"column:torrent_id"`
	Name            string         `gorm:"column:torrent_name"`
	Category_id     int            `gorm:"column:category_id"`
	Sub_category_id int            `gorm:"column:sub_category_id"`
	Status          int            `gorm:"column:status_id"`
	Hash            string         `gorm:"column:torrent_hash"`
	Date            int64          `gorm:"column:date"`
	Downloads       int            `gorm:"column:downloads"`
	Filesize        string         `gorm:"column:filesize"`
	Description     []byte         `gorm:"column:description"`
	Comments        []byte         `gorm:"column:comments"`
	Statuses        Statuses       `gorm:"ForeignKey:status_id;AssociationForeignKey:status_id"`
	Categories      Categories     `gorm:"ForeignKey:category_id;AssociationForeignKey:category_id"`
	Sub_Categories  Sub_Categories `gorm:"ForeignKey:sub_category_id;AssociationForeignKey:sub_category_id"`
}

/* We need JSON Object instead because of Magnet URL that is not in the database but generated dynamically
--------------------------------------------------------------------------------------------------------------
JSON Models Oject
--------------------------------------------------------------------------------------------------------------
*/

type CategoryJson struct {
	Id               string         `json: "id"`
	Name             string         `json: "category"`
	Torrents         []TorrentsJson `json: "torrents"`
	QueryRecordCount int            `json: "queryRecordCount"`
	TotalRecordCount int            `json: "totalRecordCount"`
}

type SubCategoryJson struct {
	Id   string `json: "id"`
	Name string `json: "category"`
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
	Id           string          `json: "id"` // Is there a need to put the ID?
	Name         string          `json: "name"`
	Status       int             `json: "status"`
	Hash         string          `json: "hash"`
	Date         string          `json: "date"`
	Filesize     string          `json: "filesize"`
	Description  template.HTML   `json: "description"`
	Comments     []CommentsJson  `json: "comments"`
	Sub_Category SubCategoryJson `json: "sub_category"`
	Category     CategoryJson    `json: "category"`
	Magnet       template.URL    `json: "magnet"`
}

/* Model Conversion to Json */

func (t *Torrents) ToJson() TorrentsJson {
	magnet := "magnet:?xt=urn:btih:" + strings.TrimSpace(t.Hash) + "&dn=" + t.Name + config.Trackers
	b := []CommentsJson{}
	_ = json.Unmarshal([]byte(util.UnZlib(t.Comments)), &b)
	res := TorrentsJson{
		Id:           strconv.Itoa(t.Id),
		Name:         html.UnescapeString(t.Name),
		Status:       t.Status,
		Hash:         t.Hash,
		Date:         time.Unix(t.Date, 0).Format(time.RFC3339),
		Filesize:     t.Filesize,
		Description:  template.HTML(util.UnZlib(t.Description)),
		Comments:     b,
		Sub_Category: t.Sub_Categories.ToJson(),
		Category:     t.Categories.ToJson(),
		Magnet:       util.Safe(magnet)}

	return res
}

func (c *Sub_Categories) ToJson() SubCategoryJson {
	return SubCategoryJson{
		Id:   strconv.Itoa(c.Id),
		Name: html.UnescapeString(c.Name)}
}

func (c *Categories) ToJson() CategoryJson {
	return CategoryJson{
		Id:   strconv.Itoa(c.Id),
		Name: html.UnescapeString(c.Name)}
}

/* Complete the functions when necessary... */
