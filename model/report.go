package model

import (
	"time"
)

// TODO Add field to specify kind of reports
// TODO Add CreatedAt field
// INFO User can be null (anonymous reports)
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

type TorrentReportJson struct {
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

func (report *TorrentReport) ToJson() TorrentReportJson {
	// FIXME: report.Torrent and report.User should never be nil
	var t TorrentJSON = TorrentJSON{}
	if report.Torrent != nil {
		t = report.Torrent.ToJSON()
	}
	var u UserJSON = UserJSON{}
	if report.User != nil {
		u = report.User.ToJSON()
	}
	json := TorrentReportJson{report.ID, getReportDescription(report.Description), t, u}
	return json
}

func TorrentReportsToJSON(reports []TorrentReport) []TorrentReportJson {
	json := make([]TorrentReportJson, len(reports))
	for i := range reports {
		json[i] = reports[i].ToJson()
	}
	return json
}
