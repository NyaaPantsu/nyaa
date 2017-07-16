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

func init() {
	if config.Get().Environment == "DEVELOPMENT" {
		View.SetDevelopmentMode(true)
		fmt.Println("Template Live Update enabled")
	}
}
func commonvariables(c *gin.Context) jet.VarMap {
	token := nosurf.Token(c.Request)
	msg := messages.GetMessages(c)
	variables := templateFunctions(make(jet.VarMap))
	variables.Set("Navigation", newNavigation())
	variables.Set("Search", newSearchForm(c))
	variables.Set("T", publicSettings.GetTfuncFromRequest(c))
	variables.Set("Theme", publicSettings.GetThemeFromRequest(c))
	variables.Set("Mascot", publicSettings.GetMascotFromRequest(c))
	variables.Set("MascotURL", publicSettings.GetMascotUrlFromRequest(c))
	variables.Set("User", getUser(c))
	variables.Set("URL", c.Request.URL)
	variables.Set("CsrfToken", token)
	variables.Set("Config", config.Get())
	variables.Set("Infos", msg.GetAllInfos())
	variables.Set("Errors", msg.GetAllErrors())
	return variables
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
	common := commonvariables(c)
	common.Set("Search", newPanelSearchForm(c))
	return common
}

func renderTemplate(c *gin.Context, templateName string, variables jet.VarMap) {
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
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = newPanelCommonVariables(c)
	} else {
		variables = commonvariables(c)
	}
	renderTemplate(c, templateName, variables)
}

func modelList(c *gin.Context, templateName string, models interface{}, nav navigation, search searchForm) {
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = newPanelCommonVariables(c)
	} else {
		variables = commonvariables(c)
	}
	variables.Set("Models", models)
	variables.Set("Navigation", nav)
	variables.Set("Search", search)
	renderTemplate(c, templateName, variables)
}

func formTemplate(c *gin.Context, templateName string, form interface{}) {
	var variables jet.VarMap
	if isAdminTemplate(templateName) {
		variables = newPanelCommonVariables(c)
	} else {
		variables = commonvariables(c)
	}
	variables.Set("Form", form)
	renderTemplate(c, templateName, variables)
}

func torrentTemplate(c *gin.Context, torrent models.TorrentJSON, rootFolder *filelist.FileListFolder, captchaID string) {
	variables := commonvariables(c)
	variables.Set("Torrent", torrent)
	variables.Set("RootFolder", rootFolder)
	variables.Set("CaptchaID", captchaID)
	renderTemplate(c, path.Join(SiteDir, "torrents/view.jet.html"), variables)
}

func userProfileEditTemplate(c *gin.Context, userProfile *models.User, userForm userValidator.UserForm, languages publicSettings.Languages) {
	variables := commonvariables(c)
	variables.Set("UserProfile", userProfile)
	variables.Set("UserForm", userForm)
	variables.Set("Languages", languages)
	renderTemplate(c, path.Join(SiteDir, "user/edit.jet.html"), variables)
}

func userProfileTemplate(c *gin.Context, userProfile *models.User) {
	variables := commonvariables(c)
	variables.Set("UserProfile", userProfile)
	renderTemplate(c, path.Join(SiteDir, "user/torrents.jet.html"), variables)
}

func userProfileNotificationsTemplate(c *gin.Context, userProfile *models.User) {
	variables := commonvariables(c)
	variables.Set("UserProfile", userProfile)
	renderTemplate(c, path.Join(SiteDir, "user/notifications.jet.html"), variables)
}
func databaseDumpTemplate(c *gin.Context, listDumps []models.DatabaseDumpJSON, GPGLink string) {
	variables := commonvariables(c)
	variables.Set("ListDumps", listDumps)
	variables.Set("GPGLink", GPGLink)
	renderTemplate(c, path.Join(SiteDir, "database/dumps.jet.html"), variables)
}
func panelAdminTemplate(c *gin.Context, torrent []models.Torrent, reports []models.TorrentReportJSON, users []models.User, comments []models.Comment) {
	variables := newPanelCommonVariables(c)
	variables.Set("Torrents", torrent)
	variables.Set("TorrentReports", reports)
	variables.Set("Users", users)
	variables.Set("Comments", comments)
	renderTemplate(c, path.Join(ModeratorDir, "index.jet.html"), variables)
}

func isAdminTemplate(templateName string) bool {
	if templateName != "" && len(templateName) > len(ModeratorDir) {
		return templateName[:5] == ModeratorDir
	}
	return false
}
