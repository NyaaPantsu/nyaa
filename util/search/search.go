package search

import (
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util/log"
	"html"
	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type SearchParam struct {
	Category string
	Order    string
	Query    string
	Max      int
	Status   string
	Sort     string
}

func SearchByQuery(r *http.Request, pagenum int) (SearchParam, []model.Torrents, int) {
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}

	if maxPerPage > 300 {
		maxPerPage = 300
	}
	search_param := SearchParam{}
	search_param.Max = maxPerPage
	search_param.Query = r.URL.Query().Get("q")
	search_param.Category = r.URL.Query().Get("c")
	search_param.Status = r.URL.Query().Get("s")
	search_param.Sort = r.URL.Query().Get("sort")
	search_param.Order = r.URL.Query().Get("order")
	userId := r.URL.Query().Get("userId")

	catsSplit := strings.Split(search_param.Category, "_")
	// need this to prevent out of index panics
	var searchCatId, searchSubCatId string
	if len(catsSplit) == 2 {

		searchCatId = html.EscapeString(catsSplit[0])
		searchSubCatId = html.EscapeString(catsSplit[1])
	}

	switch search_param.Sort {
	case "torrent_name":
		search_param.Sort = "torrent_name"
		break
	case "date":
		search_param.Sort = "date"
		break
	case "downloads":
		search_param.Sort = "downloads"
		break
	case "filesize":
		search_param.Sort = "filesize"
	case "torrent_id":
	default:
		search_param.Sort = "torrent_id"
	}

	switch search_param.Order {
	case "asc":
		search_param.Order = "asc"
		break
	case "desc":
	default:
		search_param.Order = "desc"
	}

	order_by := search_param.Sort + " " + search_param.Order

	parameters := torrentService.WhereParams{}
	conditions := []string{}
	if searchCatId != "" {
		conditions = append(conditions, "category = ?")
		parameters.Params = append(parameters.Params, searchCatId)
	}
	if searchSubCatId != "" {
		conditions = append(conditions, "sub_category = ?")
		parameters.Params = append(parameters.Params, searchSubCatId)
	}
	if userId != "" {
		conditions = append(conditions, "uploader = ?")
		parameters.Params = append(parameters.Params, userId)
	}
	if search_param.Status != "" {
		if search_param.Status == "2" {
			conditions = append(conditions, "status != ?")
		} else {
			conditions = append(conditions, "status = ?")
		}
		parameters.Params = append(parameters.Params, search_param.Status)
	}
	searchQuerySplit := strings.Fields(search_param.Query)
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
		// TODO: make this faster ?
		conditions = append(conditions, "torrent_name ILIKE ?")
		parameters.Params = append(parameters.Params, "%"+searchQuerySplit[i]+"%")
	}

	parameters.Conditions = strings.Join(conditions[:], " AND ")
	log.Infof("SQL query is :: %s\n", parameters.Conditions)
	torrents, n := torrentService.GetTorrentsOrderBy(&parameters, order_by, maxPerPage, maxPerPage*(pagenum-1))
	return search_param, torrents, n
}
