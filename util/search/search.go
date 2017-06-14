package search

import (
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/NyaaPantsu/nyaa/cache"
	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/util/log"
)

var searchOperator string
var useTSQuery bool

// Configure : initialize search
func Configure(conf *config.SearchConfig) (err error) {
	useTSQuery = false
	// Postgres needs ILIKE for case-insensitivity
	if db.ORM.Dialect().GetName() == "postgres" {
		searchOperator = "ILIKE ?"
		//useTSQuery = true
		// !!DISABLED!! because this makes search a lot stricter
		// (only matches at word borders)
	} else {
		searchOperator = "LIKE ?"
	}
	return
}

func stringIsASCII(input string) bool {
	for _, char := range input {
		if char > 127 {
			return false
		}
	}
	return true
}

// SearchByQuery : search torrents according to request without user
func SearchByQuery(r *http.Request, pagenum int) (search common.SearchParam, tor []model.Torrent, count int, err error) {
	search, tor, count, err = searchByQuery(r, pagenum, true, false, false)
	return
}

// SearchByQueryWithUser : search torrents according to request with user
func SearchByQueryWithUser(r *http.Request, pagenum int) (search common.SearchParam, tor []model.Torrent, count int, err error) {
	search, tor, count, err = searchByQuery(r, pagenum, true, true, false)
	return
}

// SearchByQueryNoCount : search torrents according to request without user and count
func SearchByQueryNoCount(r *http.Request, pagenum int) (search common.SearchParam, tor []model.Torrent, err error) {
	search, tor, _, err = searchByQuery(r, pagenum, false, false, false)
	return
}

// SearchByQueryDeleted : search deleted torrents according to request with user and count
func SearchByQueryDeleted(r *http.Request, pagenum int) (search common.SearchParam, tor []model.Torrent, count int, err error) {
	search, tor, count, err = searchByQuery(r, pagenum, true, true, true)
	return
}

// TODO Clean this up
// FIXME Some fields are not used by elasticsearch (pagenum, countAll, deleted, withUser)
// pagenum is extracted from request in .FromRequest()
// elasticsearch always provide a count to how many hits
// deleted is unused because es doesn't index deleted torrents
func searchByQuery(r *http.Request, pagenum int, countAll bool, withUser bool, deleted bool) (
	search common.SearchParam, tor []model.Torrent, count int, err error,
) {
	if db.ElasticSearchClient != nil {
		var torrentParam common.TorrentParam
		torrentParam.FromRequest(r)
		totalHits, torrents, err := torrentParam.Find(db.ElasticSearchClient)
		searchParam := common.SearchParam{
			TorrentID: uint(torrentParam.TorrentID),
			FromID:    uint(torrentParam.FromID),
			FromDate:  torrentParam.FromDate,
			ToDate:    torrentParam.ToDate,
			Order:     torrentParam.Order,
			Status:    torrentParam.Status,
			Sort:      torrentParam.Sort,
			Category:  torrentParam.Category,
			Page:      int(torrentParam.Offset),
			UserID:    uint(torrentParam.UserID),
			Max:       uint(torrentParam.Max),
			NotNull:   torrentParam.NotNull,
			Language:  torrentParam.Language,
			MinSize:   torrentParam.MinSize,
			MaxSize:   torrentParam.MaxSize,
			Query:     torrentParam.NameLike,
		}
		// Convert back to non-json torrents
		return searchParam, torrents, int(totalHits), err
	}
	log.Errorf("Unable to create elasticsearch client: %s", err)
	log.Errorf("Falling back to postgresql query")
	return searchByQueryPostgres(r, pagenum, countAll, withUser, deleted)
}

func searchByQueryPostgres(r *http.Request, pagenum int, countAll bool, withUser bool, deleted bool) (
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
	search.Language = r.URL.Query().Get("lang")
	userID, _ := strconv.Atoi(r.URL.Query().Get("userID"))
	search.UserID = uint(userID)
	fromID, _ := strconv.Atoi(r.URL.Query().Get("fromID"))
	search.FromID = uint(fromID)

	maxage, err := strconv.Atoi(r.URL.Query().Get("maxage"))
	if err != nil {
		search.FromDate = r.URL.Query().Get("fromDate")
		search.ToDate = r.URL.Query().Get("toDate")
	} else {
		search.FromDate = time.Now().AddDate(0, 0, -maxage).Format("2006-01-02")
	}

	search.Status.Parse(r.URL.Query().Get("s"))
	search.Category.Parse(r.URL.Query().Get("c"))
	search.Sort.Parse(r.URL.Query().Get("sort"))
	search.MinSize.Parse(r.URL.Query().Get("minSize"))
	search.MaxSize.Parse(r.URL.Query().Get("maxSize"))

	orderBy := search.Sort.ToDBField()
	if search.Sort == common.Date {
		search.NotNull = search.Sort.ToDBField() + " IS NOT NULL"
	}

	orderBy += " "

	switch s := r.URL.Query().Get("order"); s {
	case "true":
		search.Order = true
		orderBy += "asc"
		if db.ORM.Dialect().GetName() == "postgres" {
			orderBy += " NULLS FIRST"
		}
	default:
		orderBy += "desc"
		if db.ORM.Dialect().GetName() == "postgres" {
			orderBy += " NULLS LAST"
		}
	}

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
	if search.FromID != 0 {
		conditions = append(conditions, "torrent_id > ?")
		parameters.Params = append(parameters.Params, search.FromID)
	}
	if search.FromDate != "" {
		conditions = append(conditions, "date >= ?")
		parameters.Params = append(parameters.Params, search.FromDate)
	}
	if search.ToDate != "" {
		conditions = append(conditions, "date <= ?")
		parameters.Params = append(parameters.Params, search.ToDate)
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
	if search.Language != "" {
		conditions = append(conditions, "language "+searchOperator)
		parameters.Params = append(parameters.Params, "%"+search.Language+"%")
	}
	if search.MinSize > 0 {
		conditions = append(conditions, "filesize >= ?")
		parameters.Params = append(parameters.Params, uint64(search.MinSize))
	}
	if search.MaxSize > 0 {
		conditions = append(conditions, "filesize <= ?")
		parameters.Params = append(parameters.Params, uint64(search.MaxSize))
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

		if useTSQuery && stringIsASCII(word) {
			conditions = append(conditions, "torrent_name @@ plainto_tsquery(?)")
			parameters.Params = append(parameters.Params, word)
		} else {
			// TODO: possible to make this faster?
			conditions = append(conditions, "torrent_name "+searchOperator)
			parameters.Params = append(parameters.Params, "%"+word+"%")
		}
	}

	parameters.Conditions = strings.Join(conditions[:], " AND ")

	log.Infof("SQL query is :: %s\n", parameters.Conditions)

	tor, count, err = cache.Impl.Get(search, func() (tor []model.Torrent, count int, err error) {
		if deleted {
			tor, count, err = torrentService.GetDeletedTorrents(&parameters, orderBy, int(search.Max), int(search.Max)*(search.Page-1))
		} else if countAll && !withUser {
			tor, count, err = torrentService.GetTorrentsOrderBy(&parameters, orderBy, int(search.Max), int(search.Max)*(search.Page-1))
		} else if withUser {
			tor, count, err = torrentService.GetTorrentsWithUserOrderBy(&parameters, orderBy, int(search.Max), int(search.Max)*(search.Page-1))
		} else {
			tor, err = torrentService.GetTorrentsOrderByNoCount(&parameters, orderBy, int(search.Max), int(search.Max)*(search.Page-1))
		}
		return
	})
	return
}
