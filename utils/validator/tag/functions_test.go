package tagsValidator

import (
	"bytes"
	"net/http"
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", "..", config.Configpaths[0])
	config.Reload()
	return
}()

func TestCheck(t *testing.T) {
	tests := []struct {
		Tag      []string
		Expected bool
	}{
		{[]string{"", ""}, false},
		{[]string{"akuma06", ""}, false},
		{[]string{"videoquality", "full_hd"}, true},
		{[]string{"anidbid", ""}, false},
		{[]string{"anidbid", "20"}, true},
	}
	for _, test := range tests {
		b := Check(test.Tag[0], test.Tag[1])
		if b != test.Expected {
			t.Errorf("Error when checking tag type '%v', want '%t', got '%t'", test.Tag, test.Expected, b)
		}
	}
}

func TestBind(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		Test          string
		Expected      []CreateForm
		ExpectedEmpty []CreateForm
		Message       string
	}{
		{"", nil, nil, "Should return no CreateForm"},
		{"&", nil, nil, "Should return no CreateForm"},
		{"tag_j=d", nil, nil, "Should return no CreateForm for non-existing tagtype"},
		{"tag_videoquality=d", nil, nil, "Should return no CreateForm for non-existing options on a limited tagtype"},
		{"tag_videoquality=full_hd&tag_j=d", []CreateForm{{"full_hd", "videoquality"}}, []CreateForm{{"full_hd", "videoquality"}}, "Should filter out non-existing tagtypes and return only right ones"},
		{"tag_videoquality=", nil, []CreateForm{{"", "videoquality"}}, "Should return no CreateForm for keepEmpty=false and return the empty tagtype for keepEmpty=true"},
		{"tag_videoquality=full_hd&tag_anidbid=&tag_vndbid=", []CreateForm{{"full_hd", "videoquality"}}, []CreateForm{{"full_hd", "videoquality"}, {"", "anidbid"}, {"", "vndbid"}}, "Should keep empty tagtypes if keepEmpty is true and remove them otherwise"},
		{"tag_videoquality=full_hd&tag_anidbid=123", []CreateForm{{"full_hd", "videoquality"}, {"123", "anidbid"}}, []CreateForm{{"full_hd", "videoquality"}, {"123", "anidbid"}}, "Should be equal"},
	}
	for _, test := range tests {
		c := mockRequest(t, test.Test)
		b := Bind(c, false) // don't keep empty
		assert.Subset(b, test.Expected, test.Message)
		bEmpty := Bind(c, true) // keepEmpty
		assert.Subset(bEmpty, test.ExpectedEmpty, test.Message)
	}
}

func mockRequest(t *testing.T, params string) *gin.Context {
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(params))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		t.Fatal(err)
	}
	c := &gin.Context{Request: req}
	return c
}
