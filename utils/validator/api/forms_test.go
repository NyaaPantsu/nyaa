package apiValidator

import (
	"net/http"
	"testing"

	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/gin-gonic/gin"
)

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
		{CreateForm{}, false},
		{CreateForm{"", "f", []string{fu}, []string{}, []string{}, "", "fedr", fu, fu, fu, fu, []string{em}, ""}, false},
		{CreateForm{"", "fed", []string{fu}, []string{}, []string{}, "", "fedr", fu, fu, fu, fu, []string{em}, ""}, true},
		{CreateForm{"", "fed", []string{em}, []string{}, []string{}, "", "fedr", fu, fu, fu, fu, []string{em}, ""}, false},
	}
	for _, test := range tests {
		messages.ClearAllErrors()
		validator.ValidateForm(test.Form, messages)
		b := messages.HasErrors()
		if b != !test.Expected {
			t.Errorf("Error when validating CreateForm struct, want '%t', got '%t', please check validation arguments: %v\nand errors returned: %v", !test.Expected, b, test.Form, messages.GetAllErrors())
		}
	}
}
