package config

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"
)

// DefaultConfig : Default configuration
var DefaultConfig *Config

func getDefaultConfig() *Config {
	DefaultConfig = &Config{}
	data, err := ioutil.ReadFile(DefaultConfigPath)
	if err != nil {
		log.Printf("can't read file '%s'", DefaultConfigPath)
	}
	err = yaml.Unmarshal(data, DefaultConfig)
	if err != nil {
		log.Printf("error: %v", err)
	}
	return DefaultConfig
}
