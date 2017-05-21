package model

import (
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/bradfitz/slice"

	"fmt"
	"html/template"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	TorrentStatusNormal  = 1
	TorrentStatusRemake  = 2
	TorrentStatusTrusted = 3
	TorrentStatusAPlus   = 4
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
	DeletedAt   *time.Time

	Uploader    *User        `gorm:"AssociationForeignKey:UploaderID;ForeignKey:user_id"`
	OldUploader string       `gorm:"-"` // ???????
	OldComments []OldComment `gorm:"ForeignKey:torrent_id"`
	Comments    []Comment    `gorm:"ForeignKey:torrent_id"`

	Seeders    uint32    `gorm:"column:seeders"`
	Leechers   uint32    `gorm:"column:leechers"`
	Completed  uint32    `gorm:"column:completed"`
	LastScrape time.Time `gorm:"column:last_scrape"`
	FileList   []File    `gorm:"ForeignKey:torrent_id"`
}

// Returns the total size of memory recursively allocated for this struct
// FIXME: doesn't go have sizeof or something nicer for this?
func (t Torrent) Size() (s int) {
	s += 8 + // ints
		2*3 + // time.Time
		2 + // pointers
		4*2 + // string pointers
		// string array sizes
		len(t.Name) + len(t.Hash) + len(t.Description) + len(t.WebsiteLink) +
		2*2 // array pointers
	s *= 8 // Assume 64 bit OS

	if t.Uploader != nil {
		s += t.Uploader.Size()
	}
	for _, c := range t.OldComments {
		s += c.Size()
	}
	for _, c := range t.Comments {
		s += c.Size()
	}

	return

}

func (t Torrent) TableName() string {
	return config.TorrentsTableName
}

func (t Torrent) Identifier() string {
	return "torrent_"+strconv.Itoa(int(t.ID))
}

func (t Torrent) IsNormal() bool {
	return t.Status == TorrentStatusNormal
}

func (t Torrent) IsRemake() bool {
	return t.Status == TorrentStatusRemake
}

func (t Torrent) IsTrusted() bool {
	return t.Status == TorrentStatusTrusted
}

func (t Torrent) IsAPlus() bool {
	return t.Status == TorrentStatusAPlus
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
	UserID   int           `json:"user_id"`
	Content  template.HTML `json:"content"`
	Date     time.Time     `json:"date"`
}

type FileJSON struct {
	Path     string `json:"path"`
	Filesize int64  `json:"filesize"`
}

type TorrentJSON struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Status       int           `json:"status"`
	Hash         string        `json:"hash"`
	Date         string        `json:"date"`
	Filesize     int64         `json:"filesize"`
	Description  template.HTML `json:"description"`
	Comments     []CommentJSON `json:"comments"`
	SubCategory  string        `json:"sub_category"`
	Category     string        `json:"category"`
	Downloads    int           `json:"downloads"`
	UploaderID   uint          `json:"uploader_id"`
	UploaderName template.HTML `json:"uploader_name"`
	OldUploader  template.HTML `json:"uploader_old"`
	WebsiteLink  template.URL  `json:"website_link"`
	Magnet       template.URL  `json:"magnet"`
	TorrentLink  template.URL  `json:"torrent"`
	Seeders      uint32        `json:"seeders"`
	Leechers     uint32        `json:"leechers"`
	Completed    uint32        `json:"completed"`
	LastScrape   time.Time     `json:"last_scrape"`
	FileList     []FileJSON    `json:"file_list"`
}

// ToJSON converts a model.Torrent to its equivalent JSON structure
func (t *Torrent) ToJSON() TorrentJSON {
	magnet := util.InfoHashToMagnet(strings.TrimSpace(t.Hash), t.Name, config.Trackers...)
	commentsJSON := make([]CommentJSON, 0, len(t.OldComments)+len(t.Comments))
	for _, c := range t.OldComments {
		commentsJSON = append(commentsJSON, CommentJSON{Username: c.Username, UserID: -1, Content: template.HTML(c.Content), Date: c.Date.UTC()})
	}
	for _, c := range t.Comments {
		if c.User != nil {
			commentsJSON = append(commentsJSON, CommentJSON{Username: c.User.Username, UserID: int(c.User.ID), Content: util.MarkdownToHTML(c.Content), Date: c.CreatedAt.UTC()})
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

	uploader := ""
	if t.Uploader != nil {
		uploader = t.Uploader.Username
	}
	torrentlink := ""
	if t.ID <= config.LastOldTorrentID && len(config.TorrentCacheLink) > 0 {
		if config.IsSukebei() {
			torrentlink = "" // torrent cache doesn't have sukebei torrents
		} else {
			torrentlink = fmt.Sprintf(config.TorrentCacheLink, t.Hash)
		}
	} else if t.ID > config.LastOldTorrentID && len(config.TorrentStorageLink) > 0 {
		torrentlink = fmt.Sprintf(config.TorrentStorageLink, t.Hash)
	}
	res := TorrentJSON{
		ID:           strconv.FormatUint(uint64(t.ID), 10),
		Name:         t.Name,
		Status:       t.Status,
		Hash:         t.Hash,
		Date:         t.Date.Format(time.RFC3339),
		Filesize:     t.Filesize,
		Description:  util.MarkdownToHTML(t.Description),
		Comments:     commentsJSON,
		SubCategory:  strconv.Itoa(t.SubCategory),
		Category:     strconv.Itoa(t.Category),
		Downloads:    t.Downloads,
		UploaderID:   t.UploaderID,
		UploaderName: util.SafeText(uploader),
		OldUploader:  util.SafeText(t.OldUploader),
		WebsiteLink:  util.Safe(t.WebsiteLink),
		Magnet:       template.URL(magnet),
		TorrentLink:  util.Safe(torrentlink),
		Leechers:     t.Leechers,
		Seeders:      t.Seeders,
		Completed:    t.Completed,
		LastScrape:   t.LastScrape,
		FileList:     fileListJSON,
	}

	return res
}

/* Complete the functions when necessary... */

// Map Torrents to TorrentsToJSON without reallocations
func TorrentsToJSON(t []Torrent) []TorrentJSON {
	json := make([]TorrentJSON, len(t))
	for i := range t {
		json[i] = t[i].ToJSON()
	}
	return json
}
