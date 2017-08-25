package search

import (
	"net/http"
	"net/http/httptest"
	"testing"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/stretchr/testify/assert"
)

func TestQuery_String(t *testing.T) {
	query := Query{}
	assert := assert.New(t)

	query.Append("xd")
	assert.Equal("xd", query.String(), "Should be equal")

	query = Query{}
	assert.Equal("", query.String(), "Should be equal")
}

func TestQuery_Append(t *testing.T) {
	var query Query
	assert := assert.New(t)

	query.Append("x")
	query.Append("d")
	assert.Equal("x d", query.String(), "Should be equal")

	query.Append("")
	assert.Equal("x d", query.String(), "Should be equal")

	query.Append("d")
	assert.Equal("x d d", query.String(), "Should be equal")

	query = Query{}
	query.Append("x", 1)
	query.Append("d = ?", 2)
	search, where := query.ToDBQuery()
	assert.Equal("x = ? AND d = ?", search, "Should be equal")
	assert.Equal([]interface{}{1, 2}, where, "Should be equal")

	query.Append("", 2)
	search, where = query.ToDBQuery()
	assert.Equal("x = ? AND d = ?", search, "Should be equal")
	assert.Equal([]interface{}{1, 2}, where, "Should be equal")

	query.Append("d = true")
	search, where = query.ToDBQuery()
	assert.Equal("x = ? AND d = ? AND d = true", search, "Should be equal")
	assert.Equal([]interface{}{1, 2}, where, "Should be equal")
}

func TestQuery_Prepend(t *testing.T) {
	var query Query
	assert := assert.New(t)

	query.Prepend("x")
	query.Prepend("d")
	assert.Equal("d x", query.String(), "Should be equal")

	query.Prepend("")
	assert.Equal("d x", query.String(), "Should be equal")

	query.Prepend("d")
	assert.Equal("d d x", query.String(), "Should be equal")

	query = Query{}
	query.Prepend("x = ?", 1)
	query.Prepend("d = ?", 2)
	search, where := query.ToDBQuery()
	assert.Equal("d = ? AND x = ?", search, "Should be equal")
	assert.Equal([]interface{}{2, 1}, where, "Should be equal")

	query.Prepend("", 2)
	search, where = query.ToDBQuery()
	assert.Equal("d = ? AND x = ?", search, "Should be equal")
	assert.Equal([]interface{}{2, 1}, where, "Should be equal")

	query.Prepend("d = true")
	search, where = query.ToDBQuery()
	assert.Equal("d = true AND d = ? AND x = ?", search, "Should be equal")
	assert.Equal([]interface{}{2, 1}, where, "Should be equal")
}

func TestQuery_ToESQuery(t *testing.T) {
	assert := assert.New(t)
	handler := http.NotFound
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	}))
	defer ts.Close()

	handler = func(w http.ResponseWriter, r *http.Request) {
		resp := `{}`

		w.Write([]byte(resp))
	}

	client, err := mockService(ts.URL)
	assert.NoError(err, "Couldn't load ES Client")

	torrentParam := &TorrentParam{
		NameLike: "x",
	}
	c := mockRequest(t, "/?order=true")
	query := torrentParam.toESQuery(c)
	search, err := query.ToESQuery(client)
	assert.NoError(err, "Couldn't load ES SearchService")
	assert.NotNil(search)
}

func mockService(url string) (*elastic.Client, error) {
	client, err := elastic.NewSimpleClient(elastic.SetURL(url))
	if err != nil {
		return nil, err
	}
	return client, nil
}
