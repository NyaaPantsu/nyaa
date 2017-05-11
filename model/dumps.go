package model

import (
	"time"
	"html/template"

	"github.com/ewhal/nyaa/util"
)

type DatabaseDump struct {
	Date        time.Time
	Filesize    int64
	Name        string
	TorrentLink string
}

type DatabaseDumpJSON struct {
	Date        string `json:"date"`
	Filesize    string `json:"filesize"`
	Name        string `json:"name"`
	//Magnet       template.URL  `json:"magnet"`
	TorrentLink  template.URL  `json:"torrent"`
}

func (dump *DatabaseDump) ToJSON() DatabaseDumpJSON {
	json := DatabaseDumpJSON{
		Date:	     dump.Date.Format(time.RFC3339),
		Filesize:    util.FormatFilesize2(dump.Filesize),
		Name:        dump.Name,
		TorrentLink: template.URL(dump.TorrentLink),
	}
	return json
}
