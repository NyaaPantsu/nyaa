package search

import (
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/gin-gonic/gin"
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
func SearchByQuery(c *gin.Context, pagenum int) (search common.SearchParam, tor []model.Torrent, count int, err error) {
	search, tor, count, err = searchByQuery(c, pagenum, true, false, false, false)
	return
}

// SearchByQueryWithUser : search torrents according to request with user
func SearchByQueryWithUser(c *gin.Context, pagenum int) (search common.SearchParam, tor []model.Torrent, count int, err error) {
	search, tor, count, err = searchByQuery(c, pagenum, true, true, false, false)
	return
}

// SearchByQueryNoCount : search torrents according to request without user and count
func SearchByQueryNoCount(c *gin.Context, pagenum int) (search common.SearchParam, tor []model.Torrent, err error) {
	search, tor, _, err = searchByQuery(c, pagenum, false, false, false, false)
	return
}

// SearchByQueryDeleted : search deleted torrents according to request with user and count
func SearchByQueryDeleted(c *gin.Context, pagenum int) (search common.SearchParam, tor []model.Torrent, count int, err error) {
	search, tor, count, err = searchByQuery(c, pagenum, true, true, true, false)
	return
}

// SearchByQueryNoHidden : search torrents and filter those hidden
func SearchByQueryNoHidden(c *gin.Context, pagenum int) (search common.SearchParam, tor []model.Torrent, count int, err error) {
	search, tor, count, err = searchByQuery(c, pagenum, true, false, false, true)
	return
}

// TODO Clean this up
// FIXME Some fields are not used by elasticsearch (pagenum, countAll, deleted, withUser)
// pagenum is extracted from request in .FromRequest()
// elasticsearch always provide a count to how many hits
// deleted is unused because es doesn't index deleted torrents
func searchByQuery(c *gin.Context, pagenum int, countAll bool, withUser bool, deleted bool, hidden bool) (
	search common.SearchParam, tor []model.Torrent, count int, err error,
) {
	if db.ElasticSearchClient != nil {
		var torrentParam common.TorrentParam
		torrentParam.FromRequest(c)
		torrentParam.Offset = uint32(pagenum)
		torrentParam.Hidden = hidden
		totalHits, torrents, err := torrentParam.Find(db.ElasticSearchClient)
		searchParam := common.SearchParam{
			TorrentID: uint(torrentParam.TorrentID),
			FromID:    uint(torrentParam.FromID),
			FromDate:  torrentParam.FromDate,
			ToDate:    torrentParam.ToDate,
			Order:     torrentParam.Order,
			Hidden:    torrentParam.Hidden,
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
	return searchByQueryPostgres(c, pagenum, countAll, withUser, deleted, hidden)
}

func searchByQueryPostgres(c *gin.Context, pagenum int, countAll bool, withUser bool, deleted bool, hidden bool) (
	search common.SearchParam, tor []model.Torrent, count int, err error,
) {
	max, err := strconv.ParseUint(c.DefaultQuery("limit", "50"), 10, 32)
	if err != nil {
		max = 50 // default Value maxPerPage
	} else if max > 300 {
		max = 300
	}
	search.Max = uint(max)

	search.Page = pagenum
	search.Hidden = hidden
	search.Query = c.Query("q")
	search.Language = c.Query("lang")
	userID, _ := strconv.Atoi(c.Query("userID"))
	search.UserID = uint(userID)
	fromID, _ := strconv.Atoi(c.Query("fromID"))
	search.FromID = uint(fromID)

	maxage, err := strconv.Atoi(c.Query("maxage"))
	if err != nil {
		if c.Query("toDate") != "" {
			search.FromDate.Parse(c.Query("toDate"), c.Query("dateType"))
			search.ToDate.Parse(c.Query("fromDate"), c.Query("dateType"))
		} else {
			search.FromDate.Parse(c.Query("fromDate"), c.Query("dateType"))
		}
	} else {
		search.FromDate = common.DateFilter(time.Now().AddDate(0, 0, -maxage).Format("2006-01-02"))
	}
	search.Category = common.ParseCategories(c.Query("c"))
	search.Status.Parse(c.Query("s"))
	search.Sort.Parse(c.Query("sort"))
	search.MinSize.Parse(c.Query("minSize"), c.Query("sizeType"))
	search.MaxSize.Parse(c.Query("maxSize"), c.Query("sizeType"))

	orderBy := search.Sort.ToDBField()
	if search.Sort == common.Date {
		search.NotNull = search.Sort.ToDBField() + " IS NOT NULL"
	}

	orderBy += " "

	switch s := c.Query("order"); s {
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
	if len(search.Category) > 0 {
		conditionsOr := make([]string, len(search.Category))
		for key, val := range search.Category {
			conditionsOr[key] = "(category = ? AND sub_category = ?)"
			parameters.Params = append(parameters.Params, val.Main)
			parameters.Params = append(parameters.Params, val.Sub)
		}
		conditions = append(conditions, strings.Join(conditionsOr, " OR "))
	}

	if search.UserID != 0 {
		conditions = append(conditions, "uploader = ?")
		parameters.Params = append(parameters.Params, search.UserID)
	}
	if search.Hidden {
		conditions = append(conditions, "hidden = ?")
		parameters.Params = append(parameters.Params, false)
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
}
