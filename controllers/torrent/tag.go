package torrentController

import (
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/tags"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/tags"
	"github.com/gin-gonic/gin"
)

func postTag(c *gin.Context, torrent *models.Torrent, user *models.User) {
	messages := msg.GetMessages(c)
	tagForm := &tagsValidator.CreateForm{}

	c.Bind(tagForm)
	validator.ValidateForm(tagForm, messages)

	for _, tag := range user.Tags {
		if tag.Tag == tagForm.Tag {
			return // already a tag by the user, don't add one more
		}
	}

	tags.Create(tagForm.Tag, tagForm.Type, torrent, user) // Add a tag to the db
	tags.Filter(tagForm.Tag, tagForm.Type, torrent.ID)    // Check if we have a tag reaching the maximum weight, if yes, deletes every tag and add only the one accepted
}
