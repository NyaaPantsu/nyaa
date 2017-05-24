package common

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	elastic "gopkg.in/olivere/elastic.v5"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/util/log"
)

// TorrentParam defines all parameters that can be provided when searching for a torrent
type TorrentParam struct {
	All       bool // True means ignore everything but Max and Offset
	Full      bool // True means load all members
	Order     bool // True means ascending
	Status    Status
	Sort      SortMode
	Category  Category
	Max       uint32
	Offset    uint32
	UserID    uint32
	TorrentID uint32
	NotNull   string // csv
	Null      string // csv
	NameLike  string // csv
}

// TODO Should probably return an error ?
func (p *TorrentParam) FromRequest(r *http.Request) {
	var err error

	nameLike := r.URL.Query().Get("q")
	if nameLike == "" {
		nameLike = "*"
	}

	page := mux.Vars(r)["page"]
	pagenum, err := strconv.ParseUint(page, 10, 32)
	if err != nil {
		pagenum = 1
	}

	max, err := strconv.ParseUint(r.URL.Query().Get("max"), 10, 32)
	if err != nil {
		max = config.TorrentsPerPage
	} else if max > config.MaxTorrentsPerPage {
		max = config.MaxTorrentsPerPage
	}

	// FIXME 0 means no userId defined
	userId, err := strconv.ParseUint(r.URL.Query().Get("userID"), 10, 32)
	if err != nil {
		userId = 0
	}

	var status Status
	status.Parse(r.URL.Query().Get("s"))

	var category Category
	category.Parse(r.URL.Query().Get("c"))

	var sortMode SortMode
	sortMode.Parse(r.URL.Query().Get("sort"))

	ascending := false
	if r.URL.Query().Get("order") == "true" {
		ascending = true
	}

	p.NameLike = nameLike
	p.Offset = uint32(pagenum)
	p.Max = uint32(max)
	p.UserID = uint32(userId)
	// TODO Use All
	p.All = false
	// TODO Use Full
	p.Full = false
	p.Order = ascending
	p.Status = status
	p.Sort = sortMode
	p.Category = category
	// FIXME 0 means no TorrentId defined
	// Do we even need that ?
	p.TorrentID = 0
}

// Builds a query string with for es query string query defined here
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html
func (p *TorrentParam) ToFilterQuery() string {
	// Don't set sub category unless main category is set
	query := ""
	if p.Category.IsMainSet() {
		query += "category:" + strconv.FormatInt(int64(p.Category.Main), 10)
		if p.Category.IsSubSet() {
			query += " sub_category:" + strconv.FormatInt(int64(p.Category.Sub), 10)
		}
	}

	if p.UserID != 0 {
		query += "uploader_id:" + strconv.FormatInt(int64(p.UserID), 10)
	}

	if p.Status != ShowAll {
		query += " status:" + p.Status.ToString()
	}
	return query
}

// Uses elasticsearch to find the torrents based on TorrentParam
func (p *TorrentParam) Find(client *elastic.Client) (int64, []model.TorrentJSON, error) {
	// TODO Why is it needed, what does it do ?
	ctx := context.Background()

	query := elastic.NewSimpleQueryStringQuery(p.NameLike).
		Field("name").
		Analyzer(config.DefaultElasticsearchAnalyzer).
		DefaultOperator("AND")

	// TODO Find a better way to keep in sync with mapping in ansible
	search := client.Search().
		Index(config.DefaultElasticsearchIndex).
		Query(query).
		Type(config.DefaultElasticsearchType).
		From(int((p.Offset - 1) * p.Max)).
		Size(int(p.Max)).
		Sort("_score", false). // TODO Do we want to sort by score first ?
		Sort(p.Sort.ToESField(), p.Order)

	filterQueryString := p.ToFilterQuery()
	if filterQueryString != "" {
		filterQuery := elastic.NewQueryStringQuery(filterQueryString).
			DefaultOperator("AND")
		search = search.PostFilter(filterQuery)
	}


	result, err := search.Do(ctx)
	if err != nil {
		return 0, nil, err
	}

	log.Infof("Query '%s' took %d milliseconds.", p.NameLike, result.TookInMillis)
	log.Infof("Amount of results %d.", result.TotalHits())

	torrents := make([]model.TorrentJSON, len(result.Hits.Hits))
	for i, hit := range result.Hits.Hits {
		var torrent model.TorrentJSON
		err := json.Unmarshal(*hit.Source, &torrent)
		if err == nil {
			torrents[i] = torrent
		} else {
			log.Errorf("Failed to decode TorrentJSON from '%v'", *hit.Source)
		}
	}

	return result.TotalHits(), torrents, nil
}

func (p *TorrentParam) Clone() TorrentParam {
	return TorrentParam{
		Order:     p.Order,
		Status:    p.Status,
		Sort:      p.Sort,
		Category:  p.Category,
		Max:       p.Max,
		Offset:    p.Offset,
		UserID:    p.UserID,
		TorrentID: p.TorrentID,
		NotNull:   p.NotNull,
		Null:      p.Null,
		NameLike:  p.NameLike,
	}
}
