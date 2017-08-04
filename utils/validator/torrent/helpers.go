package torrentValidator

import (
	"net/url"
	"strings"
)

// CheckTrackers : Check if there is good trackers in torrent
func CheckTrackers(trackers []string) []string {
	// TODO: move to runtime configuration
	var deadTrackers = []string{ // substring matches!
		"://open.nyaatorrents.info:6544",
		"://tracker.openbittorrent.com:80",
		"://tracker.publicbt.com:80",
		"://stats.anisource.net:2710",
		"://exodus.desync.com",
		"://open.demonii.com:1337",
		"://tracker.istole.it:80",
		"://tracker.ccc.de:80",
		"://bt2.careland.com.cn:6969",
		"://announce.torrentsmd.com:8080",
		"://open.demonii.com:1337",
		"://tracker.btcake.com",
		"://tracker.prq.to",
		"://bt.rghost.net"}

	var trackerRet []string
	for _, t := range trackers {
		urlTracker, err := url.Parse(t)
		if err == nil {
			good := true
			for _, check := range deadTrackers {
				if strings.Contains(t, check) {
					good = false
					break // No need to continue the for loop
				}
			}
			if good {
				trackerRet = append(trackerRet, urlTracker.String())
			}
		}
	}
	return trackerRet
}
