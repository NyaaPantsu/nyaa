package router

import (
	"net/http"

	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/service/user"
	"github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/gorilla/mux"
)

// Getting View User Registration
func UserRegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	_, errorUser := userService.CurrentUser(r)
	if errorUser != nil {
		b := form.RegistrationForm{}
		modelHelper.BindValueForm(&b, r)
		b.CaptchaID = captcha.GetID()
		languages.SetTranslationFromRequest(viewRegisterTemplate, r, "en-us")
		htv := UserRegisterTemplateVariables{b, form.NewErrors(), NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}
		err := viewRegisterTemplate.ExecuteTemplate(w, "index.html", htv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		HomeHandler(w, r)
	}
}

// Getting View User Login
func UserLoginFormHandler(w http.ResponseWriter, r *http.Request) {
	b := form.LoginForm{}
	modelHelper.BindValueForm(&b, r)

	languages.SetTranslationFromRequest(viewLoginTemplate, r, "en-us")
	htv := UserLoginFormVariables{b, form.NewErrors(), NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}

	err := viewLoginTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Getting User Profile
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil {
		currentUser := GetUser(r)
		view := r.URL.Query()["edit"]
		deleteVar := r.URL.Query()["delete"]
		if (view != nil) && (userPermission.CurrentOrAdmin(currentUser, userProfile.ID)) {
			b := form.UserForm{}
			modelHelper.BindValueForm(&b, r)
			languages.SetTranslationFromRequest(viewProfileEditTemplate, r, "en-us")
			htv := UserProfileEditVariables{&userProfile, b, form.NewErrors(), form.NewInfos(), NewSearchForm(), Navigation{}, currentUser, r.URL, mux.CurrentRoute(r)}

			err := viewProfileEditTemplate.ExecuteTemplate(w, "index.html", htv)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		} else if (deleteVar != nil) && (userPermission.CurrentOrAdmin(currentUser, userProfile.ID)) {
			err := form.NewErrors()
			_, errUser := userService.DeleteUser(w, currentUser, id)
			if errUser != nil {
				err["errors"] = append(err["errors"], errUser.Error())
			}
			languages.SetTranslationFromRequest(viewUserDeleteTemplate, r, "en-us")
			htv := UserVerifyTemplateVariables{err, NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}
			errorTmpl := viewUserDeleteTemplate.ExecuteTemplate(w, "index.html", htv)
			if errorTmpl != nil {
				http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
			}
		} else {
			languages.SetTranslationFromRequest(viewProfileTemplate, r, "en-us")
			htv := UserProfileVariables{&userProfile, NewSearchForm(), Navigation{}, currentUser, r.URL, mux.CurrentRoute(r)}

			err := viewProfileTemplate.ExecuteTemplate(w, "index.html", htv)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		searchForm := NewSearchForm()
		searchForm.HideAdvancedSearch = true

		languages.SetTranslationFromRequest(notFoundTemplate, r, "en-us")
		err := notFoundTemplate.ExecuteTemplate(w, "index.html", NotFoundTemplateVariables{Navigation{}, searchForm, GetUser(r), r.URL, mux.CurrentRoute(r)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Getting View User Profile Update
func UserProfileFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := GetUser(r)
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil {
		if userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
			b := form.UserForm{}
			err := form.NewErrors()
			infos := form.NewInfos()
			T := languages.SetTranslationFromRequest(viewProfileEditTemplate, r, "en-us")
			if len(r.PostFormValue("email")) > 0 {
				_, err = form.EmailValidation(r.PostFormValue("email"), err)
			}
			if len(r.PostFormValue("username")) > 0 {
				_, err = form.ValidateUsername(r.PostFormValue("username"), err)
			}
			if len(err) == 0 {
				modelHelper.BindValueForm(&b, r)
				err = modelHelper.ValidateForm(&b, err)
				if len(err) == 0 {
					userProfile, _, errorUser = userService.UpdateUser(w, &b, currentUser, id)
					if errorUser != nil {
						err["errors"] = append(err["errors"], errorUser.Error())
					}
					if len(err) == 0 {
						infos["infos"] = append(infos["infos"], T("profile_updated"))
					}
				}
			}
			htv := UserProfileEditVariables{&userProfile, b, err, infos, NewSearchForm(), Navigation{}, currentUser, r.URL, mux.CurrentRoute(r)}
			errorTmpl := viewProfileEditTemplate.ExecuteTemplate(w, "index.html", htv)
			if errorTmpl != nil {
				http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
			}
		} else {
			searchForm := NewSearchForm()
			searchForm.HideAdvancedSearch = true

			languages.SetTranslationFromRequest(notFoundTemplate, r, "en-us")
			err := notFoundTemplate.ExecuteTemplate(w, "index.html", NotFoundTemplateVariables{Navigation{}, searchForm, GetUser(r), r.URL, mux.CurrentRoute(r)})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		searchForm := NewSearchForm()
		searchForm.HideAdvancedSearch = true

		languages.SetTranslationFromRequest(notFoundTemplate, r, "en-us")
		err := notFoundTemplate.ExecuteTemplate(w, "index.html", NotFoundTemplateVariables{Navigation{}, searchForm, GetUser(r), r.URL, mux.CurrentRoute(r)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Post Registration controller, we do some check on the form here, the rest on user service
func UserRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	// Check same Password
	b := form.RegistrationForm{}
	err := form.NewErrors()
	if !captcha.Authenticate(captcha.Extract(r)) {
		err["errors"] = append(err["errors"], "Wrong captcha!")
	}
	if len(err) == 0 {
		if len(r.PostFormValue("email")) > 0 {
			_, err = form.EmailValidation(r.PostFormValue("email"), err)
		}
		_, err = form.ValidateUsername(r.PostFormValue("username"), err)
		if len(err) == 0 {
			modelHelper.BindValueForm(&b, r)
			err = modelHelper.ValidateForm(&b, err)
			if len(err) == 0 {
				_, errorUser := userService.CreateUser(w, r)
				if errorUser != nil {
					err["errors"] = append(err["errors"], errorUser.Error())
				}
				if len(err) == 0 {
					languages.SetTranslationFromRequest(viewRegisterSuccessTemplate, r, "en-us")
					u := model.User{
						Email: r.PostFormValue("email"), // indicate whether user had email set
					}
					htv := UserRegisterTemplateVariables{b, err, NewSearchForm(), Navigation{}, &u, r.URL, mux.CurrentRoute(r)}
					errorTmpl := viewRegisterSuccessTemplate.ExecuteTemplate(w, "index.html", htv)
					if errorTmpl != nil {
						http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
					}
				}
			}
		}
	}
	if len(err) > 0 {
		b.CaptchaID = captcha.GetID()
		languages.SetTranslationFromRequest(viewRegisterTemplate, r, "en-us")
		htv := UserRegisterTemplateVariables{b, err, NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}
		errorTmpl := viewRegisterTemplate.ExecuteTemplate(w, "index.html", htv)
		if errorTmpl != nil {
			http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
		}
	}
}

func UserVerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	err := form.NewErrors()
	_, errEmail := userService.EmailVerification(token, w)
	if errEmail != nil {
		err["errors"] = append(err["errors"], errEmail.Error())
	}
	languages.SetTranslationFromRequest(viewVerifySuccessTemplate, r, "en-us")
	htv := UserVerifyTemplateVariables{err, NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}
	errorTmpl := viewVerifySuccessTemplate.ExecuteTemplate(w, "index.html", htv)
	if errorTmpl != nil {
		http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
	}
}

// Post Login controller
func UserLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	b := form.LoginForm{}
	modelHelper.BindValueForm(&b, r)
	err := form.NewErrors()
	err = modelHelper.ValidateForm(&b, err)
	if len(err) == 0 {
		_, errorUser := userService.CreateUserAuthentication(w, r)
		if errorUser != nil {
			err["errors"] = append(err["errors"], errorUser.Error())
			languages.SetTranslationFromRequest(viewLoginTemplate, r, "en-us")
			htv := UserLoginFormVariables{b, err, NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}
			errorTmpl := viewLoginTemplate.ExecuteTemplate(w, "index.html", htv)
			if errorTmpl != nil {
				http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
			}
		} else {
			url, _ := Router.Get("home").URL()
			http.Redirect(w, r, url.String(), http.StatusSeeOther)
		}
	}
}

// Logout
func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = userService.ClearCookie(w)
	url, _ := Router.Get("home").URL()
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}
