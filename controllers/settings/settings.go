package settingsController

import (
	"net/http"
	"net/url"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/templates"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/timeHelper"
	"github.com/gin-gonic/gin"
)

// SeePublicSettingsHandler : Controller to view the languages and themes
func SeePublicSettingsHandler(c *gin.Context) {
	_, Tlang := publicSettings.GetTfuncAndLanguageFromRequest(c)
	availableLanguages := publicSettings.GetAvailableLanguages()
	languagesJSON := templates.LanguagesJSONResponse{Tlang.Tag, availableLanguages}
	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "application/json" {
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, languagesJSON)
	} else {
		templates.Form(c, "site/user/public/settings.jet.html", languagesJSON)
	}
}

// ChangePublicSettingsHandler : Controller for changing the current language and theme
func ChangePublicSettingsHandler(c *gin.Context) {
	theme := c.PostForm("theme")
	lang := c.PostForm("language")
	mascot := c.PostForm("mascot")
	mascotURL := c.PostForm("mascot_url")
	altColors := c.PostForm("altColors")
	hideAds := c.PostForm("hideAds")

	messages := msg.GetMessages(c)

	availableLanguages := publicSettings.GetAvailableLanguages()

	if !availableLanguages.Exist(lang) {
		messages.AddErrorT("errors", "language_not_available")
	}
	// FIXME Are the settings actually sanitized?
	// Limit the mascot URL, so base64-encoded images aren't valid
	if len(mascotURL) > 256 {
		messages.AddErrorT("errors", "mascot_url_too_long")
	}

	_, err := url.Parse(mascotURL)
	if err != nil {
		messages.AddErrorTf("errors", "mascor_url_parse_error", err.Error())
	}

	// If logged in, update user settings.
	user := router.GetUser(c)
	if user.ID > 0 {
		user.Language = lang
		user.Theme = theme
		user.Mascot = mascot
		user.MascotURL = mascotURL
		user.AltColors = altColors
		user.HideAds = hideAds
		user.UpdateRaw()
	}
	// Set cookie with http and not gin for expires (maxage not supported in <IE8)
	http.SetCookie(c.Writer, &http.Cookie{Name: "lang", Value: lang, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme", Value: theme, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "mascot", Value: mascot, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "mascot_url", Value: mascotURL, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "altColors", Value: altColors, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "hideAds", Value: hideAds, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})

	c.Redirect(http.StatusSeeOther, "/")
}
func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
