package modelHelper

import (
	"net/http"
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.ConfigPath = path.Join("..", "..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", "..", config.DefaultConfigPath)
	config.Parse()
	return
}()

type TestForm struct {
	DefaultVal  int    `form:"default" default:"3" notnull:"true"`
	ConfirmVal  string `form:"confirm" needed:"true" equalInput:"ConfirmeVal" len_min:"7" len_max:"8"`
	ConfirmeVal string `form:"confirme" needed:"true"`
}

func TestValidateForm(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	messages := msg.GetMessages(c)
	testform := TestForm{}
	ValidateForm(&testform, messages)
	if !messages.HasErrors() {
		t.Errorf("No errors when parsing empty invalid form: %v", testform)
	}
	messages.ClearAllErrors()
	testform.DefaultVal, testform.ConfirmVal, testform.ConfirmeVal = 1, "testingl", "testingl"
	ValidateForm(&testform, messages)
	if messages.HasErrors() {
		t.Errorf("Errors when parsing valid form %v\n with errors %v", testform, messages.GetAllErrors())
	}
	messages.ClearAllErrors()
	testform.ConfirmVal = "test"
	testform.ConfirmeVal = "test"
	ValidateForm(&testform, messages)
	if len(messages.GetErrors("confirm")) == 0 {
		t.Errorf("No errors on minimal length test when parsing invalid form: %v", testform)
	}
	messages.ClearAllErrors()
	testform.ConfirmVal, testform.ConfirmeVal = "testing", "testind"
	ValidateForm(&testform, messages)
	if len(messages.GetErrors("confirm")) == 0 {
		t.Errorf("No errors on equal test when parsing invalid form: %v", testform)
	}
	messages.ClearAllErrors()
	testform.ConfirmVal, testform.ConfirmeVal = "", "testing"
	ValidateForm(&testform, messages)
	if len(messages.GetErrors("confirm")) == 0 {
		t.Errorf("No errors on needed test when parsing invalid form: %v", testform)
	}
	messages.ClearAllErrors()
	testform.ConfirmVal, testform.ConfirmeVal = "azertyuid", "azertyuid"
	ValidateForm(&testform, messages)
	if len(messages.GetErrors("confirm")) == 0 {
		t.Errorf("No errors on maximal length test when parsing invalid form %v", testform)
	}
	messages.ClearAllErrors()
	testform.DefaultVal = 0
	ValidateForm(&testform, messages)
	if testform.DefaultVal == 0 {
		t.Errorf("Default value are not assigned on int with notnull specified: %v", testform)
	}
	messages.ClearAllErrors()
	testform.DefaultVal = 1
	ValidateForm(&testform, messages)
	if testform.DefaultVal != 1 {
		t.Errorf("Default value are assigned on int with non null value: %v", testform)
	}
	messages.ClearAllErrors()
}
