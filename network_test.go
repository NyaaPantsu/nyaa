package main

import (
	"net/http"
	"testing"

	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"sync"

	"github.com/NyaaPantsu/nyaa/config"
)

type MyHandler struct {
	sync.Mutex
	count int
}

func (h *MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var count int
	h.Lock()
	h.count++
	count = h.count
	h.Unlock()

	fmt.Fprintf(w, "Visitor count: %d.", count)
}

func TestCreateHTTPListener(t *testing.T) {
	// Set up server,
	conf := config.Get()
	testListener, err := CreateHTTPListener(conf)
	if err != nil {
		t.Error(err)
	}
	srv := httptest.NewUnstartedServer(&MyHandler{})
	srv.Listener = testListener
	defer srv.Close()

	srv.Start()
	for _, i := range []int{1, 2} {
		resp, err := http.Get(srv.URL)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Received non-200 response: %d\n", resp.StatusCode)
		}
		expected := fmt.Sprintf("Visitor count: %d.", i)
		actual, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}
		if expected != string(actual) {
			t.Errorf("Expected the message '%s', got '%s'\n", expected, string(actual))
		}
	}
}
