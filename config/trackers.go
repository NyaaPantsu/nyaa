package config

// TODO: Update FAQ template to use this variable

// Trackers : Default trackers supported
var Trackers = []string{
	"udp://tracker.doko.moe:6969",
	"udp://tracker.coppersurfer.tk:6969",
	"udp://tracker.zer0day.to:1337/announce",
	"udp://tracker.leechers-paradise.org:6969",
	"udp://explodie.org:6969",
	"udp://tracker.opentrackr.org:1337",
	"udp://tracker.internetwarriors.net:1337/announce",
	"http://mgtracker.org:6969/announce",
	"udp://ipv6.leechers-paradise.org:6969/announce"}

// NeededTrackers : Array indexes of Trackers for needed tracker in a torrent file
var NeededTrackers = []int{
	0,
}
