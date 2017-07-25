package apiValidator

import (
	"reflect"
	"testing"

	"github.com/NyaaPantsu/nyaa/models"
)

func TestCreateForm_Bind(t *testing.T) {
	fu := "http://nyaa.cat"
	em := "cop@cat.fe"
	tests := []CreateForm{
		CreateForm{},
		CreateForm{"", "f", []string{fu}, []string{}, []string{}, "", "fedr", fu, fu, fu, fu, []string{em}, ""},
		CreateForm{"", "fed", []string{fu}, []string{}, []string{}, "", "fedr", fu, fu, fu, fu, []string{em}, ""},
		CreateForm{"", "fed", []string{em}, []string{}, []string{}, "", "fedr", fu, fu, fu, fu, []string{em}, ""},
	}
	for _, test := range tests {
		d := &models.OauthClient{}
		nd := test.Bind(d)
		b := reflect.DeepEqual(nd, d)
		if !b {
			t.Errorf("The bind hasn't modified the original variable got '%v', wanted '%v'!", nd, d)
		}
	}

}
