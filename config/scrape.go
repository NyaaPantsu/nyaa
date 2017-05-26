package config

// ScrapeConfig : Config struct for Scraping
type ScrapeConfig struct {
	URL             string `json:"scrape_url"`
	Name            string `json:"name"`
	IntervalSeconds int64  `json:"interval"`
}

// ScraperConfig :  Config struct for Scraper
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
		{
			URL:  "udp://tracker.coppersurfer.tk:6969/",
			Name: "coppersurfer.tk",
		},
	},
}
