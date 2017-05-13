package uploadService

import (
	"strings"
)

func CheckTrackers(trackers []string) bool {
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
		"://announce.torrentsmd.com:8080"}

	var numGood int
	for _, t := range trackers {
		good := true
		for _, check := range deadTrackers {
			if strings.Contains(t, check) {
				good = false
			}
		}
		if good {
			numGood++
		}
	}
	return numGood > 0
}
