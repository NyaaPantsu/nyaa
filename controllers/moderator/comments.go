package moderatorController

import (
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/models/comments"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/gin-gonic/gin"
)

// CommentsListPanel : Controller for listing comments, can accept pages and userID
func CommentsListPanel(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	userid := c.Query("userid")
	username := c.Query("user")
	var err error
	messages := msg.GetMessages(c)
	deleted := c.Request.URL.Query()["deleted"]
	if deleted != nil {
		messages.AddInfoTf("infos", "comment_deleted")
	}
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	var conditions string
	var values []interface{}
	searchForm := templates.NewSearchForm(c)
	// if there is a username in url
	if username != "" {
		conditions = "user = ?"
		values = append(values, username)
		searchForm.UserName = username
	// else we look if there is a userid
	} else if userid != "" {
		id, err := strconv.Atoi(userid)
		if err == nil {
			conditions = "user_id = ?"
			values = append(values, id)
			searchForm.UserID = uint32(id)
		}
	}


	comments, nbComments := comments.FindAll(offset, (pagenum-1)*offset, conditions, values...)
	nav := templates.Navigation{nbComments, offset, pagenum, "mod/comments/p"}
	templates.ModelList(c, "admin/commentlist.jet.html", comments, nav, searchForm)
}

// CommentDeleteModPanel : Controller for deleting a comment
func CommentDeleteModPanel(c *gin.Context) {
	id, err := strconv.ParseInt(c.PostForm("id"), 10, 32)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/mod/comments")
		return
	}
	
	comment, _, err := comments.Delete(uint(id))
	if err == nil {
		username := "れんちょん"
		if comment.UserID != 0 {
			username = comment.User.Username
		}
		activities.Log(&models.User{}, comment.Identifier(), "delete", "comment_deleted_by", strconv.Itoa(int(comment.ID)), username, router.GetUser(c).Username)
	}

	c.Redirect(http.StatusSeeOther, "/mod/comments?deleted")
}
