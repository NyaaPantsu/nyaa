package userService

import (
	"errors"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	formStruct "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

func Token(r *http.Request) (string, error) {
	var token string
	cookie, err := r.Cookie("session")
	if err != nil {
		return token, err
	}
	cookieValue := make(map[string]string)
	err = cookieHandler.Decode("session", cookie.Value, &cookieValue)
	if err != nil {
		return token, err
	}
	token = cookieValue["token"]
	if len(token) == 0 {
		return token, errors.New("token is empty")
	}
	return token, nil
}

// SetCookie sets a cookie.
func SetCookie(w http.ResponseWriter, token string) (int, error) {
	value := map[string]string{
		"token": token,
	}
	encoded, err := cookieHandler.Encode("session", value)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	cookie := &http.Cookie{
		Name:  "session",
		Value: encoded,
		Path:  "/",
	}
	http.SetCookie(w, cookie)
	return http.StatusOK, nil
}

// ClearCookie clears a cookie.
func ClearCookie(w http.ResponseWriter) (int, error) {
	cookie := &http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	return http.StatusOK, nil
}

// SetCookieHandler sets a cookie with email and password.
func SetCookieHandler(w http.ResponseWriter, email string, pass string) (int, error) {
	if email != "" && pass != "" {
		var user model.User
		isValidEmail, _ := formStruct.EmailValidation(email, formStruct.NewErrors())
		if isValidEmail {
			log.Debug("User entered valid email.")
			if db.ORM.Where("email = ?", email).First(&user).RecordNotFound() {
				return http.StatusNotFound, errors.New("User not found")
			}
		} else {
			log.Debug("User entered username.")
			if db.ORM.Where("username = ?", email).First(&user).RecordNotFound() {
				return http.StatusNotFound, errors.New("User not found")
			}
		}
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))
		if err != nil {
			return http.StatusUnauthorized, errors.New("Password incorrect")
		}
		if user.Status == -1 {
			return http.StatusUnauthorized, errors.New("Account banned")
		}
		status, err := SetCookie(w, user.Token)
		if err != nil {
			return status, err
		}
		w.Header().Set("X-Auth-Token", user.Token)
		return http.StatusOK, nil
	}
	return http.StatusNotFound, errors.New("user not found")
}

// RegisterHanderFromForm sets cookie from a RegistrationForm.
func RegisterHanderFromForm(w http.ResponseWriter, registrationForm formStruct.RegistrationForm) (int, error) {
	username := registrationForm.Username // email isn't set at this point
	pass := registrationForm.Password
	return SetCookieHandler(w, username, pass)
}

// RegisterHandler sets a cookie when user registered.
func RegisterHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	var registrationForm formStruct.RegistrationForm
	modelHelper.BindValueForm(&registrationForm, r)
	return RegisterHanderFromForm(w, registrationForm)
}

// CurrentUser get a current user.
func CurrentUser(r *http.Request) (model.User, error) {
	var user model.User
	var token string
	var err error
	token = r.Header.Get("X-Auth-Token")
	if len(token) > 0 {
		log.Debug("header token exists")
	} else {
		token, err = Token(r)
		log.Debug("header token does not exist")
		if err != nil {
			return user, err
		}
	}
	if db.ORM.Where("api_token = ?", token).First(&user).RecordNotFound() {
		return user, errors.New("user not found")
	}
	err = db.ORM.Model(&user).Error
	return user, err
}
