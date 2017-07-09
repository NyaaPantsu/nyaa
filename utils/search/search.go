package search

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/utils/cache"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/gin-gonic/gin"
)

var searchOperator string
var useTSQuery bool

// Configure : initialize search
func Configure(conf *config.SearchConfig) (err error) {
	useTSQuery = false
	// Postgres needs ILIKE for case-insensitivity
	if models.ORM.Dialect().GetName() == "postgres" {
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

// ByQuery : search torrents according to request without user
func ByQuery(c *gin.Context, pagenum int) (search structs.TorrentParam, tor []models.Torrent, count int, err error) {
	search, tor, count, err = byQuery(c, pagenum, true, false, false, false)
	return
}

// ByQueryWithUser : search torrents according to request with user
func ByQueryWithUser(c *gin.Context, pagenum int) (search structs.TorrentParam, tor []models.Torrent, count int, err error) {
	search, tor, count, err = byQuery(c, pagenum, true, true, false, false)
	return
}

// ByQueryNoCount : search torrents according to request without user and count
func ByQueryNoCount(c *gin.Context, pagenum int) (search structs.TorrentParam, tor []models.Torrent, err error) {
	search, tor, _, err = byQuery(c, pagenum, false, false, false, false)
	return
}

// ByQueryDeleted : search deleted torrents according to request with user and count
func ByQueryDeleted(c *gin.Context, pagenum int) (search structs.TorrentParam, tor []models.Torrent, count int, err error) {
	search, tor, count, err = byQuery(c, pagenum, true, true, true, false)
	return
}

// ByQueryNoHidden : search torrents and filter those hidden
func ByQueryNoHidden(c *gin.Context, pagenum int) (search structs.TorrentParam, tor []models.Torrent, count int, err error) {
	search, tor, count, err = byQuery(c, pagenum, true, false, false, true)
	return
}

// TODO Clean this up
// FIXME Some fields are not used by elasticsearch (pagenum, countAll, deleted, withUser)
// pagenum is extracted from request in .FromRequest()
// elasticsearch always provide a count to how many hits
// deleted is unused because es doesn't index deleted torrents
func byQuery(c *gin.Context, pagenum int, countAll bool, withUser bool, deleted bool, hidden bool) (structs.TorrentParam, []models.Torrent, int, error) {
	var err error
	if models.ElasticSearchClient != nil {
		var torrentParam structs.TorrentParam
		torrentParam.FromRequest(c)
		torrentParam.Offset = uint32(pagenum)
		torrentParam.Hidden = hidden
		torrentParam.Full = withUser
		if found, ok := cache.C.Get(torrentParam.Identifier()); ok {
			torrentCache := found.(*structs.TorrentCache)
			return torrentParam, torrentCache.Torrents, torrentCache.Count, nil
		}
		totalHits, tor, err := torrentParam.Find(models.ElasticSearchClient)
		cache.C.Set(torrentParam.Identifier(), &structs.TorrentCache{tor, int(totalHits)}, 5*time.Minute)
		// Convert back to non-json torrents
		return torrentParam, tor, int(totalHits), err
	}
	log.Errorf("Unable to create elasticsearch client: %s", err)
	log.Errorf("Falling back to postgresql query")
	return byQueryPostgres(c, pagenum, countAll, withUser, deleted, hidden)
}

func byQueryPostgres(c *gin.Context, pagenum int, countAll bool, withUser bool, deleted bool, hidden bool) (
	search structs.TorrentParam, tor []models.Torrent, count int, err error,
) {
	search.FromRequest(c)
	search.Offset = uint32(pagenum)
	search.Hidden = hidden
	search.Full = withUser

	orderBy := search.Sort.ToDBField()
	if search.Sort == structs.Date {
		search.NotNull = search.Sort.ToDBField() + " IS NOT NULL"
	}
	if found, ok := cache.C.Get(search.Identifier()); ok {
		torrentCache := found.(*structs.TorrentCache)
		tor = torrentCache.Torrents
		count = torrentCache.Count
		return
	}

	orderBy += " "

	switch search.Order {
	case true:
		orderBy += "asc"
		if models.ORM.Dialect().GetName() == "postgres" {
			orderBy += " NULLS FIRST"
		}
	case false:
		orderBy += "desc"
		if models.ORM.Dialect().GetName() == "postgres" {
			orderBy += " NULLS LAST"
		}
	}

	parameters := structs.WhereParams{
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
	if len(search.Languages) > 0 {
		langs := ""
		for key, val := range search.Languages {
			langs += val.Code
			if key+1 < len(search.Languages) {
				langs += ","
			}
		}
		conditions = append(conditions, "language "+searchOperator)
		parameters.Params = append(parameters.Params, "%"+langs+"%")
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
		if search.Status == structs.FilterRemakes {
			conditions = append(conditions, "status <> ?")
		} else {
			conditions = append(conditions, "status >= ?")
		}
		parameters.Params = append(parameters.Params, strconv.Itoa(int(search.Status)+1))
	}
	if len(search.NotNull) > 0 {
		conditions = append(conditions, search.NotNull)
	}
	if search.MinSize > 0 {
		conditions = append(conditions, "filesize >= ?")
		parameters.Params = append(parameters.Params, uint64(search.MinSize))
	}
	if search.MaxSize > 0 {
		conditions = append(conditions, "filesize <= ?")
		parameters.Params = append(parameters.Params, uint64(search.MaxSize))
	}

	querySplit := strings.Fields(search.NameLike)
	for _, word := range querySplit {
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
		tor, count, err = torrents.FindDeleted(&parameters, orderBy, int(search.Max), int(search.Max*(search.Offset-1)))
	} else if countAll && !withUser {
		tor, count, err = torrents.FindOrderBy(&parameters, orderBy, int(search.Max), int(search.Max*(search.Offset-1)))
	} else if withUser {
		tor, count, err = torrents.FindWithUserOrderBy(&parameters, orderBy, int(search.Max), int(search.Max*(search.Offset-1)))
	} else {
		tor, err = torrents.FindOrderByNoCount(&parameters, orderBy, int(search.Max), int(search.Max*(search.Offset-1)))
	}
	cache.C.Set(search.Identifier(), &structs.TorrentCache{tor, count}, 5*time.Minute)
	return
}
