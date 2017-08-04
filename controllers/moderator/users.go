package moderatorController

import (
	"html"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/gin-gonic/gin"
)

// UsersListPanel : Controller for listing users, can accept pages
func UsersListPanel(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	var err error

	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	users, nbUsers := users.FindUsersForAdmin(offset, (pagenum-1)*offset)
	nav := templates.Navigation{nbUsers, offset, pagenum, "mod/users/p"}
	templates.ModelList(c, "admin/userlist.jet.html", users, nav, templates.NewSearchForm(c))
}
