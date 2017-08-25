package search

import (
	"errors"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"gopkg.in/olivere/elastic.v5"
)

// Query Type for outputing
type Query struct {
	text         []string
	TorrentParam *TorrentParam
	where        []interface{}
}

// String convert a query to string
func (q *Query) String() string {
	return strings.Join(q.text, " ")
}

// ToESQuery convert a query to use with ES search
// TODO add a test for this function
func (q *Query) ToESQuery(client *elastic.Client) (*elastic.SearchService, error) {
	if q.TorrentParam == nil {
		return nil, errors.New("You have to define TorrentParam before using ToESQuery")
	}
	var query elastic.Query
	if q.TorrentParam.NameLike == "" {
		query = elastic.NewMatchAllQuery()
	}
	query = elastic.NewSimpleQueryStringQuery(q.TorrentParam.NameLike).
		Field("name").
		Analyzer(config.Get().Search.ElasticsearchAnalyzer).
		DefaultOperator("AND")

	search := client.Search().
		Index(config.Get().Search.ElasticsearchIndex).
		Query(query).
		Type(config.Get().Search.ElasticsearchType).
		From(int((q.TorrentParam.Offset-1)*uint32(q.TorrentParam.Max))).
		Size(int(q.TorrentParam.Max)).
		Sort(q.TorrentParam.Sort.ToESField(), q.TorrentParam.Order).
		Sort("_score", false)

	if len(q.text) > 0 {
		filterQuery := elastic.NewQueryStringQuery(strings.Join(q.text, " ")).
			DefaultOperator("AND")
		search = search.PostFilter(filterQuery)
	}
	return search, nil
}

// ToDBQuery convert a query to use with PG search
func (q *Query) ToDBQuery() (string, []interface{}) {
	return strings.Join(q.text, " AND "), q.where
}

// Append a string to query
func (q *Query) Append(s string, args ...interface{}) {
	if s != "" {
		if len(args) > 0 {
			q.where = append(q.where, args...)
			if !strings.Contains(s, "?") && !strings.Contains(s, "=") {
				q.text = append(q.text, s+" = ?")
				return
			}
		}
		q.text = append(q.text, s)
	}
}

// Prepend a string to query
func (q *Query) Prepend(s string, args ...interface{}) {
	if s != "" {
		if len(args) > 0 {
			q.where = append(args, q.where...)
			if !strings.Contains(s, "?") && !strings.Contains(s, "=") {
				sarr := []string{s + " = ?"}
				q.text = append(sarr, q.text...)
				return
			}
		}
		sarr := []string{s}
		q.text = append(sarr, q.text...)
	}
}
