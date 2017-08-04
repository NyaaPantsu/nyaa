package config

import (
	"reflect"
	"testing"

	"flag"

	"github.com/jinzhu/configor"
)

var _ = func() (_ struct{}) {
	Configpaths = []string{"config.yml", "default_config.yml"}
	Reload()
	return
}()

func TestBindFlags(t *testing.T) {
	oldConf := *Get()
	if !reflect.DeepEqual(oldConf, *Get()) {
		t.Error("Couldn't copy Conf *Config to make the test")
	}

	Get().Host = "something"
	if reflect.DeepEqual(oldConf, *Get()) {
		t.Error("Couldn't overwrtie property on Conf *Config")
	}

	Get().Host = oldConf.Host // Reset value
	if !reflect.DeepEqual(oldConf, *Get()) {
		t.Error("Couldn't reset value on Conf *Config")
	}

	flag.StringVar(&Get().DBType, "dbtype", "dfdf", "database backend")
	configor.Load(Get())
	if reflect.DeepEqual(oldConf, *Get()) {
		t.Error("Couldn't overwrtie property with flags on Conf *Config")
	}

	Get().DBType = oldConf.DBType // Reset
	if !reflect.DeepEqual(oldConf, *Get()) {
		t.Error("Couldn't reset value on Conf *Config")
	}

	flag.StringVar(&Get().Host, "host", Get().Host, "binding address of the server")
	configor.Load(Get())
	if !reflect.DeepEqual(oldConf, *Get()) {
		t.Error("Couldn't overwrtie by default property on Conf *Config")
	}
}
