package router

import (
	"encoding/json"
	"net/http"

	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/NyaaPantsu/nyaa/util/timeHelper"
)

type LanguagesJSONResponse struct {
	Current   string            `json:"current"`
	Languages map[string]string `json:"languages"`
}

func SeeLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	_, Tlang := languages.GetTfuncAndLanguageFromRequest(r)
	availableLanguages := languages.GetAvailableLanguages()

	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(LanguagesJSONResponse{Tlang.Tag, availableLanguages})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		clv := ChangeLanguageVariables{
			CommonTemplateVariables: NewCommonVariables(r),
			Language:   Tlang.Tag,
			Languages:  availableLanguages,
		}
		err := changeLanguageTemplate.ExecuteTemplate(w, "index.html", clv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func ChangeLanguageHandler(w http.ResponseWriter, r *http.Request) {
	lang := r.FormValue("language")

	availableLanguages := languages.GetAvailableLanguages()
	if _, exists := availableLanguages[lang]; !exists {
		http.Error(w, "Language not available", http.StatusInternalServerError)
		return
	}

	// If logged in, update user language; if not, set cookie.
	user, _ := userService.CurrentUser(r)
	if user.ID > 0 { 
		user.Language = lang
		// I don't know if I should use this...
		userService.UpdateUserCore(&user)
	}
	http.SetCookie(w, &http.Cookie{Name: "lang", Value: lang, Expires: timeHelper.FewDaysLater(365)})

	url, _ := Router.Get("home").URL()
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}
