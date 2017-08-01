package models

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"

	"net/http"
	"net/url"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/cache"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/sanitize"
	"github.com/bradfitz/slice"
	"github.com/fatih/structs"
)

const (
	// TorrentStatusNormal Int for Torrent status normal
	TorrentStatusNormal = 1
	// TorrentStatusRemake Int for Torrent status remake
	TorrentStatusRemake = 2
	// TorrentStatusTrusted Int for Torrent status trusted
	TorrentStatusTrusted = 3
	// TorrentStatusAPlus Int for Torrent status a+
	TorrentStatusAPlus = 4
	// TorrentStatusBlocked Int for Torrent status locked
	TorrentStatusBlocked = 5
)

// Torrent model
type Torrent struct {
	ID          uint      `gorm:"column:torrent_id;primary_key"`
	Name        string    `gorm:"column:torrent_name"`
	Hash        string    `gorm:"column:torrent_hash;unique"`
	Category    int       `gorm:"column:category"`
	SubCategory int       `gorm:"column:sub_category"`
	Status      int       `gorm:"column:status"`
	Hidden      bool      `gorm:"column:hidden"`
	Date        time.Time `gorm:"column:date"`
	UploaderID  uint      `gorm:"column:uploader"`
	Stardom     int       `gorm:"column:stardom"`
	Filesize    int64     `gorm:"column:filesize"`
	Description string    `gorm:"column:description"`
	WebsiteLink string    `gorm:"column:website_link"`
	DbID        string    `gorm:"column:db_id"`
	Trackers    string    `gorm:"column:trackers"`
	// Indicates the language of the torrent's content (eg. subs, dubs, raws, manga TLs)
	Language  string `gorm:"column:language"`
	DeletedAt *time.Time

	Uploader    *User        `gorm:"AssociationForeignKey:UploaderID;ForeignKey:user_id"`
	OldUploader string       `gorm:"-"` // ???????
	OldComments []OldComment `gorm:"ForeignKey:torrent_id"`
	Comments    []Comment    `gorm:"ForeignKey:torrent_id"`
	Tags        Tags          `gorm:"-"`
	Scrape      *Scrape      `gorm:"AssociationForeignKey:ID;ForeignKey:torrent_id"`
	FileList    []File       `gorm:"ForeignKey:torrent_id"`
	Languages   []string     `gorm:"-"` // This is parsed when retrieved from db
}

/* We need a JSON object instead of a Gorm structure because magnet URLs are
   not in the database and have to be generated dynamically */

// TorrentJSON for torrent model in json for api
type TorrentJSON struct {
	ID           uint          `json:"id"`
	Name         string        `json:"name"`
	Status       int           `json:"status"`
	Hidden       bool          `json:"-"`
	Hash         string        `json:"hash"`
	Date         string        `json:"date"`
	Filesize     int64         `json:"filesize"`
	Description  template.HTML `json:"description"`
	Comments     []CommentJSON `json:"comments"`
	SubCategory  string        `json:"sub_category"`
	Category     string        `json:"category"`
	DbID         string        `json:"db_id"`
	UploaderID   uint          `json:"uploader_id"`
	UploaderName template.HTML `json:"uploader_name"`
	OldUploader  template.HTML `json:"uploader_old"`
	WebsiteLink  template.URL  `json:"website_link"`
	Languages    []string      `json:"languages"`
	Magnet       template.URL  `json:"magnet"`
	TorrentLink  template.URL  `json:"torrent"`
	Seeders      uint32        `json:"seeders"`
	Leechers     uint32        `json:"leechers"`
	Completed    uint32        `json:"completed"`
	LastScrape   time.Time     `json:"last_scrape"`
	FileList     []FileJSON    `json:"file_list"`
	Tags         Tags           `json:"-"` // not needed in json to reduce db calls
}

// Size : Returns the total size of memory recursively allocated for this struct
// FIXME: Is it deprecated?
func (t Torrent) Size() (s int) {
	s = int(reflect.TypeOf(t).Size())
	return
}

// TableName : Return the table name of torrents table
func (t Torrent) TableName() string {
	return config.Get().Models.TorrentsTableName
}

