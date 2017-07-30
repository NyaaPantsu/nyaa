package torrentController

import (
	"errors"

	"fmt"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/tags"
	"github.com/gin-gonic/gin"
	"github.com/NyaaPantsu/nyaa/models/tag"
)

func postTag(c *gin.Context, torrent *models.Torrent, user *models.User) {
	messages := msg.GetMessages(c)
	tagForm := &tagsValidator.CreateForm{}

	c.Bind(tagForm)
	validator.ValidateForm(tagForm, messages)

	// We check that the tag type sent is one enabled in config.yml
	if !tagsValidator.CheckTagType(tagForm.Type) {
		messages.ErrorT(errors.New("wrong_tag_type"))
		return
	}

	for _, tag := range user.Tags {
		if tag.Tag == tagForm.Tag {
			return // already a tag by the user, don't add one more
		}
	}

	tags.Create(tagForm.Tag, tagForm.Type, torrent, user) // Add a tag to the db
	tags.Filter(tagForm.Tag, tagForm.Type, torrent.ID)    // Check if we have a tag reaching the maximum weight, if yes, deletes every tag and add only the one accepted
}

func ViewFormTag(c *gin.Context) {
	messages := msg.GetMessages(c)
	user := router.GetUser(c)
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	// Retrieve the torrent
	torrent, err := torrents.FindByID(uint(id))

	// If torrent not found, display 404
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// We load tags for user and torrents
	user.LoadTags(torrent)
	torrent.LoadTags()

	// We add a tag if posted
	if c.PostForm("tag") != "" && user.ID > 0 {
		postTag(c, torrent, user)
		if !messages.HasErrors() {
			if _, ok := c.GetQuery("json"); ok {
				c.JSON(http.StatusOK, struct {
					Ok bool
				}{true})
				return
			}
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/view/%d", id))
		}
		if _, ok := c.GetQuery("json"); ok {
			c.JSON(http.StatusOK, struct {
				Ok bool
			}{false})
			return
		}
	}
	tagForm := &tagsValidator.CreateForm{}
	c.Bind(tagForm)

	templates.Form(c, "/site/torrents/tag.jet.html", tagForm)
}
