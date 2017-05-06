package router

import (
	"net/url"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/templates"
	"github.com/gorilla/mux"
)

/* Each Page should have an object to pass to their own template
 * Therefore, we put them in a separate file for better maintenance
 *
 * MAIN Template Variables
 */

type FaqTemplateVariables struct {
	Navigation templates.Navigation
	Search     templates.SearchForm
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type ViewTemplateVariables struct {
	Torrent    model.TorrentsJson
	Search     templates.SearchForm
	Navigation templates.Navigation
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type HomeTemplateVariables struct {
	ListTorrents   []model.TorrentsJson
	Search         templates.SearchForm
	Navigation     templates.Navigation
	URL            *url.URL   // For parsing Url in templates
	Route          *mux.Route // For getting current route in templates
}

type UploadTemplateVariables struct {
	Upload     UploadForm
	Search     templates.SearchForm
	Navigation templates.Navigation
	URL        *url.URL
	Route      *mux.Route
}
