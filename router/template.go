package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/model"
	userForms "github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/util/filelist"
	"github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"

	"fmt"

	"github.com/CloudyKit/jet"
)

// TemplateDir : Variable to the template directory
var TemplateDir = "templates" // FIXME: Need to be a constant!

// ModeratorDir : Variable to the admin template sub directory
const ModeratorDir = "admin"

// View : Jet Template Renderer
var View = jet.NewHTMLSet("./templates")
var vars = templateFunctions(make(jet.VarMap))

func init() {
	if config.Conf.Environment == "DEVELOPMENT" {
		View.SetDevelopmentMode(true)
		fmt.Println("Template Live Update enabled")
	}
}
func commonVars(c *gin.Context) jet.VarMap {
	msg := messages.GetMessages(c)
	vars.Set("Navigation", newNavigation())
	vars.Set("Search", newSearchForm(c))
	vars.Set("T", publicSettings.GetTfuncFromRequest(c))
	vars.Set("Theme", publicSettings.GetThemeFromRequest(c))
	vars.Set("Mascot", publicSettings.GetMascotFromRequest(c))
	vars.Set("MascotURL", publicSettings.GetMascotUrlFromRequest(c))
	vars.Set("User", getUser(c))
	vars.Set("URL", c.Request.URL)
	vars.Set("CsrfToken", nosurf.Token(c.Request))
	vars.Set("Config", config.Conf)
	vars.Set("Infos", msg.GetAllInfos())
	vars.Set("Errors", msg.GetAllErrors())
	return vars
}

// newPanelSearchForm : Helper that creates a search form without items/page field
// these need to be used when the templateVariables don't include `navigation`
func newPanelSearchForm(c *gin.Context) searchForm {
	form := newSearchForm(c)
	form.ShowItemsPerPage = false
	return form
}

//
func newPanelCommonVariables(c *gin.Context) jet.VarMap {
	common := commonVars(c)
	common.Set("Search", newPanelSearchForm(c))
	return common
}

func renderTemplate(c *gin.Context, templateName string, vars jet.VarMap) {
	t, err := View.GetTemplate(templateName)
	if err != nil {
		fmt.Println("404")
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err = t.Execute(c.Writer, vars, nil); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func httpError(c *gin.Context, errorCode int) {
	if errorCode == http.StatusNotFound {
		c.Status(http.StatusNotFound)
		staticTemplate(c, "404.jet.html")
		return
	}
	c.Status(errorCode)
}

func staticTemplate(c *gin.Context, templateName string) {
	var vars jet.VarMap
	if isAdminTemplate(templateName) {
		vars = newPanelCommonVariables(c)
	} else {
		vars = commonVars(c)
	}
	renderTemplate(c, templateName, vars)
}

func modelList(c *gin.Context, templateName string, models interface{}, nav navigation, search searchForm) {
	var vars jet.VarMap
	if isAdminTemplate(templateName) {
		vars = newPanelCommonVariables(c)
	} else {
		vars = commonVars(c)
	}
	vars.Set("Models", models)
	vars.Set("Navigation", nav)
	vars.Set("Search", search)
	renderTemplate(c, templateName, vars)
}

func formTemplate(c *gin.Context, templateName string, form interface{}) {
	var vars jet.VarMap
	if isAdminTemplate(templateName) {
		vars = newPanelCommonVariables(c)
	} else {
		vars = commonVars(c)
	}
	vars.Set("Form", form)
	renderTemplate(c, templateName, vars)
}

func torrentTemplate(c *gin.Context, torrent model.TorrentJSON, rootFolder *filelist.FileListFolder, captchaID string) {
	vars := commonVars(c)
	vars.Set("Torrent", torrent)
	vars.Set("RootFolder", rootFolder)
	vars.Set("CaptchaID", captchaID)
	renderTemplate(c, "view.jet.html", vars)
}

func userProfileEditTemplate(c *gin.Context, userProfile *model.User, userForm userForms.UserForm, languages map[string]string) {
	vars := commonVars(c)
	vars.Set("UserProfile", userProfile)
	vars.Set("UserForm", userForm)
	vars.Set("Languages", languages)
	renderTemplate(c, "user/profile_edit.jet.html", vars)
}

func userProfileTemplate(c *gin.Context, userProfile *model.User) {
	vars := commonVars(c)
	vars.Set("UserProfile", userProfile)
	renderTemplate(c, "user/profile.jet.html", vars)
}
func databaseDumpTemplate(c *gin.Context, listDumps []model.DatabaseDumpJSON, GPGLink string) {
	vars := commonVars(c)
	vars.Set("ListDumps", listDumps)
	vars.Set("GPGLink", GPGLink)
	renderTemplate(c, "dumps.jet.html", vars)
}
func changeLanguageTemplate(c *gin.Context, language string, languages map[string]string) {
	vars := commonVars(c)
	vars.Set("Language", language)
	vars.Set("Languages", languages)
	renderTemplate(c, "user/public_settings.jet.html", vars)
}

func panelAdminTemplate(c *gin.Context, torrent []model.Torrent, reports []model.TorrentReportJSON, users []model.User, comments []model.Comment) {
	vars := newPanelCommonVariables(c)
	vars.Set("Torrent", torrent)
	vars.Set("TorrentReports", reports)
	vars.Set("Users", users)
	vars.Set("Comments", comments)
	renderTemplate(c, "admin/index.jet.html", vars)
}

func isAdminTemplate(templateName string) bool {
	if templateName != "" && len(templateName) > len(ModeratorDir) {
		return templateName[:5] == ModeratorDir
	}
	return false
}
