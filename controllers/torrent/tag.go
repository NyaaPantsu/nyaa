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

// postTags is a function used by controllers to post tags
func postTags(c *gin.Context, torrent *models.Torrent, user *models.User) []models.Tag {
	var newTags []models.Tag
	messages := msg.GetMessages(c)
	// Bind already check if a tag is valid or not
	tagsForm := tagsValidator.Bind(c)
	if len(tagsForm) == 0 {
		messages.ErrorT(errors.New("wrong_tag_type"))
		return nil
	}
	for _, tagForm := range tagsForm {
		// We validate the tag CreateForm
		validator.ValidateForm(tagForm, messages)
		if messages.HasErrors() { // if there are errors while validating the form, we skip it
			continue
		}

		if len(user.Tags) == 0 { // In case we didn't call userLoadTags before calling this function
			user.LoadTags(torrent)
		}

		newTag := models.Tag{Tag: tagForm.Tag, Type: tagForm.Type, UserID: user.ID, TorrentID: torrent.ID, Weight: user.Pantsu}
		if user.Tags.Contains(newTag) {
			// We check if the user has already submitted this tag, if he has, we prevent him to vote twice for it
			log.Info("User has already tagged the type for the torrent")
			return nil
		}

		// Add a tag to the db and filter them if needed
		// Filtering means that we sum up all the tag with the same type/value
		// and compare the sum with the maximum value (of votes) a tag can have
		// if the value is greater than the maximum, we don't add the tag as a simple vote
		// we add it directly in torrent model as an accepted tag and remove other tags with the same type
		if tags.FilterOrCreate(&newTag, torrent, user) {
			newTags = append(newTags, newTag)
		}

	}
	return newTags
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

	// We add a tag if posted
	if user.ID > 0 {

		// We load tags for user so we can check if they have them
		user.LoadTags(torrent)
		tag := postTags(c, torrent, user)
		if _, ok := c.GetQuery("json"); ok {
			apiUtils.ResponseHandler(c, tag)
			return
		}
		if !messages.HasErrors() {
			c.Redirect(http.StatusSeeOther, fmt.Sprintf("/view/%d", id))
		}
	}
	tagForm := &tagsValidator.CreateForm{}
	c.Bind(tagForm)

	templates.Form(c, "/site/torrents/tag.jet.html", tagForm)
}

// AddTag is a controller to add a tag
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

	if c.Query("tag") != "" && user.ID > 0 {
		tagForm := &tagsValidator.CreateForm{c.Query("tag"), c.Query("type")}
		validator.ValidateForm(tagForm, messages)
		if !messages.HasErrors() {
			// We load tags for user and torrents
			user.LoadTags(torrent)
			postTags(c, torrent, user)
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

	if c.Query("tag") != "" && user.ID > 0 {
		tagForm := &tagsValidator.CreateForm{c.Query("tag"), c.Query("type")}

		validator.ValidateForm(tagForm, messages)

		if !messages.HasErrors() {
			for _, tag := range user.Tags {
				if tag.Tag == tagForm.Tag && tag.Type == tagForm.Type {
					_, err := tag.Delete()
					log.CheckError(err)
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
