package config

// SearchConfig : Config struct for search
// Is it deprecated?
type SearchConfig struct {
}

// DefaultSearchConfig : Default config for search
var DefaultSearchConfig = SearchConfig{}

const (
	DefaultElasticsearchAnalyzer = "nyaapantsu_analyzer"
	DefaultElasticsearchIndex = "nyaapantsu"
	DefaultElasticsearchType = "torrents" // Name of the type in the es mapping
)
