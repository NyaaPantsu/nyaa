package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sync"
	)

type Config struct {
	Host   string `json: "host"`
	Port   int    `json: "port"`
	DBType string `json: "db_type"`
	// This will be directly passed to Gorm, and its internal
	// structure depends on the dialect for each db type
	DBParams string `json: "db_type"`
}

var Defaults = Config{"localhost", 9999, "sqlite3", "./nyaa.db?cache_size=50"}
var PrintDefaults *bool
var allowedDatabaseTypes = map[string]bool{
	"sqlite3":  true,
	"postgres": true,
	"mysql":    true,
	"mssql":    true,
}

var instance *Config
var once sync.Once

func GetInstance() *Config {
    once.Do(func() {
        instance = &Config{}
		instance.Host = Defaults.Host
		instance.Port = Defaults.Port
		instance.DBType = Defaults.DBType
		instance.DBParams = Defaults.DBParams
		instance.BindFlags()

    })
    return instance
}

type processFlags func() error

func (config *Config) BindFlags() error {
	// This function returns a function which is to be used after
	// flag.Parse to check and copy the flags' values to the config instance.

	conf_file := flag.String("conf", "", "path to the configuration file")
	db_type := flag.String("dbtype", Defaults.DBType, "database backend")
	host := flag.String("host", Defaults.Host, "binding address of the server")
	port := flag.Int("port", Defaults.Port, "port of the server")
	db_params := flag.String("dbparams", Defaults.DBParams, "parameters to open the database (see Gorm's doc)")
	PrintDefaults = flag.Bool("print-defaults", false, "print the default configuration file on stdout")
	flag.Parse()
	err := config.HandleConfFileFlag(*conf_file)
	if err != nil {
		return err
	}
	// You can override fields in the config file with flags.
	config.Host = *host
	config.Port = *port
	config.DBParams = *db_params
	err = config.SetDBType(*db_type)
	return err
}

func (config *Config) HandleConfFileFlag(path string) error {
	if path != "" {
		file, err := os.Open(path)
		if err != nil {
			return errors.New(fmt.Sprintf("Can't read file '%s'.", path))
		}

		err = config.Read(bufio.NewReader(file))
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to parse file '%s' (%s).", path, err))
		}
	}
	return nil
}

func (config *Config) SetDBType(db_type string) error {
	if !allowedDatabaseTypes[db_type] {
		return errors.New(fmt.Sprintf("Unknown database backend '%s'.", db_type))
	}
	config.DBType = db_type
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
