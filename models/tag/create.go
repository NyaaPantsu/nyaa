package tags

import (
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cache"
	"github.com/pkg/errors"
)

// Create a new tag based on string inputs
func Create(tag string, tagType string, torrent *models.Torrent, user *models.User) (*models.Tag, error) {
	newTag := &models.Tag{
		Tag:       tag,
		Type:      tagType,
		TorrentID: torrent.ID,
		UserID:    user.ID,
		Weight:    user.Pantsu,
	}
	return New(newTag, torrent)
}

// New is the low level functions that actually create a tag in db
func New(tag *models.Tag, torrent *models.Torrent) (*models.Tag, error) {
	if torrent.ID == 0 {
		return tag, errors.New("Can't add a tag to no torrents")
	}
	tag.TorrentID = torrent.ID
	if err := models.ORM.Create(tag).Error; err != nil {
		return tag, err
	}
	cache.C.Delete(torrent.Identifier())
	return tag, nil
}
