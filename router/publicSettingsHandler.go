package router

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/NyaaPantsu/nyaa/util/timeHelper"
)

// LanguagesJSONResponse : Structure containing all the languages to parse it as a JSON response
type LanguagesJSONResponse struct {
	Current   string            `json:"current"`
	Languages map[string]string `json:"languages"`
}

// SeePublicSettingsHandler : Controller to view the languages and themes
func SeePublicSettingsHandler(w http.ResponseWriter, r *http.Request) {
	_, Tlang := publicSettings.GetTfuncAndLanguageFromRequest(r)
	availableLanguages := publicSettings.GetAvailableLanguages()
	defer r.Body.Close()
	contentType := r.Header.Get("Content-Type")
	if contentType == "application/json" {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(LanguagesJSONResponse{Tlang.Tag, availableLanguages})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		clv := changeLanguageVariables{
			commonTemplateVariables: newCommonVariables(r),
			Language:                Tlang.Tag,
			Languages:               availableLanguages,
		}
		err := changePublicSettingsTemplate.ExecuteTemplate(w, "index.html", clv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// ChangePublicSettingsHandler : Controller for changing the current language and theme
func ChangePublicSettingsHandler(w http.ResponseWriter, r *http.Request) {
	theme := r.FormValue("theme")
	lang := r.FormValue("language")
	mascot := r.FormValue("mascot")
	mascotURL := r.FormValue("mascot_url")

	availableLanguages := publicSettings.GetAvailableLanguages()
	defer r.Body.Close()
	if _, exists := availableLanguages[lang]; !exists {
		http.Error(w, "Language not available", http.StatusInternalServerError)
		return
	}
	// FIXME Are the settings actually sanitized?
	// Limit the mascot URL, so base64-encoded images aren't valid
	if len(mascotURL) > 256 {
		http.Error(w, "Mascot URL is too long (max is 255 chars)", http.StatusInternalServerError)
		return
	}

	_, err := url.Parse(mascotURL)
	if err != nil {
		http.Error(w, "Error parsing mascot URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// If logged in, update user settings.
	user, _ := userService.CurrentUser(r)
	if user.ID > 0 {
		user.Language = lang
		user.Theme = theme
		user.Mascot = mascot
		user.MascotURL = mascotURL
		// I don't know if I should use this...
		userService.UpdateRawUser(&user)
	}
	// Set cookie
	http.SetCookie(w, &http.Cookie{Name: "lang", Value: lang, Domain: getDomainName(), Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(w, &http.Cookie{Name: "theme", Value: theme, Domain: getDomainName(), Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(w, &http.Cookie{Name: "mascot", Value: mascot, Domain: getDomainName(), Expires: timeHelper.FewDaysLater(365)})
	http.SetCookie(w, &http.Cookie{Name: "mascot_url", Value: mascotURL, Domain: getDomainName(), Expires: timeHelper.FewDaysLater(365)})

	url, _ := Router.Get("home").URL()
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}
func getDomainName() string {
	domain := config.Conf.Cookies.DomainName
	if config.Conf.Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}
