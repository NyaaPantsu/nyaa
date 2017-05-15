package postgres

import (
	"database/sql"
	"github.com/ewhal/nyaa/common"
)

func (db *Database) RecordScrapes(scrape []common.ScrapeResult) (err error) {
	if len(scrape) > 0 {
		var tx *sql.Tx
		tx, err = db.conn.Begin()
		if err == nil {
			st := tx.Stmt(db.getPrepared(queryInsertScrape))
			for idx := range scrape {
				_, err = st.Exec(scrape[idx].Seeders, scrape[idx].Leechers, scrape[idx].Completed, scrape[idx].Date, scrape[idx].TorrentID)
				if err != nil {
					break
				}
			}
		}
		if err == nil {
			err = tx.Commit()
		} else {
			tx.Rollback()
		}
	}
	return
}
