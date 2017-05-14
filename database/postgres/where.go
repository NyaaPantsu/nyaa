package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/model"
)

// build sql query from SearchParam for torrent search
func searchParamToTorrentQuery(param *common.TorrentParam) (q sqlQuery) {
	counter := 1
	q.query = fmt.Sprintf("SELECT %s FROM %s ", torrentSelectColumnsFull, tableTorrents)
	if param.Category.IsSet() {
		q.query += fmt.Sprintf("WHERE category = $%d AND sub_category = $%d ", counter, counter+1)
		q.params = append(q.params, param.Category.Main, param.Category.Sub)
		counter += 2
	}
	if counter > 1 {
		q.query += "AND "
	} else {
		q.query += "WHERE "
	}

	q.query += fmt.Sprintf("status >= $%d ", counter)
	q.params = append(q.params, param.Status)
	counter++
	if param.UserID > 0 {
		q.query += fmt.Sprintf("AND uploader = $%d ", counter)
		q.params = append(q.params, param.UserID)
		counter++
	}

	notnulls := strings.Split(param.NotNull, ",")
	for idx := range notnulls {
		k := strings.ToLower(strings.TrimSpace(notnulls[idx]))
		switch k {
		case "date":
		case "downloads":
		case "filesize":
		case "website_link":
		case "deleted_at":
		case "seeders":
		case "leechers":
		case "completed":
		case "last_scrape":
			q.query += fmt.Sprintf("AND %s IS NOT NULL ", k)
			break
		default:
			break
		}
	}

	nulls := strings.Split(param.Null, ",")
	for idx := range nulls {
		k := strings.ToLower(strings.TrimSpace(nulls[idx]))
		switch k {
		case "date":
		case "downloads":
		case "filesize":
		case "website_link":
		case "deleted_at":
		case "seeders":
		case "leechers":
		case "completed":
		case "last_scrape":
			q.query += fmt.Sprintf("AND %s IS NULL ", k)
			break
		default:
			break
		}
	}

	nameLikes := strings.Split(param.NameLike, ",")
	for idx := range nameLikes {
		q.query += fmt.Sprintf("AND torrent_name ILIKE $%d", counter)
		q.params = append(q.params, strings.TrimSpace(nameLikes[idx]))
		counter++
	}

	var sort string
	switch param.Sort {
	case common.Name:
		sort = "torrent_name"
		break
	case common.Date:
		sort = "date"
		break
	case common.Downloads:
		sort = "downloads"
		break
	case common.Size:
		sort = "filesize"
		break
	case common.Seeders:
		sort = "seeders"
		break
	case common.Leechers:
		sort = "leechers"
		break
	case common.Completed:
		sort = "completed"
		break
	case common.ID:
	default:
		sort = "torrent_id"
	}

	q.query += fmt.Sprintf("ORDER BY %s ", sort)
	if param.Order {
		q.query += "ASC "
	} else {
		q.query += "DESC "
	}

	if param.Max > 0 {
		q.query += fmt.Sprintf("LIMIT $%d ", counter)
		q.params = append(q.params, param.Max)
		counter++
	}
	if param.Offset > 0 {
		q.query += fmt.Sprintf("OFFSET $%d ", counter)
		q.params = append(q.params, param.Offset)
		counter++
	}
	return
}

func (db *Database) GetTorrentsWhere(param *common.TorrentParam) (torrents []model.Torrent, err error) {
	if param.TorrentID > 0 {
		var torrent model.Torrent
		var has bool
		torrent, has, err = db.GetTorrentByID(param.TorrentID)
		if has {
			torrents = append(torrents, torrent)
		}
		return
	}

	if param.All {
		torrents, err = db.GetAllTorrents(param.Offset, param.Max)
		return
	}

	q := searchParamToTorrentQuery(param)
	err = q.Query(db.conn, func(rows *sql.Rows) error {
		if param.Max == 0 {
			torrents = make([]model.Torrent, 0, 128)
		} else {
			torrents = make([]model.Torrent, 0, param.Max)
		}

		for rows.Next() {
			var t model.Torrent
			scanTorrentColumnsFull(rows, &t)
			torrents = append(torrents, t)
			// grow as needed
			if len(torrents) >= cap(torrents) {
				newtorrents := make([]model.Torrent, cap(torrents), cap(torrents)*3/2) // XXX: adjust as needed
				copy(newtorrents, torrents)
				torrents = newtorrents
			}
		}
		return nil
	})
	return
}
