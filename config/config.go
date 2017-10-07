package config

import (
	"flag"
	"io"
	"sync"

	"fmt"

	"github.com/jinzhu/configor"
	yaml "gopkg.in/yaml.v2"
)

var config *Config
var once sync.Once

// Configpaths default configuration file paths
var Configpaths = []string{"config/config.yml", "config/default_config.yml"}

// Get config variable
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

// DefaultTheme : Return the default theme or default dark theme
func DefaultTheme(dark bool) string {
	if Get().DefaultTheme.Forced != "" {
		return Get().DefaultTheme.Forced
	}
	if !dark {
		return Get().DefaultTheme.Theme
	} else {
		return Get().DefaultTheme.Dark
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
	fmt.Println("Config reload")
	fmt.Println(Configpaths)
	newConf := &Config{
		DBType:    Get().DBType,
		Host:      Get().Host,
		Port:      Get().Port,
		DBParams:  Get().DBParams,
		DBLogMode: Get().DBLogMode,
	}
	config = newConf
	configor.Load(config, Configpaths...)
	fmt.Printf("Port: %d", Get().Port)
}

// BindFlags returns a function which is to be used after
// flag.Parse to check and copy the flags' values to the Config instance.
func BindFlags() func() {
	confFile := flag.String("conf", Configpaths[1], "path to the configuration file")
	flag.StringVar(&Get().DBType, "dbtype", Get().DBType, "database backend")
	flag.StringVar(&Get().Host, "host", Get().Host, "binding address of the server")
	flag.IntVar(&Get().Port, "port", Get().Port, "port of the server")
	flag.StringVar(&Get().DBParams, "dbparams", Get().DBParams, "parameters to open the database (see Gorm's doc)")
	flag.StringVar(&Get().DBLogMode, "dblogmode", Get().DBLogMode, "database log verbosity (errors only by default)")
	return func() {
		if *confFile != "" && *confFile != Configpaths[1] {
			Configpaths = append([]string{*confFile}, Configpaths...)
			Reload()
		}
	}
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
