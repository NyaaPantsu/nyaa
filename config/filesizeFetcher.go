package config

type FilesizeFetcherConfig struct {
	QueueSize      int `json:"queue_size"`
	Timeout        int `json:"timeout"`
	MaxDays        int `json:"max_days"`
	WakeUpInterval int `json:"wake_up_interval"`
}

var DefaultFilesizeFetcherConfig = FilesizeFetcherConfig{
	QueueSize: 10,
	Timeout: 120, // 2 min
	MaxDays: 90,
	WakeUpInterval: 300, // 5 min
}

