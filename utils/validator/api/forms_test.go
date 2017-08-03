package apiValidator

import (
	"net/http"
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/gin-gonic/gin"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", "..", config.Configpaths[0])
	config.Reload()
	config.Get().DBType = models.SqliteType
	config.Get().DBParams = ":memory:?cache=shared&mode=memory"

	models.ORM, _ = models.GormInit(models.DefaultLogger)
	return
}()

func TestForms(t *testing.T) {
	fu := "http://nyaa.cat"
	em := "cop@cat.fe"
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &gin.Context{Request: req}
	messages := msg.GetMessages(c)
	tests := []struct {
		Form     CreateForm
		Expected bool
	}{
		{CreateForm{}, true},
		{CreateForm{"", "f", []string{fu}, []string{}, []string{}, "", "fedr", fu, fu, fu, fu, []string{em}, ""}, true},
		{CreateForm{"", "fed", []string{fu}, []string{}, []string{}, "", "fedr", "", "", fu, "", []string{em}, ""}, false},
		{CreateForm{"", "fed", []string{em}, []string{}, []string{}, "", "fedr", "", "", fu, "", []string{em}, ""}, true},
	}
	for _, test := range tests {
		messages.ClearAllErrors()
		validator.ValidateForm(test.Form, messages)
		b := messages.HasErrors()
		if b != test.Expected {
			t.Errorf("Error when validating CreateForm struct, want '%t', got '%t', please check validation arguments: %v\nand errors returned: %v", test.Expected, b, test.Form, messages.GetAllErrors())
		}
	}
}
