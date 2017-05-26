package config

const (
	// AccessLogFilePath : Path to logs access
	AccessLogFilePath = "log/access"
	// AccessLogFileExtension : Extension for log file
	AccessLogFileExtension = ".txt"
	// AccessLogMaxSize : Size max for a log file in megabytes
	AccessLogMaxSize = 5
	// AccessLogMaxBackups : Number of file for logs
	AccessLogMaxBackups = 7
	// AccessLogMaxAge : Number of days that we keep logs
	AccessLogMaxAge = 30
	// ErrorLogFilePath : Path to log errors
	ErrorLogFilePath = "log/error"
	// ErrorLogFileExtension : Extension for log file
	ErrorLogFileExtension = ".json"
	// ErrorLogMaxSize : Size max for a log file in megabytes
	ErrorLogMaxSize = 10
	// ErrorLogMaxBackups : Number of file for logs
	ErrorLogMaxBackups = 7
	// ErrorLogMaxAge : Number of days that we keep logs
	ErrorLogMaxAge = 30
)
