package config

type SearchConfig struct {
}

var DefaultSearchConfig = SearchConfig{}

const (
	DefaultElasticsearchAnalyzer = "nyaapantsu_analyzer"
	DefaultElasticsearchIndex = "nyaapantsu"
	DefaultElasticsearchType = "torrents" // Name of the type in the es mapping
)
