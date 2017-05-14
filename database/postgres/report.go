package postgres

import (
	"database/sql"
	"fmt"

	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/model"
)

func (db *Database) InsertTorrentReport(report *model.TorrentReport) (err error) {
	_, err = db.getPrepared(queryInsertTorrentReport).Exec(report.Type, report.TorrentID, report.UserID, report.CreatedAt)
	return
}

func reportParamToQuery(param *common.ReportParam) (q sqlQuery) {
	q.query += fmt.Sprintf("SELECT %s FROM %s WHERE created_at IS NOT NULL ", torrentReportSelectColumnsFull, tableTorrentReports)

	counter := 1

	if !param.AllTime {
		q.query += fmt.Sprintf("AND created_at < $%d AND created_at > $%d", counter, counter+1)
		q.params = append(q.params, param.Before, param.After)
		counter += 2
	}

	if param.Limit > 0 {
		q.query += fmt.Sprintf("LIMIT $%d ", counter)
		q.params = append(q.params, param.Limit)
		counter++
	}
	if param.Offset > 0 {
		q.query += fmt.Sprintf("OFFSET $%d ", counter)
		q.params = append(q.params, param.Offset)
		counter++
	}
	return
}

func (db *Database) GetTorrentReportsWhere(param *common.ReportParam) (reports []model.TorrentReport, err error) {
	q := reportParamToQuery(param)
	err = q.Query(db.conn, func(rows *sql.Rows) error {
		for rows.Next() {
			var r model.TorrentReport
			scanTorrentReportColumnsFull(rows, &r)
			reports = append(reports, r)
		}
		return nil
	})
	return
}

func (db *Database) DeleteTorrentReportByID(id uint32) (err error) {
	_, err = db.getPrepared(queryDeleteTorrentReportByID).Exec(id)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

func (db *Database) DeleteTorrentReportsWhere(param *common.ReportParam) (deleted uint32, err error) {

	return
}
