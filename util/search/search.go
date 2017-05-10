package search

import (
	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ewhal/nyaa/cache"
	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util/log"
)

func SearchByQuery(r *http.Request, pagenum int) (search common.SearchParam, tor []model.Torrent, count int, err error) {
	search, tor, count, err = searchByQuery(r, pagenum, true)
	return
}

func SearchByQueryNoCount(r *http.Request, pagenum int) (search common.SearchParam, tor []model.Torrent, err error) {
	search, tor, _, err = searchByQuery(r, pagenum, false)
	return
}

func searchByQuery(r *http.Request, pagenum int, countAll bool) (
	search common.SearchParam, tor []model.Torrent, count int, err error,
) {
	max, err := strconv.ParseUint(r.URL.Query().Get("max"), 10, 32)
	if err != nil {
		max = 50 // default Value maxPerPage
	} else if max > 300 {
		max = 300
	}
	search.Max = uint(max)

	search.Page = pagenum
	search.Query = r.URL.Query().Get("q")

	switch s := r.URL.Query().Get("s"); s {
	case "1":
		search.Status = common.FilterRemakes
	case "2":
		search.Status = common.Trusted
	case "3":
		search.Status = common.APlus
	}

	catString := r.URL.Query().Get("c")
	if s := catString; len(s) > 1 && s != "_" {
		var tmp uint64
		tmp, err = strconv.ParseUint(string(s[0]), 10, 8)
		if err != nil {
			return
		}
		search.Category.Main = uint8(tmp)

		if len(s) == 3 {
			tmp, err = strconv.ParseUint(string(s[2]), 10, 8)
			if err != nil {
				return
			}
			search.Category.Sub = uint8(tmp)
		}
	}

	orderBy := ""

	switch s := r.URL.Query().Get("sort"); s {
	case "1":
		search.Sort = common.Name
		orderBy += "torrent_name"
	case "2":
		search.Sort = common.Date
		orderBy += "date"
	case "3":
		search.Sort = common.Downloads
		orderBy += "downloads"
	case "4":
		search.Sort = common.Size
		orderBy += "filesize"
	default:
		orderBy += "torrent_id"
	}

	orderBy += " "

	switch s := r.URL.Query().Get("order"); s {
	case "true":
		search.Order = true
		orderBy += "asc"
	default:
		orderBy += "desc"
	}

	tor, count, err = cache.Get(search, func() (tor []model.Torrent, count int, err error) {
		parameters := serviceBase.WhereParams{
			Params: make([]interface{}, 0, 64),
		}
		conditions := make([]string, 0, 64)
		if search.Category.Main != 0 {
			conditions = append(conditions, "category = ?")
			parameters.Params = append(parameters.Params, string(catString[0]))
		}
		if search.Category.Sub != 0 {
			conditions = append(conditions, "sub_category = ?")
			parameters.Params = append(parameters.Params, string(catString[2]))
		}
		if search.Status != 0 {
			if search.Status == 3 {
				conditions = append(conditions, "status != ?")
			} else {
				conditions = append(conditions, "status = ?")
			}
			parameters.Params = append(parameters.Params, strconv.Itoa(int(search.Status)+1))
		}

		searchQuerySplit := strings.Fields(search.Query)
		for i, word := range searchQuerySplit {
			firstRune, _ := utf8.DecodeRuneInString(word)
			if len(word) == 1 && unicode.IsPunct(firstRune) {
				// some queries have a single punctuation character
				// which causes a full scan instead of using the index
				// and yields no meaningful results.
				// due to len() == 1 we're just looking at 1-byte/ascii
				// punctuation characters.
				continue
			}

			// SQLite has case-insensitive LIKE, but no ILIKE
			var operator string
			if db.ORM.Dialect().GetName() == "sqlite3" {
				operator = "LIKE ?"
			} else {
				operator = "ILIKE ?"
			}

			// TODO: make this faster ?
			conditions = append(conditions, "torrent_name "+operator)
			parameters.Params = append(parameters.Params, "%"+searchQuerySplit[i]+"%")
		}

		parameters.Conditions = strings.Join(conditions[:], " AND ")
		log.Infof("SQL query is :: %s\n", parameters.Conditions)
		if countAll {
			tor, count, err = torrentService.GetTorrentsOrderBy(&parameters, orderBy, int(search.Max), int(search.Max)*(search.Page-1))
		} else {
			tor, err = torrentService.GetTorrentsOrderByNoCount(&parameters, orderBy, int(search.Max), int(search.Max)*(search.Page-1))
		}

		return
	})

	return
}
