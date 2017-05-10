package config

type ScrapeConfig struct {
	URL             string `json:"scrape_url"`
	Name            string `json:"name"`
	IntervalSeconds int64  `json:"interval"`
}

type ScraperConfig struct {
	Addr            string         `json:"bind"`
	NumWorkers      int            `json:"workers"`
	IntervalSeconds int64          `json:"default_interval"`
	Trackers        []ScrapeConfig `json:"trackers"`
}

// DefaultScraperConfig is the default config for bittorrent scraping
var DefaultScraperConfig = ScraperConfig{
	Addr: ":9999",
	// TODO: query system?
	NumWorkers: 4,
	// every hour
	IntervalSeconds: 60 * 60,
	Trackers: []ScrapeConfig{
		ScrapeConfig{
			URL:  "udp://tracker.doko.moe:6969/",
			Name: "doko.moe",
		},
	},
}
