package torrentValidator

import (
	"net/url"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/anacrolix/torrent/metainfo"
)

// CheckTrackers : Check if there is good trackers in torrent
func CheckTrackers(t *metainfo.MetaInfo) []string {
	// TODO: move to runtime configuration
	var deadTrackers = config.Get().Torrents.Trackers.DeadTrackers

	trackerRet := []string{}
	tempList := metainfo.AnnounceList{}
	for _, group := range t.AnnounceList {
		var trackers []string
		for _, tracker := range group {
			urlTracker, err := url.ParseRequestURI(tracker)
			if err == nil {
				good := true
				for _, check := range deadTrackers {
					// the tracker is part of the deadtracker list
					// we don't keep it
					if strings.Contains(tracker, check) {
						good = false
						break // No need to continue the for loop
					}
				}
				if good {
					// We only keep the good trackers
					trackers = append(trackers, urlTracker.String())
				}
			}
		}
		if len(trackers) > 0 {
			tempList = append(tempList, trackers)
			trackerRet = append(trackerRet, trackers...)
		}
	}
	t.AnnounceList = tempList
	defaultTracker := config.Get().Torrents.Trackers.GetDefault()
	if defaultTracker != "" {
		t.Announce = defaultTracker
	}

	for _, key := range config.Get().Torrents.Trackers.NeededTrackers {
		inside := false
		if key < len(config.Get().Torrents.Trackers.Default) {
			tracker := config.Get().Torrents.Trackers.Default[key]
			for _, tr := range trackerRet {
				if tr == tracker {
					inside = true
				}
			}
			if !inside {
				trackerRet = append(trackerRet, tracker)
				t.AnnounceList = append(t.AnnounceList, []string{tracker})
			}
		}
	}
	return trackerRet
}
