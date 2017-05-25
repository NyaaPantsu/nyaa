package router

import (
	"net/http"

	"github.com/NyaaPantsu/nyaa/service/user"
	//"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/NyaaPantsu/nyaa/util/timeHelper"
)

func SeeThemesHandler(w http.ResponseWriter, r *http.Request) {
	//Just render a static page for now
	ctv := ChangeThemeVariables{
		commonTemplateVariables: newCommonVariables(r),
	}
	changeThemeTemplate.ExecuteTemplate(w, "index.html", ctv)
	return
}

func ChangeThemeHandler(w http.ResponseWriter, r *http.Request) {
	theme := r.FormValue("theme")
        user, _ := userService.CurrentUser(r)
	// Update User if he's logged in
	if user.ID > 0 {
		user.Theme = theme
		// I don't know if I should use this...
		userService.UpdateUserCore(&user)
	}
	// Set cookie
        http.SetCookie(w, &http.Cookie{Name: "theme", Value: theme, Expires: timeHelper.FewDaysLater(365)})
	url, _ := Router.Get("home").URL()
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
	return
}