// Identifier : Return the identifier of a torrent
func (t *Torrent) Identifier() string {
	return fmt.Sprintf("torrent_%d", t.ID)
}

// IsNormal : Return if a torrent status is normal
func (t *Torrent) IsNormal() bool {
	return t.Status == TorrentStatusNormal
}

// IsRemake : Return if a torrent status is remake
func (t *Torrent) IsRemake() bool {
	return t.Status == TorrentStatusRemake
}

// IsTrusted : Return if a torrent status is trusted
func (t *Torrent) IsTrusted() bool {
	return t.Status == TorrentStatusTrusted
}

// IsAPlus : Return if a torrent status is a+
func (t *Torrent) IsAPlus() bool {
	return t.Status == TorrentStatusAPlus
}

// IsBlocked : Return if a torrent status is locked
func (t *Torrent) IsBlocked() bool {
	return t.Status == TorrentStatusBlocked
}

// IsDeleted : Return if a torrent status is deleted
func (t *Torrent) IsDeleted() bool {
	return t.DeletedAt != nil
}

// AddToESIndex : Adds a torrent to Elastic Search
func (t Torrent) AddToESIndex(client *elastic.Client) error {
	ctx := context.Background()
	torrentJSON := t.ToJSON()
	_, err := client.Index().
		Index(config.Get().Search.ElasticsearchIndex).
		Type(config.Get().Search.ElasticsearchType).
		Id(strconv.FormatUint(uint64(torrentJSON.ID), 10)).
		BodyJson(torrentJSON).
		Refresh("true").
		Do(ctx)
	return err
}

// DeleteFromESIndex : Removes a torrent from Elastic Search
func (t *Torrent) DeleteFromESIndex(client *elastic.Client) error {
	ctx := context.Background()
	_, err := client.Delete().
		Index(config.Get().Search.ElasticsearchIndex).
		Type(config.Get().Search.ElasticsearchType).
		Id(strconv.FormatInt(int64(t.ID), 10)).
		Do(ctx)
	return err
}

// ParseTrackers : Takes an array of trackers, adds needed trackers and parse it to url string
func (t *Torrent) ParseTrackers(trackers []string) {
	v := url.Values{}
	if len(config.Get().Torrents.Trackers.NeededTrackers) > 0 { // if we have some needed trackers configured
		if len(trackers) == 0 {
			trackers = config.Get().Torrents.Trackers.Default
		} else {
			for _, id := range config.Get().Torrents.Trackers.NeededTrackers {
				found := false
				for _, tracker := range trackers {
					if tracker == config.Get().Torrents.Trackers.Default[id] {
						found = true
						break
					}
				}
				if !found {
					trackers = append(trackers, config.Get().Torrents.Trackers.Default[id])
				}
			}
		}
	}
	v["tr"] = trackers
	t.Trackers = v.Encode()
}

func (t *Torrent) ParseLanguages() {
	t.Languages = strings.Split(t.Language, ",")
}

func (t *Torrent) EncodeLanguages() {
	t.Language = strings.Join(t.Languages, ",")
}

// GetTrackersArray : Convert trackers string to Array
func (t *Torrent) GetTrackersArray() (trackers []string) {
	v, _ := url.ParseQuery(t.Trackers)
	trackers = v["tr"]
	return
}

