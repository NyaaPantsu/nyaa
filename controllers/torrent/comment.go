package torrentController

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models/comments"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/sanitize"
	"github.com/gin-gonic/gin"
)

// PostCommentHandler : Controller for posting a comment
func PostCommentHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	torrent, err := torrents.FindByID(uint(id))
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	currentUser := router.GetUser(c)
	messages := msg.GetMessages(c)

	if currentUser.NeedsCaptcha() {
		userCaptcha := captcha.Extract(c)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddErrorT("errors", "bad_captcha")
		}
	}
	if currentUser.IsBanned() {
	    messages.AddErrorT("errors", "account_banned")
	}
	content := sanitize.Sanitize(c.PostForm("comment"), "comment")
	
	userID := currentUser.ID
	if c.PostForm("anonymous") == "true" {
		userID = 0
	}

	if strings.TrimSpace(content) == "" {
		messages.AddErrorT("errors", "comment_empty")
	}
	if len(content) > config.Get().CommentLength {
		messages.AddErrorT("errors", "comment_toolong")
	}
	if !messages.HasErrors() {

		_, err := comments.Create(content, torrent, userID)
		if err != nil {
			messages.Error(err)
		}
	}
	ViewHandler(c)
}
