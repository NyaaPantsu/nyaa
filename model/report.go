package model

import (
	"time"

	"github.com/NyaaPantsu/nyaa/config"
)

// TorrentReport model
// User can be null (anonymous reports)
// FIXME  can't preload field Torrents for model.TorrentReport
type TorrentReport struct {
	ID          uint   `gorm:"column:torrent_report_id;primary_key"`
	Description string `gorm:"column:type"`
	TorrentID   uint   `gorm:"column:torrent_id"`
	UserID      uint   `gorm:"column:user_id"`

	CreatedAt time.Time `gorm:"column:created_at"`

	Torrent *Torrent `gorm:"AssociationForeignKey:TorrentID;ForeignKey:torrent_id"`
	User    *User    `gorm:"AssociationForeignKey:UserID;ForeignKey:user_id"`
}

// TableName : Return the name of torrent report table
func (report TorrentReport) TableName() string {
	return config.ReportsTableName
}

// TorrentReportJSON : Json struct of torrent report model
type TorrentReportJSON struct {
	ID          uint        `json:"id"`
	Description string      `json:"description"`
	Torrent     TorrentJSON `json:"torrent"`
	User        UserJSON    `json:"user"`
}

/* Model Conversion to Json */

func getReportDescription(d string) string {
	if d == "illegal" {
		return "Illegal content"
	} else if d == "spam" {
		return "Spam / Garbage"
	} else if d == "wrongcat" {
		return "Wrong category"
	} else if d == "dup" {
		return "Duplicate / Deprecated"
	}
	return "???"
}

// ToJSON : conversion to json of a torrent report
func (report *TorrentReport) ToJSON() TorrentReportJSON {
	t := TorrentJSON{}
	if report.Torrent != nil { // FIXME: report.Torrent should never be nil
		t = report.Torrent.ToJSON()
	}
	u := UserJSON{}
	if report.User != nil {
		u = report.User.ToJSON()
	}
	json := TorrentReportJSON{report.ID, getReportDescription(report.Description), t, u}
	return json
}

// TorrentReportsToJSON : Conversion of multiple reports to json
func TorrentReportsToJSON(reports []TorrentReport) []TorrentReportJSON {
	json := make([]TorrentReportJSON, len(reports))
	for i := range reports {
		json[i] = reports[i].ToJSON()
	}
	return json
}
