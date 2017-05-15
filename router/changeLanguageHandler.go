package router;

import (
	"encoding/json"
	"net/http"

	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/service/user"
	"github.com/gorilla/mux"
)

type LanguagesJSONResponse struct {
	Current   string `json:"current"`
	Languages map[string]string `json:"languages"`
}

func SeeLanguagesHandler(w http.ResponseWriter, r *http.Request) {
	_, Tlang := languages.GetTfuncAndLanguageFromRequest(r)
	availableLanguages := languages.GetAvailableLanguages()

	format := r.URL.Query().Get("format")
	if format == "json" {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(LanguagesJSONResponse{Tlang.Tag, availableLanguages})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		clv := ChangeLanguageVariables{NewSearchForm(), Navigation{}, Tlang.Tag, availableLanguages, GetUser(r), r.URL, mux.CurrentRoute(r)}
		languages.SetTranslationFromRequest(changeLanguageTemplate, r)
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
	user, err := userService.CurrentUser(r)
	if err == nil {
		user.Language = lang
		// I don't know if I should use this...
		userService.UpdateUserCore(&user)
	} else {
		http.SetCookie(w, &http.Cookie{Name: "lang", Value: lang})
	}

	url, _ := Router.Get("home").URL()
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}

