package models

import (
	"errors"
	"net/http"

	"github.com/fatih/structs"
)

// Tag model for a torrent vote system
type Tag struct {
	TorrentID uint    `gorm:"column:torrent_id" json:"torrent_id"`
	UserID    uint    `gorm:"column:user_id" json:"user_id"`
	Tag       string  `gorm:"column:tag" json:"tag"`
	Type      string  `gorm:"column:type" json:"type"`
	Weight    float64 `gorm:"column:weight" json:"weight"`
	Accepted  bool    `gorm:"column:accepted" json:"accepted"`
	Total     float64 `gorm:"-" json:"total"`
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

// HasAccepted check if a tag has been accepted in the tags map
func (ts Tags) HasAccepted() bool {
	for _, tag := range ts {
		if tag.Accepted {
			return true
		}
	}
	return false
}
