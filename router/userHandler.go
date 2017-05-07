package router

import (
	"html/template"
	"net/http"

	"github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/service/user"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/gorilla/mux"
)

var viewRegisterTemplate = template.Must(template.New("userRegister").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/user/register.html"))
var viewLoginTemplate = template.Must(template.New("userLogin").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/user/login.html"))
var viewRegisterSuccessTemplate = template.Must(template.New("userRegisterSuccess").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/user/signup_success.html"))
//var viewTemplate = template.Must(template.New("view").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/view.html"))
//var viewTemplate = template.Must(template.New("view").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/view.html"))

func init() {
	template.Must(viewRegisterTemplate.ParseGlob("templates/_*.html"))
	template.Must(viewLoginTemplate.ParseGlob("templates/_*.html"))
	template.Must(viewRegisterSuccessTemplate.ParseGlob("templates/_*.html"))
}

// Getting View User Registration

func UserRegisterFormHandler(w http.ResponseWriter, r *http.Request) {
	b := form.RegistrationForm{}
	modelHelper.BindValueForm(&b, r)
	languages.SetTranslation("en-us", viewRegisterTemplate)
	htv := UserRegisterTemplateVariables{b, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
	err := viewRegisterTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Getting View User Login
func UserLoginFormHandler(w http.ResponseWriter, r *http.Request) {
	b := form.LoginForm{}
	modelHelper.BindValueForm(&b, r)
	languages.SetTranslation("en-us", viewLoginTemplate)
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

// Post Registration controller
func UserRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	// Check same Password
	if ((r.PostFormValue("password") == r.PostFormValue("password_confirm"))&&(r.PostFormValue("password") != "")) {
		if ((form.EmailValidation(r.PostFormValue("email")))&&(form.ValidateUsername(r.PostFormValue("username")))) {
			_ , err := userService.CreateUser(w, r)
			if (err == nil) {
				b := form.RegistrationForm{}
				htv := UserRegisterTemplateVariables{b, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
				err = viewRegisterSuccessTemplate.ExecuteTemplate(w, "index.html", htv)
				if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			} else {
				UserRegisterFormHandler(w, r)
			}
		} else {
			UserRegisterFormHandler(w, r)
		}
	} else {
		UserRegisterFormHandler(w, r)
	}
}

// Post Login controller
func UserLoginPostHandler(w http.ResponseWriter, r *http.Request) {

}

// Post Profule Update controller
func UserProfilePostHandler(w http.ResponseWriter, r *http.Request) {

}
