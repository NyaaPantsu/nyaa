package userService

import (
	"errors"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	formStruct "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/ewhal/nyaa/util/timeHelper"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

const CookieName = "session"

// If you want to keep login cookies between restarts you need to make these permanent
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

// Encoding & Decoding of the cookie value
func DecodeCookie(cookie_value string) (uint, error) {
	value := make(map[string]string)
	err := cookieHandler.Decode(CookieName, cookie_value, &value)
	if err != nil {
		return 0, err
	}
	time_int, _ := strconv.ParseInt(value["t"], 10, 0)
	if timeHelper.IsExpired(time.Unix(time_int, 0)) {
		return 0, errors.New("Cookie is expired")
	}
	ret, err := strconv.ParseUint(value["u"], 10, 0)
	return uint(ret), err
}

func EncodeCookie(user_id uint) (string, error) {
	validUntil := timeHelper.FewDaysLater(7) // 1 week
	value := map[string]string{
		"u": strconv.FormatUint(uint64(user_id), 10),
		"t": strconv.FormatInt(validUntil.Unix(), 10),
	}
	return cookieHandler.Encode(CookieName, value)
}

func ClearCookie(w http.ResponseWriter) (int, error) {
	cookie := &http.Cookie{
		Name:   CookieName,
		Value:  "",
		Path:   "/",
		HttpOnly: true,
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	return http.StatusOK, nil
}

// SetCookieHandler sets the authentication cookie
func SetCookieHandler(w http.ResponseWriter, email string, pass string) (int, error) {
	if email == "" || pass == "" {
		return http.StatusNotFound, errors.New("No username/password entered")
	}

	var user model.User
	// search by email or username
	isValidEmail, _ := formStruct.EmailValidation(email, formStruct.NewErrors())
	if isValidEmail {
		if db.ORM.Where("email = ?", email).First(&user).RecordNotFound() {
			return http.StatusNotFound, errors.New("User not found")
		}
	} else {
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

	encoded, err := EncodeCookie(user.ID)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	cookie := &http.Cookie{
		Name:  CookieName,
		Value: encoded,
		Path:  "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
	// also set response header for convenience
	w.Header().Set("X-Auth-Token", encoded)
	return http.StatusOK, nil
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

// CurrentUser determines the current user from the request
func CurrentUser(r *http.Request) (model.User, error) {
	var user model.User
	var encoded string

	encoded = r.Header.Get("X-Auth-Token")
	if len(encoded) == 0 {
		// check cookie instead
		cookie, err := r.Cookie(CookieName)
		if err != nil {
			return user, err
		}
		encoded = cookie.Value
	}
	user_id, err := DecodeCookie(encoded)
	if err != nil {
		return user, err
	}
	if db.ORM.Where("user_id = ?", user_id).First(&user).RecordNotFound() {
		return user, errors.New("User not found")
	}

	if user.Status == -1 {
		// recheck as user might've been banned in the meantime
		return user, errors.New("Account banned")
	}
	return user, nil
}
