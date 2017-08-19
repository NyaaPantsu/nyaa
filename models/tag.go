package models

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/fatih/structs"
)

// Tag model for a torrent vote system
type Tag struct {
	TorrentID uint    `gorm:"column:torrent_id" json:"torrent_id"`
	UserID    uint    `gorm:"column:user_id" json:"user_id"`
	Tag       string  `gorm:"column:tag" json:"tag"`
	Type      string  `gorm:"column:type" json:"type"`
	Weight    float64 `gorm:"column:weight" json:"weight"`
	Total     float64 `gorm:"-" json:"total"`
	Accepted  bool    `gorm:"-" json:"accepted"`
}

// Update a tag
func (ta *Tag) Update() (int, error) {
	if ORM.Model(ta).UpdateColumn(ta.toMap()).Error != nil {
		return http.StatusInternalServerError, errors.New("Tag was not updated")
	}
	return http.StatusOK, nil
}

// Delete : delete a tag based on id
func (ta *Tag) Delete() (int, error) {
	if ORM.Where("tag = ? AND type = ? AND torrent_id = ? AND user_id = ?", ta.Tag, ta.Type, ta.TorrentID, ta.UserID).Delete(ta).Error != nil {
		return http.StatusInternalServerError, errors.New("tag_not_deleted")
	}

	return http.StatusOK, nil
}

// toMap : convert the model to a map of interface
func (ta *Tag) toMap() map[string]interface{} {
	return structs.Map(ta)
}

// GetName : get the translated name
func (ta *Tag) GetName() string {
	tagtype := config.Get().Torrents.Tags.Types.Get(ta.Type)
	if len(tagtype.Defaults) > 0 && tagtype.Defaults[0] != "db" {
		return "tagvalue_" + ta.Tag
	}
	return ta.Tag
}

type Tags []Tag

// Contains check if the tag map has the same tag in it (tag value + tag type)
func (ts Tags) Contains(tag Tag) bool {
	for _, ta := range ts {
		if ta.Tag == tag.Tag && ta.Type == tag.Type {
			return true
		}
	}
	return false
}

// HasType check if the tag map has the same tag type in it and returns its index
func (ts Tags) HasType(tagtype string) int {
	for i, ta := range ts {
		if ta.Type == tagtype {
			return i
		}
	}
	return -1
}

// Get in the tag map the same tag type
func (ts Tags) Get(tagtype string) Tag {
	i := ts.HasType(tagtype)
	if i == -1 {
		return Tag{}
	}
	return ts[i]
}

// DeleteType remove all tags from the map that have the same tag type in it
func (ts *Tags) DeleteType(tagtype string) {
	var newTs Tags
	for _, ta := range *ts {
		if ta.Type == tagtype {
			continue
		}
		newTs = append(newTs, ta)
	}
	ts = &newTs
}

// HasAccepted check if a tag has been accepted in the tags map
func (ts Tags) HasAccepted() bool {
	for _, tag := range ts {
		if tag.Accepted {
			return true
		}
	}
	return false
}

// Replace a tag in map of tags
func (ts Tags) Replace(index int, tag *Tag) {
	if index >= 0 && index < len(ts) {
		ts[index].Delete()
		ts[index] = *tag
		ts[index].Update()
	}
}

// ToJSON convert tags map to a json map and can exclud non accepted tags
func (ts Tags) ToJSON() string {
	toParse := ts

	b, err := json.Marshal(toParse)
	if err != nil {
		log.Infof("Couldn't parse to json the tags %v", toParse)
		return ""
	}
	return string(b)
}
