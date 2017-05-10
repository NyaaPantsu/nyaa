package model

// TODO Add field to specify kind of reports
// TODO Add CreatedAt field
// INFO User can be null (anonymous reports)
// FIXME  can't preload field Torrents for model.TorrentReport
type TorrentReport struct {
	ID          uint   `gorm:"column:torrent_report_id;primary_key"`
	Description string `gorm:"column:type"`
	TorrentID   uint   `gorm:"column:torrent_id"`
	UserID      uint   `gorm:"column:user_id"`

	Torrent Torrent `gorm:"AssociationForeignKey:TorrentID;ForeignKey:torrent_id"`
	User    User    `gorm:"AssociationForeignKey:UserID;ForeignKey:ID"`
}

type TorrentReportJson struct {
	ID          uint        `json:"id"`
	Description string      `json:"description"`
	Torrent     TorrentJSON `json:"torrent"`
	User        UserJSON    `json:"user"`
}

/* Model Conversion to Json */

func (report *TorrentReport) ToJson() TorrentReportJson {
	json := TorrentReportJson{report.ID, report.Description, report.Torrent.ToJSON(), report.User.ToJSON()}
	return json
}

func TorrentReportsToJSON(reports []TorrentReport) []TorrentReportJson {
	json := make([]TorrentReportJson, len(reports))
	for i := range reports {
		json[i] = reports[i].ToJson()
	}
	return json
}
