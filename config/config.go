package config

import (
	"flag"
	"io"
	"sync"

	"github.com/jinzhu/configor"
	yaml "gopkg.in/yaml.v2"
)

var (
	// DefaultConfigPath : path to the default config file (please do not change it)
	DefaultConfigPath = "config/default_config.yml"
	// ConfigPath : path to the user specific config file (please do not change it)
	ConfigPath = "config/config.yml"
)

var config *Config
var once sync.Once

func Get() *Config {
	once.Do(func() {
		config = &Config{}
	})
	return config
}

// IsSukebei : Tells if we are on the sukebei website
func IsSukebei() bool {
	return Get().Models.TorrentsTableName == "sukebei_torrents"
}

// WebAddress : Returns web address for current site
func WebAddress() string {
	if IsSukebei() {
		return Get().WebAddress.Sukebei
	} else {
		return Get().WebAddress.Nyaa
	}
}

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

func init() {
	Reload()
}

// Reload the configuration from the files provided in the config variables
func Reload() {
	configor.Load(Get(), DefaultConfigPath, ConfigPath)
}

// BindFlags returns a function which is to be used after
// flag.Parse to check and copy the flags' values to the Config instance.
func BindFlags() {
	confFile := flag.String("conf", ConfigPath, "path to the configuration file")
	flag.StringVar(&Get().DBType, "dbtype", Get().DBType, "database backend")
	flag.StringVar(&Get().Host, "host", Get().Host, "binding address of the server")
	flag.IntVar(&Get().Port, "port", Get().Port, "port of the server")
	flag.StringVar(&Get().DBParams, "dbparams", Get().DBParams, "parameters to open the database (see Gorm's doc)")
	flag.StringVar(&Get().DBLogMode, "dblogmode", Get().DBLogMode, "database log verbosity (errors only by default)")
	configor.Load(Get(), DefaultConfigPath, ConfigPath, *confFile)
}

// Pretty : Write config json in a file
func (config *Config) Pretty(output io.Writer) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	_, err = output.Write(data)
	return err
}
