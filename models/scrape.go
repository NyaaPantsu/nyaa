package models

import (
	"net/http"
	"errors"
	"time"
	"fmt"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/cache"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/fatih/structs"
)

// Scrape model
type Scrape struct {
	TorrentID  uint      `gorm:"column:torrent_id;primary_key"`
	Seeders    uint32    `gorm:"column:seeders"`
	Leechers   uint32    `gorm:"column:leechers"`
	Completed  uint32    `gorm:"column:completed"`
	LastScrape time.Time `gorm:"column:last_scrape"`
}

// TableName : return the table name of the scrape table
func (t Scrape) TableName() string {
	return config.Get().Models.ScrapeTableName
}

// Update : Update scrape data based on Scrape model
func (s *Scrape) Update(unscope bool) (int, error) {
	db := ORM
	if unscope {
		db = ORM.Unscoped()
	}

	if db.Model(s).UpdateColumn(s.toMap()).Error != nil {
		return http.StatusInternalServerError, errors.New("Scrape data was not updated")
	}

	// We only flush cache after update
	cache.C.Delete(s.Identifier())
	cache.C.Flush()

	return http.StatusOK, nil
}


// toMap : convert the model to a map of interface
func (s *Scrape) toMap() map[string]interface{} {
	return structs.Map(s)
}

// Identifier : Return the identifier of a torrent
func (s *Scrape) Identifier() string {
	return fmt.Sprintf("torrent_%d", s.TorrentID)
}

//Create a Scrape entry in the DB
func (s *Scrape) Create(torrentid uint, seeders uint32, leechers uint32, completed uint32, lastscrape time.Time) (*Scrape) {
	ScrapeData := Scrape{
		TorrentID:      torrentid,
		Seeders:    	seeders,
		Leechers: 	leechers,
		Completed:      completed,
		LastScrape:     lastscrape
	}

	err := ORM.Create(&ScrapeData).Error
	log.Infof("Scrape data ID %d created!\n", ScrapeData.TorrentID)
	if err != nil {
		log.CheckErrorWithMessage(err, "ERROR_SCRAPE_CREATE: Cannot create a scrape data")
	}

	ScrapeData.Update(false)

	return &ScrapeData
}
