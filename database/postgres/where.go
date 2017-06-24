package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
)

// build sql query from SearchParam for torrent search
func searchParamToTorrentQuery(param *common.TorrentParam) (q sqlQuery) {
	counter := 1
	q.query = fmt.Sprintf("SELECT %s FROM %s ", torrentSelectColumnsFull, tableTorrents)
	if len(param.Category) > 0 {
		conditionsOr := make([]string, len(param.Category))
		for key, val := range param.Category {
			conditionsOr[key] = fmt.Sprintf("(category = $%d AND sub_category = $%d)", counter, counter+1)
			q.params = append(q.params, val.Main, val.Sub)
			counter += 2
		}
		q.query += "WHERE " + strings.Join(conditionsOr, " OR ")
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
		default:
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
		default:
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
	case common.Date:
		sort = "date"
	case common.Downloads:
		sort = "downloads"
	case common.Size:
		sort = "filesize"
	case common.Seeders:
		sort = "seeders"
	case common.Leechers:
		sort = "leechers"
	case common.Completed:
		sort = "completed"
	case common.ID:
	default:
		sort = config.Conf.Torrents.Order
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
