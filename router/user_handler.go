package router

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/notifier"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/languages"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
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
	messages := msg.GetMessages(r)
	registrationForm := form.RegistrationForm{}
	modelHelper.BindValueForm(&registrationForm, r)
	registrationForm.CaptchaID = captcha.GetID()
	urtv := UserRegisterTemplateVariables{
		CommonTemplateVariables: NewCommonVariables(r),
		RegistrationForm:        registrationForm,
		FormErrors:              messages.GetAllErrors(),
	}
	err := viewRegisterTemplate.ExecuteTemplate(w, "index.html", urtv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Getting View User Login
func UserLoginFormHandler(w http.ResponseWriter, r *http.Request) {
	_, errorUser := userService.CurrentUser(r)
	// User is already connected, redirect to home
	if errorUser == nil {
		HomeHandler(w, r)
		return
	}

	loginForm := form.LoginForm{}
	modelHelper.BindValueForm(&loginForm, r)
	messages := msg.GetMessages(r)
	ulfv := UserLoginFormVariables{
		CommonTemplateVariables: NewCommonVariables(r),
		LoginForm:               loginForm,
		FormErrors:              messages.GetAllErrors(),
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
	messages := msg.GetMessages(r)

	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil {
		currentUser := GetUser(r)
		follow := r.URL.Query()["followed"]
		unfollow := r.URL.Query()["unfollowed"]
		deleteVar := r.URL.Query()["delete"]

		if (deleteVar != nil) && (userPermission.CurrentOrAdmin(currentUser, userProfile.ID)) {
			_, errUser := userService.DeleteUser(w, currentUser, id)
			if errUser != nil {
				messages.ImportFromError("errors", errUser)
			}
			htv := UserVerifyTemplateVariables{NewCommonVariables(r), messages.GetAllErrors()}
			errorTmpl := viewUserDeleteTemplate.ExecuteTemplate(w, "index.html", htv)
			if errorTmpl != nil {
				http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
			}
		} else {
			if follow != nil {
				messages.AddInfof("infos", string(T("user_followed_msg")), userProfile.Username)
			}
			if unfollow != nil {
				messages.AddInfof("infos", string(T("user_unfollowed_msg")), userProfile.Username)
			}
			htv := UserProfileVariables{NewCommonVariables(r), &userProfile, messages.GetAllInfos()}

			err := viewProfileTemplate.ExecuteTemplate(w, "index.html", htv)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		NotFoundHandler(w, r)
	}
}

//Getting User Profile Details View
func UserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := GetUser(r)
	messages := msg.GetMessages(r)

	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil && userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
		if userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
			b := form.UserForm{}
			modelHelper.BindValueForm(&b, r)
			availableLanguages := languages.GetAvailableLanguages()
			htv := UserProfileEditVariables{NewCommonVariables(r), &userProfile, b, messages.GetAllErrors(), messages.GetAllInfos(), availableLanguages}
			err := viewProfileEditTemplate.ExecuteTemplate(w, "index.html", htv)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		NotFoundHandler(w, r)
	}
}

// Getting View User Profile Update
func UserProfileFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := GetUser(r)
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser != nil || !userPermission.CurrentOrAdmin(currentUser, userProfile.ID) || userProfile.ID == 0 {
		NotFoundHandler(w, r)
		return
	}
	messages := msg.GetMessages(r)
	userForm := form.UserForm{}
	userSettingsForm := form.UserSettingsForm{}


	T := languages.GetTfuncFromRequest(r)
	if len(r.PostFormValue("email")) > 0 {
		form.EmailValidation(r.PostFormValue("email"), &messages)
	}
	if len(r.PostFormValue("username")) > 0 {
		form.ValidateUsername(r.PostFormValue("username"), &messages)
	}

	if !messages.HasErrors() {
		modelHelper.BindValueForm(&userForm, r)
		modelHelper.BindValueForm(&userSettingsForm, r)
		if !userPermission.HasAdmin(currentUser) {
			userForm.Username = userProfile.Username
			userForm.Status = userProfile.Status
		} else {
			if userProfile.Status != userForm.Status && userForm.Status == 2 {
				messages.AddError("errors", "Elevating status to moderator is prohibited")
			}
		}
		modelHelper.ValidateForm(&userForm, &messages)
		if !messages.HasErrors() {
			if userForm.Email != userProfile.Email {
				userService.SendVerificationToUser(*currentUser, userForm.Email)
				messages.AddInfof("infos", string(T("email_changed")), userForm.Email)
				userForm.Email = userProfile.Email // reset, it will be set when user clicks verification
			}
			userProfile, _, errorUser = userService.UpdateUser(w, &userForm, currentUser, id)
			if errorUser != nil {
				messages.ImportFromError("errors", errorUser)
			} else {
				messages.AddInfo("infos", string(T("profile_updated")))
			}
		}
	}
	availableLanguages := languages.GetAvailableLanguages()
	upev := UserProfileEditVariables{
		CommonTemplateVariables: NewCommonVariables(r),
		UserProfile:             &userProfile,
		UserForm:                userForm,
		FormErrors:              messages.GetAllErrors(),
		FormInfos:               messages.GetAllInfos(),
		Languages:               availableLanguages,
	}
	errorTmpl := viewProfileEditTemplate.ExecuteTemplate(w, "index.html", upev)
	if errorTmpl != nil {
		http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
	}
}

