package router

import (
	"path"
	"testing"
)

// run before router/init.go:init()
var _ = func() (_ struct{}) {
	TemplateDir = path.Join("..", TemplateDir)
	return
}()

func TestReloadTemplates(t *testing.T) {
	ReloadTemplates()
}