// ToTorrent :
// TODO: Need to get rid of TorrentJSON altogether and have only one true Torrent
//       model
func (t *TorrentJSON) ToTorrent() Torrent {
	category, err := strconv.ParseInt(t.Category, 10, 64)
	if err != nil {
		category = 0
	}
	subCategory, err := strconv.ParseInt(t.SubCategory, 10, 64)
	if err != nil {
		subCategory = 0
	}
	// Need to add +00:00 at the end because ES doesn't store it by default
	dateFixed := t.Date
	if len(dateFixed) > 6 && dateFixed[len(dateFixed)-6] != '+' {
		dateFixed += "Z"
	}
	date, err := time.Parse(time.RFC3339, dateFixed)
	if err != nil {
		log.Errorf("Problem parsing date '%s' from ES: %s", dateFixed, err)
	}
	torrent := Torrent{
		ID:          t.ID,
		Name:        t.Name,
		Hash:        t.Hash,
		Category:    int(category),
		SubCategory: int(subCategory),
		Status:      t.Status,
		Date:        date,
		UploaderID:  t.UploaderID,
		//Stardom: t.Stardom,
		Filesize:    t.Filesize,
		Description: string(t.Description),
		Hidden:      t.Hidden,
		//WebsiteLink: t.WebsiteLink,
		//Trackers: t.Trackers,
		//DeletedAt: t.DeletedAt,
		// Uploader: TODO
		//OldUploader: t.OldUploader,
		//OldComments: TODO
		// Comments: TODO
		// LastScrape not stored in ES, counts won't show without a value however
		Scrape:    &Scrape{Seeders: t.Seeders, Leechers: t.Leechers, Completed: t.Completed, LastScrape: time.Now()},
		Languages: t.Languages,
		//FileList: TODO
	}
	torrent.EncodeLanguages()
	return torrent
}

// ToJSON converts a models.Torrent to its equivalent JSON structure
func (t *Torrent) ToJSON() TorrentJSON {
	var trackers []string
	if t.Trackers == "" {
		trackers = config.Get().Torrents.Trackers.Default
	} else {
		trackers = t.GetTrackersArray()
	}
	magnet := format.InfoHashToMagnet(strings.TrimSpace(t.Hash), t.Name, trackers...)
	commentsJSON := make([]CommentJSON, 0, len(t.OldComments)+len(t.Comments))
	for _, c := range t.OldComments {
		commentsJSON = append(commentsJSON, CommentJSON{Username: c.Username, UserID: -1, Content: template.HTML(c.Content), Date: c.Date.UTC()})
	}
	for _, c := range t.Comments {
		if c.User != nil {
			commentsJSON = append(commentsJSON, CommentJSON{Username: c.User.Username, UserID: int(c.User.ID), Content: sanitize.MarkdownToHTML(c.Content), Date: c.CreatedAt.UTC(), UserAvatar: c.User.MD5})
		} else {
			commentsJSON = append(commentsJSON, CommentJSON{})
		}
	}

	// Sort comments by date
	slice.Sort(commentsJSON, func(i, j int) bool {
		return commentsJSON[i].Date.Before(commentsJSON[j].Date)
	})

	fileListJSON := make([]FileJSON, 0, len(t.FileList))
	for _, f := range t.FileList {
		fileListJSON = append(fileListJSON, FileJSON{
			Path:     filepath.Join(f.Path()...),
			Filesize: f.Filesize,
		})
	}

	// Sort file list by lowercase filename
	slice.Sort(fileListJSON, func(i, j int) bool {
		return strings.ToLower(fileListJSON[i].Path) < strings.ToLower(fileListJSON[j].Path)
	})

	uploader := "れんちょん" // by default
	var uploaderID uint
	if t.UploaderID > 0 && t.Uploader != nil {
		uploader = t.Uploader.Username
		uploaderID = t.UploaderID
	} else if t.OldUploader != "" {
		uploader = t.OldUploader
	}
	torrentlink := ""
	if t.ID <= config.Get().Models.LastOldTorrentID && len(config.Get().Torrents.CacheLink) > 0 {
		if config.IsSukebei() {
			torrentlink = "" // torrent cache doesn't have sukebei torrents
		} else {
			torrentlink = fmt.Sprintf(config.Get().Torrents.CacheLink, t.Hash)
		}
	} else if t.ID > config.Get().Models.LastOldTorrentID && len(config.Get().Torrents.StorageLink) > 0 {
		torrentlink = fmt.Sprintf(config.Get().Torrents.StorageLink, t.Hash)
	}
	scrape := Scrape{}
	if t.Scrape != nil {
		scrape = *t.Scrape
	}
	t.ParseLanguages()
	res := TorrentJSON{
		ID:           t.ID,
		Name:         t.Name,
		Status:       t.Status,
		Hidden:       t.Hidden,
		Hash:         t.Hash,
		Date:         t.Date.Format(time.RFC3339),
		Filesize:     t.Filesize,
		Description:  sanitize.MarkdownToHTML(t.Description),
		Comments:     commentsJSON,
		SubCategory:  strconv.Itoa(t.SubCategory),
		Category:     strconv.Itoa(t.Category),
		UploaderID:   uploaderID,
		UploaderName: sanitize.SafeText(uploader),
		WebsiteLink:  sanitize.Safe(t.WebsiteLink),
		Languages:    t.Languages,
		Magnet:       template.URL(magnet),
		TorrentLink:  sanitize.Safe(torrentlink),
		Leechers:     scrape.Leechers,
		Seeders:      scrape.Seeders,
		Completed:    scrape.Completed,
		LastScrape:   scrape.LastScrape,
		FileList:     fileListJSON,
		Tags:         t.Tags,
	}

	return res
}

