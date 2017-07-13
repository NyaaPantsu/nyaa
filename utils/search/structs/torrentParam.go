package structs

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/gin-gonic/gin"
)

// FromRequest : parse a request in torrent param
// TODO Should probably return an error ?
func (p *TorrentParam) FromRequest(c *gin.Context) {
	var err error

	nameLike := strings.TrimSpace(c.Query("q"))
	max, err := strconv.ParseUint(c.DefaultQuery("limit", strconv.Itoa(config.Get().Navigation.TorrentsPerPage)), 10, 32)
	if err != nil {
		max = uint64(config.Get().Navigation.TorrentsPerPage)
	} else if max > uint64(config.Get().Navigation.MaxTorrentsPerPage) {
		max = uint64(config.Get().Navigation.MaxTorrentsPerPage)
	}

	// FIXME 0 means no userId defined
	userID, err := strconv.ParseUint(c.Query("userID"), 10, 32)
	if err != nil {
		userID = 0
	}

	// FIXME 0 means no userId defined
	fromID, err := strconv.ParseUint(c.Query("fromID"), 10, 32)
	if err != nil {
		fromID = 0
	}

	var status Status
	status.Parse(c.Query("s"))

	maxage, err := strconv.Atoi(c.Query("maxage"))
	fromDate, toDate := DateFilter(""), DateFilter("")
	if err != nil {
		// if to xxx is not provided, fromDate is equal to from xxx
		if c.Query("toDate") != "" {
			fromDate.Parse(c.Query("toDate"), c.Query("dateType"))
			toDate.Parse(c.Query("fromDate"), c.Query("dateType"))
		} else {
			fromDate.Parse(c.Query("fromDate"), c.Query("dateType"))
		}
	} else {
		fromDate = DateFilter(time.Now().AddDate(0, 0, -maxage).Format("2006-01-02"))
	}

	categories := ParseCategories(c.Query("c"))

	var sortMode SortMode
	sortMode.Parse(c.Query("sort"))

	var minSize SizeBytes
	var maxSize SizeBytes

	minSize.Parse(c.Query("minSize"), c.Query("sizeType"))
	maxSize.Parse(c.Query("maxSize"), c.Query("sizeType"))

	ascending := false
	if c.Query("order") == "true" {
		ascending = true
	}
	langs := c.QueryArray("lang")

	language := ParseLanguages(strings.Join(langs, ","))

	p.NameLike = nameLike
	p.Max = uint32(max)
	p.UserID = uint32(userID)
	p.Order = ascending
	p.Status = status
	p.Sort = sortMode
	p.Category = categories
	p.Languages = language
	p.FromDate = fromDate
	p.ToDate = toDate
	p.MinSize = minSize
	p.MaxSize = maxSize
	// FIXME 0 means no TorrentId defined
	// Do we even need that ?
	p.TorrentID = 0
	// Needed to display result after a certain torrentID
	p.FromID = uint32(fromID)
}

// ToFilterQuery : Builds a query string with for es query string query defined here
// https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html
func (p *TorrentParam) ToFilterQuery() string {
	// Don'p set sub category unless main category is set
	query := ""
	if len(p.Category) > 0 {
		conditionsOr := make([]string, len(p.Category))
		for key, val := range p.Category {
			if val.IsSubSet() {
				conditionsOr[key] = "(category: " + strconv.FormatInt(int64(val.Main), 10) + " AND sub_category: " + strconv.FormatInt(int64(val.Sub), 10) + ")"
			} else {
				conditionsOr[key] = "(category: " + strconv.FormatInt(int64(val.Main), 10) + ")"
			}
		}
		query += "(" + strings.Join(conditionsOr, " OR ") + ")"
	}

	if p.UserID != 0 {
		query += " uploader_id:" + strconv.FormatInt(int64(p.UserID), 10)
	}

	if p.Hidden {
		query += " hidden:false"
	}

	if p.Status != ShowAll {
		if p.Status != FilterRemakes {
			query += " status:" + p.Status.ToString()
		} else {
			/* From the old nyaa behavior, FilterRemake means everything BUT
			 * remakes
			 */
			query += " !status:" + p.Status.ToString()
		}
	}

	if p.FromID != 0 {
		query += " id:>" + strconv.FormatInt(int64(p.FromID), 10)
	}

	if p.FromDate != "" && p.ToDate != "" {
		query += " date: [" + string(p.FromDate) + " " + string(p.ToDate) + "]"
	} else if p.FromDate != "" {
		query += " date: [" + string(p.FromDate) + " *]"
	} else if p.ToDate != "" {
		query += " date: [* " + string(p.ToDate) + "]"
	}

	sMinSize := strconv.FormatUint(uint64(p.MinSize), 10)
	sMaxSize := strconv.FormatUint(uint64(p.MaxSize), 10)
	if p.MinSize > 0 && p.MaxSize > 0 {
		query += " filesize: [" + sMinSize + " " + sMaxSize + "]"
	} else if p.MinSize > 0 {
		query += " filesize: [" + sMinSize + " *]"
	} else if p.MaxSize > 0 {
		query += " filesize: [* " + sMaxSize + "]"
	}

	if len(p.Languages) > 0 {
		for _, val := range p.Languages {
			query += " language: " + val.Code
		}
	}

	return query
}

// Find :
/* Uses elasticsearch to find the torrents based on TorrentParam
 * We decided to fetch only the ids from ES and then query these ids to the
 * database
 */
func (p *TorrentParam) Find(client *elastic.Client) (int64, []models.Torrent, error) {
	// TODO Why is it needed, what does it do ?
	ctx := context.Background()

	var query elastic.Query
	if p.NameLike == "" {
		query = elastic.NewMatchAllQuery()
	} else {
		query = elastic.NewSimpleQueryStringQuery(p.NameLike).
			Field("name").
			Analyzer(config.Get().Search.ElasticsearchAnalyzer).
			DefaultOperator("AND")
	}

	// TODO Find a better way to keep in sync with mapping in ansible
	search := client.Search().
		Index(config.Get().Search.ElasticsearchIndex).
		Query(query).
		Type(config.Get().Search.ElasticsearchType).
		From(int((p.Offset-1)*p.Max)).
		Size(int(p.Max)).
		Sort(p.Sort.ToESField(), p.Order).
		Sort("_score", false) // Don'p put _score before the field sort, it messes with the sorting

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

	var torrents []models.Torrent
	var torrentCount int64
	torrentCount = 0
	//torrents := make([]models.Torrent, len(result.Hits.Hits))
	if len(result.Hits.Hits) <= 0 {
		return 0, nil, nil
	}
	for _, hit := range result.Hits.Hits {
		var tJSON models.TorrentJSON
		err := json.Unmarshal(*hit.Source, &tJSON)
		if err == nil {
			torrents = append(torrents, tJSON.ToTorrent())
			torrentCount++
		} else {
			log.Infof("Cannot unmarshal elasticsearch torrent: %s", err)
		}
	}
	return result.TotalHits(), torrents, nil

}

// Clone : To clone a torrent params
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
		FromID:    p.FromID,
		FromDate:  p.FromDate,
		ToDate:    p.ToDate,
		NotNull:   p.NotNull,
		NameLike:  p.NameLike,
		Languages: p.Languages,
		MinSize:   p.MinSize,
		MaxSize:   p.MaxSize,
	}
}
