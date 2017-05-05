package main

import (
	"github.com/gorilla/mux"
	"net/url"
)

/* Each Page should have an object to pass to their own template
 * Therefore, we put them in a separate file for better maintenance
 */

type FaqTemplateVariables struct {
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
	Query      string
	Status     string
	Category   string
	Navigation Navigation
}

type HomeTemplateVariables struct {
	ListTorrents   []TorrentsJson
	ListCategories []Categories
	Query          string
	Status         string
	Category       string
	Navigation     Navigation
	URL            *url.URL   // For parsing Url in templates
	Route          *mux.Route // For getting current route in templates
}

type Navigation struct {
	TotalItem      int
	MaxItemPerPage int
	CurrentPage    int
	Route          string
}
