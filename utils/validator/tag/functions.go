package tagsValidator

import (
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/gin-gonic/gin"
)

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

// Bind a post request to tags
func Bind(c *gin.Context) []CreateForm {
	var tags []CreateForm
	for key, defaults := range config.Get().Torrents.Tags.Types {
		if value := c.PostForm("tag_" + key); value != "" {
			if len(defaults) > 0 && defaults[0] != "db" && !defaults.Contains(value) {
				continue
			}
			tags = append(tags, CreateForm{Tag: value, Type: key})
		}
	}
	return tags
}
