package search

import(
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util/log"
	"net/http"
	"html"
	"strconv"
	"strings"
)

type SearchParam struct {
	Category   string
	Order      string
	Query      string
	Max        int
	Status     string
	Sort       string
}

func SearchByQuery(r *http.Request, pagenum int) (SearchParam, []model.Torrents, int) {
	maxPerPage, errConv := strconv.Atoi(r.URL.Query().Get("max"))
	if errConv != nil {
		maxPerPage = 50 // default Value maxPerPage
	}
	
	search_param := SearchParam{}
	search_param.Max = maxPerPage
	search_param.Query = r.URL.Query().Get("q")
	search_param.Category = r.URL.Query().Get("c")
	search_param.Status = r.URL.Query().Get("s")
	search_param.Sort = r.URL.Query().Get("sort")
	search_param.Order = r.URL.Query().Get("order")

	catsSplit := strings.Split(search_param.Category, "_")
	// need this to prevent out of index panics
	var searchCatId, searchSubCatId string
	if len(catsSplit) == 2 {

		searchCatId = html.EscapeString(catsSplit[0])
		searchSubCatId = html.EscapeString(catsSplit[1])
	}
	if search_param.Sort == "" {
		search_param.Sort = "torrent_id"
	}
	if search_param.Order == "" {
		search_param.Order = "desc"
	}
	order_by := search_param.Sort + " " + search_param.Order

	parameters := torrentService.WhereParams{}
	conditions := []string{}
	if searchCatId != "" {
		conditions = append(conditions, "category_id = ?")
		parameters.Params = append(parameters.Params, searchCatId)
	}
	if searchSubCatId != "" {
		conditions = append(conditions, "sub_category_id = ?")
		parameters.Params = append(parameters.Params, searchSubCatId)
	}
	if search_param.Status != "" {
		if search_param.Status == "2" {
			conditions = append(conditions, "status_id != ?")
		} else {
			conditions = append(conditions, "status_id = ?")
		}
		parameters.Params = append(parameters.Params, search_param.Status)
	}
	searchQuerySplit := strings.Split(search_param.Query, " ")
	for i, _ := range searchQuerySplit {
		conditions = append(conditions, "torrent_name LIKE ?")
		parameters.Params = append(parameters.Params, "%"+searchQuerySplit[i]+"%")
	}

	parameters.Conditions = strings.Join(conditions[:], " AND ")
	log.Infof("SQL query is :: %s\n", parameters.Conditions)
	torrents, n := torrentService.GetTorrentsOrderBy(&parameters, order_by, maxPerPage, maxPerPage*(pagenum-1))
	return search_param, torrents, n
}
