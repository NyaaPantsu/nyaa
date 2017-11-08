package settingsController

import (
	"net/http"
	"net/url"
	"time"

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
	lang := c.PostForm("lang")
	mascot := c.PostForm("mascot")
	mascotURL := c.PostForm("mascot_url")
	altColors := c.PostForm("altColors")
	oldNav := c.PostForm("oldNav")

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
		user.OldNav = oldNav
		user.UpdateRaw()
	}
	
	if getDomainName() != "" {
		//Clear every cookie from current domain so that users that old cookies from current domain do not interfere with new ones, should new one be shared within multiple subdomains
		http.SetCookie(c.Writer, &http.Cookie{Name: "lang", Value: "", Domain: "", Path: "/", Expires: time.Now().AddDate(-1, -1, -1)})
		http.SetCookie(c.Writer, &http.Cookie{Name: "theme", Value: "", Domain: "", Path: "/", Expires: time.Now().AddDate(-1, -1, -1)})
		http.SetCookie(c.Writer, &http.Cookie{Name: "mascot", Value: "", Domain: "", Path: "/", Expires: time.Now().AddDate(-1, -1, -1)})
		http.SetCookie(c.Writer, &http.Cookie{Name: "mascot_url", Value: "", Domain: "", Path: "/", Expires: time.Now().AddDate(-1, -1, -1)})
		http.SetCookie(c.Writer, &http.Cookie{Name: "oldNav", Value: "", Domain: "", Path: "/", Expires: time.Now().AddDate(-1, -1, -1)})
	}
	
	
	// Set cookie with http and not gin for expires (maxage not supported in <IE8)
	http.SetCookie(c.Writer, &http.Cookie{Name: "lang", Value: lang, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme", Value: theme, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "mascot", Value: mascot, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "mascot_url", Value: mascotURL, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "oldNav", Value: oldNav, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})

	c.Redirect(http.StatusSeeOther, "/")
}
func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
