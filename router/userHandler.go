package router

import (
	"net/http"

	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/service/user"
	"github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/gorilla/mux"
)

//var viewTemplate = template.Must(template.New("view").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/view.html"))
//var viewTemplate = template.Must(template.New("view").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/view.html"))

// Getting View User Registration
func UserRegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	_, errorUser := userService.CurrentUser(r)
	if (errorUser != nil) {
		b := form.RegistrationForm{}
		modelHelper.BindValueForm(&b, r)
		b.CaptchaID = captcha.GetID()
		languages.SetTranslationFromRequest(viewRegisterTemplate, r, "en-us")
		htv := UserRegisterTemplateVariables{b, form.NewErrors(), NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
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
	htv := UserLoginFormVariables{b, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
	err := viewLoginTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Getting User Profile
func UserProfileHandler(w http.ResponseWriter, r *http.Request) {

}

// Getting View User Profile Update
func UserProfileFormHandler(w http.ResponseWriter, r *http.Request) {

}

// Post Registration controller, we do some check on the form here, the rest on user service
func UserRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	// Check same Password
	b := form.RegistrationForm{}
	err := form.NewErrors()
	if !captcha.Authenticate(captcha.Extract(r)) {
		err["errors"] = append(err["errors"], "Wrong captcha!")
	}
	if (len(err) == 0) {
		_, err = form.EmailValidation(r.PostFormValue("email"), err)
		_, err = form.ValidateUsername(r.PostFormValue("username"), err)
		if (len(err) == 0) {
			modelHelper.BindValueForm(&b, r)
			err = modelHelper.ValidateForm(&b, err)
			if (len(err) == 0) {
				_, errorUser := userService.CreateUser(w, r)
				if (errorUser != nil) {
					err["errors"] = append(err["errors"], errorUser.Error())
				}
				if (len(err) == 0) {
					b := form.RegistrationForm{}
					languages.SetTranslationFromRequest(viewRegisterSuccessTemplate, r, "en-us")
					htv := UserRegisterTemplateVariables{b, err, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
					errorTmpl := viewRegisterSuccessTemplate.ExecuteTemplate(w, "index.html", htv)
					if errorTmpl != nil {
						http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
					}
				} 
			} 
		} 
	}
	if (len(err) > 0) {
		b.CaptchaID = captcha.GetID()
		languages.SetTranslationFromRequest(viewRegisterTemplate, r, "en-us")
		htv := UserRegisterTemplateVariables{b, err, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
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
	if (errEmail != nil) {
		err["errors"] = append(err["errors"], errEmail.Error())
	}
	languages.SetTranslationFromRequest(viewVerifySuccessTemplate, r, "en-us")
	htv := UserVerifyTemplateVariables{err, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
	errorTmpl := viewVerifySuccessTemplate.ExecuteTemplate(w, "index.html", htv)
	if errorTmpl != nil {
		http.Error(w, errorTmpl.Error(), http.StatusInternalServerError)
	}
}

// Post Login controller
func UserLoginPostHandler(w http.ResponseWriter, r *http.Request) {

}

// Post Profule Update controller
func UserProfilePostHandler(w http.ResponseWriter, r *http.Request) {

}
