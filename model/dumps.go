package model

import (
	"html/template"
	"time"

	"github.com/NyaaPantsu/nyaa/util"
)

// DatabaseDump model
type DatabaseDump struct {
	Date        time.Time
	Filesize    int64
	Name        string
	TorrentLink string
}

// DatabaseDumpJSON : Json format of DatabaseDump model
type DatabaseDumpJSON struct {
	Date     string `json:"date"`
	Filesize string `json:"filesize"`
	Name     string `json:"name"`
	//Magnet       template.URL  `json:"magnet"`
	TorrentLink template.URL `json:"torrent"`
}

// ToJSON : convert to JSON DatabaseDump model
func (dump *DatabaseDump) ToJSON() DatabaseDumpJSON {
	json := DatabaseDumpJSON{
		Date:        dump.Date.Format(time.RFC3339),
		Filesize:    util.FormatFilesize(dump.Filesize),
		Name:        dump.Name,
		TorrentLink: template.URL(dump.TorrentLink),
	}
	return json
}
