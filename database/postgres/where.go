package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/log"
)

// build sql query from SearchParam for torrent search
func searchParamToTorrentQuery(param *common.TorrentParam) (q sqlQuery) {
	log.Debugf("build query: %v", param)
	counter := 1
	q.query = fmt.Sprintf("SELECT %s FROM %s ", torrentSelectColumnsFull, tableTorrents)
	if param.Category.IsSet() {
		q.query += fmt.Sprintf("WHERE category = $%d AND sub_category = $%d ", counter, counter+1)
		q.params = append(q.params, param.Category.Main, param.Category.Sub)
		counter += 2
	}
	s := "WHERE"
	if counter > 1 {
		s = "AND "
	}
	if param.Status > 0 {
		q.query += fmt.Sprintf("%s status >= $%d ", s, counter)
		q.params = append(q.params, param.Status)
		counter++
		s = "AND"
	}
	if param.UserID > 0 {
		q.query += fmt.Sprintf("%s uploader = $%d ", s, counter)
		q.params = append(q.params, param.UserID)
		counter++
		s = "AND"
	}

	for idx := range param.NotNull {
		k := strings.ToLower(strings.TrimSpace(param.NotNull[idx]))
		switch k {
		case "date":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		case "downloads":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		case "filesize":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		case "website_link":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		case "deleted_at":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		case "seeders":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		case "leechers":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		case "completed":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		case "last_scrape":
			q.query += fmt.Sprintf("%s %s IS NOT NULL ", s, k)
			s = "AND"
		}
	}

	for idx := range param.Null {
		k := strings.ToLower(strings.TrimSpace(param.Null[idx]))
		switch k {
		case "date":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		case "downloads":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		case "filesize":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		case "website_link":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		case "deleted_at":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		case "seeders":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		case "leechers":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		case "completed":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		case "last_scrape":
			q.query += fmt.Sprintf("%s %s IS NULL ", s, k)
			s = "AND"
		}
	}
	for idx := range param.NameLike {
		name := param.NameLike[idx]
		if len(name) > 0 {
			q.query += fmt.Sprintf("%s torrent_name ILIKE $%d ", s, counter)
			q.params = append(q.params, "%"+name+"%")
			counter++
			s = "AND"
		}
	}
	sort := "torrent_id"
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
	err = q.Query(db.conn, func(rows *sql.Rows) (e error) {
		if param.Max == 0 {
			torrents = make([]model.Torrent, 0, 128)
		} else {
			torrents = make([]model.Torrent, 0, param.Max)
		}

		for rows.Next() {
			var t model.Torrent
			e = scanTorrentColumnsFull(rows, &t)
			if e != nil {
				break
			}
			torrents = append(torrents, t)
		}
		return
	})
	return
}
