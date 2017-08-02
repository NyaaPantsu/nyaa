package format

import (
	"path"
	"testing"

	"math"

	"github.com/NyaaPantsu/nyaa/config"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", config.Configpaths[0])
	config.Reload()
	return
}()

func TestFileSize(t *testing.T) {
	format := FileSize(0)
	if format != "0.0 B" {
		t.Fatalf("Format of 0 bytes gives %s, expected '0.0 B'", format)
	}
	format = FileSize(int64(math.Exp2(0)))
	if format != "1.0 B" {
		t.Fatalf("Format of 1byte gives %s, expected '1.0 B'", format)
	}
	format = FileSize(int64(math.Exp2(10)))
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
