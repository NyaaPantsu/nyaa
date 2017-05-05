package main

import (
	"github.com/gorilla/mux"
	"net/url"
)

/* Each Page should have an object to pass to their own template
 * Therefore, we put them in a separate file for better maintenance
 *
 * MAIN Template Variables
 */

type FaqTemplateVariables struct {
	Navigation Navigation
	Search 		SearchForm
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type ViewTemplateVariables struct {
	Torrent 		TorrentsJson
	Search 			SearchForm
	Navigation     Navigation
	URL            *url.URL   // For parsing Url in templates
	Route          *mux.Route // For getting current route in templates
}

type HomeTemplateVariables struct {
	ListTorrents   []TorrentsJson
	ListCategories []Categories
	Search 			SearchForm
	Navigation     Navigation
	URL            *url.URL   // For parsing Url in templates
	Route          *mux.Route // For getting current route in templates
}


/*
 * Variables used by the upper ones
 */
type Navigation struct {
	TotalItem      int
	MaxItemPerPage int
	CurrentPage    int
	Route          string
}

type SearchForm struct {
	Query      string
	Status     string
	Category   string
	Sort       string
	Order      string
}

// Some Default Values to ease things out
func NewSearchForm(params ...string) SearchForm {
	searchForm := SearchForm{}
	if (len(params) > 1) { 
		searchForm.Category = params[0]
	} else { 
		searchForm.Category = "_" 
	}
	if (len(params) > 2) { 
		searchForm.Sort = params[0]
	} else { 
		searchForm.Sort = "torrent_id" 
	}
	if (len(params) > 3) { 
		searchForm.Order = params[0]
	} else { 
		searchForm.Order = "DESC" 
	}
	return searchForm
}