package torrentController

import (
	"errors"

	"fmt"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/tag"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/api"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/tag"
	"github.com/gin-gonic/gin"
)

// postTag is a function used by controllers to post a tag
func postTag(c *gin.Context, torrent *models.Torrent, user *models.User) *models.Tag {
	messages := msg.GetMessages(c)
	tagForm := &tagsValidator.CreateForm{}

	c.Bind(tagForm)
	validator.ValidateForm(tagForm, messages)

	// We check that the tag type sent is one enabled in config.yml
	if !tagsValidator.CheckTagType(tagForm.Type) {
		messages.ErrorT(errors.New("wrong_tag_type"))
		return nil
	}

	if user.Tags.Contains(models.Tag{Tag: tagForm.Tag, Type: tagForm.Type}) {
		log.Info("User has already tagged the type for the torrent")
		return nil
	}

	tag, _ := tags.Create(tagForm.Tag, tagForm.Type, torrent, user) // Add a tag to the db
	tags.Filter(tagForm.Tag, tagForm.Type, torrent)                 // Check if we have a tag reaching the maximum weight, if yes, deletes every tag and add only the one accepted
	return tag
}

// ViewFormTag is a controller displaying a form to add a tag to a torrent
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
		tag := postTag(c, torrent, user)
		if !messages.HasErrors() {
			if _, ok := c.GetQuery("json"); ok {
				apiUtils.ResponseHandler(c, tag)
				return
			}
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/view/%d", id))
		}
		if _, ok := c.GetQuery("json"); ok {
			apiUtils.ResponseHandler(c)
			return
		}
	}
	tagForm := &tagsValidator.CreateForm{}
	c.Bind(tagForm)

	templates.Form(c, "/site/torrents/tag.jet.html", tagForm)
}

// AddTag is a controller to add a
func AddTag(c *gin.Context) {
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

	if c.Query("tag") != "" && user.ID > 0 {
		tagForm := &tagsValidator.CreateForm{c.Query("tag"), c.Query("type")}

		validator.ValidateForm(tagForm, messages)

		if !messages.HasErrors() {
			tag := postTag(c, torrent, user)
			if !messages.HasErrors() {
				if _, ok := c.GetQuery("json"); ok {
					apiUtils.ResponseHandler(c, tag)
					return
				}
			}
		}
	}
	if _, ok := c.GetQuery("json"); ok {
		apiUtils.ResponseHandler(c)
		return
	}
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/view/%d", id))
}

// DeleteTag is a controller to delete a user tag
func DeleteTag(c *gin.Context) {
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

	if c.Query("tag") != "" && user.ID > 0 {
		tagForm := &tagsValidator.CreateForm{c.Query("tag"), c.Query("type")}

		validator.ValidateForm(tagForm, messages)

		if !messages.HasErrors() {
			for _, tag := range user.Tags {
				if tag.Tag == tagForm.Tag && tag.Type == tagForm.Type {
					tagRef := &models.Tag{tag.TorrentID, tag.UserID, tag.Tag, tag.Type, tag.Weight, tag.Accepted, tag.Total}
					_, err := tag.Delete()
					if err != nil {
						log.CheckError(err)
						break
					}
					if _, ok := c.GetQuery("json"); ok {
						apiUtils.ResponseHandler(c, tagRef)
						return
					}
					break
				}
			}
		}
	}
	if _, ok := c.GetQuery("json"); ok {
		apiUtils.ResponseHandler(c)
		return
	}
	c.Redirect(http.StatusSeeOther, fmt.Sprintf("/view/%d", id))
}
