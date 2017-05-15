package postgres

import (
	"github.com/ewhal/nyaa/model"

	"database/sql"
)

func (db *Database) GetAllTorrents(offset, limit uint32) (torrents []model.Torrent, err error) {
	err = db.queryWithPrepared(queryGetAllTorrents, func(rows *sql.Rows) error {
		torrents = make([]model.Torrent, 0, limit)
		var idx uint64
		for rows.Next() {
			rows.Scan(torrents[idx])
		}
		return nil
	}, offset, limit)
	return
}

func (db *Database) GetTorrentByID(id uint32) (torrent model.Torrent, has bool, err error) {
	err = db.queryWithPrepared(queryGetTorrentByID, func(rows *sql.Rows) error {
		rows.Next()
		scanTorrentColumnsFull(rows, &torrent)
		has = true
		return nil
	}, id)
	return
}

func (db *Database) UpsertTorrent(t *model.Torrent) (err error) {
	return
}
