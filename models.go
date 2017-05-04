package main

import (
	"github.com/jinzhu/gorm"
	"html"
	"html/template"
	"strconv"
	"net/url"
	"errors"
)

type Categories struct {
	Category_id      int    
	Category_name    string    
	Torrents          	[]Torrents `gorm:"ForeignKey:category_id;AssociationForeignKey:category_id"`
	Sub_category    []Sub_Categories `gorm:"ForeignKey:category_id;AssociationForeignKey:parent_id"`
}

type Sub_Categories struct {
	Sub_category_id     int    
	Sub_category_name   string  
	Parent_id           int  
	Torrents          	[]Torrents `gorm:"ForeignKey:sub_category_id;AssociationForeignKey:sub_category_id"`
}

type Torrents struct {
	gorm.Model
	Id      		int        `gorm:"column:torrent_id"` 
	Name   			string     `gorm:"column:torrent_name"` 
	Category_id 	int        `gorm:"column:category_id"`  
	Sub_category_id int         
	Status 			int        `gorm:"column:status_id"` 
	Hash  			string     `gorm:"column:torrent_hash"` 
	Categories        Categories `gorm:"ForeignKey:category_id;AssociationForeignKey:category_id"`
	Sub_Categories    Sub_Categories `gorm:"ForeignKey:sub_category_id;AssociationForeignKey:sub_category_id"`
}



/* We need JSON Object instead because of Magnet URL that is not in the database but generated dynamically
--------------------------------------------------------------------------------------------------------------
JSON Models Oject
--------------------------------------------------------------------------------------------------------------
*/

type CategoryJson struct {
	Category         string    `json: "category"`
	Torrents          []TorrentsJson `json: "torrents"`
	QueryRecordCount int       `json: "queryRecordCount"`
	TotalRecordCount int       `json: "totalRecordCount"`
}

type TorrentsJson struct {
	Id     string       `json: "id"`  // Is there a need to put the ID?
	Name   string       `json: "name"`
	Status int          `json: "status"`
	Hash   string       `json: "hash"`
	Magnet template.URL `json: "magnet"`
}

type WhereParams struct {
	conditions string // Ex : name LIKE ? AND category_id LIKE ?
	params []interface{} 
}


/* Each Page should have an object to pass to their own template */

type HomeTemplateVariables struct {
	ListTorrents   []TorrentsJson
	ListCategories []Categories
	Query string
	Category string
	QueryRecordCount int   
	TotalRecordCount int 
}


/* Function to interact with Models
 *
 * Get the torrents with where clause
 *
*/

func getTorrentById(id string) (Torrents, error) {
	var torrent Torrents

	if db.Order("torrent_id DESC").First(&torrent, "id = ?", html.EscapeString(id)).RecordNotFound() {
		return torrent, errors.New("Article is not found.")
	}

	return torrent, nil
}

func getTorrents(parameters WhereParams, limit int, offset int) ([]Torrents) {
	var torrents []Torrents
	db.Limit(limit).Offset(offset).Order("torrent_id DESC").Where(parameters.conditions, parameters.params...).Preload("Categories").Preload("Sub_Categories").Find(&torrents)
	return torrents
}

func getTorrentsDB(parameters WhereParams) ([]Torrents) {
	var torrents []Torrents
	db.Where(parameters.conditions, parameters.params...).Order("torrent_id DESC").Preload("Categories").Preload("Sub_Categories").Find(&torrents)
	return torrents
}

/* Function to get all torrents 
*/

func getAllTorrents(limit int, offset int) ([]Torrents) {
	var torrents []Torrents
	db.Model(&torrents).Limit(limit).Offset(offset).Order("torrent_id DESC").Preload("Categories").Preload("Sub_Categories").Find(&torrents)

	return torrents
}

func getAllTorrentsDB() ([]Torrents) {
	var torrents []Torrents
	db.Order("torrent_id DESC").Preload("Categories").Preload("Sub_Categories").Find(&torrents)
	return torrents
}

/* Function to get all categories with/without torrents (bool)
*/
func getAllCategories(populatedWithTorrents bool) ([]Categories) {
	var categories []Categories
	if populatedWithTorrents {
		db.Preload("Torrents").Preload("Sub_Categories").Find(&categories)
	} else {
		db.Preload("Sub_Categories").Find(&categories)
	}
	return categories
}



func (t *Torrents) toJson() (TorrentsJson) {
	magnet := "magnet:?xt=urn:btih:" + t.Hash + "&dn=" + url.QueryEscape(t.Name) + trackers
	res := TorrentsJson{
		Id:     strconv.Itoa(t.Id),
		Name:   t.Name,
		Status: t.Status,
		Hash:   t.Hash,
		Magnet: safe(magnet)}
		return res;
}

func createWhereParams(conditions string, params ...string) (WhereParams) {
	whereParams := WhereParams{}
	whereParams.conditions = conditions
	for i, _ := range params {
		whereParams.params = append(whereParams.params, params[i])
	}

	return whereParams
}
/* Complete the functions when necessary... */