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
	path := "config/default_config.yml"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("can't read file '%s'", path)
	}
	err = yaml.Unmarshal(data, DefaultConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	log.Printf("Config by default: %v", DefaultConfig)
	return DefaultConfig
}
