package tags

import (
	"reflect"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/log"
)

// FilterOrCreate check if a tag type has reached the maximal votes and removes the other tag of the same type
// Filtering means that we sum up all the tag with the same type/value
// and compare the sum with the maximum value (of votes) a tag can have
// if the value is greater than the maximum, we don't add the tag as a simple vote
// we add it directly in torrent model as an accepted tag and remove other tags with the same type
// This function return true if it has added/filtered the tags and false if errors were encountered
func FilterOrCreate(tag *models.Tag, torrent *models.Torrent, currentUser *models.User) bool {
	if torrent.ID == 0 || tag.Type == "" || tag.Tag == "" {
		return false
	}
	tagSum := models.Tag{}
	if !tag.Accepted {
		if models.ORM.Select("torrent_id, tag, type, SUM(weight) as total").Where("torrent_id = ? AND tag = ? AND type = ?", torrent.ID, tag.Tag, tag.Type).Group("type, tag").Find(&tagSum).Error != nil {
			return false
		}
	} else {
		// if the tag given is an accepted one, tagSum is the tag given
		tagSum = *tag
	}

	// if the total sum is equal or lesser than the maximum set in config
	if (tagSum.Total+tag.Weight) <= config.Get().Torrents.Tags.MaxWeight && !tagSum.Accepted {
		// We only add the tag
		_, err := New(tag, torrent)
		if err != nil {
			log.CheckErrorWithMessage(err, "TAG_NOT_CREATED: Couldn't create tag: %s")
			return false
		}
		return true
	}
	// if the total sum is greater than the maximum set in config
	// or if the tag is accepted (when owner uploads/edit torrent details)
	// we can select all the tags of the same type
	tags, err := FindAll(tag.Type, torrent.ID)
	if err != nil {
		return false
	}
	// delete them and decrease/increase pantsu of the users who have voted wrongly/rightly
	for _, toDelete := range tags {
		// find the user who has voted for the tag
		user, _, err := users.FindRawByID(toDelete.UserID)
		if err != nil {
			log.CheckErrorWithMessage(err, "USER_NOT_FOUND: Couldn't update pantsu points!")
		}
		// if the user exist
		if user.ID > 0 {
			// and if he has voted for the right tag value
			if toDelete.Tag == tag.Tag {
				// we increase his pantsu
				user.IncreasePantsu()
			} else {
				// else we decrease them
				user.DecreasePantsu()
			}
			// and finally we update the user so the changes take effect
			user.Update()
		}
		// Not forget to delete the tag
		toDelete.Delete()
	}
	if currentUser.ID > 0 {
		// Same as for the current user, we increase his pantsus and update
		currentUser.IncreasePantsu()
		currentUser.Update() // we do it here since we didn't save the tag previously and didn't increase his pantsu
	}
	callbackOnType(&tagSum, torrent) // This callback will make different action depending on the tag type
	return true
}

/// callbackOnType is a function which will perform different action depending on the tag type
func callbackOnType(tag *models.Tag, torrent *models.Torrent) {
	switch tag.Type {
	case config.Get().Torrents.Tags.Default:
		if tag.TorrentID > 0 && torrent.ID > 0 {
			// We check if the torrent has already accepted tags
			if torrent.AcceptedTags != "" {
				// if yes we append to it a comma before inserting the tag
				torrent.AcceptedTags += ","
			}
			// We finally add the tag to the column
			torrent.AcceptedTags += tag.Tag
			// and update the torrent
			torrent.Update(false)
		}
	case "anidbid":
		// TODO: Perform a check that anidbid is in anidb database
		if tag.TorrentID > 0 && torrent.ID > 0 {
			torrent.AnidbID = tag.Tag
			// and update the torrent
			torrent.Update(false)
		}
	default:
		// Some tag type can have default values that you have to choose from
		// We, here, check that the tag is one of them
		for _, tagConf := range config.Get().Torrents.Tags.Types {
			// We look for the tag type in config
			if tagConf.Name == tag.Type {
				// and then check that the value is in his defaults if defaults are set
				if len(tagConf.Defaults) > 0 && tagConf.Defaults[0] != "db" && !tagConf.Defaults.Contains(tag.Tag) {
					// if not we return the function
					return
				}
				// We overwrite the tag type in the torrent model
				reflect.ValueOf(torrent).Elem().FieldByName(tagConf.Field).SetString(tag.Tag)
				// if it's good, we break of the loop
				break
			}
		}
		torrent.Update(false)
	}
}
