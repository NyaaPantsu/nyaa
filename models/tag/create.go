package tags

import "github.com/NyaaPantsu/nyaa/models"

func Create(tag string, tagType string, torrent *models.Torrent, user *models.User) (*models.Tag, error) {
	newTag := &models.Tag{
		Tag:       tag,
		Type:      tagType,
		TorrentID: torrent.ID,
		UserID:    user.ID,
		Weight:    user.Pantsu,
		Accepted:  false,
	}

	if err := models.ORM.Create(newTag).Error; err != nil {
		return newTag, err
	}
	return newTag, nil
}
