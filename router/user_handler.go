package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/NyaaPantsu/nyaa/util/modelHelper"
	"github.com/gorilla/mux"
)

// Getting View User Registration
func UserRegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	_, errorUser := userService.CurrentUser(r)
	// User is already connected, redirect to home
	if errorUser == nil {
		HomeHandler(w, r)
		return
	}
	registrationForm := form.RegistrationForm{}
	modelHelper.BindValueForm(&registrationForm, r)
	registrationForm.CaptchaID = captcha.GetID()
	urtv := UserRegisterTemplateVariables{
		RegistrationForm: registrationForm,
		FormErrors:       form.NewErrors(),
		Search:           NewSearchForm(),
		Navigation:       NewNavigation(),
		T:                languages.GetTfuncFromRequest(r),
		User:             GetUser(r),
		URL:              r.URL,
		Route:            mux.CurrentRoute(r),
	}
	err := viewRegisterTemplate.ExecuteTemplate(w, "index.html", urtv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Getting View User Login
func UserLoginFormHandler(w http.ResponseWriter, r *http.Request) {
	loginForm := form.LoginForm{}
	modelHelper.BindValueForm(&loginForm, r)

	ulfv := UserLoginFormVariables{
		LoginForm:  loginForm,
		FormErrors: form.NewErrors(),
		Search:     NewSearchForm(),
		Navigation: NewNavigation(),
		T:          languages.GetTfuncFromRequest(r),
		User:       GetUser(r),
		URL:        r.URL,
		Route:      mux.CurrentRoute(r),
	}

	err := viewLoginTemplate.ExecuteTemplate(w, "index.html", ulfv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Getting User Profile
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	T := languages.GetTfuncFromRequest(r)
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil {
		currentUser := GetUser(r)
		follow := r.URL.Query()["followed"]
		unfollow := r.URL.Query()["unfollowed"]
		infosForm := form.NewInfos()
		deleteVar := r.URL.Query()["delete"]

		if (deleteVar != nil) && (userPermission.CurrentOrAdmin(currentUser, userProfile.ID)) {
			err := form.NewErrors()
			_, errUser := userService.DeleteUser(w, currentUser, id)
			if errUser != nil {
				err["errors"] = append(err["errors"], errUser.Error())
			}
			htv := UserVerifyTemplateVariables{err, NewSearchForm(), NewNavigation(), T, GetUser(r), r.URL, mux.CurrentRoute(r)}
			errorTmpl := viewUserDeleteTemplate.ExecuteTemplate(w, "index.html", htv)
			if errorTmpl != nil {
				http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
			}
		} else {
			if follow != nil {
				infosForm["infos"] = append(infosForm["infos"], fmt.Sprintf(string(T("user_followed_msg")), userProfile.Username))
			}
			if unfollow != nil {
				infosForm["infos"] = append(infosForm["infos"], fmt.Sprintf(string(T("user_unfollowed_msg")), userProfile.Username))
			}
			htv := UserProfileVariables{&userProfile, infosForm, NewSearchForm(), NewNavigation(), T, currentUser, r.URL, mux.CurrentRoute(r)}

			err := viewProfileTemplate.ExecuteTemplate(w, "index.html", htv)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		err := notFoundTemplate.ExecuteTemplate(w, "index.html", NotFoundTemplateVariables{NewNavigation(), NewSearchForm(), T, GetUser(r), r.URL, mux.CurrentRoute(r)})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

//Getting User Profile Details View
func UserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := GetUser(r)
	T := languages.GetTfuncFromRequest(r)
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil && userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
		if userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
			b := form.UserForm{}
			modelHelper.BindValueForm(&b, r)
			availableLanguages := languages.GetAvailableLanguages()
			htv := UserProfileEditVariables{&userProfile, b, form.NewErrors(), form.NewInfos(), availableLanguages, NewSearchForm(), NewNavigation(), T, currentUser, r.URL, mux.CurrentRoute(r)}
			err := viewProfileEditTemplate.ExecuteTemplate(w, "index.html", htv)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		err := notFoundTemplate.ExecuteTemplate(w, "index.html", NotFoundTemplateVariables{NewNavigation(), NewSearchForm(), T, GetUser(r), r.URL, mux.CurrentRoute(r)})
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
	if errorUser != nil || !userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
		NotFoundHandler(w, r)
		return
	}

	userForm := form.UserForm{}
	err := form.NewErrors()
	infos := form.NewInfos()
	T := languages.GetTfuncFromRequest(r)
	if len(r.PostFormValue("email")) > 0 {
		_, err = form.EmailValidation(r.PostFormValue("email"), err)
	}
	if len(r.PostFormValue("username")) > 0 {
		_, err = form.ValidateUsername(r.PostFormValue("username"), err)
	}

	if len(err) == 0 {
		modelHelper.BindValueForm(&userForm, r)
		if !userPermission.HasAdmin(currentUser) {
			userForm.Username = userProfile.Username
			userForm.Status = userProfile.Status
		} else {
			if userProfile.Status != userForm.Status && userForm.Status == 2 {
				err["errors"] = append(err["errors"], "Elevating status to moderator is prohibited")
			}
		}
		err = modelHelper.ValidateForm(&userForm, err)
		if len(err) == 0 {
			if userForm.Email != userProfile.Email {
				userService.SendVerificationToUser(*currentUser, userForm.Email)
				infos["infos"] = append(infos["infos"], fmt.Sprintf(string(T("email_changed")), userForm.Email))
				userForm.Email = userProfile.Email // reset, it will be set when user clicks verification
			}
			userProfile, _, errorUser = userService.UpdateUser(w, &userForm, currentUser, id)
			if errorUser != nil {
				err["errors"] = append(err["errors"], errorUser.Error())
			} else {
				infos["infos"] = append(infos["infos"], string(T("profile_updated")))
			}
		}
	}
	availableLanguages := languages.GetAvailableLanguages()
	upev := UserProfileEditVariables{
		UserProfile: &userProfile,
		UserForm:    userForm,
		FormErrors:  err,
		FormInfos:   infos,
		Languages:   availableLanguages,
		Search:      NewSearchForm(),
		Navigation:  NewNavigation(),
		T:           T,
		User:        currentUser,
		URL:         r.URL,
		Route:       mux.CurrentRoute(r),
	}
	errorTmpl := viewProfileEditTemplate.ExecuteTemplate(w, "index.html", upev)
	if errorTmpl != nil {
		http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
	}
}

// Post Registration controller, we do some check on the form here, the rest on user service
func UserRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	T := languages.GetTfuncFromRequest(r)
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
					u := model.User{
						Email: r.PostFormValue("email"), // indicate whether user had email set
					}
					htv := UserRegisterTemplateVariables{b, err, NewSearchForm(), NewNavigation(), T, &u, r.URL, mux.CurrentRoute(r)}
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
		htv := UserRegisterTemplateVariables{b, err, NewSearchForm(), NewNavigation(), T, GetUser(r), r.URL, mux.CurrentRoute(r)}
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
	T := languages.GetTfuncFromRequest(r)
	htv := UserVerifyTemplateVariables{err, NewSearchForm(), NewNavigation(), T, GetUser(r), r.URL, mux.CurrentRoute(r)}
	errorTmpl := viewVerifySuccessTemplate.ExecuteTemplate(w, "index.html", htv)
	if errorTmpl != nil {
		http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
	}
}

// Post Login controller
func UserLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	T := languages.GetTfuncFromRequest(r)
	b := form.LoginForm{}
	modelHelper.BindValueForm(&b, r)
	err := form.NewErrors()
	err = modelHelper.ValidateForm(&b, err)
	if len(err) == 0 {
		_, errorUser := userService.CreateUserAuthentication(w, r)
		if errorUser != nil {
			err["errors"] = append(err["errors"], errorUser.Error())
			htv := UserLoginFormVariables{b, err, NewSearchForm(), NewNavigation(), T, GetUser(r), r.URL, mux.CurrentRoute(r)}
			errorTmpl := viewLoginTemplate.ExecuteTemplate(w, "index.html", htv)
			if errorTmpl != nil {
				http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
			}
			return
		} else {
			url, _ := Router.Get("home").URL()
			http.Redirect(w, r, url.String(), http.StatusSeeOther)
		}
	}
	if len(err) > 0 {
		htv := UserLoginFormVariables{b, err, NewSearchForm(), NewNavigation(), T, GetUser(r), r.URL, mux.CurrentRoute(r)}
		errorTmpl := viewLoginTemplate.ExecuteTemplate(w, "index.html", htv)
		if errorTmpl != nil {
			http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
		}
	}
}

// Logout
func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = userService.ClearCookie(w)
	url, _ := Router.Get("home").URL()
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}

func UserFollowHandler(w http.ResponseWriter, r *http.Request) {
	var followAction string
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := GetUser(r)
	user, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil {
		if !userPermission.IsFollower(&user, currentUser) {
			followAction = "followed"
			userService.SetFollow(&user, currentUser)
		} else {
			followAction = "unfollowed"
			userService.RemoveFollow(&user, currentUser)
		}
	}
	url, _ := Router.Get("user_profile").URL("id", strconv.Itoa(int(user.ID)), "username", user.Username)
	http.Redirect(w, r, url.String()+"?"+followAction, http.StatusSeeOther)
}
