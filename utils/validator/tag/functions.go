package tagsValidator

import "github.com/NyaaPantsu/nyaa/config"

func CheckTagType(tagType string) bool {
	return config.Get().Torrents.Tags.Types.Contains(tagType)
}
