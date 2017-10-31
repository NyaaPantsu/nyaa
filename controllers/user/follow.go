package userController

import (
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/gin-gonic/gin"
)

// UserFollowHandler : Controller to follow/unfollow users, need user id to follow
func UserFollowHandler(c *gin.Context) {
	var followAction string
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	currentUser := router.GetUser(c)
	user, _, errorUser := users.FindForAdmin(uint(id))
	if errorUser == nil && user.ID > 0 && currentUser.ID > 0 && user.ID != currentUser.ID  {
		if !currentUser.IsFollower(uint(id)) {
			followAction = "followed"
			currentUser.SetFollow(user)
		} else {
			followAction = "unfollowed"
			currentUser.RemoveFollow(user)
		}
	}
	url := "/user/" + strconv.Itoa(int(user.ID)) + "/" + user.Username + "?" + followAction
	if c.Query("id") != "" {
		url = "/view/" + c.Query("id") + "?" + followAction
	}
	if currentUser.ID == 0 {
		url = "/login"
	}
	c.Redirect(http.StatusSeeOther, url)
}
