package structs

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/gin-gonic/gin"
)

// TorrentParam defines all parameters that can be provided when searching for a torrent
type TorrentParam struct {
	Full      bool // True means load all members
	Order     bool // True means ascending
	Hidden    bool // True means filter hidden torrents
	Deleted   bool // False means filter deleted torrents
	Status    Status
	Sort      SortMode
	Category  Categories
	Max       maxType
	Offset    uint32
	UserID    uint32
	TorrentID []uint32
	FromID    uint32
	FromDate  DateFilter
	ToDate    DateFilter
	NotNull   string // csv
	NameLike  string // csv
	Languages publicSettings.Languages
	MinSize   SizeBytes
	MaxSize   SizeBytes
}

// Identifier returns a unique identifier for the struct
func (p *TorrentParam) Identifier() string {
	cats := ""
	for _, v := range p.Category {
		cats += fmt.Sprintf("%d%d", v.Main, v.Sub)
	}
	languages := ""
	for _, v := range p.Languages {
		languages += fmt.Sprintf("%s%s", v.Code, v.Name)
	}
	ids := ""
	for _, v := range p.TorrentID {
		ids += fmt.Sprintf("%d", v)
	}
	return fmt.Sprintf("%s%s%s%d%d%d%d%d%d%d%s%s%d%s%s%t%t%t%t", p.NameLike, p.NotNull, languages, p.Max, p.Offset, p.FromID, p.MinSize, p.MaxSize, p.Status, p.Sort, p.FromDate, p.ToDate, p.UserID, ids, cats, p.Full, p.Order, p.Hidden, p.Deleted)
}

// FromRequest : parse a request in torrent param
// TODO Should probably return an error ?
func (p *TorrentParam) FromRequest(c *gin.Context) {
	var err error

	// We take the search arguments from "q" in url
	nameLike := strings.TrimSpace(c.Query("q"))
	var max maxType
	// We take the maxximum results to display from "limit" in url
	max.Parse(c.Query("limit"))

	// Get the user id from the url
	userID, err := strconv.ParseUint(c.Query("userID"), 10, 32)
	if err != nil {
		// if you can't convert it, you set it to 0
		userID = 0
	}

	// Get the torrent ID to limit the results to the ones after this torrent
	fromID, err := strconv.ParseUint(c.Query("fromID"), 10, 32)
	if err != nil {
		// if you can't convert it, you set it to 0
		fromID = 0
	}

	var status Status
	// helper to parse status from the "s" argument in url
	status.Parse(c.Query("s"))

	// maxage is an int parameter limiting the results to the last "x" days (old nyaa behavior)
	fromDate, toDate := backwardCompatibility(c.Query("maxage"), c.Query("fromDate"), c.Query("toDate"), c.Query("dateType"))

	// Parse the categories from the "c" argument in url
	categories := ParseCategories(c.Query("c"))

	var sortMode SortMode
	// Parse the sorting mode of the result from the "sort" argument in url
	sortMode.Parse(c.Query("sort"))

	var minSize SizeBytes
	var maxSize SizeBytes

	// Parsing minimum and maximum size from the sizeType given (minSize & maxSize & sizeType in url)
	minSize.Parse(c.Query("minSize"), c.Query("sizeType"))
	maxSize.Parse(c.Query("maxSize"), c.Query("sizeType"))

	// Getting the order from the "order" argument in url, we default to descending order
	ascending := false
	if c.Query("order") == "true" {
		ascending = true
	}

	// We get the languages filtering the results from the "lang" argument in url
	language := ParseLanguages(c.QueryArray("lang"))

	ids := c.QueryArray("id")

	for _, id := range ids {
		idInt, err := strconv.Atoi(id)
		if err == nil {
			p.TorrentID = append(p.TorrentID, uint32(idInt))
		}
	}
	// Search by name
	p.NameLike = nameLike
	// Maximum results returned
	p.Max = max
	// Limit search to one user
	p.UserID = uint32(userID)
	// Order to return the results
	p.Order = ascending
	// Limit to some status the results
	p.Status = status
	// Sort the results
	p.Sort = sortMode
	// Category in which you have to search
	p.Category = categories
	// Languages filter of the torrents
	p.Languages = language
	// From which date you need to search
	p.FromDate = fromDate
	// To which date you need to search
	p.ToDate = toDate
	// Minimum size to search
	p.MinSize = minSize
	// Maximum size to search
	p.MaxSize = maxSize
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
				if val.Main > 0 {
					conditionsOr[key] = "(category: " + strconv.FormatInt(int64(val.Main), 10) + ")"
				}
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
			query += " status:" + p.Status.String()
		} else {
			/* From the old nyaa behavior, FilterRemake means everything BUT
			 * remakes
			 */
			query += " !status:" + p.Status.String()
		}
	}

	if p.FromID != 0 {
		query += " id:>" + strconv.FormatInt(int64(p.FromID), 10)
	}

	if len(p.TorrentID) > 0 {
		for _, id := range p.TorrentID {
			query += fmt.Sprintf(" id:%d", id)
		}
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
		From(int((p.Offset-1)*uint32(p.Max))).
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
