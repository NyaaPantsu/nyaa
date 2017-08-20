package torrentValidator

import (
	"path"
	"testing"

	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/validator/tag"
	"github.com/stretchr/testify/assert"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", "..", config.Configpaths[0])
	config.Reload()
	config.Get().I18n.Directory = path.Join("..", "..", "..", config.Get().I18n.Directory)
	categories.InitCategories()
	return
}()

func TestValidateName(t *testing.T) {
	r := TorrentRequest{}
	tests := []struct {
		Test     string
		Expected error
	}{
		{"", errTorrentNameInvalid},
		{"something", nil},
		{"fjr*$é)à_'", nil},
	}
	for _, test := range tests {
		r.Name = test.Test
		err := r.ValidateName()
		if err != test.Expected {
			t.Errorf("Validation of torrent name for '%s' doesn't give the expected result, have '%v', wants '%v'", test.Test, err, test.Expected)
		}
	}
}

func TestValidateDescription(t *testing.T) {
	r := TorrentRequest{}
	config.Get().DescriptionLength = 5
	tests := []struct {
		Test     string
		Expected error
	}{
		{"", nil},
		{"something", errTorrentDescInvalid},
		{"fed", nil},
	}
	for _, test := range tests {
		r.Description = test.Test
		err := r.ValidateDescription()
		if err != test.Expected {
			t.Errorf("Validation of torrent description for '%s' doesn't give the expected result, have '%v', wants '%v'", test.Test, err, test.Expected)
		}
	}
}

func TestValidateMagnet(t *testing.T) {
	r := TorrentRequest{}
	config.Get().DescriptionLength = 5
	tests := []struct {
		Test     string
		Expected error
	}{
		{"", errTorrentMagnetInvalid},
		{"something", errTorrentMagnetInvalid},
		{"magnet:?xt=urn:btih:2BCE960D3CF61462DFB68C10C68D20ED56133BAD&dn=The+King%27s+Avatar+%5BQuan+Zhi+Gao+Shou%5D+-+07+-+%5B1080P%5D+-+Vostfr+-+Fastsub+-+BS.mkv&tr=http://nyaa.tracker.wf:7777/announce&tr=http://nyaa.tracker.wf:7777/announce&tr=udp://tracker.opentrackr.org:1337/announce&tr=http://anidex.moe:6969/announce&tr=http://tracker.anirena.com:80/announce&tr=http://tracker.t411.al&tr=udp://tracker.doko.moe:6969", nil},
	}
	for _, test := range tests {
		r.Magnet = test.Test
		err := r.ValidateMagnet()
		if err != test.Expected {
			t.Errorf("Validation of torrent magnet for '%s' doesn't give the expected result, have '%v', wants '%v'", test.Test, err, test.Expected)
		}
	}
}

func TestValidateWebsiteLink(t *testing.T) {
	r := TorrentRequest{}
	tests := []struct {
		Test     string
		Expected error
	}{
		{"", nil},
		{"something", errTorrentURIInvalid},
		{"https://kkk.cd", nil},
		{"http://kkk.cd/xd.?lol=eds", nil},
		{"ircs://kkk.cd", nil},
		{"irc://kkk.cd/lol", nil},
	}
	for _, test := range tests {
		r.WebsiteLink = test.Test
		err := r.ValidateWebsiteLink()
		if err != test.Expected {
			t.Errorf("Validation of torrent uri for '%s' doesn't give the expected result, have '%v', wants '%v'", test.Test, err, test.Expected)
		}
	}
}

func TestValidateHash(t *testing.T) {
	r := TorrentRequest{}
	tests := []struct {
		Test     string
		Expected error
	}{
		{"", errTorrentHashInvalid},
		{"something", errTorrentHashInvalid},
		{"2BCE960D3CF61462DFB68C10C68D20ED56133BAD", nil},
		{"2BCE960D3CF61462DFB68C10C68D20ED56133BADE", errTorrentHashInvalid},
	}
	for _, test := range tests {
		r.Infohash = test.Test
		err := r.ValidateHash()
		if err != test.Expected {
			t.Errorf("Validation of torrent hash for '%s' doesn't give the expected result, have '%v', wants '%v'", test.Test, err, test.Expected)
		}
	}
}

