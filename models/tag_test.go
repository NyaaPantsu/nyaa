package models

import (
	"reflect"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/stretchr/testify/assert"
)

func TestTagsType(t *testing.T) {
	assert := assert.New(t)
	torrent := &Torrent{}
	for _, tagConf := range config.Get().Torrents.Tags.Types {
		assert.NotEmpty(tagConf.Name, "A name must be provided for a tag type")
		assert.NotEmpty(tagConf.Field, "You need to provide a field from the torrent Struct for '%s' tag type in config", tagConf.Name)
		if tagConf.Defaults != nil && len(tagConf.Defaults) == 0 {
			t.Log("For '%s', you provided a defaults attributes but it is empty. You should remove the attribute or fill it", tagConf.Name)
		}
		field := reflect.ValueOf(torrent).Elem().FieldByName(tagConf.Field)
		assert.True(field.IsValid(), "The field '%s' provided for '%s' doesn't exist", tagConf.Field, tagConf.Name)
		assert.Equal(reflect.String, field.Type().Kind(), "Only string tag types are supported, if you want to make an array do a string splitted by commas for '%s'", tagConf.Name)
	}
}
