package tagsValidator

import "github.com/NyaaPantsu/nyaa/config"

// Check if a tag type exist and if it does, if the tag value is part of the defaults
func Check(tagType string, tag string) bool {
	for key, defaults := range config.Get().Torrents.Tags.Types {
		// We look for the tag type in config
		if key == tagType {
			// and then check that the value is in his defaults if defaults are set
			if len(defaults) > 0 && !defaults.Contains(tag) {
				return false
			}
			return true
		}
	}
	return false
}