func TestExtractCategory(t *testing.T) {
	r := TorrentRequest{}
	tests := []struct {
		Test     string
		Expected error
	}{
		{"", errTorrentCatInvalid},
		{"something", errTorrentCatInvalid},
		{"33_5", errTorrentCatInvalid},
		{"3_", errTorrentCatInvalid},
		{"3_12", nil},
		{"_12", errTorrentCatInvalid},
	}
	for _, test := range tests {
		r.Category = test.Test
		err := r.ExtractCategory()
		if err != test.Expected {
			t.Errorf("Validation of torrent category for '%s' doesn't give the expected result, have '%v', wants '%v'", test.Test, err, test.Expected)
		}
	}
}

func TestExtractLanguage(t *testing.T) {
	var retriever publicSettings.UserRetriever // not required during initialization
	err := publicSettings.InitI18n(config.Get().I18n, retriever)
	if err != nil {
		t.Errorf("failed to initialize language translations: %v", err)
	}
	r := TorrentRequest{}
	tests := []struct {
		Test     string
		Expected error
	}{
		{"", nil},
		{"something,fr-fr", errTorrentLangInvalid},
		{"fr-ems", errTorrentLangInvalid},
		{"fr-fr,en-us", nil},
		{"es,fr", nil},
		{"es", nil},
		{"es-es", nil},
		{"ca-es", nil},
	}
	for _, test := range tests {
		r.Languages = strings.Split(test.Test, ",")
		err := r.ExtractLanguage()
		if err != test.Expected {
			t.Errorf("Validation of torrent language for '%s' doesn't give the expected result, have '%v', wants '%v'", test.Test, err, test.Expected)
		}
	}
}

func TestValidateTags(t *testing.T) {
	r := TorrentRequest{}
	assert := assert.New(t)
	tests := []struct {
		Test     TagsRequest
		Expected TagsRequest
	}{
		{TagsRequest{}, TagsRequest{}},
		{TagsRequest{{Tag: "", Type: ""}}, TagsRequest{}},
		{TagsRequest{{Tag: "xx", Type: "lol"}}, TagsRequest{{Tag: "xx", Type: "lol"}}},
		{TagsRequest{{Tag: "xx", Type: "lol"}, {Tag: "xxs", Type: "lol"}}, TagsRequest{{Tag: "xx", Type: "lol"}}},
	}
	for _, test := range tests {
		r.Tags = test.Test
		r.ValidateTags()
		assert.Equal(test.Expected, r.Tags, "Validation of torrent tags for '%v' doesn't give the expected result, have '%v', wants '%v'", test.Test, r.Tags, test.Expected)
	}

}

func TestTagsRequest_Bind(t *testing.T) {
	r := TorrentRequest{}
	assert := assert.New(t)
	tests := []struct {
		Test     *models.Torrent
		Expected TagsRequest
		Error    error
	}{
		{&models.Torrent{}, nil, nil},
		{&models.Torrent{AnidbID: 1}, TagsRequest{tagsValidator.CreateForm{Tag: "1", Type: "anidbid"}}, nil},
		{&models.Torrent{AnidbID: 1, VndbID: 2}, TagsRequest{tagsValidator.CreateForm{Tag: "1", Type: "anidbid"}, tagsValidator.CreateForm{Tag: "1", Type: "anidbid"}, tagsValidator.CreateForm{Tag: "2", Type: "vndbid"}}, nil},
		{&models.Torrent{AnidbID: 1, VndbID: 2, VgmdbID: 3, Dlsite: 4, AcceptedTags: "ddd,ddd,ddd", VideoQuality: "full_hd"}, TagsRequest{tagsValidator.CreateForm{Tag: "1", Type: "anidbid"}, tagsValidator.CreateForm{Tag: "1", Type: "anidbid"}, tagsValidator.CreateForm{Tag: "2", Type: "vndbid"}, tagsValidator.CreateForm{Tag: "1", Type: "anidbid"}, tagsValidator.CreateForm{Tag: "2", Type: "vndbid"}, tagsValidator.CreateForm{Tag: "4", Type: "dlsite"}, tagsValidator.CreateForm{Tag: "3", Type: "vgmdbid"}, tagsValidator.CreateForm{Tag: "full_hd", Type: "videoquality"}}, nil},
	}
	for _, test := range tests {
		err := r.Tags.Bind(test.Test)
		assert.Equal(test.Error, err, "Validation of torrent tags for '%v' doesn't give the expected error, have '%v', wants '%v'", test.Test, err, test.Error)
		assert.Equal(test.Expected, r.Tags, "Validation of torrent tags for '%v' doesn't give the expected result, have '%v', wants '%v'", test.Test, r.Tags, test.Expected)
	}
}
