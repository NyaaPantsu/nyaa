package controllers

import (
	"net/http"
	"path"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/filelist"
	"github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/gin-gonic/gin"
	"github.com/justinas/nosurf"

	"fmt"

	"github.com/CloudyKit/jet"
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
var vars = templateFunctions(make(jet.VarMap))

func init() {
	if config.Conf.Environment == "DEVELOPMENT" {
		View.SetDevelopmentMode(true)
		fmt.Println("Template Live Update enabled")
	}
}
func commonVars(c *gin.Context) jet.VarMap {
	token := nosurf.Token(c.Request)
	msg := messages.GetMessages(c)
	vars.Set("Navigation", newNavigation())
	vars.Set("Search", newSearchForm(c))
	vars.Set("T", publicSettings.GetTfuncFromRequest(c))
	vars.Set("Theme", publicSettings.GetThemeFromRequest(c))
	vars.Set("Mascot", publicSettings.GetMascotFromRequest(c))
	vars.Set("MascotURL", publicSettings.GetMascotUrlFromRequest(c))
	vars.Set("User", getUser(c))
	vars.Set("URL", c.Request.URL)
	vars.Set("CsrfToken", token)
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
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if err = t.Execute(c.Writer, vars, nil); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

func httpError(c *gin.Context, errorCode int) {
	switch errorCode {
	case http.StatusNotFound:
		staticTemplate(c, path.Join(ErrorsDir, "404.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	case http.StatusBadRequest:
		staticTemplate(c, path.Join(ErrorsDir, "400.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	case http.StatusInternalServerError:
		staticTemplate(c, path.Join(ErrorsDir, "500.jet.html"))
		c.AbortWithStatus(errorCode)
		return
	}
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

func torrentTemplate(c *gin.Context, torrent models.TorrentJSON, rootFolder *filelist.FileListFolder, captchaID string) {
	vars := commonVars(c)
	vars.Set("Torrent", torrent)
	vars.Set("RootFolder", rootFolder)
	vars.Set("CaptchaID", captchaID)
	renderTemplate(c, path.Join(SiteDir, "torrents/view.jet.html"), vars)
}

func userProfileEditTemplate(c *gin.Context, userProfile *models.User, userForm userValidator.UserForm, languages publicSettings.Languages) {
	vars := commonVars(c)
	vars.Set("UserProfile", userProfile)
	vars.Set("UserForm", userForm)
	vars.Set("Languages", languages)
	renderTemplate(c, path.Join(SiteDir, "user/edit.jet.html"), vars)
}

func userProfileTemplate(c *gin.Context, userProfile *models.User) {
	vars := commonVars(c)
	vars.Set("UserProfile", userProfile)
	renderTemplate(c, path.Join(SiteDir, "user/torrents.jet.html"), vars)
}

func userProfileNotificationsTemplate(c *gin.Context, userProfile *models.User) {
	vars := commonVars(c)
	vars.Set("UserProfile", userProfile)
	renderTemplate(c, path.Join(SiteDir, "user/notifications.jet.html"), vars)
}
func databaseDumpTemplate(c *gin.Context, listDumps []models.DatabaseDumpJSON, GPGLink string) {
	vars := commonVars(c)
	vars.Set("ListDumps", listDumps)
	vars.Set("GPGLink", GPGLink)
	renderTemplate(c, path.Join(SiteDir, "database/dumps.jet.html"), vars)
}
func panelAdminTemplate(c *gin.Context, torrent []models.Torrent, reports []models.TorrentReportJSON, users []models.User, comments []models.Comment) {
	vars := newPanelCommonVariables(c)
	vars.Set("Torrents", torrent)
	vars.Set("TorrentReports", reports)
	vars.Set("Users", users)
	vars.Set("Comments", comments)
	renderTemplate(c, path.Join(ModeratorDir, "index.jet.html"), vars)
}

func isAdminTemplate(templateName string) bool {
	if templateName != "" && len(templateName) > len(ModeratorDir) {
		return templateName[:5] == ModeratorDir
	}
	return false
}
