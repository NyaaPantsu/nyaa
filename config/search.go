package config

// SearchConfig : Config struct for search
// Is it deprecated?
type SearchConfig struct {
}

// DefaultSearchConfig : Default config for search
var DefaultSearchConfig = SearchConfig{}

const (
	// DefaultElasticsearchAnalyzer : default analyzer for ES
	DefaultElasticsearchAnalyzer = "nyaapantsu_search_analyzer"
	// DefaultElasticsearchIndex : default search index for ES
	DefaultElasticsearchIndex = "nyaapantsu"
	// DefaultElasticsearchType :  Name of the type in the es mapping
	DefaultElasticsearchType = "torrents"
)
