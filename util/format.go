package util

import (
	"fmt"
)

func FormatFilesize(bytes int64) string {
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

func FormatFilesize2(bytes int64) string {
	if bytes == 0 { // this is what gorm returns for NULL
		return "Unknown"
	}
	return FormatFilesize(bytes)
}
