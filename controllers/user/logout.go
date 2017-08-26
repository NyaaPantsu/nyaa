package userController

import (
	"net/http"
	"strings"

	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/gin-gonic/gin"
)

// UserLogoutHandler : Controller to logout users
func UserLogoutHandler(c *gin.Context) {
	logout := c.PostForm("logout")
	if logout != "" {
		cookies.Clear(c)
		url := c.DefaultPostForm("redirectTo", "/")
		if strings.Contains(url, "/mod/") {
			url = "/"
		}
		c.Redirect(http.StatusSeeOther, url)
	} else {
		c.Status(http.StatusNotFound)
	}
}
