package config

type MetainfoFetcherConfig struct {
	QueueSize      int `json:"queue_size"`
	Timeout        int `json:"timeout"`
	MaxDays        int `json:"max_days"`
	FailCooldown   int `json:"fail_cooldown"`
	WakeUpInterval int `json:"wake_up_interval"`

	UploadRateLimiter   int `json:"upload_rate_limiter"`
	DownloadRateLimiter int `json:"download_rate_limiter"`
}

var DefaultMetainfoFetcherConfig = MetainfoFetcherConfig{
	QueueSize:      10,
	Timeout:        120, // 2 min
	MaxDays:        90,
	FailCooldown:   30 * 60, // in seconds, when failed torrents will be able to be fetched again.
	WakeUpInterval: 300, // 5 min

	UploadRateLimiter:   1024, // kbps
	DownloadRateLimiter: 1024,
}

