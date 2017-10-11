package themeToggleController

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/utils/timeHelper"
	"github.com/gin-gonic/gin"
	
)

// toggleThemeHandler : Controller to switch between theme1 & theme2
func toggleThemeHandler(c *gin.Context) {

	theme, err := c.Cookie("theme")
	if err != nil {
		theme = config.DefaultTheme(false)
	}
	theme2, err := c.Cookie("theme2")
	if err != nil {
		theme2 = config.DefaultTheme(true)
	}
	if theme != config.DefaultTheme(true) && theme2 != config.DefaultTheme(true) {
		//None of the themes are dark ones, force the second one as the dark one
		theme2 = config.DefaultTheme(true)
	} else if  theme == theme2 {
		//Both theme are dark ones, force the second one as the default (light) theme
		theme2 = config.DefaultTheme(false)
	}
	//Get theme1 & theme2 value
	
	// If logged in, update user theme (will not work otherwise)
	user := router.GetUser(c)
	if user.ID > 0 {
		user.Theme = theme2
		user.UpdateRaw()
	}
	
	//Switch theme & theme2 value
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme", Value: theme2, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme2", Value: theme, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})	
	
	//Redirect user to page he was in beforehand
	c.Redirect(http.StatusSeeOther, c.Param("redirect") + "#footer")
	return
}

func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
