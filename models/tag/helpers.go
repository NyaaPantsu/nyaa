package tags

import (
	"fmt"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/log"
)

func Filter(tag string, tagType string, torrentID uint) bool {
	if torrentID == 0 || tagType == "" || tag == "" {
		return false
	}
	tagSum := models.Tag{}
	if err := models.ORM.Select("torrent_id, tag, type, accepted, SUM(weight) as total").Where("torrent_id = ? AND tag = ? AND type = ?", torrentID, tag, tagType).Group("type, tag").Find(&tagSum).Error; err == nil {
		fmt.Println(tagSum)
		if tagSum.Total > config.Get().Torrents.Tags.MaxWeight {
			tags, err := FindAll(tagType, torrentID)
			if err != nil {
				return false
			}
			for _, toDelete := range tags {
				user, _, err := users.FindRawByID(toDelete.UserID)
				if err != nil {
					log.CheckErrorWithMessage(err, "USER_NOT_FOUND: Couldn't update pantsu points!")
				}
				if toDelete.Tag == tag {
					user.IncreasePantsu()
				} else {
					user.DecreasePantsu()
				}
				user.Update()
				toDelete.Delete()
			}
			/* err := DeleteAllType(tagType, torrentID) // We can also delete them in batch
			log.CheckError(err) */
			tagSum.Accepted = true
			tagSum.UserID = 0                                    // System ID
			tagSum.Weight = config.Get().Torrents.Tags.MaxWeight // Overriden to the maximal value
			models.ORM.Save(&tagSum)                             // We only add back the tag accepted
			return true
		}
	}
	return false
}
