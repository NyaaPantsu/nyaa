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

func (p *TorrentParam) FromRequest(r *http.Request) error {
	var err error

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

	// FIXME Use maxint to determine whether the query should include userID
	userId, err := strconv.ParseUint(r.URL.Query().Get("userID"), 10, 32)
	if err != nil {
		userId = uint64(^uint32(0)) // Bit-wise NOT to find maximum value
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

	p.NameLike = r.URL.Query().Get("q")
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
	// FIXME Don't use Max Uint32 value to determine if we use TorrentID
	p.TorrentID = ^uint32(0) // Bit-wise NOT to find maximum value
	return nil
}

// TODO Query filter with Sort, Status, Category, UserID
func (p *TorrentParam) Find(client elastic.Client) ([]model.TorrentJSON, error) {
	// TODO Don't create a new client for every search
	// TODO Why is it needed, what does it do ?
	ctx := context.Background()

	simpleQuery := elastic.NewSimpleQueryStringQuery(p.NameLike).
		Analyzer(config.DefaultElasticsearchAnalyzer).
		Flags("AND|OR|NOT|PREFIX|PHRASE|PRECEDENCE|WHITESPACE").
		Field("name").
		DefaultOperator("AND")

		// Specify all field that should be returned by Elasticsearch
		// TODO Add seeders, leechers, completed, etc
		// TODO Find a better way to keep in sync with mapping in ansible
	fsc := elastic.NewFetchSourceContext(true).
		Include("id", "name", "category", "sub_category", "status", "hash",
			"date", "uploader_id", "downloads", "filesize")

	search := client.Search().
		Index(config.DefaultElasticsearchIndex).
		Query(simpleQuery).
		Type(config.DefaultElasticsearchType).
		From(int(p.Offset)).
		Size(int(p.Max)).
		Sort("_score", false).
		Sort("date", false).
		FetchSourceContext(fsc)

	result, err := search.Do(ctx)
	if err != nil {
		return nil, err
	}

	log.Infof("Query '%s' took %d milliseconds.\n", p.NameLike, result.TookInMillis)
	log.Infof("Amount of results %d.\n", result.TotalHits())

	torrents := make([]model.TorrentJSON, result.Hits.TotalHits)
	for _, hit := range result.Hits.Hits {
		var torrent model.TorrentJSON
		err := json.Unmarshal(*hit.Source, &torrent)
		if err == nil {
			torrents = append(torrents, torrent)
		} else {
			log.Errorf("Failed to decode TorrentJSON from '%v'", *hit.Source)
		}
	}

	return torrents, nil
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
