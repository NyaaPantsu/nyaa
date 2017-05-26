package config

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	// LastOldTorrentID is the highest torrent ID
	// that was copied from the original Nyaa
	LastOldTorrentID = 923000
	// TorrentsTableName : Name of torrent table in DB
	TorrentsTableName = "torrents"
	// ReportsTableName : Name of torrent report table in DB
	ReportsTableName = "torrent_reports"
	// CommentsTableName : Name of comments table in DB
	CommentsTableName = "comments"
	// UploadsOldTableName : Name of uploads table in DB
	UploadsOldTableName = "user_uploads_old"
	// FilesTableName : Name of files table in DB
	FilesTableName = "files"
	// NotificationTableName : Name of notifications table in DB
	NotificationTableName = "notifications"

	// for sukebei:
	//LastOldTorrentID    = 2303945
	//TorrentsTableName   = "sukebei_torrents"
	//ReportsTableName    = "sukebei_torrent_reports"
	//CommentsTableName   = "sukebei_comments"
	//UploadsOldTableName = "sukebei_user_uploads_old"
	//FilesTableName      = "sukebei_files"
)

// IsSukebei : Tells if we are on the sukebei website
func IsSukebei() bool {
	return TorrentsTableName == "sukebei_torrents"
}

// Config : Configuration for DB, I2P, Fetcher, Go Server and Translation
type Config struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	DBType string `json:"db_type"`
	// DBParams will be directly passed to Gorm, and its internal
	// structure depends on the dialect for each db type
	DBParams  string `json:"db_params"`
	DBLogMode string `json:"db_logmode"`
	// tracker scraper config (required)
	Scrape ScraperConfig `json:"scraper"`
	// cache config
	Cache CacheConfig `json:"cache"`
	// search config
	Search SearchConfig `json:"search"`
	// optional i2p configuration
	I2P *I2PConfig `json:"i2p"`
	// filesize fetcher config
	MetainfoFetcher MetainfoFetcherConfig `json:"metainfo_fetcher"`
	// internationalization config
	I18n I18nConfig `json:"i18n"`
}

// Defaults : Configuration by default
var Defaults = Config{"localhost", 9999, "sqlite3", "./nyaa.db?cache_size=50", "default", DefaultScraperConfig, DefaultCacheConfig, DefaultSearchConfig, nil, DefaultMetainfoFetcherConfig, DefaultI18nConfig}

var allowedDatabaseTypes = map[string]bool{
	"sqlite3":  true,
	"postgres": true,
	"mysql":    true,
	"mssql":    true,
}

var allowedDBLogModes = map[string]bool{
	"default":  true, // errors only
	"detailed": true,
	"silent":   true,
}

// New : Construct a new config variable
func New() *Config {
	var config Config
	config.Host = Defaults.Host
	config.Port = Defaults.Port
	config.DBType = Defaults.DBType
	config.DBParams = Defaults.DBParams
	config.DBLogMode = Defaults.DBLogMode
	config.Scrape = Defaults.Scrape
	config.Cache = Defaults.Cache
	config.MetainfoFetcher = Defaults.MetainfoFetcher
	config.I18n = Defaults.I18n
	return &config
}

// BindFlags returns a function which is to be used after
// flag.Parse to check and copy the flags' values to the Config instance.
func (config *Config) BindFlags() func() error {
	confFile := flag.String("conf", "", "path to the configuration file")
	dbType := flag.String("dbtype", Defaults.DBType, "database backend")
	host := flag.String("host", Defaults.Host, "binding address of the server")
	port := flag.Int("port", Defaults.Port, "port of the server")
	dbParams := flag.String("dbparams", Defaults.DBParams, "parameters to open the database (see Gorm's doc)")
	dbLogMode := flag.String("dblogmode", Defaults.DBLogMode, "database log verbosity (errors only by default)")

	return func() error {
		// You can override fields in the config file with flags.
		config.Host = *host
		config.Port = *port
		config.DBParams = *dbParams
		err := config.SetDBType(*dbType)
		if err != nil {
			return err
		}
		err = config.SetDBLogMode(*dbLogMode)
		if err != nil {
			return err
		}
		err = config.HandleConfFileFlag(*confFile)
		return err
	}
}

// HandleConfFileFlag : Read the config from a file
func (config *Config) HandleConfFileFlag(path string) error {
	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("can't read file '%s'", path)
		}

		err = config.Read(bufio.NewReader(file))
		if err != nil {
			return fmt.Errorf("failed to parse file '%s' (%s)", path, err)
		}
	}
	return nil
}

// SetDBType : Set the DataBase type in config
func (config *Config) SetDBType(dbType string) error {
	if !allowedDatabaseTypes[dbType] {
		return fmt.Errorf("unknown database backend '%s'", dbType)
	}
	config.DBType = dbType
	return nil
}

// SetDBLogMode : Set the log mode in config
func (config *Config) SetDBLogMode(dbLogmode string) error {
	if !allowedDBLogModes[dbLogmode] {
		return fmt.Errorf("unknown database log mode '%s'", dbLogmode)
	}
	config.DBLogMode = dbLogmode
	return nil
}

// Read : Decode config from json to config
func (config *Config) Read(input io.Reader) error {
	return json.NewDecoder(input).Decode(config)
}

// Write : Encode config from json to config
func (config *Config) Write(output io.Writer) error {
	return json.NewEncoder(output).Encode(config)
}

// Pretty : Write config json in a file
func (config *Config) Pretty(output io.Writer) error {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	data = append(data, []byte("\n")...)
	_, err = output.Write(data)
	return err
}
