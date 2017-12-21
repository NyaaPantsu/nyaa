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

	var trackerRet []string
	tempList := t.AnnounceList
	for kgroup, group := range tempList {
		for ktracker, tracker := range group {
			urlTracker, err := url.Parse(tracker)
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
					trackerRet = append(trackerRet, urlTracker.String())
				} else {
					// if the tracker is no good, we remove it from the group
					t.AnnounceList[kgroup] = append(t.AnnounceList[kgroup][:ktracker], t.AnnounceList[kgroup][ktracker+1:]...)
				}
			}
		}

		// We need to update the group of the trackers
		// if there is no good trackers in this group, we remove the group
		if len(t.AnnounceList[kgroup]) == 0 {
			t.AnnounceList = append(t.AnnounceList[:kgroup], t.AnnounceList[kgroup+1:]...)
		}
	}
	defaultTracker := config.Get().Torrents.Trackers.GetDefault()
	if defaultTracker != "" {
		t.Announce = defaultTracker
	}
	return trackerRet
}
