package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/notifier"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/crypto"
	"github.com/NyaaPantsu/nyaa/util/languages"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/modelHelper"
	"github.com/gorilla/mux"
)

// UserRegisterFormHandler : Getting View User Registration
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
	urtv := formTemplateVariables{
		commonTemplateVariables: newCommonVariables(r),
		Form:       registrationForm,
		FormErrors: messages.GetAllErrors(),
	}
	err := viewRegisterTemplate.ExecuteTemplate(w, "index.html", urtv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UserLoginFormHandler : Getting View User Login
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
	ulfv := formTemplateVariables{
		commonTemplateVariables: newCommonVariables(r),
		Form:       loginForm,
		FormErrors: messages.GetAllErrors(),
	}

	err := viewLoginTemplate.ExecuteTemplate(w, "index.html", ulfv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// UserProfileHandler :  Getting User Profile
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	Ts, _ := languages.GetTfuncAndLanguageFromRequest(r)
	messages := msg.GetMessages(r)

	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil {
		currentUser := getUser(r)
		follow := r.URL.Query()["followed"]
		unfollow := r.URL.Query()["unfollowed"]
		deleteVar := r.URL.Query()["delete"]

		if (deleteVar != nil) && (userPermission.CurrentOrAdmin(currentUser, userProfile.ID)) {
			_, errUser := userService.DeleteUser(w, currentUser, id)
			if errUser != nil {
				messages.ImportFromError("errors", errUser)
			}
			htv := userVerifyTemplateVariables{newCommonVariables(r), messages.GetAllErrors()}
			errorTmpl := viewUserDeleteTemplate.ExecuteTemplate(w, "index.html", htv)
			if errorTmpl != nil {
				http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
			}
		} else {
			if follow != nil {
				messages.AddInfof("infos", Ts("user_followed_msg"), userProfile.Username)
			}
			if unfollow != nil {
				messages.AddInfof("infos", Ts("user_unfollowed_msg"), userProfile.Username)
			}
			userProfile.ParseSettings()
			htv := userProfileVariables{newCommonVariables(r), &userProfile, messages.GetAllInfos()}

			err := viewProfileTemplate.ExecuteTemplate(w, "index.html", htv)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		NotFoundHandler(w, r)
	}
}

// UserDetailsHandler : Getting User Profile Details View
func UserDetailsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := getUser(r)
	messages := msg.GetMessages(r)

	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser == nil && userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
		if userPermission.CurrentOrAdmin(currentUser, userProfile.ID) {
			b := form.UserForm{}
			modelHelper.BindValueForm(&b, r)
			availableLanguages := languages.GetAvailableLanguages()
			userProfile.ParseSettings()
			htv := userProfileEditVariables{newCommonVariables(r), &userProfile, b, messages.GetAllErrors(), messages.GetAllInfos(), availableLanguages}
			err := viewProfileEditTemplate.ExecuteTemplate(w, "index.html", htv)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
	} else {
		NotFoundHandler(w, r)
	}
}

// UserProfileFormHandler : Getting View User Profile Update
func UserProfileFormHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := getUser(r)
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser != nil || !userPermission.CurrentOrAdmin(currentUser, userProfile.ID) || userProfile.ID == 0 {
		NotFoundHandler(w, r)
		return
	}
	userProfile.ParseSettings()
	messages := msg.GetMessages(r)
	userForm := form.UserForm{}
	userSettingsForm := form.UserSettingsForm{}

	Ts, _ := languages.GetTfuncAndLanguageFromRequest(r)

	if len(r.PostFormValue("email")) > 0 {
		form.EmailValidation(r.PostFormValue("email"), messages)
	}
	if len(r.PostFormValue("username")) > 0 {
		form.ValidateUsername(r.PostFormValue("username"), messages)
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
		modelHelper.ValidateForm(&userForm, messages)
		if !messages.HasErrors() {
			if userForm.Email != userProfile.Email {
				userService.SendVerificationToUser(*currentUser, userForm.Email)
				messages.AddInfof("infos", Ts("email_changed"), userForm.Email)
				userForm.Email = userProfile.Email // reset, it will be set when user clicks verification
			}
			userProfile, _, errorUser = userService.UpdateUser(w, &userForm, &userSettingsForm, currentUser, id)
			if errorUser != nil {
				messages.ImportFromError("errors", errorUser)
			} else {
				messages.AddInfo("infos", Ts("profile_updated"))
			}
		}
	}
	availableLanguages := languages.GetAvailableLanguages()
	upev := userProfileEditVariables{
		commonTemplateVariables: newCommonVariables(r),
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

// UserRegisterPostHandler : Post Registration controller, we do some check on the form here, the rest on user service
func UserRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	b := form.RegistrationForm{}
	messages := msg.GetMessages(r)

	if !captcha.Authenticate(captcha.Extract(r)) {
		messages.AddError("errors", "Wrong captcha!")
	}
	if !messages.HasErrors() {
		if len(r.PostFormValue("email")) > 0 {
			form.EmailValidation(r.PostFormValue("email"), messages)
		}
		form.ValidateUsername(r.PostFormValue("username"), messages)
		if !messages.HasErrors() {
			modelHelper.BindValueForm(&b, r)
			modelHelper.ValidateForm(&b, messages)
			if !messages.HasErrors() {
				_, errorUser := userService.CreateUser(w, r)
				if errorUser != nil {
					messages.ImportFromError("errors", errorUser)
				}
				if !messages.HasErrors() {
					common := newCommonVariables(r)
					common.User = &model.User{
						Email: r.PostFormValue("email"), // indicate whether user had email set
					}
					htv := formTemplateVariables{common, b, messages.GetAllErrors(), messages.GetAllInfos()}
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

// UserVerifyEmailHandler : Controller when verifying email, needs a token
func UserVerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token := vars["token"]
	messages := msg.GetMessages(r)

	_, errEmail := userService.EmailVerification(token, w)
	if errEmail != nil {
		messages.ImportFromError("errors", errEmail)
	}
	htv := userVerifyTemplateVariables{newCommonVariables(r), messages.GetAllErrors()}
	errorTmpl := viewVerifySuccessTemplate.ExecuteTemplate(w, "index.html", htv)
	if errorTmpl != nil {
		http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
	}
}

// UserLoginPostHandler : Post Login controller
func UserLoginPostHandler(w http.ResponseWriter, r *http.Request) {
	b := form.LoginForm{}
	modelHelper.BindValueForm(&b, r)
	messages := msg.GetMessages(r)

	modelHelper.ValidateForm(&b, messages)
	if !messages.HasErrors() {
		_, errorUser := userService.CreateUserAuthentication(w, r)
		if errorUser != nil {
			messages.ImportFromError("errors", errorUser)
			htv := formTemplateVariables{newCommonVariables(r), b, messages.GetAllErrors(), messages.GetAllInfos()}
			errorTmpl := viewLoginTemplate.ExecuteTemplate(w, "index.html", htv)
			if errorTmpl != nil {
				http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
			}
			return
		}
		messages.ImportFromError("errors", errorUser)
	}
	if messages.HasErrors() {
		UserLoginFormHandler(w, r)
	}
}

// UserLogoutHandler : Controller to logout users
func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = userService.ClearCookie(w)
	url, _ := Router.Get("home").URL()
	http.Redirect(w, r, url.String(), http.StatusSeeOther)
}

// UserFollowHandler : Controller to follow/unfollow users, need user id to follow
func UserFollowHandler(w http.ResponseWriter, r *http.Request) {
	var followAction string
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := getUser(r)
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

// UserNotificationsHandler : Controller to show user notifications
func UserNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := getUser(r)
	if currentUser.ID > 0 {
		messages := msg.GetMessages(r)
		Ts, _ := languages.GetTfuncAndLanguageFromRequest(r)
		if r.URL.Query()["clear"] != nil {
			notifierService.DeleteAllNotifications(currentUser.ID)
			messages.AddInfo("infos", Ts("notifications_cleared"))
			currentUser.Notifications = []model.Notification{}
		}
		htv := userProfileNotifVariables{newCommonVariables(r), messages.GetAllInfos()}
		err := viewProfileNotifTemplate.ExecuteTemplate(w, "index.html", htv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		NotFoundHandler(w, r)
	}
}

// UserAPIKeyResetHandler : Controller to reset user api key
func UserAPIKeyResetHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	currentUser := getUser(r)

	Ts, _ := languages.GetTfuncAndLanguageFromRequest(r)
	messages := msg.GetMessages(r)
	userProfile, _, errorUser := userService.RetrieveUserForAdmin(id)
	if errorUser != nil || !userPermission.CurrentOrAdmin(currentUser, userProfile.ID) || userProfile.ID == 0 {
		NotFoundHandler(w, r)
		return
	}
	userProfile.ApiToken, _ = crypto.GenerateRandomToken32()
	userProfile.ApiTokenExpiry = time.Unix(0, 0)
	_, errorUser = userService.UpdateUserCore(&userProfile)
	if errorUser != nil {
		messages.ImportFromError("errors", errorUser)
	} else {
		messages.AddInfo("infos", Ts("profile_updated"))
	}

}
