package main

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"html/template"
	"strconv"
	"strings"
)

type Categories struct {
	Id   int `gorm:"column:category_id"`
	Name string `gorm:"column:category_name"`
	Torrents      []Torrents       `gorm:"ForeignKey:category_id;AssociationForeignKey:category_id"`
	Sub_Categories  []Sub_Categories `gorm:"ForeignKey:parent_id;AssociationForeignKey:category_id"`
}

type Sub_Categories struct {
	Id   int          `gorm:"column:sub_category_id"`
	Name string       `gorm:"column:Sub_category_name"`
	Parent_id         int          `gorm:"column:parent_id"`
	Torrents          []Torrents `gorm:"ForeignKey:sub_category_id;AssociationForeignKey:sub_category_id"`
}

type Torrents struct {
	gorm.Model
	Id              int            `gorm:"column:torrent_id"`
	Name            string         `gorm:"column:torrent_name"`
	Category_id     int            `gorm:"column:category_id"`
	Sub_category_id int            `gorm:"column:sub_category_id"`
	Status          int            `gorm:"column:status_id"`
	Hash            string         `gorm:"column:torrent_hash"`
	Categories      Categories     `gorm:"ForeignKey:category_id;AssociationForeignKey:category_id"`
	Sub_Categories  Sub_Categories `gorm:"ForeignKey:sub_category_id;AssociationForeignKey:sub_category_id"`
}

/* We need JSON Object instead because of Magnet URL that is not in the database but generated dynamically
--------------------------------------------------------------------------------------------------------------
JSON Models Oject
--------------------------------------------------------------------------------------------------------------
*/

type CategoryJson struct {
	Id 				 string  		`json: "id"`
	Name         string         `json: "category"`
	Torrents         []TorrentsJson `json: "torrents"`
	QueryRecordCount int            `json: "queryRecordCount"`
	TotalRecordCount int            `json: "totalRecordCount"`
}

type SubCategoryJson struct {
	Id 				string  		`json: "id"`
	Name     string         	`json: "category"`
}

type TorrentsJson struct {
	Id     string       `json: "id"` // Is there a need to put the ID?
	Name   string       `json: "name"`
	Status int          `json: "status"`
	Hash   string       `json: "hash"`
	Sub_Category SubCategoryJson `json: "sub_category"`
	Category CategoryJson `json: "category"`
	Magnet template.URL `json: "magnet"`
}

type WhereParams struct {
	conditions string // Ex : name LIKE ? AND category_id LIKE ?
	params     []interface{}
}


/* Function to interact with Models
 *
 * Get the torrents with where clause
 *
 */

func getTorrentById(id string) (Torrents, error) {
	var torrent Torrents

	if db.Where("torrent_id = ?", id).Find(&torrent).RecordNotFound() {
		return torrent, errors.New("Article is not found.")
	}

	return torrent, nil
}

func getTorrentsOrderBy(parameters *WhereParams, orderBy string, limit int, offset int) ([]Torrents, int) {
	var torrents []Torrents
	var dbQuery *gorm.DB
	var count int
	if (parameters != nil) {  // if there is where parameters
		db.Model(&torrents).Where(parameters.conditions, parameters.params...).Count(&count)
		dbQuery = db.Model(&torrents).Where(parameters.conditions, parameters.params...)
	} else {
		db.Model(&torrents).Count(&count)
		dbQuery = db.Model(&torrents)
	} 
	
	if (orderBy == "") { orderBy = "torrent_id DESC" } // Default OrderBy
	if limit != 0 || offset != 0 { // if limits provided
		dbQuery = dbQuery.Limit(limit).Offset(offset)
	}
		dbQuery.Order(orderBy).Preload("Categories").Preload("Sub_Categories").Find(&torrents)
	return torrents, count
}

/* Functions to simplify the get parameters of the main function 
 * 
 * Get Torrents with where parameters and limits, order by default
 */
func getTorrents(parameters WhereParams, limit int, offset int) ([]Torrents, int) {
	return getTorrentsOrderBy(&parameters, "", limit, offset)
}

/* Get Torrents with where parameters but no limit and order by default (get all the torrents corresponding in the db)
 */
func getTorrentsDB(parameters WhereParams) ([]Torrents, int) {
	return getTorrentsOrderBy(&parameters, "", 0, 0)
}

/* Function to get all torrents
 */

func getAllTorrentsOrderBy(orderBy string, limit int, offset int) ([]Torrents, int) {
	return getTorrentsOrderBy(nil, orderBy, limit, offset)
}

func getAllTorrents(limit int, offset int) ([]Torrents, int) {
	return getTorrentsOrderBy(nil, "", limit, offset)
}

func getAllTorrentsDB() ([]Torrents, int) {
	return getTorrentsOrderBy(nil, "", 0, 0)
}

/* Function to get all categories with/without torrents (bool)
 */
func getAllCategories(populatedWithTorrents bool) []Categories {
	var categories []Categories
	if populatedWithTorrents {
		db.Preload("Torrents").Preload("Sub_Categories").Find(&categories)
	} else {
		db.Preload("Sub_Categories").Find(&categories)
	}
	return categories
}

func createWhereParams(conditions string, params ...string) WhereParams {
	whereParams := WhereParams{}
	whereParams.conditions = conditions
	for i, _ := range params {
		whereParams.params = append(whereParams.params, params[i])
	}

	return whereParams
}

/* Model Conversion to Json */

func (t *Torrents) toJson() TorrentsJson {
	magnet := "magnet:?xt=urn:btih:" + strings.TrimSpace(t.Hash) + "&dn=" + t.Name + trackers
	res := TorrentsJson{
		Id:     strconv.Itoa(t.Id),
		Name:   html.UnescapeString(t.Name),
		Status: t.Status,
		Hash:   t.Hash,
		Sub_Category: t.Sub_Categories.toJson(),
		Category: t.Categories.toJson(),
		Magnet: safe(magnet)}
	return res
}

func (c *Sub_Categories) toJson() SubCategoryJson {
	return SubCategoryJson{
		Id:     strconv.Itoa(c.Id),
		Name:   html.UnescapeString(c.Name)}
}

func (c *Categories) toJson() CategoryJson {
	return CategoryJson{
		Id:     strconv.Itoa(c.Id),
		Name:   html.UnescapeString(c.Name)}
}

/* Complete the functions when necessary... */