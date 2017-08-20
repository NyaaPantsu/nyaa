package themeToggleController

import (
	"net/http"
	"log"

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
	//Get theme1 & theme2 value, set g.css & tomorrow.css by default

	//Switch theme & theme2 value
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme", Value: theme2, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(c.Writer, &http.Cookie{Name: "theme2", Value: theme, Domain: getDomainName(), Path: "/", Expires: timeHelper.FewDaysLater(365)})	
	
	//redirect user to page he was beforehand
	c.Redirect(http.StatusSeeOther, c.Param("redirect"))
	return
}

func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
