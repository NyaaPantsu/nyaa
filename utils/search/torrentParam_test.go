package search

import (
	"net/http"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTorrentParam_Identifier(t *testing.T) {
	torrentParam := &TorrentParam{}
	assert := assert.New(t)
	assert.Equal("MDAwMDAwMDAwMDBmYWxzZWZhbHNlZmFsc2VmYWxzZWZhbHNl", torrentParam.Identifier(), "It should be empty")
	torrentParam = &TorrentParam{
		NameLike: "test",
		NotNull:  "IS NULL",
		Hidden:   false,
	}
	assert.Equal("dGVzdElTIE5VTEwwMDAwMDAwMDAwMGZhbHNlZmFsc2VmYWxzZWZhbHNlZmFsc2U=", torrentParam.Identifier(), "It should be empty")
}

func TestTorrentParam_FromRequest(t *testing.T) {
	torrentParam := &TorrentParam{}
	assert := assert.New(t)
	defTorrent := &TorrentParam{Sort: 2, Max: maxType(config.Get().Navigation.TorrentsPerPage), NotNull: "date IS NOT NULL"}

	c := mockRequest(t, "/?")
	torrentParam.FromRequest(c)
	assert.Equal(defTorrent, torrentParam)

	c = mockRequest(t, "/?fromID=3&q=xx&c=_")
	torrentParam.FromRequest(c)
	defTorrent.FromID, defTorrent.NameLike = 3, "xx"
	defTorrent.NameSearch = "xx "
	assert.Equal(defTorrent, torrentParam)
}

func TestTorrentParam_Clone(t *testing.T) {
	torrentParam := TorrentParam{
		NameLike: "xx",
		ToDate:   DateFilter("2017-08-01"),
	}
	clone := torrentParam.Clone()
	assert.Equal(t, torrentParam, clone, "Should be equal")
}

// TODO implement this by asking a json response
func TestTorrentParam_FindES(t *testing.T) {

}

func TestTorrentParam_ToESQuery(t *testing.T) {
	assert := assert.New(t)
	c := mockRequest(t, "/?fromID=3")
	tests := []struct {
		Test     TorrentParam
		Expected string
	}{
		{TorrentParam{}, "!status:5"},
		{TorrentParam{NameLike: "lol"}, "!status:5"},
		{TorrentParam{NameLike: "lol", FromID: 12}, "!status:5 id:>12"},
		{TorrentParam{NameLike: "lol", FromID: 12, FromDate: DateFilter("2017-08-01"), ToDate: DateFilter("2017-08-05")}, "!status:5 id:>12 date: [2017-08-01 2017-08-05]"},
		{TorrentParam{NameLike: "lol", FromID: 12, ToDate: DateFilter("2017-08-05")}, "!status:5 id:>12 date: [* 2017-08-05]"},
		{TorrentParam{NameLike: "lol", FromID: 12, FromDate: DateFilter("2017-08-01")}, "!status:5 id:>12 date: [2017-08-01 *]"},
		{TorrentParam{NameLike: "lol", FromID: 12, Category: Categories{&Category{3, 12}}}, "(category: 3 AND sub_category: 12) !status:5 id:>12"},
		{TorrentParam{NameLike: "lol", FromID: 12, Category: Categories{&Category{3, 12}, &Category{3, 12}}}, "((category: 3 AND sub_category: 12) OR (category: 3 AND sub_category: 12)) !status:5 id:>12"},
	}

	for _, test := range tests {
		assert.Equal(test.Expected, test.Test.toESQuery(c).String())
	}
}

func TestParseUInt(t *testing.T) {
	assert := assert.New(t)

	c := mockRequest(t, "/?userID=3")
	userID := parseUInt(c, "userID")
	assert.Equal(uint32(3), userID, "Should be equal to 3")

	c = mockRequest(t, "/?userID=")
	userID = parseUInt(c, "userID")
	assert.Empty(userID, "Should be empty")

	c = mockRequest(t, "/?userID=lol")
	userID = parseUInt(c, "userID")
	assert.Empty(userID, "Should be empty")
}

func TestParseOrder(t *testing.T) {
	assert := assert.New(t)

	c := mockRequest(t, "/?order=true")
	order := parseOrder(c)
	assert.Equal(true, order, "Should be true")

	c = mockRequest(t, "/?order=")
	order = parseOrder(c)
	assert.Equal(false, order, "Should be false")

	c = mockRequest(t, "/?order=lol")
	order = parseOrder(c)
	assert.Equal(false, order, "Should be false")
}

func TestParseTorrentID(t *testing.T) {
	assert := assert.New(t)

	c := mockRequest(t, "/?fromID=3")
	fromID, torrentIDs := parseTorrentID(c)
	assert.Equal(uint32(3), fromID, "Should be equal to 3")
	assert.Empty(torrentIDs, "Should be empty")

	c = mockRequest(t, "/?fromID=")
	fromID, torrentIDs = parseTorrentID(c)
	assert.Empty(fromID, "Should be empty")
	assert.Empty(torrentIDs, "Should be empty")

	c = mockRequest(t, "/?fromID=lol")
	fromID, torrentIDs = parseTorrentID(c)
	assert.Empty(fromID, "Should be empty")
	assert.Empty(torrentIDs, "Should be empty")

	c = mockRequest(t, "/?fromID=3&id=2")
	fromID, torrentIDs = parseTorrentID(c)
	assert.Equal(uint32(3), fromID, "Should be equal to 3")
	assert.Equal([]uint32{2}, torrentIDs, "Should be 2")

	c = mockRequest(t, "/?fromID=3&id=2&id=3&id=4")
	fromID, torrentIDs = parseTorrentID(c)
	assert.Equal(uint32(3), fromID, "Should be equal to 3")
	assert.Equal([]uint32{2, 3, 4}, torrentIDs, "Should be 2,3,4")

	c = mockRequest(t, "/?fromID=&id=")
	fromID, torrentIDs = parseTorrentID(c)
	assert.Empty(fromID, "Should be empty")
	assert.Empty(torrentIDs, "Should be empty")

	c = mockRequest(t, "/?fromID=lol&id=lol")
	fromID, torrentIDs = parseTorrentID(c)
	assert.Empty(fromID, "Should be empty")
	assert.Empty(torrentIDs, "Should be empty")
}

func mockRequest(t *testing.T, url string) *gin.Context {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &gin.Context{Request: req}
	return c
}
