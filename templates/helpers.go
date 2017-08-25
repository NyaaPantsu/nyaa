package templates

import (
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/gin-gonic/gin"
)

// LanguagesJSONResponse : Structure containing all the languages to parse it as a JSON response
type LanguagesJSONResponse struct {
	Current   string                   `json:"current"`
	Languages publicSettings.Languages `json:"language"`
}

// Navigation is used to display navigation links to pages on list view
type Navigation struct {
	TotalItem      int
	MaxItemPerPage int // FIXME: shouldn't this be in SearchForm?
	CurrentPage    int
	Route          string
}

// SearchForm struct used to display the search form
type SearchForm struct {
	search.TorrentParam
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

// NewNavigation return a navigation struct with
// Some Default Values to ease things out
func NewNavigation() Navigation {
	return Navigation{
		MaxItemPerPage: 50,
	}
}

// NewSearchForm return a searchForm struct with
// Some Default Values to ease things out
func NewSearchForm(c *gin.Context) SearchForm {
	sizeType := c.DefaultQuery("sizeType", "m")
	return SearchForm{
		Category:         "_",
		ShowItemsPerPage: true,
		ShowRefine:       false,
		SizeType:         sizeType,
		DateType:         c.Query("dateType"),
		MinSize:          c.Query("minSize"),  // We need to overwrite the value here, since size are formatted
		MaxSize:          c.Query("maxSize"),  // We need to overwrite the value here, since size are formatted
		FromDate:         c.Query("fromDate"), // We need to overwrite the value here, since we can have toDate instead and date are formatted
		ToDate:           c.Query("toDate"),   // We need to overwrite the value here, since date are formatted
	}
}
