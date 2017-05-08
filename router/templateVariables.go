package router

import (
	"net/url"
	"net/http"

	"github.com/ewhal/nyaa/model"
	userForms "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/service/user"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/gorilla/mux"
)

/* Each Page should have an object to pass to their own template
 * Therefore, we put them in a separate file for better maintenance
 *
 * MAIN Template Variables
 */

type FaqTemplateVariables struct {
	Navigation Navigation
	Search     SearchForm
	User       model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type NotFoundTemplateVariables struct {
	Navigation Navigation
	Search     SearchForm
	User       model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type ViewTemplateVariables struct {
	Torrent    model.TorrentsJson
	Captcha    captcha.Captcha
	Search     SearchForm
	Navigation Navigation
	User       model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type UserRegisterTemplateVariables struct {
	RegistrationForm userForms.RegistrationForm
	FormErrors 		 map[string][]string
	Search           SearchForm
	Navigation       Navigation
	User      		 model.User
	URL              *url.URL   // For parsing Url in templates
	Route            *mux.Route // For getting current route in templates
}

type UserVerifyTemplateVariables struct {
	FormErrors 		 map[string][]string
	Search           SearchForm
	Navigation       Navigation
	User       		 model.User
	URL              *url.URL   // For parsing Url in templates
	Route            *mux.Route // For getting current route in templates
}

type UserLoginFormVariables struct {
	LoginForm userForms.LoginForm
	FormErrors 		 map[string][]string
	Search           SearchForm
	Navigation       Navigation
	User      		 model.User
	URL              *url.URL   // For parsing Url in templates
	Route            *mux.Route // For getting current route in templates
}

type HomeTemplateVariables struct {
	ListTorrents []model.TorrentsJson
	Search       SearchForm
	Navigation   Navigation
	User      	 model.User
	URL          *url.URL   // For parsing Url in templates
	Route        *mux.Route // For getting current route in templates
}

type UploadTemplateVariables struct {
	Upload     UploadForm
	Search     SearchForm
	Navigation Navigation
	User       model.User
	URL        *url.URL
	Route      *mux.Route
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
	Query              string
	Status             string
	Category           string
	Sort               string
	Order              string
	HideAdvancedSearch bool
}

// Some Default Values to ease things out
func NewSearchForm(params ...string) (searchForm SearchForm) {
	if len(params) > 1 {
		searchForm.Category = params[0]
	} else {
		searchForm.Category = "_"
	}
	if len(params) > 2 {
		searchForm.Sort = params[1]
	} else {
		searchForm.Sort = "torrent_id"
	}
	if len(params) > 3 {
		order := params[2]
		if order == "DESC" {
			searchForm.Order = order
		} else if order == "ASC" {
			searchForm.Order = order
		} else {
			// TODO: handle invalid value (?)
		}
	} else {
		searchForm.Order = "DESC"
	}
	return
}

func GetUser(r *http.Request) model.User {
	user, _ , _ := userService.RetrieveCurrentUser(r)
	return user
}