// Post Registration controller, we do some check on the form here, the rest on user service
func UserRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	b := form.RegistrationForm{}
	messages = msg.GetMessages(r)

	if !captcha.Authenticate(captcha.Extract(r)) {
		messages.AddError("errors", "Wrong captcha!")
	}
	if !messages.HasErrors() {
		if len(r.PostFormValue("email")) > 0 {
			form.EmailValidation(r.PostFormValue("email"), &messages)
		}
		form.ValidateUsername(r.PostFormValue("username"), &messages)
		if !messages.HasErrors() {
			modelHelper.BindValueForm(&b, r)
			modelHelper.ValidateForm(&b, &messages)
			if !messages.HasErrors() {
				_, errorUser := userService.CreateUser(w, r)
				if errorUser != nil {
					messages.ImportFromError("errors", errorUser)
				}
				if !messages.HasErrors() {
					common := NewCommonVariables(r)
					common.User = &model.User{
						Email: r.PostFormValue("email"), // indicate whether user had email set
					}
					htv := UserRegisterTemplateVariables{common, b, messages.GetAllErrors()}
					errorTmpl := viewRegisterSuccessTemplate.ExecuteTemplate(w, "index.html", htv)
					if errorTmpl != nil {
						http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
					}
				}
			}
		}
	}
	if messages.HasErrors() {
		UserRegisterFormHandler(w, r)
	}
}

func UserVerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	messages := msg.GetMessages(r)

	_, errEmail := userService.EmailVerification(token, w)
	if errEmail != nil {
		messages.ImportFromError("errors", errEmail)
	}
	htv := UserVerifyTemplateVariables{NewCommonVariables(r), messages.GetAllErrors()}
	errorTmpl := viewVerifySuccessTemplate.ExecuteTemplate(w, "index.html", htv)
	if errorTmpl != nil {
		http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
	}
}

// Post Login controller
func UserLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	b := form.LoginForm{}
	modelHelper.BindValueForm(&b, r)
	messages := msg.GetAllErrors()

	modelHelper.ValidateForm(&b, &messages)
	if !messages.HasErrors() {
		_, errorUser := userService.CreateUserAuthentication(w, r)
		if errorUser != nil {
			messages.ImportFromError("errors", errorUser)
			htv := UserLoginFormVariables{NewCommonVariables(r), b, messages.GetAllErrors()}
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
	if messages.HasErrors() {
		UserLoginFormHandler(w,r)
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
	if errorUser == nil && user.ID > 0 {
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

func UserNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if currentUser.ID > 0 {
		messages := msg.GetMessages(r)
		Ts, _ := languages.GetTfuncAndLanguageFromRequest(r)
		if r.URL.Query()["clear"] != nil {
			notifierService.DeleteAllNotifications(currentUser.ID)
			messages.AddInfo("infos", Ts("notifications_cleared"))
			currentUser.Notifications = []model.Notification{}
		}
		htv := UserProfileNotifVariables{NewCommonVariables(r), messages.GetAllInfos()}
		err := viewProfileNotifTemplate.ExecuteTemplate(w, "index.html", htv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		NotFoundHandler(w, r)
	}
}