/* Complete the functions when necessary... */

// TorrentsToJSON : Map Torrents to TorrentsToJSON without reallocations
func TorrentsToJSON(t []Torrent) []TorrentJSON {
	json := make([]TorrentJSON, len(t))
	for i := range t {
		json[i] = t[i].ToJSON()
	}
	return json
}

// Update : Update a torrent based on model
func (t *Torrent) Update(unscope bool) (int, error) {
	db := ORM
	if unscope {
		db = ORM.Unscoped()
	}
	t.EncodeLanguages() // Need to transform array into single string

	if db.Model(t).UpdateColumn(t.toMap()).Error != nil {
		return http.StatusInternalServerError, errors.New("Torrent was not updated")
	}

	// TODO Don't create a new client for each request
	if ElasticSearchClient != nil {
		err := t.AddToESIndex(ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully updated torrent to ES index.")
		} else {
			log.Errorf("Unable to update torrent to ES index: %s", err)
		}
	}
	// We only flush cache after update
	cache.C.Delete(t.Identifier())

	return http.StatusOK, nil
}

// UpdateUnscope : Update a torrent based on model
func (t *Torrent) UpdateUnscope() (int, error) {
	return t.Update(true)
}

// Delete : delete a torrent based on id
func (t *Torrent) Delete(definitely bool) (*Torrent, int, error) {
	if t.ID == 0 {
		err := errors.New("ERROR: Tried to delete a torrent with ID 0")
		log.CheckErrorWithMessage(err, "ERROR_IMPORTANT: ")
		return t, http.StatusBadRequest, err
	}
	db := ORM
	if definitely {
		db = ORM.Unscoped()
	}
	if db.Delete(t).Error != nil {
		return t, http.StatusInternalServerError, errors.New("torrent_not_deleted")
	}

	if ElasticSearchClient != nil {
		err := t.DeleteFromESIndex(ElasticSearchClient)
		if err == nil {
			log.Infof("Successfully deleted torrent to ES index.")
		} else {
			log.Errorf("Unable to delete torrent to ES index: %s", err)
		}
	}
	// We flush cache only after delete
	cache.C.Flush()
	return t, http.StatusOK, nil
}

// DefinitelyDelete : deletes definitely a torrent based on id
func (t *Torrent) DefinitelyDelete() (*Torrent, int, error) {
	return t.Delete(true)

}

// toMap : convert the model to a map of interface
func (t *Torrent) toMap() map[string]interface{} {
	return structs.Map(t)
}

// LoadTags : load all the unique tags with summed up weight from the database in torrent
func (t *Torrent) LoadTags() {
	// Only load if necessary
	if len(t.Tags) == 0 {
		// Should output a query like this: SELECT tag, type, accepted, SUM(weight) as total FROM tags WHERE torrent_id=923000 GROUP BY type, tag ORDER BY type, total DESC
		err := ORM.Select("tag, type, accepted, SUM(weight) as total").Where("torrent_id = ?", t.ID).Group("type, tag").Order("type ASC, total DESC").Find(&t.Tags).Error
		log.CheckErrorWithMessage(err, "LOAD_TAGS_ERROR: Couldn't load tags!")
	}
}
