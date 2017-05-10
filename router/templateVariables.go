package router

import (
	"net/http"
	"net/url"

	"github.com/ewhal/nyaa/common"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/service/user"
	userForms "github.com/ewhal/nyaa/service/user/form"
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
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type NotFoundTemplateVariables struct {
	Navigation Navigation
	Search     SearchForm
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type ViewTemplateVariables struct {
	Torrent    model.TorrentJSON
	Captcha    captcha.Captcha
	Search     SearchForm
	Navigation Navigation
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type UserRegisterTemplateVariables struct {
	RegistrationForm userForms.RegistrationForm
	FormErrors       map[string][]string
	Search           SearchForm
	Navigation       Navigation
	User             *model.User
	URL              *url.URL   // For parsing Url in templates
	Route            *mux.Route // For getting current route in templates
}

type UserProfileEditVariables struct {
	UserProfile *model.User
	UserForm    userForms.UserForm
	FormErrors  map[string][]string
	FormInfos   map[string][]string
	Search      SearchForm
	Navigation  Navigation
	User        *model.User
	URL         *url.URL   // For parsing Url in templates
	Route       *mux.Route // For getting current route in templates
}

type UserVerifyTemplateVariables struct {
	FormErrors map[string][]string
	Search     SearchForm
	Navigation Navigation
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type UserLoginFormVariables struct {
	LoginForm  userForms.LoginForm
	FormErrors map[string][]string
	Search     SearchForm
	Navigation Navigation
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type UserProfileVariables struct {
	UserProfile *model.User
	FormInfos   map[string][]string
	Search      SearchForm
	Navigation  Navigation
	User        *model.User
	URL         *url.URL   // For parsing Url in templates
	Route       *mux.Route // For getting current route in templates
}

type HomeTemplateVariables struct {
	ListTorrents []model.TorrentJSON
	Search       SearchForm
	Navigation   Navigation
	User         *model.User
	URL          *url.URL   // For parsing Url in templates
	Route        *mux.Route // For getting current route in templates
}

type UploadTemplateVariables struct {
	Upload     UploadForm
	Search     SearchForm
	Navigation Navigation
	User       *model.User
	URL        *url.URL
	Route      *mux.Route
}

/* MODERATION Variables */

type PanelIndexVbs struct {
	Torrents       []model.Torrent
	TorrentReports []model.TorrentReport
	Users          []model.User
	Comments       []model.Comment
	Search         SearchForm
	User           *model.User
	URL            *url.URL // For parsing Url in templates
}

type PanelTorrentListVbs struct {
	Torrents   []model.Torrent
	Search     SearchForm
	Navigation Navigation
	User       *model.User
	URL        *url.URL // For parsing Url in templates
}
type PanelUserListVbs struct {
	Users      []model.User
	Search     SearchForm
	Navigation Navigation
	User       *model.User
	URL        *url.URL // For parsing Url in templates
}
type PanelCommentListVbs struct {
	Comments   []model.Comment
	Search     SearchForm
	Navigation Navigation
	User       *model.User
	URL        *url.URL // For parsing Url in templates
}
type PanelTorrentEdVbs struct {
	Torrent model.Torrent
	Search  SearchForm
	User    *model.User
}

type PanelTorrentReportListVbs struct {
	TorrentReports []model.TorrentReportJson
	Search         SearchForm
	Navigation     Navigation
	User           *model.User
	URL            *url.URL // For parsing Url in templates
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
	common.SearchParam
	Category           string
	HideAdvancedSearch bool
}

// Some Default Values to ease things out
func NewSearchForm() SearchForm {
	return SearchForm{
		Category: "_",
	}
}

func GetUser(r *http.Request) *model.User {
	user, _, _ := userService.RetrieveCurrentUser(r)
	return &user
}
