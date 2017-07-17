package activitiesController

import (
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/gin-gonic/gin"
)

// ActivityListHandler : Show a list of activity
func ActivityListHandler(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	userid := c.Query("userid")
	filter := c.Query("filter")

	var err error
	currentUser := router.GetUser(c)
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	var conditions []string
	var values []interface{}
	if userid != "" && currentUser.HasAdmin() {
		conditions = append(conditions, "user_id = ?")
		values = append(values, userid)
	}
	if filter != "" {
		conditions = append(conditions, "filter = ?")
		values = append(values, filter)
	}

	activity, nbActivities := activities.FindAll(offset, (pagenum-1)*offset, strings.Join(conditions, " AND "), values...)

	nav := templates.Navigation{nbActivities, offset, pagenum, "activities/p"}
	templates.ModelList(c, "site/torrents/activities.jet.html", activity, nav, templates.NewSearchForm(c))
}
