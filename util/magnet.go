package util

import (
	"fmt"
)

// InfoHashToMagnet : convert a binary infohash to a magnet uri given a display name and tracker urls
func InfoHashToMagnet(ih string, name string, trackers ...string) (str string) {
	str = fmt.Sprintf("magnet:?xt=urn:btih:%s", ih)
	if len(name) > 0 {
		str += fmt.Sprintf("&dn=%s", name)
	}
	for idx := range trackers {
		str += fmt.Sprintf("&tr=%s", trackers[idx])
	}
	return
}
