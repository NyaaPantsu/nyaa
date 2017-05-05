package main

import (
	"errors"
	"github.com/jinzhu/gorm"
	"html"
	"html/template"
	"strconv"
	"strings"
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
	Name      string     `gorm:"column:Sub_category_name"`
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
	Date            int            `gorm:"column:date"`
	Downloads       int            `gorm:"column:downloads"`
	Filesize        string         `gorm:"column:filesize"`
	Description     []byte         `gorm:"column:description"`
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

type TorrentsJson struct {
	Id           string          `json: "id"` // Is there a need to put the ID?
	Name         string          `json: "name"`
	Status       int             `json: "status"`
	Hash         string          `json: "hash"`
	Date         int             `json: "date"`
	Filesize     string          `json: "filesize"`
	Description  template.HTML   `json: "description"`
	Sub_Category SubCategoryJson `json: "sub_category"`
	Category     CategoryJson    `json: "category"`
	Magnet       template.URL    `json: "magnet"`
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

// don't need raw SQL once we get MySQL
func getFeeds() []Feed {
	var result []Feed
	rows, err := db.DB().
		Query(
			"SELECT `torrent_id` AS `id`, `torrent_name` AS `name`, `torrent_hash` AS `hash`, `timestamp` FROM `torrents` " +
				"ORDER BY `timestamp` desc LIMIT 50")
	if err == nil {
		for rows.Next() {
			item := Feed{}
			rows.Scan(&item.Id, &item.Name, &item.Hash, &item.Timestamp)
			magnet := "magnet:?xt=urn:btih:" + strings.TrimSpace(item.Hash) + "&dn=" + item.Name + trackers
			item.Magnet = magnet
			// memory hog
			result = append(result, item)
		}
		rows.Close()
	}
	return result
}

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
	conditions := "torrent_hash is not null" //filter out broken entries
	var params []interface{}
	if parameters != nil { // if there is where parameters
		conditions += " AND " + parameters.conditions
		params = parameters.params
	}
	db.Model(&torrents).Where(conditions, params...).Count(&count)
	dbQuery = db.Model(&torrents).Where(conditions, params...)

	if orderBy == "" {
		orderBy = "torrent_id DESC"
	} // Default OrderBy
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
		Id:           strconv.Itoa(t.Id),
		Name:         html.UnescapeString(t.Name),
		Status:       t.Status,
		Hash:         t.Hash,
		Date:         t.Date,
		Filesize:     t.Filesize,
		Description:  template.HTML(unZlib(t.Description)),
		Sub_Category: t.Sub_Categories.toJson(),
		Category:     t.Categories.toJson(),
		Magnet:       safe(magnet)}

	return res
}

func (c *Sub_Categories) toJson() SubCategoryJson {
	return SubCategoryJson{
		Id:   strconv.Itoa(c.Id),
		Name: html.UnescapeString(c.Name)}
}

func (c *Categories) toJson() CategoryJson {
	return CategoryJson{
		Id:   strconv.Itoa(c.Id),
		Name: html.UnescapeString(c.Name)}
}

/* Complete the functions when necessary... */
