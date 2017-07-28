package format

import (
	"fmt"
	"net/url"
	"strings"
)

// FileSize : format file size
func FileSize(bytes int64) string {
	var unit string
	var value float64
	if bytes >= 1024*1024*1024*1024 {
		unit = "TiB"
		value = float64(bytes) / (1024 * 1024 * 1024 * 1024)
	} else if bytes >= 1024*1024*1024 {
		unit = "GiB"
		value = float64(bytes) / (1024 * 1024 * 1024)
	} else if bytes >= 1024*1024 {
		unit = "MiB"
		value = float64(bytes) / (1024 * 1024)
	} else if bytes >= 1024 {
		unit = "KiB"
		value = float64(bytes) / (1024)
	} else {
		unit = "B"
		value = float64(bytes)
	}
	return fmt.Sprintf("%.1f %s", value, unit)
}

// GetHostname : Returns the host of a URL, without any scheme or port number.
func GetHostname(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl
	}
	return u.Hostname()
}

// SplitNonEmpty is a special case of strings.Split
// which returns an empty slice if string is empty
func SplitNonEmpty(s, sep string) []string {
	if s == "" {
		return nil
	}

	return strings.Split(s, sep)
}
