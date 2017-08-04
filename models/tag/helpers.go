package tags

import (
	"fmt"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/log"
)

// Filter check if a tag type has reached the maximal votes and removes the other tag of the same type
func Filter(tag string, tagType string, torrent *models.Torrent) bool {
	if torrent.ID == 0 || tagType == "" || tag == "" {
		return false
	}
	tagSum := models.Tag{}
	if err := models.ORM.Select("torrent_id, tag, type, accepted, SUM(weight) as total").Where("torrent_id = ? AND tag = ? AND type = ?", torrent.ID, tag, tagType).Group("type, tag").Find(&tagSum).Error; err == nil {
		fmt.Println(tagSum)
		if tagSum.Total > config.Get().Torrents.Tags.MaxWeight {
			tags, err := FindAll(tagType, torrent.ID)
			if err != nil {
				return false
			}
			for _, toDelete := range tags {
				user, _, err := users.FindRawByID(toDelete.UserID)
				if err != nil {
					log.CheckErrorWithMessage(err, "USER_NOT_FOUND: Couldn't update pantsu points!")
				}
				if user.ID > 0 {
					if toDelete.Tag == tag {
						user.IncreasePantsu()
					} else {
						user.DecreasePantsu()
					}
					user.Update()
				}
				toDelete.Delete()
			}
			/* err := DeleteAllType(tagType, torrent.ID) // We can also delete them in batch
			log.CheckError(err) */
			tagSum.Accepted = true
			tagSum.UserID = 0                                    // System ID
			tagSum.Weight = config.Get().Torrents.Tags.MaxWeight // Overriden to the maximal value
			models.ORM.Save(&tagSum)                             // We only add back the tag accepted
			callbackOnType(&tagSum, torrent)
			return true
		}
	}
	return false
}

func callbackOnType(tag *models.Tag, torrent *models.Torrent) {
	switch tag.Type {
	case "anidbid", "vndbid":
		if tag.Accepted && tag.TorrentID > 0 && torrent.ID > 0 {
			torrent.DbID = tag.Tag
			torrent.Update(false)
		}
	}
}
