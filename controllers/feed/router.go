package feedController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Any("/feed", RSSHandler)
	router.Get().Any("/feed/p/:page", RSSHandler)
	router.Get().Any("/feed/magnet", RSSMagnetHandler)
	router.Get().Any("/feed/magnet/p/:page", RSSMagnetHandler)
	router.Get().Any("/feed/torznab", RSSTorznabHandler)
	router.Get().Any("/feed/torznab/api", RSSTorznabHandler)
	router.Get().Any("/feed/torznab/p/:page", RSSTorznabHandler)
	router.Get().Any("/feed/eztv", RSSEztvHandler)
	router.Get().Any("/feed/eztv/p/:page", RSSEztvHandler)
}
