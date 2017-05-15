package search

import (
	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/ewhal/nyaa/cache"
	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/database"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service"
	"github.com/ewhal/nyaa/util/log"
)

var searchOperator string
var useTSQuery bool

func Configure(conf *config.SearchConfig) (err error) {
	return
}

func stringIsAscii(input string) bool {
	for _, char := range input {
		if char > 127 {
			return false
		}
	}
	return true
}

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
	var param common.TorrentParam
	max, err := strconv.ParseUint(r.URL.Query().Get("max"), 10, 32)
	if err != nil {
		max = 50 // default Value maxPerPage
	} else if max > 300 {
		max = 300
	}
	search.Max = uint(max)
	param.Max = uint32(max)

	search.Page = pagenum
	param.Offset = param.Max * uint32(pagenum-1)
	search.Query = r.URL.Query().Get("q")
	userID, _ := strconv.Atoi(r.URL.Query().Get("userID"))
	search.UserID = uint(userID)
	param.UserID = uint32(userID)

	switch s := r.URL.Query().Get("s"); s {
	case "1":
		search.Status = common.FilterRemakes
	case "2":
		search.Status = common.Trusted
	case "3":
		search.Status = common.APlus
	}

	param.Status = search.Status

	catString := r.URL.Query().Get("c")
	if s := catString; len(s) > 1 && s != "_" {
		var tmp uint64
		tmp, err = strconv.ParseUint(string(s[0]), 10, 8)
		if err != nil {
			return
		}
		search.Category.Main = uint8(tmp)

		if len(s) > 2 && len(s) < 5 {
			tmp, err = strconv.ParseUint(s[2:], 10, 8)
			if err != nil {
				return
			}
			search.Category.Sub = uint8(tmp)
		}
	}

	param.Category = search.Category

	orderBy := ""

	switch s := r.URL.Query().Get("sort"); s {
	case "1":
		search.Sort = common.Name
		orderBy += "torrent_name"
		break
	case "2":
		search.Sort = common.Date
		orderBy += "date"
		search.NotNull = "date IS NOT NULL"
		param.NotNull = append(param.NotNull, "date")
		break
	case "3":
		search.Sort = common.Downloads
		orderBy += "downloads"
		break
	case "4":
		search.Sort = common.Size
		orderBy += "filesize"
		// avoid sorting completely breaking on postgres
		search.NotNull = "filesize IS NOT NULL"
		param.NotNull = append(param.NotNull, "filesize")
		break
	case "5":
		search.Sort = common.Seeders
		orderBy += "seeders"
		search.NotNull = "seeders IS NOT NULL"
		param.NotNull = append(param.NotNull, "seeders")
		break
	case "6":
		search.Sort = common.Leechers
		orderBy += "leechers"
		search.NotNull = "leechers IS NOT NULL"
		param.NotNull = append(param.NotNull, "leechers")
		break
	case "7":
		search.Sort = common.Completed
		orderBy += "completed"
		search.NotNull = "completed IS NOT NULL"
		param.NotNull = append(param.NotNull, "completed")
		break
	default:
		search.Sort = common.ID
		orderBy += "torrent_id"
	}

	param.Sort = search.Sort

	orderBy += " "

	switch s := r.URL.Query().Get("order"); s {
	case "true":
		search.Order = true
		orderBy += "asc"
	default:
		orderBy += "desc"
	}

	param.Order = search.Order

	parameters := serviceBase.WhereParams{
		Params: make([]interface{}, 0, 64),
	}
	conditions := make([]string, 0, 64)

	if search.Category.Main != 0 {
		conditions = append(conditions, "category = ?")
		parameters.Params = append(parameters.Params, search.Category.Main)
	}
	if search.UserID != 0 {
		conditions = append(conditions, "uploader = ?")
		parameters.Params = append(parameters.Params, search.UserID)
	}
	if search.Category.Sub != 0 {
		conditions = append(conditions, "sub_category = ?")
		parameters.Params = append(parameters.Params, search.Category.Sub)
	}
	if search.Status != 0 {
		if search.Status == common.FilterRemakes {
			conditions = append(conditions, "status <> ?")
		} else {
			conditions = append(conditions, "status >= ?")
		}
		parameters.Params = append(parameters.Params, strconv.Itoa(int(search.Status)+1))
	}
	if len(search.NotNull) > 0 {
		conditions = append(conditions, search.NotNull)
	}

	searchQuerySplit := strings.Fields(search.Query)
	for _, word := range searchQuerySplit {
		firstRune, _ := utf8.DecodeRuneInString(word)
		if len(word) == 1 && unicode.IsPunct(firstRune) {
			// some queries have a single punctuation character
			// which causes a full scan instead of using the index
			// and yields no meaningful results.
			// due to len() == 1 we're just looking at 1-byte/ascii
			// punctuation characters.
			continue
		}

		if useTSQuery && stringIsAscii(word) {
			conditions = append(conditions, "torrent_name @@ plainto_tsquery(?)")
			parameters.Params = append(parameters.Params, word)
		} else {
			// TODO: possible to make this faster?
			conditions = append(conditions, "torrent_name "+searchOperator)
			parameters.Params = append(parameters.Params, "%"+word+"%")
			param.NameLike = append(param.NameLike, word)
		}
	}

	parameters.Conditions = strings.Join(conditions[:], " AND ")

	tor, count, err = cache.Impl.Get(search, func() (tor []model.Torrent, count int, err error) {
		tor, err = database.Impl.GetTorrentsWhere(&param)
		count = len(tor)
		count += int(param.Offset)
		count += int(param.Max) * 20
		log.Debugf("%d", count)
		return
	})
	return
}
