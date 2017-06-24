package config

// Config : Configuration for DB, I2P, Fetcher, Go Server and Translation
type Config struct {
	Host                   string `json:"host" yaml:"host,omitempty"`
	Port                   int    `json:"port" yaml:"port,omitempty"`
	DBType                 string `json:"db_type" yaml:"db_type,omitempty"`
	Environment            string `json:"environment" yaml:"environment,omitempty"`
	AuthTokenExpirationDay int    `json:"auth_token_expiration" yaml:"auth_token_expiration,omitempty"`
	EnableSecureCSRF       bool   `json:"enable_secure_csrf" yaml:"enable_secure_csrf,omitempty"`
	DescriptionLength      int    `json:"description_length" yaml:"description_length,omitempty"`
	CommentLength          int    `json:"comment_length" yaml:"comment_length,omitempty"`
	// DBParams will be directly passed to Gorm, and its internal
	// structure depends on the dialect for each db type
	DBParams  string `json:"db_params" yaml:"db_params,omitempty"`
	DBLogMode string `json:"db_logmode" yaml:"db_logmode,omitempty"`
	Version   string `json:"version" yaml:"version,omitempty"`
	Build     string `yaml:"-"`
	// web address config
	WebAddress WebAddressConfig `yaml:"web_address,flow,omitempty"`
	// cookies config
	Cookies CookiesConfig `yaml:"cookies,flow,omitempty"`
	// tracker scraper config (required)
	Scrape ScraperConfig `json:"scraper" yaml:"scraper,flow,omitempty"`
	// cache config
	Cache CacheConfig `json:"cache" yaml:"cache,flow,omitempty"`
	// search config
	Search SearchConfig `json:"search" yaml:"search,flow,omitempty"`
	// optional i2p configuration
	I2P *I2PConfig `json:"i2p" yaml:"i2p,flow"`
	// filesize fetcher config
	MetainfoFetcher MetainfoFetcherConfig `json:"metainfo_fetcher" yaml:"metainfo_fetcher,flow,omitempty"`
	// internationalization config
	I18n I18nConfig `json:"i18n" yaml:"i18n,flow,omitempty"`
	// torrents config
	Torrents TorrentsConfig `yaml:"torrents,flow,omitempty"`
	// user config
	Users UsersConfig `yaml:"users,flow,omitempty"`
	// navigation config
	Navigation NavigationConfig `yaml:"navigation,flow,omitempty"`
	// log config
	Log LogConfig `yaml:"log,flow,omitempty"`
	// email config
	Email EmailConfig `yaml:"email,flow,omitempty"`
	// models config
	Models ModelsConfig `yaml:"models,flow,omitempty"`
}

// WebAddressConfig : Config struct for web addresses
type WebAddressConfig struct {
	Nyaa    string `yaml:"nyaa,omitempty"`
	Sukebei string `yaml:"sukebei,omitempty"`
	Status  string `yaml:"status,omitempty"`
}

// CookiesConfig : Config struct for session cookies
type CookiesConfig struct {
	DomainName    string `yaml:"domain_name,omitempty"`
	MaxAge        int    `yaml:"max_age,omitempty"`
	HashKey       string `yaml:"hash_key,omitempty"`
	EncryptionKey string `yaml:"encryption_key,omitempty"`
}

// CacheConfig is config struct for caching strategy
type CacheConfig struct {
	Dialect string  `yaml:"dialect,omitempty"`
	URL     string  `yaml:"url,omitempty"`
	Size    float64 `yaml:"size,omitempty"`
}

// I2PConfig : Config struct for I2P
type I2PConfig struct {
	Name    string `json:"name" yaml:"name,omitempty"`
	Addr    string `json:"samaddr" yaml:"addr,omitempty"`
	Keyfile string `json:"keyfile" yaml:"keyfile,omitempty"`
}

// I18nConfig : Config struct for translation
type I18nConfig struct {
	Directory       string `json:"translations_directory" yaml:"directory,omitempty"`
	DefaultLanguage string `json:"default_language" yaml:"default_language,omitempty"`
}

// ScrapeConfig : Config struct for Scraping
type ScrapeConfig struct {
	URL             string `json:"scrape_url" yaml:"url,omitempty"`
	Name            string `json:"name"  yaml:"name,omitempty"`
	IntervalSeconds int64  `json:"interval" yaml:"interval,omitempty"`
}

// ScraperConfig :  Config struct for Scraper
type ScraperConfig struct {
	Addr            string         `json:"bind" yaml:"addr,omitempty"`
	NumWorkers      int            `json:"workers" yaml:"workers,omitempty"`
	IntervalSeconds int64          `json:"default_interval" yaml:"default_interval,omitempty"`
	Trackers        []ScrapeConfig `json:"trackers" yaml:"trackers,omitempty"`
}

// TrackersConfig ; Config struct for Trackers
type TrackersConfig struct {
	Default        []string `yaml:"default,flow,omitempty"`
	NeededTrackers []int    `yaml:"needed,flow,omitempty"`
}

