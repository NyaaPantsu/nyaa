package validator

import (
	"net/http"
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/gin-gonic/gin"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", config.Configpaths[0])
	config.Reload()
	return
}()

type TestForm struct {
	DefaultVal  int    `validate:"default=3,required"`
	ConfirmVal  string `validate:"eqfield=ConfirmeVal,min=7,max=8,required"`
	ConfirmeVal string `validate:"required"`
	Optional    string
}

func TestValidateForm(t *testing.T) {
	t.Parallel()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &gin.Context{Request: req}
	messages := msg.GetMessages(c)
	testform := TestForm{}
	ValidateForm(&testform, messages)
	if !messages.HasErrors() {
		t.Errorf("No errors when parsing empty invalid form: %v", testform)
	}
	messages.ClearAllErrors()
	testform.DefaultVal, testform.ConfirmVal, testform.ConfirmeVal, testform.Optional = 1, "testingl", "testingl", "test"
	ValidateForm(&testform, messages)
	if messages.HasErrors() {
		t.Errorf("Errors when parsing valid form %v\n with errors %v", testform, messages.GetAllErrors())
	}
	messages.ClearAllErrors()
	testform.Optional = ""
	ValidateForm(&testform, messages)
	if messages.HasErrors() {
		t.Errorf("Errors when testing an empty optional field in form %v\n with errors %v", testform, messages.GetAllErrors())
	}
	messages.ClearAllErrors()
	testform.ConfirmVal = "test"
	testform.ConfirmeVal = "test"
	ValidateForm(&testform, messages)
	if len(messages.GetErrors("ConfirmVal")) == 0 {
		t.Errorf("No errors on minimal length test when parsing invalid form: %v", testform)
	}
	messages.ClearAllErrors()
	testform.ConfirmVal, testform.ConfirmeVal = "testing", "testind"
	ValidateForm(&testform, messages)
	if len(messages.GetErrors("ConfirmVal")) == 0 {
		t.Errorf("No errors on equal test when parsing invalid form: %v", testform)
	}
	messages.ClearAllErrors()
	testform.ConfirmVal, testform.ConfirmeVal = "", "testing"
	ValidateForm(&testform, messages)
	if len(messages.GetErrors("ConfirmVal")) == 0 {
		t.Errorf("No errors on needed test when parsing invalid form: %v", testform)
	}
	messages.ClearAllErrors()
	testform.ConfirmVal, testform.ConfirmeVal = "azertyuid", "azertyuid"
	ValidateForm(&testform, messages)
	if len(messages.GetErrors("ConfirmVal")) == 0 {
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
func TestIsUTFLetterNumeric(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", true},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", true},
		{"abc〩", true},
		{"abc", true},
		{"소주", true},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", true},
		{"소", true},
		{"달기&Co.", false},
		{"〩Hours", true},
		{"\ufff0", false},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", true},
		{"〥〩", true},
		{"모자", true},
		{"ix", true},
		{"۳۵۶۰", true},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}
	for _, test := range tests {
		actual := IsUTFLetterNumeric(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsUTFLetterNumeric(%q) to be %v, got %v", test.param, test.expected, actual)
		}
	}
}
