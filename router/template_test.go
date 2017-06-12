package router

import (
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
)

// run before router/init.go:init()
var _ = func() (_ struct{}) {
	TemplateDir = path.Join("..", TemplateDir)
	config.ConfigPath = path.Join("..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", config.DefaultConfigPath)
	config.Parse()
	return
}()

func TestReloadTemplates(t *testing.T) {
	ReloadTemplates()
}
