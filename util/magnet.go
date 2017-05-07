package util

import (
	"fmt"
	"net/url"
)

// convert a binary infohash to a magnet uri given a display name and tracker urls
func InfoHashToMagnet(ih string, name string, trackers ...url.URL) (str string) {
	str = fmt.Sprintf("magnet:?xt=urn:btih:%s", ih)
	if len(name) > 0 {
		str += fmt.Sprintf("&dn=%s", name)
	}
	for idx := range trackers {
		str += fmt.Sprintf("&tr=%s", trackers[idx].String())
	}
	return
}
