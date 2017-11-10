package feedController

import (
	"errors"
	"html"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/feeds"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/gin-gonic/gin"
)

func getTorrentList(c *gin.Context) (torrents []models.Torrent, createdAsTime time.Time, title string, err error) {
	page := c.Param("page")
	userID := c.Param("id")
	cat := c.Query("cat")
	offset := 0
	currentUser := router.GetUser(c)
	if c.Query("offset") != "" {
		offset, err = strconv.Atoi(html.EscapeString(c.Query("offset")))
		if err != nil {
			return
		}
	}

	createdAsTime = time.Now()

	if len(torrents) > 0 {
		createdAsTime = torrents[0].Date
	}

	title = "Nyaa Pantsu"
	if config.IsSukebei() {
		title = "Sukebei Pantsu"
	}

	pagenum := 1
	if page == "" && offset > 0 { // first page for offset is 0
		pagenum = offset + 1
	} else if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if err != nil {
			return
		}
	}
	if pagenum <= 0 {
		err = errors.New("Page number is invalid")
		return
	}

	if userID != "" {
		var userIDnum int
		userIDnum, err = strconv.Atoi(html.EscapeString(userID))
		// Should we have a feed for anonymous uploads?
		if err != nil || userIDnum == 0 {
			return
		}

		_, _, err = users.FindForAdmin(uint(userIDnum))
		if err != nil {
			return
		}
		// Set the user ID on the request, so that SearchByQuery finds it.
		query := c.Request.URL.Query()
		query.Set("userID", userID)
		c.Request.URL.RawQuery = query.Encode()
	}

	if cat != "" {
		query := c.Request.URL.Query()
		catConv := nyaafeeds.ConvertToCat(cat)
		if catConv == "" {
			return
		}
		query.Set("c", catConv)
		c.Request.URL.RawQuery = query.Encode()
	}

	user, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		user = 0
	}

	_, torrents, _, err = search.AuthorizedQuery(c, pagenum, currentUser.CurrentOrJanitor(uint(user)))

	return
}
