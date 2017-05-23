package router

import (
	"net/http"
	"net/url"

	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/user"
	userForms "github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/gorilla/mux"
)

/* Each Page should have an object to pass to their own template
* Therefore, we put them in a separate file for better maintenance
*
* MAIN Template Variables
 */

type FaqTemplateVariables struct {
	CommonTemplateVariables
}

type NotFoundTemplateVariables struct {
	CommonTemplateVariables
}

type ViewTemplateVariables struct {
	CommonTemplateVariables
	Torrent    model.TorrentJSON
	CaptchaID  string
	FormErrors  map[string][]string
	Infos   map[string][]string
}

type UserRegisterTemplateVariables struct {
	CommonTemplateVariables
	RegistrationForm userForms.RegistrationForm
	FormErrors       map[string][]string
}

type UserProfileEditVariables struct {
	CommonTemplateVariables
	UserProfile *model.User
	UserForm    userForms.UserForm
	FormErrors  map[string][]string
	FormInfos   map[string][]string
	Languages   map[string]string
}

type UserVerifyTemplateVariables struct {
	CommonTemplateVariables
	FormErrors map[string][]string
}

type UserLoginFormVariables struct {
	CommonTemplateVariables
	LoginForm  userForms.LoginForm
	FormErrors map[string][]string
}

type UserProfileVariables struct {
	CommonTemplateVariables
	UserProfile *model.User
	FormInfos   map[string][]string
}

type UserProfileNotifVariables struct {
	CommonTemplateVariables
	Infos   map[string][]string
}

type UserTorrentEdVbs struct {
	CommonTemplateVariables
	Upload     UploadForm
	FormErrors map[string][]string
	FormInfos  map[string][]string
}

type HomeTemplateVariables struct {
	CommonTemplateVariables
	ListTorrents []model.TorrentJSON
	Infos   map[string][]string
}

type DatabaseDumpTemplateVariables struct {
	CommonTemplateVariables
	ListDumps  []model.DatabaseDumpJSON
	GPGLink    string
}

type UploadTemplateVariables struct {
	CommonTemplateVariables
	Upload     UploadForm
	FormErrors  map[string][]string
}

type ChangeLanguageVariables struct {
	CommonTemplateVariables
	Language   string
	Languages  map[string]string
}

/* MODERATION Variables */

type PanelIndexVbs struct {
	CommonTemplateVariables
	Torrents       []model.Torrent
	TorrentReports []model.TorrentReportJson
	Users          []model.User
	Comments       []model.Comment
}

type PanelTorrentListVbs struct {
	CommonTemplateVariables
	Torrents   []model.Torrent
	Errors map[string][]string
	Infos  map[string][]string
}
type PanelUserListVbs struct {
	CommonTemplateVariables
	Users      []model.User
}
type PanelCommentListVbs struct {
	CommonTemplateVariables
	Comments   []model.Comment
}

type PanelTorrentEdVbs struct {
	CommonTemplateVariables
	Upload     UploadForm
	FormErrors map[string][]string
	FormInfos  map[string][]string
}

type PanelTorrentReportListVbs struct {
	CommonTemplateVariables
	TorrentReports []model.TorrentReportJson
}

type PanelTorrentReassignVbs struct {
	CommonTemplateVariables
	Reassign   ReassignForm
	FormErrors map[string][]string
	FormInfos  map[string][]string
}

/*
* Variables used by the upper ones
 */

type CommonTemplateVariables struct {
	Navigation Navigation
	Search     SearchForm
	T          languages.TemplateTfunc
	User       *model.User
	URL        *url.URL // for parsing URL in templates
    Route      *mux.Route // for getting current route in templates
}

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

func NewCommonVariables(r *http.Request) CommonTemplateVariables {
	return CommonTemplateVariables{
		Navigation: NewNavigation(),
		Search:     NewSearchForm(),
		T:          languages.GetTfuncFromRequest(r),
		User:       GetUser(r),
		URL:        r.URL,
		Route:      mux.CurrentRoute(r),
	}
}

