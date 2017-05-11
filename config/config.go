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
)

type Config struct {
	Host   string `json:"host"`
	Port   int    `json:"port"`
	DBType string `json:"db_type"`
	// DBParams will be directly passed to Gorm, and its internal
	// structure depends on the dialect for each db type
	DBParams string `json:"db_params"`
	DBLogMode string `json:"db_logmode"`
	// tracker scraper config (required)
	Scrape ScraperConfig `json:"scraper"`
	// cache config
	Cache CacheConfig `json:"cache"`
	// search config
	Search SearchConfig `json:"search"`
	// optional i2p configuration
	I2P *I2PConfig `json:"i2p"`
}

var Defaults = Config{"localhost", 9999, "sqlite3", "./nyaa.db?cache_size=50", "default", DefaultScraperConfig, DefaultCacheConfig, DefaultSearchConfig, nil}

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

func New() *Config {
	var config Config
	config.Host = Defaults.Host
	config.Port = Defaults.Port
	config.DBType = Defaults.DBType
	config.DBParams = Defaults.DBParams
	config.DBLogMode = Defaults.DBLogMode
	config.Scrape = Defaults.Scrape
	config.Cache = Defaults.Cache
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

func (config *Config) SetDBType(db_type string) error {
	if !allowedDatabaseTypes[db_type] {
		return fmt.Errorf("unknown database backend '%s'", db_type)
	}
	config.DBType = db_type
	return nil
}

func (config *Config) SetDBLogMode(db_logmode string) error {
	if !allowedDBLogModes[db_logmode] {
		return fmt.Errorf("unknown database log mode '%s'", db_logmode)
	}
	config.DBLogMode = db_logmode
	return nil
}

func (config *Config) Read(input io.Reader) error {
	return json.NewDecoder(input).Decode(config)
}

func (config *Config) Write(output io.Writer) error {
	return json.NewEncoder(output).Encode(config)
}

func (config *Config) Pretty(output io.Writer) error {
	data, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	data = append(data, []byte("\n")...)
	_, err = output.Write(data)
	return err
}
