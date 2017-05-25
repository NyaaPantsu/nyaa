package router

import (
	"net/http"
	"net/url"

	"github.com/NyaaPantsu/nyaa/common"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/user"
	userForms "github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/util/filelist"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/gorilla/mux"
)

/* Each Page should have an object to pass to their own template
* Therefore, we put them in a separate file for better maintenance
*
* MAIN Template Variables
 */

type viewTemplateVariables struct {
	commonTemplateVariables
	Torrent    model.TorrentJSON
	RootFolder *filelist.FileListFolder // used for tree view
	CaptchaID  string
	FormErrors map[string][]string
	Infos      map[string][]string
}

type formTemplateVariables struct {
	commonTemplateVariables
	Form       interface{}
	FormErrors map[string][]string
	FormInfos  map[string][]string
}

type userProfileEditVariables struct {
	commonTemplateVariables
	UserProfile *model.User
	UserForm    userForms.UserForm
	FormErrors  map[string][]string
	FormInfos   map[string][]string
	Languages   map[string]string
}

type userVerifyTemplateVariables struct {
	commonTemplateVariables
	FormErrors map[string][]string
}

type userProfileVariables struct {
	commonTemplateVariables
	UserProfile *model.User
	FormInfos   map[string][]string
}

type userProfileNotifVariables struct {
	commonTemplateVariables
	Infos map[string][]string
}

type databaseDumpTemplateVariables struct {
	commonTemplateVariables
	ListDumps []model.DatabaseDumpJSON
	GPGLink   string
}

type changeLanguageVariables struct {
	commonTemplateVariables
	Language  string
	Languages map[string]string
}

/* MODERATION Variables */

type panelIndexVbs struct {
	commonTemplateVariables
	Torrents       []model.Torrent
	TorrentReports []model.TorrentReportJson
	Users          []model.User
	Comments       []model.Comment
}

type modelListVbs struct {
	commonTemplateVariables
	Models interface{}
	Errors map[string][]string
	Infos  map[string][]string
}

/*
* Variables used by the upper ones
 */

type commonTemplateVariables struct {
	Navigation navigation
	Search     searchForm
	T          languages.TemplateTfunc
	User       *model.User
	URL        *url.URL   // for parsing URL in templates
	Route      *mux.Route // for getting current route in templates
}

type navigation struct {
	TotalItem      int
	MaxItemPerPage int // FIXME: shouldn't this be in SearchForm?
	CurrentPage    int
	Route          string
}

type searchForm struct {
	common.SearchParam
	Category         string
	ShowItemsPerPage bool
}

// Some Default Values to ease things out
func newNavigation() navigation {
	return navigation{
		MaxItemPerPage: 50,
	}
}

func newSearchForm() searchForm {
	return searchForm{
		Category:         "_",
		ShowItemsPerPage: true,
	}
}

func getUser(r *http.Request) *model.User {
	user, _, _ := userService.RetrieveCurrentUser(r)
	return &user
}

func newCommonVariables(r *http.Request) commonTemplateVariables {
	return commonTemplateVariables{
		Navigation: newNavigation(),
		Search:     newSearchForm(),
		T:          languages.GetTfuncFromRequest(r),
		User:       getUser(r),
		URL:        r.URL,
		Route:      mux.CurrentRoute(r),
	}
}
