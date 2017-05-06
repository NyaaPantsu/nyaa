package util

import (
	"encoding/hex"
	"fmt"
	"net/url"
)

// convert a binary infohash to a magnet uri given a display name and tracker urls
func InfoHashToMagnet(ih [20]byte, name string, trackers ...url.URL) (str string) {
	str = hex.EncodeToString(ih[:])
	str = fmt.Sprintf("magnet:?xt=urn:btih:%s", str)
	if len(name) > 0 {
		str += fmt.Sprintf("&dn=%s", name)
	}
	for idx := range trackers {
		str += fmt.Sprintf("&tr=%s", trackers[idx].String())
	}
	return
}
