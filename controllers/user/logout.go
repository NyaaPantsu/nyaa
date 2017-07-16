package userController

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/gin-gonic/gin"
)

// UserLogoutHandler : Controller to logout users
func UserLogoutHandler(c *gin.Context) {
	logout := c.PostForm("logout")
	if logout != "" {
		cookies.Clear(c)
		url := c.DefaultPostForm("redirectTo", "/")
		c.Redirect(http.StatusSeeOther, url)
	} else {
		c.Status(http.StatusNotFound)
	}
}
