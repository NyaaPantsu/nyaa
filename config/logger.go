package config

const (
	AccessLogFilePath      = "log/access"
	AccessLogFileExtension = ".txt"
	AccessLogMaxSize       = 5 // megabytes
	AccessLogMaxBackups    = 7
	AccessLogMaxAge        = 30 //days
	ErrorLogFilePath       = "log/error"
	ErrorLogFileExtension  = ".json"
	ErrorLogMaxSize        = 10 // megabytes
	ErrorLogMaxBackups     = 7
	ErrorLogMaxAge         = 30 //days
)

type LogConfig struct {
	Environment string `json:"env"`
}

var DefaultLogConfig = LogConfig{
	Environment: "PRODUCTION",
}
