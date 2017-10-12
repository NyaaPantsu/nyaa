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
	if userid != "" {
		conditions = "user_id = ?"
		values = append(values, userid)
	}

	comments, nbComments := comments.FindAll(offset, (pagenum-1)*offset, conditions, values...)
	nav := templates.Navigation{nbComments, offset, pagenum, "mod/comments/p"}
	templates.ModelList(c, "admin/commentlist.jet.html", comments, nav, templates.NewSearchForm(c))
}

// CommentDeleteModPanel : Controller for deleting a comment
func CommentDeleteModPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.PostForm("id"), 10, 32)
	comment, _, err := comments.Delete(uint(id))
	if err == nil {
		activities.Log(&models.User{}, comment.Identifier(), "delete", "comment_deleted_by", strconv.Itoa(int(comment.ID)), comment.User.Username, router.GetUser(c).Username)
	}

	c.Redirect(http.StatusSeeOther, "/mod/comments?deleted")
}
