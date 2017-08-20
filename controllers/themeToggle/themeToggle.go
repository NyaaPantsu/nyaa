package themeToggleController

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/timeHelper"
	"github.com/gin-gonic/gin"
	
)

// toggleThemeHandler : Controller to switch between theme1 & theme2
func toggleThemeHandler(c *gin.Context) {

	theme, err := c.Cookie("theme")
	if err != nil {
		theme = "g"
	}
	theme2, err := c.Cookie("theme2")
	if err != nil {
		theme2 = "tomorrow"
	}
	if theme == theme2 {
		if theme == "tomorrow" {
			theme2 = "g"
		}
		if theme != "tomorrow" {
			theme2 = "tomorrow"	
		}
	}
	//Get theme & theme2 value, if not set by default is g.css & tomorrow.css
	//if both theme are identical, which can happen, we fix it
	
	//Set value of redirect url, add #footer to get user back to same page position
	redirectUrl = "https://nyaa.pantsu.cat/"
	fmt.Sprintf(redirectUrl, "%s#footer", redirectUrl)

	//Switch theme & thele2 value
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme", Value: theme2, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme2", Value: theme, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})	
	
	//redirect user to page
	c.Redirect(http.StatusSeeOther, redirectUrl)
	return
}

func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