// TorrentsConfig : Config struct for Torrents
type TorrentsConfig struct {
	Status                        []bool            `yaml:"status,omitempty,omitempty"`
	SukebeiCategories             map[string]string `yaml:"sukebei_categories,omitempty"`
	CleanCategories               map[string]string `yaml:"clean_categories,omitempty"`
	EnglishOnlyCategories         []string          `yaml:"english_only_categories,omitempty"`
	NonEnglishOnlyCategories      []string          `yaml:"non_english_only_categories,omitempty"`
	AdditionalLanguages           []string          `yaml:"additional_languages,omitempty"`
	FileStorage                   string            `yaml:"filestorage,omitempty"`
	StorageLink                   string            `yaml:"storage_link,omitempty"`
	CacheLink                     string            `yaml:"cache_link,omitempty"`
	UploadsDisabled               bool              `yaml:"uploads_disabled,omitempty"`
	AdminsAreStillAllowedTo       bool              `yaml:"admins_are_still_allowed_to,omitempty"`
	TrustedUsersAreStillAllowedTo bool              `yaml:"trusted_users_are_still_allowed_to,omitempty"`
	Trackers                      TrackersConfig    `yaml:"trackers,flow,omitempty"`
	Order                         string            `yaml:"order,omitempty"`
	Sort                          string            `yaml:"sort,omitempty"`
}

// UsersConfig : Config struct for Users
type UsersConfig struct {
	DefaultUserSettings map[string]bool `yaml:"default_notifications_settings,flow,omitempty"`
}

// NavigationConfig : Config struct for Navigation
type NavigationConfig struct {
	TorrentsPerPage    int `yaml:"torrents_per_page,omitempty"`
	MaxTorrentsPerPage int `yaml:"max_torrents_per_page,omitempty"`
}

// MetainfoFetcherConfig : Config struct for metainfo fetcher
type MetainfoFetcherConfig struct {
	QueueSize        int `json:"queue_size" yaml:"queue_size,omitempty"`
	Timeout          int `json:"timeout" yaml:"timeout,omitempty"`
	MaxDays          int `json:"max_days" yaml:"max_days,omitempty"`
	BaseFailCooldown int `json:"base_fail_cooldown" yaml:"base_fail_cooldown,omitempty"`
	MaxFailCooldown  int `json:"max_fail_cooldown" yaml:"max_fail_cooldown,omitempty"`
	WakeUpInterval   int `json:"wake_up_interval" yaml:"wake_up_interval,omitempty"`

	UploadRateLimitKiB   int `json:"upload_rate_limit" yaml:"upload_rate_limit,omitempty"`
	DownloadRateLimitKiB int `json:"download_rate_limit" yaml:"download_rate_limit,omitempty"`

	FetchNewTorrentsOnly bool `json:"fetch_new_torrents_only" yaml:"fetch_new_torrents_only,omitempty"`
}

// LogConfig : Config struct for Logs
type LogConfig struct {
	AccessLogFilePath      string `yaml:"access_log_filepath,omitempty"`
	AccessLogFileExtension string `yaml:"access_log_fileextension,omitempty"`
	AccessLogMaxSize       int    `yaml:"access_log_max_size,omitempty"`
	AccessLogMaxBackups    int    `yaml:"access_log_max_backups,omitempty"`
	AccessLogMaxAge        int    `yaml:"access_log_max_age,omitempty"`
	ErrorLogFilePath       string `yaml:"error_log_filepath,omitempty"`
	ErrorLogFileExtension  string `yaml:"error_log_fileextension,omitempty"`
	ErrorLogMaxSize        int    `yaml:"error_log_max_size,omitempty"`
	ErrorLogMaxBackups     int    `yaml:"error_log_max_backups,omitempty"`
	ErrorLogMaxAge         int    `yaml:"error_log_max_age,omitempty"`
}

// EmailConfig : Config struct for email
type EmailConfig struct {
	SendEmail bool   `yaml:"send_email,omitempty"`
	From      string `yaml:"from,omitempty"`
	TestTo    string `yaml:"test_to,omitempty"`
	Host      string `yaml:"host,omitempty"`
	Username  string `yaml:"username,omitempty"`
	Password  string `yaml:"password,omitempty"`
	Port      int    `yaml:"port,omitempty"`
	Timeout   int    `yaml:"timeout,omitempty"`
}

// ModelsConfig : Config struct for models
type ModelsConfig struct {
	LastOldTorrentID       uint   `yaml:"last_old_torrent_id,omitempty"`
	TorrentsTableName      string `yaml:"torrents_table_name,omitempty"`
	ReportsTableName       string `yaml:"reports_table_name,omitempty"`
	CommentsTableName      string `yaml:"comments_table_name,omitempty"`
	UploadsOldTableName    string `yaml:"uploads_old_table_name,omitempty"`
	FilesTableName         string `yaml:"files_table_name,omitempty"`
	NotificationsTableName string `yaml:"notifications_table_name,omitempty"`
	ActivityTableName      string `yaml:"activities_table_name,omitempty"`
	ScrapeTableName        string `yaml:"scrape_table_name,omitempty"`
}

// SearchConfig : Config struct for search
type SearchConfig struct {
	ElasticsearchAnalyzer string `yaml:"es_analyze,omitempty"`
	ElasticsearchIndex    string `yaml:"es_index,omitempty"`
	ElasticsearchType     string `yaml:"es_type,omitempty"`
}
