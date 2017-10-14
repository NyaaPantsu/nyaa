package searchController

import (
	"html"
	"net/http"
	"strconv"
	"fmt"

	"math"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)
// SearchHandler : Controller for displaying search result page, accepting common search arguments
func SearchHandler(c *gin.Context) {
	var err error
	// TODO Don't create a new client for each request
	// TODO Fallback to postgres search if es is down

	page := c.Param("page")
	currentUser := router.GetUser(c)
	// db params url
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if err != nil {
			c.AbortWithError(http.StatusNotFound, err)
			return
		}
		if pagenum <= 0 {
			c.AbortWithError(http.StatusNotFound, errors.New("Can't find a page with negative value"))
			return
		}
	}
	
	searchForm := templates.NewSearchForm(c)
	
	if c.Param("id") != "" {
		query := c.Request.URL.Query()
		query.Set("userID", c.Param("id"))
		c.Request.URL.RawQuery = query.Encode()
		searchForm.SearchURL = fmt.Sprintf("/user/%s/%s/search", c.Param("id"), c.Param("username"))
		searchForm.UserName = c.Param("username") //Only add username if user search route
	}
	
	userID, err := strconv.ParseUint(c.Query("userID"), 10, 32)
	if err != nil {
		userID = 0
	}
	
	if userID == 0 && c.Param("id") != "" && c.Param("id") != "0" {
		c.Redirect(http.StatusSeeOther, fmt.Sprintf("/user/%s/%s", c.Param("id"), c.Param("username")))
		//User is trying to use the user search route with an inexisting user
		//Must redirect him to user search instead of simply showing "no torrents found!"
	}
	
	
	searchParam, torrents, nbTorrents, err := search.AuthorizedQuery(c, pagenum, currentUser.CurrentOrAdmin(uint(userID)))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	
	// Convert back to strings for now.
	category := ""
	if len(searchParam.Category) > 0 {
		category = searchParam.Category[0].String()
	}
	nav := templates.Navigation{int(nbTorrents), int(searchParam.Max), int(searchParam.Offset), "search"}
	
	searchForm.TorrentParam, searchForm.Category = searchParam, category

	if c.Query("refine") == "1" {
		searchForm.ShowRefine = true
	}

	maxPages := math.Ceil(float64(nbTorrents) / float64(searchParam.Max))
	if pagenum > int(maxPages) {
		variables := templates.Commonvariables(c)
		variables.Set("Search", searchForm)
		templates.Render(c, "errors/no_results.jet.html", variables)
		return
	}

	templates.ModelList(c, "site/torrents/listing.jet.html", models.TorrentsToJSON(torrents), nav, searchForm)
}
