package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Config struct {
	Host   string `json: "host"`
	Port   int    `json: "port"`
	DBType string `json: "db_type"`
	// This will be directly passed to Gorm, and its internal
	// structure depends on the dialect for each db type
	DBParams string `json: "db_type"`
}

var Defaults = Config{"localhost", 9999, "sqlite3", "./nyaa.db"}

var allowedDatabaseTypes = map[string]bool{
	"sqlite3":  true,
	"postgres": true,
	"mysql":    true,
	"mssql":    true,
}

func NewConfig() *Config {
	var config Config
	config.Host = Defaults.Host
	config.Port = Defaults.Port
	config.DBType = Defaults.DBType
	config.DBParams = Defaults.DBParams
	return &config
}

type processFlags func() error

func (config *Config) BindFlags() processFlags {
	// This function returns a function which is to be used after
	// flag.Parse to check and copy the flags' values to the Config instance.

	conf_file := flag.String("conf", "", "path to the configuration file")
	db_type := flag.String("dbtype", Defaults.DBType, "database backend")
	host := flag.String("host", Defaults.Host, "binding address of the server")
	port := flag.Int("port", Defaults.Port, "port of the server")
	db_params := flag.String("dbparams", Defaults.DBParams, "parameters to open the database (see Gorm's doc)")

	return func() error {
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
	data = append(data, []byte("\n")...)
	if err != nil {
		return err
	}
	_, err = output.Write(data)
	return err
}
