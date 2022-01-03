package tags

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

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
	// In case we are adding a default tag, we need to split the tag. Since the input form ask for comma separated tags
	if tag.Type == config.Get().Torrents.Tags.Default && strings.Contains(tag.Tag, ",") {
		tagsToAdd := strings.Split(tag.Tag, ",")
		for _, tagToAdd := range tagsToAdd {
			tagCopy := *tag
			tagCopy.Tag = strings.TrimSpace(tagToAdd)
			FilterOrCreate(&tagCopy, torrent, currentUser)
		}
		return true
	}
	if tag.Type != config.Get().Torrents.Tags.Default {
		tagConf := config.Get().Torrents.Tags.Types.Get(tag.Type)
		if tagConf.Name == "" {
			return false
		}
		var oldValue = fmt.Sprint(reflect.ValueOf(torrent).Elem().FieldByName(tagConf.Field).Interface())
		// If the tag is already accepted in torrent, don't need to create it again or modify it
		if oldValue == tag.Tag || (oldValue == "0" && tag.Tag == "") || tag.Tag == "0" {
			return true
		}
	}
	tagSum := models.Tag{}
	if !tag.Accepted {
		if torrent.ID == 0 { // We can't search tags in dv for non existing torrent
			return false
		}
		// Here we only sum the tags of the same type, same value for the torrent specified, we don't handle errors since if no tags found, it returns an error
		models.ORM.Select("torrent_id, tag, type, SUM(weight) as total").Where("torrent_id = ? AND tag = ? AND type = ?", torrent.ID, tag.Tag, tag.Type).Group("type, tag").Find(&tagSum)
	} else {
		// if the tag given is an accepted one, tagSum is the tag given
		tagSum = *tag
	}

	// if the total sum is equal or lesser than the maximum set in config
	if (tagSum.Total+tag.Weight) <= config.Get().Torrents.Tags.MaxWeight && !tagSum.Accepted {
		// We only add the tag
		_, err := New(tag, torrent)
		if err != nil {
			log.CheckErrorWithMessage(err, "TAG_NOT_CREATED: Couldn't create tag")
			return false
		}
		return true
	}
	// if the total sum is greater than the maximum set in config
	// or if the tag is accepted (when owner uploads/edit torrent details)
	// we can select all the tags of the same type
	if torrent.ID > 0 { // We can't filter tags for non existing torrent
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
	}
	callbackOnType(&tagSum, torrent) // This callback will make different action depending on the tag type
	return true
}

/// callbackOnType is a function which will perform different action depending on the tag type
func callbackOnType(tag *models.Tag, torrent *models.Torrent) {
	switch tag.Type {
	case config.Get().Torrents.Tags.Default:
		// We need to check that the tag doesn't actually exist
		tags := strings.Split(torrent.AcceptedTags, ",")
		for _, tagComp := range tags {
			// if it exists, we return
			if tag.Tag == tagComp {
				return
			}
		}
		// if it doesn't exist
		// We check if the torrent has already accepted tags and that the tag is not empty
		if torrent.AcceptedTags != "" && tag.Tag != "" {
			// if yes we append to it a comma before inserting the tag
			torrent.AcceptedTags += ","
		}
		// We finally add the tag to the column
		torrent.AcceptedTags += tag.Tag
	case "anidbid", "vndbid", "vgmdbid":
		u64, _ := strconv.ParseUint(trimNonNumbers(tag.Tag), 10, 32)
		// TODO: Perform a check that anidbid is in anidb database
		tagConf := config.Get().Torrents.Tags.Types.Get(tag.Type)
		reflect.ValueOf(torrent).Elem().FieldByName(tagConf.Field).SetUint(u64)
	case "dlsite":
		// Since DLSite has a particular format, we don't use the default behavior but we check the format
		tagConf := config.Get().Torrents.Tags.Types.Get(tag.Type)
		if len(tag.Tag) != 8 { // eg RJ001001
			return
		}
		var validID = regexp.MustCompile(`^[A-Za-z]{2}[0-9]{6}$`)
		if validID.MatchString(tag.Tag) {
			reflect.ValueOf(torrent).Elem().FieldByName(tagConf.Field).SetString(tag.Tag)
		}
	default:
		// Some tag type can have default values that you have to choose from
		// We, here, check that the tag is one of them
		tagConf := config.Get().Torrents.Tags.Types.Get(tag.Type)
		// We look for the tag type in config
		if tagConf.Name != "" {
			// and then check that the value is in his defaults if defaults are set
			if len(tagConf.Defaults) > 0 && tagConf.Defaults[0] != "db" && tag.Tag != "" && !tagConf.Defaults.Contains(tag.Tag) {
				// if not we return the function
				return
			}
			// We overwrite the tag type in the torrent model
			reflect.ValueOf(torrent).Elem().FieldByName(tagConf.Field).SetString(tag.Tag)
		}
	}
}

func trimNonNumbers(source string) string {
	output := ""

	for i := 0; i < len(source); i++ {
		if source[i] < 58 {
			output += source[i : i+1]
		}
	}
	return output
}
