package controllers

import (
	"html"
	"net/http"
	"strconv"

	"math"

	"github.com/NyaaPantsu/nyaa/models"
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

	searchParam, torrents, nbTorrents, err := search.ByQuery(c, pagenum)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	maxPages := math.Ceil(float64(nbTorrents) / float64(searchParam.Max))
	if pagenum > int(maxPages) {
		c.AbortWithError(http.StatusNotFound, errors.New("Page superior to the maximum number of pages"))
		return
	}

	// Convert back to strings for now.
	category := ""
	if len(searchParam.Category) > 0 {
		category = searchParam.Category[0].String()
	}
	nav := navigation{int(nbTorrents), int(searchParam.Max), int(searchParam.Offset), "search"}
	searchForm := newSearchForm(c)
	searchForm.TorrentParam, searchForm.Category = searchParam, category

	if c.Request.URL.Path == "/" {
		searchForm.ShowRefine = false
		//pls change that condition to check if url has REFINE get parameter
	}

	modelList(c, "site/torrents/listing.jet.html", models.TorrentsToJSON(torrents), nav, searchForm)
}
