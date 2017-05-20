package router

import (
	"net/http"
	"net/url"

	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/user"
	userForms "github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"
)

/* Each Page should have an object to pass to their own template
* Therefore, we put them in a separate file for better maintenance
*
* MAIN Template Variables
 */

type FaqTemplateVariables struct {
	Navigation Navigation
	Search     SearchForm
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type NotFoundTemplateVariables struct {
	Navigation Navigation
	Search     SearchForm
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type ViewTemplateVariables struct {
	Torrent    model.TorrentJSON
	CaptchaID  string
	FormErrors  map[string][]string
	Infos   map[string][]string
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type UserRegisterTemplateVariables struct {
	RegistrationForm userForms.RegistrationForm
	FormErrors       map[string][]string
	Search           SearchForm
	Navigation       Navigation
	T                i18n.TranslateFunc
	User             *model.User
	URL              *url.URL   // For parsing Url in templates
	Route            *mux.Route // For getting current route in templates
}

type UserProfileEditVariables struct {
	UserProfile *model.User
	UserForm    userForms.UserForm
	FormErrors  map[string][]string
	FormInfos   map[string][]string
	Languages   map[string]string
	Search      SearchForm
	Navigation  Navigation
	T           i18n.TranslateFunc
	User        *model.User
	URL         *url.URL   // For parsing Url in templates
	Route       *mux.Route // For getting current route in templates
}

type UserVerifyTemplateVariables struct {
	FormErrors map[string][]string
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type UserLoginFormVariables struct {
	LoginForm  userForms.LoginForm
	FormErrors map[string][]string
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type UserProfileVariables struct {
	UserProfile *model.User
	FormInfos   map[string][]string
	Search      SearchForm
	Navigation  Navigation
	T           i18n.TranslateFunc
	User        *model.User
	URL         *url.URL   // For parsing Url in templates
	Route       *mux.Route // For getting current route in templates
}

type HomeTemplateVariables struct {
	ListTorrents []model.TorrentJSON
	Search       SearchForm
	Navigation   Navigation
	T            i18n.TranslateFunc
	User         *model.User
	URL          *url.URL   // For parsing Url in templates
	Route        *mux.Route // For getting current route in templates
}

type DatabaseDumpTemplateVariables struct {
	ListDumps  []model.DatabaseDumpJSON
	GPGLink    string
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL   // For parsing Url in templates
	Route      *mux.Route // For getting current route in templates
}

type UploadTemplateVariables struct {
	Upload     UploadForm
	FormErrors  map[string][]string
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL
	Route      *mux.Route
}

type ChangeLanguageVariables struct {
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	Language   string
	Languages  map[string]string
	User       *model.User
	URL        *url.URL
	Route      *mux.Route
}

/* MODERATION Variables */

type PanelIndexVbs struct {
	Torrents       []model.Torrent
	TorrentReports []model.TorrentReportJson
	Users          []model.User
	Comments       []model.Comment
	Search         SearchForm
	T              i18n.TranslateFunc
	User           *model.User
	URL            *url.URL // For parsing Url in templates
}

type PanelTorrentListVbs struct {
	Torrents   []model.Torrent
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	User       *model.User
	Errors map[string][]string
	Infos  map[string][]string
	URL        *url.URL // For parsing Url in templates
}
type PanelUserListVbs struct {
	Users      []model.User
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL // For parsing Url in templates
}
type PanelCommentListVbs struct {
	Comments   []model.Comment
	Search     SearchForm
	Navigation Navigation
	T          i18n.TranslateFunc
	User       *model.User
	URL        *url.URL // For parsing Url in templates
}

type PanelTorrentEdVbs struct {
	Upload     UploadForm
	Search     SearchForm
	T          i18n.TranslateFunc
	User       *model.User
	FormErrors map[string][]string
	FormInfos  map[string][]string
	URL        *url.URL // For parsing Url in templates
}

type PanelTorrentReportListVbs struct {
	TorrentReports []model.TorrentReportJson
	Search         SearchForm
	Navigation     Navigation
	T              i18n.TranslateFunc
	User           *model.User
	URL            *url.URL // For parsing Url in templates
}

type PanelTorrentReassignVbs struct {
	Reassign   ReassignForm
	Search     SearchForm  // unused?
	T          i18n.TranslateFunc
	User       *model.User // unused?
	FormErrors map[string][]string
	FormInfos  map[string][]string
	URL        *url.URL // For parsing Url in templates
}

/*
* Variables used by the upper ones
 */
type Navigation struct {
	TotalItem      int
	MaxItemPerPage int // FIXME: shouldn't this be in SearchForm?
	CurrentPage    int
	Route          string
}

type SearchForm struct {
	common.SearchParam
	Category         string
	ShowItemsPerPage bool
}

// Some Default Values to ease things out
func NewNavigation() Navigation {
	return Navigation{
		MaxItemPerPage: 50,
	}
}

func NewSearchForm() SearchForm {
	return SearchForm{
		Category:         "_",
		ShowItemsPerPage: true,
	}
}

func GetUser(r *http.Request) *model.User {
	user, _, _ := userService.RetrieveCurrentUser(r)
	return &user
}
