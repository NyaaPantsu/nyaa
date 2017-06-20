package util

import (
	"path"
	"testing"

	"math"

	"github.com/NyaaPantsu/nyaa/config"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.ConfigPath = path.Join("..", config.ConfigPath)
	config.DefaultConfigPath = path.Join("..", config.DefaultConfigPath)
	config.Parse()
	return
}()

func TestFormatFilesize(t *testing.T) {
	format := FormatFilesize(0)
	if format != "0.0 B" {
		t.Fatalf("Format of 0 bytes gives %s, expected '0.0 B'", format)
	}
	format = FormatFilesize(int64(math.Exp2(0)))
	if format != "1.0 B" {
		t.Fatalf("Format of 1byte gives %s, expected '1.0 B'", format)
	}
	format = FormatFilesize(int64(math.Exp2(10)))
	if format != "1.0 KiB" {
		t.Fatalf("Format of 1024 bytes gives %s, expected '1.0 KiB'", format)
	}
}

func TestGetHostname(t *testing.T) {
	hostname := GetHostname("http://exemple.com")
	if hostname != "exemple.com" {
		t.Fatalf("Hostname gives %s, expected 'exemple.com'", hostname)
	}
	hostname = GetHostname("ircs://exemple.com/ddf/?kfkf=http://something.net&id=")
	if hostname != "exemple.com" {
		t.Fatalf("Hostname gives %s, expected 'exemple.com'", hostname)
	}
	hostname = GetHostname("")
	if hostname != "" {
		t.Fatalf("Hostname gives %s, expected ''", hostname)
	}
}

func TestInfoHashToMagnet(t *testing.T) {
	magnetExpected := "magnet:?xt=urn:btih:213d354dd354d534d&dn=Test&tr=udp://tracker.doko.moe:6969&tr=udp://tracker.zer0day.to:1337/announce"

	magnet := InfoHashToMagnet("213d354dd354d534d", "Test", "udp://tracker.doko.moe:6969", "udp://tracker.zer0day.to:1337/announce")
	if magnetExpected != magnet {
		t.Fatalf("Magnet URL parsed doesn't give the expected result, have this '%s', want this '%s'", magnet, magnetExpected)
	}
}

func TestSafe(t *testing.T) {
	safeString := map[string]string{
		"'":                  "&#39;",
		"&":                  "&amp;",
		"http://exemple.com": "http://exemple.com",
	}
	for key, val := range safeString {
		safe := Safe(key)
		if string(safe) != val {
			t.Errorf("Safe doesn't escape the right values, expected result %s, got %s", key, val)
		}
	}
}

func TestSafeText(t *testing.T) {
	safeString := map[string]string{
		"'":                                    "&#39;",
		"&":                                    "&amp;",
		"http://exemple.com":                   "http://exemple.com",
		"<em>test</em><script>lol();</script>": "&lt;em&gt;test&lt;/em&gt;&lt;script&gt;lol();&lt;/script&gt;",
	}
	for key, val := range safeString {
		safe := Safe(key)
		if string(safe) != val {
			t.Errorf("Safe doesn't escape the right values, expected result %s, got %s", key, val)
		}
	}
}
