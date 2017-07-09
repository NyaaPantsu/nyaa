package controllers

import (
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/gin-gonic/gin"
)

type navigation struct {
	TotalItem      int
	MaxItemPerPage int // FIXME: shouldn't this be in SearchForm?
	CurrentPage    int
	Route          string
}

type searchForm struct {
	structs.TorrentParam
	Category         string
	ShowItemsPerPage bool
	ShowRefine       bool
	SizeType         string
	DateType         string
	MinSize          string
	MaxSize          string
	FromDate         string
	ToDate           string
}

// Some Default Values to ease things out
func newNavigation() navigation {
	return navigation{
		MaxItemPerPage: 50,
	}
}

func newSearchForm(c *gin.Context) searchForm {
	sizeType := c.DefaultQuery("sizeType", "m")
	return searchForm{
		Category:         "_",
		ShowItemsPerPage: true,
		ShowRefine:       true,
		SizeType:         sizeType,
		DateType:         c.Query("dateType"),
		MinSize:          c.Query("minSize"),  // We need to overwrite the value here, since size are formatted
		MaxSize:          c.Query("maxSize"),  // We need to overwrite the value here, since size are formatted
		FromDate:         c.Query("fromDate"), // We need to overwrite the value here, since we can have toDate instead and date are formatted
		ToDate:           c.Query("toDate"),   // We need to overwrite the value here, since date are formatted
	}
}
func getUser(c *gin.Context) *models.User {
	user, _, _ := cookies.CurrentUser(c)
	if user == nil {
		return &models.User{}
	}
	return user
}
