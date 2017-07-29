package tags

import "github.com/NyaaPantsu/nyaa/models"

// DeleteAllType : Deletes all tag from a tag type and torrent ID
func DeleteAllType(tagType string, torrentID uint) error {
	if err := models.ORM.Model(&models.Tag{}).Where("torrent_id = ? AND type = ?", torrentID, tagType).Delete(&models.Tag{}).Error; err != nil {
		return err
	}

	return nil
}
