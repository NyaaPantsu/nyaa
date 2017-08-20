package tagsValidator

import (
	"github.com/NyaaPantsu/nyaa/config"
	"github.com/gin-gonic/gin"
)

// Check if a tag type exist and if it does, if the tag value is part of the defaults
func Check(tagType string, tag string) bool {
	if tagType == "" || tag == "" {
		return false
	}
	// We look for the tag type in config
	if tagConf := config.Get().Torrents.Tags.Types.Get(tagType); tagConf.Name != "" {
		// and then check that the value is in his defaults if defaults are set
		if len(tagConf.Defaults) > 0 && tagConf.Defaults[0] != "db" && !tagConf.Defaults.Contains(tag) {
			return false
		}
		return true
	}
	return false
}

// Bind a post request to tags
func Bind(c *gin.Context, keepEmpty bool) []CreateForm {
	var tags []CreateForm
	for _, tagConf := range config.Get().Torrents.Tags.Types {
		if value, ok := c.GetPostForm("tag_" + tagConf.Name); ok && (value != "" || keepEmpty) {
			if len(tagConf.Defaults) > 0 && tagConf.Defaults[0] != "db" && !tagConf.Defaults.Contains(value) && value != "" {
				continue
			}
			tags = append(tags, CreateForm{Tag: value, Type: tagConf.Name})
		}
	}
	return tags
}
