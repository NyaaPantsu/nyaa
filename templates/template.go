package templates

import (
	"net/http"
	"path"
	"strconv"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/filelist"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"

	"fmt"

	"github.com/CloudyKit/jet"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
)

// TemplateDir : Variable to the template directory
const TemplateDir = "./templates" // FIXME: Need to be a constant!

// ModeratorDir : Variable to the admin template sub directory
const ModeratorDir = "admin"

// SiteDir : Variable pointing to the site page templates
const SiteDir = "site"

// ErrorsDir : Variable pointing to the errors page templates
const ErrorsDir = "errors"

// View : Jet Template Renderer
var View = jet.NewHTMLSet("./templates")

func init() {
	if config.Get().Environment == "DEVELOPMENT" {
		View.SetDevelopmentMode(true)
		fmt.Println("Template Live Update enabled")
	}
}

// Commonvariables return a jet.VarMap variable containing the necessary variables to run index layouts
func Commonvariables(c *gin.Context) jet.VarMap {
	token := nosurf.Token(c.Request)
	messages := msg.GetMessages(c)
	user, _, _ := cookies.CurrentUser(c)
	variables := templateFunctions(make(jet.VarMap))
	variables.Set("Navigation", NewNavigation())
	variables.Set("Search", NewSearchForm(c))
	variables.Set("T", publicSettings.GetTfuncFromRequest(c))
	variables.Set("Theme", publicSettings.GetThemeFromRequest(c))
	variables.Set("AltColors", publicSettings.GetAltColorsFromRequest(c))
	variables.Set("OldNav", publicSettings.GetOldNavFromRequest(c))
	variables.Set("Mascot", publicSettings.GetMascotFromRequest(c))
	variables.Set("MascotURL", publicSettings.GetMascotURLFromRequest(c))
	variables.Set("User", user)
	variables.Set("URL", c.Request.URL)
	variables.Set("CsrfToken", token)
	variables.Set("EUCookieLaw", publicSettings.GetEUCookieFromRequest(c))
	variables.Set("Config", config.Get())
	variables.Set("Infos", messages.GetAllInfos())
	variables.Set("Errors", messages.GetAllErrors())
	return variables
}

// NewPanelSearchForm : Helper that creates a search form without items/page field
// these need to be used when the templateVariables don't include `navigation`
func NewPanelSearchForm(c *gin.Context) SearchForm {
	form := NewSearchForm(c)
	form.ShowItemsPerPage = false
	return form
}

// NewPanelCommonvariables return a jet.VarMap variable containing the necessary variables to run index admin layouts
func NewPanelCommonvariables(c *gin.Context) jet.VarMap {
	common := Commonvariables(c)
	common.Set("Search", NewPanelSearchForm(c))
	return common
}

// Render is a function rendering a template
func Render(c *gin.Context, templateName string, variables jet.VarMap) {
	t, err := View.GetTemplate(templateName)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err = t.Execute(c.Writer, variables, nil); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

// HttpError render an error template
func HttpError(c *gin.Context, errorCode int) {
	switch errorCode {
	case http.StatusNotFound:
		Static(c, path.Join(ErrorsDir, "404.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	case http.StatusBadRequest:
		Static(c, path.Join(ErrorsDir, "400.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	case http.StatusInternalServerError:
		Static(c, path.Join(ErrorsDir, "500.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	}
}

// Static render static templates
func Static(c *gin.Context, templateName string) {
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = NewPanelCommonvariables(c)
	} else {
		variables = Commonvariables(c)
	}
	Render(c, templateName, variables)
}

// ModelList render list models templates
func ModelList(c *gin.Context, templateName string, models interface{}, nav Navigation, search SearchForm) {
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = NewPanelCommonvariables(c)
	} else {
		variables = Commonvariables(c)
	}
	variables.Set("Models", models)
	variables.Set("Navigation", nav)
	variables.Set("Search", search)
	Render(c, templateName, variables)
}

// Form render a template form
func Form(c *gin.Context, templateName string, form interface{}) {
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = NewPanelCommonvariables(c)
	} else {
		variables = Commonvariables(c)
	}
	variables.Set("Form", form)
	Render(c, templateName, variables)
}

// Torrent render a torrent view template
func Torrent(c *gin.Context, torrent models.TorrentJSON, rootFolder *filelist.FileListFolder, captchaID string) {
	variables := Commonvariables(c)
	variables.Set("Torrent", torrent)
	variables.Set("RootFolder", rootFolder)
	variables.Set("CaptchaID", captchaID)
	Render(c, path.Join(SiteDir, "torrents", "view.jet.html"), variables)
}

// userProfilBase render the base for user profile
func userProfileBase(c *gin.Context, templateName string, userProfile *models.User, variables jet.VarMap) {
	currentUser, _, _ := cookies.CurrentUser(c)
	query := c.Request.URL.Query()
	query.Set("userID", strconv.Itoa(int(userProfile.ID)))
	query.Set("limit", "20")
	c.Request.URL.RawQuery = query.Encode()
	nbTorrents := 0
	if userProfile.ID > 0 && currentUser.CurrentOrAdmin(userProfile.ID) {
		_, userProfile.Torrents, nbTorrents, _ = search.ByQuery(c, 1, true, false, false)
	} else {
		_, userProfile.Torrents, nbTorrents, _ = search.ByQuery(c, 1, true, false, true)
	}

	variables.Set("UserProfile", userProfile)
	variables.Set("NbTorrents", nbTorrents)
	Render(c, path.Join(SiteDir, "user", templateName), variables)
}

// UserProfileEdit render a form to edit a profile
func UserProfileEdit(c *gin.Context, userProfile *models.User, userForm userValidator.UserForm, languages publicSettings.Languages) {
	variables := Commonvariables(c)
	variables.Set("UserForm", userForm)
	variables.Set("Languages", languages)
	userProfileBase(c, "edit.jet.html", userProfile, variables)
}

// UserProfile render a user profile
func UserProfile(c *gin.Context, userProfile *models.User) {
	userProfileBase(c, "torrents.jet.html", userProfile, Commonvariables(c))
}

// UserProfileNotifications render a user profile notifications
func UserProfileNotifications(c *gin.Context, userProfile *models.User) {
	userProfileBase(c, "notifications.jet.html", userProfile, Commonvariables(c))
}

// DatabaseDump render the list of database dumps template
func DatabaseDump(c *gin.Context, listDumps []models.DatabaseDumpJSON, GPGLink string) {
	variables := Commonvariables(c)
	variables.Set("ListDumps", listDumps)
	variables.Set("GPGLink", GPGLink)
	Render(c, path.Join(SiteDir, "database", "dumps.jet.html"), variables)
}

// PanelAdmin render the panel admin template index
func PanelAdmin(c *gin.Context, torrent []models.Torrent, reports []models.TorrentReportJSON, users []models.User, comments []models.Comment) {
	variables := NewPanelCommonvariables(c)
	variables.Set("Torrents", torrent)
	variables.Set("TorrentReports", reports)
	variables.Set("Users", users)
	variables.Set("Comments", comments)
	Render(c, path.Join(ModeratorDir, "index.jet.html"), variables)
}

func isAdminTemplate(templateName string) bool {
	if templateName != "" && len(templateName) > len(ModeratorDir) {
		return templateName[:5] == ModeratorDir
	}
	return false
}
