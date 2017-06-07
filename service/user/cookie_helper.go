package userService

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	formStruct "github.com/NyaaPantsu/nyaa/service/user/form"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/modelHelper"
	"github.com/NyaaPantsu/nyaa/util/timeHelper"
	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
	"golang.org/x/crypto/bcrypt"
)

const (
	// CookieName : Name of cookie
	CookieName = "session"

	// UserContextKey : key for user context
	UserContextKey = "nyaapantsu.user"
)

// If you want to keep login cookies between restarts you need to make these permanent
var cookieHandler = securecookie.New(
	getOrGenerateKey(config.Conf.Cookies.HashKey, 64),
	getOrGenerateKey(config.Conf.Cookies.EncryptionKey, 32))

func getOrGenerateKey(key string, requiredLen int) []byte {
	data := []byte(key)
	if len(data) == 0 {
		data = securecookie.GenerateRandomKey(requiredLen)
	} else if len(data) != requiredLen {
		panic(fmt.Sprintf("failed to load cookie key. required key length is %d bytes and the provided key length is %d bytes.", requiredLen, len(data)))
	}
	return data
}

// DecodeCookie : Encoding & Decoding of the cookie value
func DecodeCookie(cookieValue string) (uint, error) {
	value := make(map[string]string)
	err := cookieHandler.Decode(CookieName, cookieValue, &value)
	if err != nil {
		return 0, err
	}
	timeInt, _ := strconv.ParseInt(value["t"], 10, 0)
	if timeHelper.IsExpired(time.Unix(timeInt, 0)) {
		return 0, errors.New("Cookie is expired")
	}
	ret, err := strconv.ParseUint(value["u"], 10, 0)
	return uint(ret), err
}

// EncodeCookie : Encoding of the cookie value
func EncodeCookie(userID uint, validUntil time.Time) (string, error) {
	value := map[string]string{
		"u": strconv.FormatUint(uint64(userID), 10),
		"t": strconv.FormatInt(validUntil.Unix(), 10),
	}
	return cookieHandler.Encode(CookieName, value)
}

// ClearCookie : Erase cookie session
func ClearCookie(w http.ResponseWriter) (int, error) {
	cookie := &http.Cookie{
		Name:     CookieName,
		Domain:   getDomainName(),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	}
	http.SetCookie(w, cookie)
	return http.StatusOK, nil
}

// SetCookieHandler sets the authentication cookie
func SetCookieHandler(w http.ResponseWriter, r *http.Request, email string, pass string) (int, error) {
	if email == "" || pass == "" {
		return http.StatusNotFound, errors.New("No username/password entered")
	}

	var user model.User
	messages := msg.GetMessages(r)
	// search by email or username
	isValidEmail := formStruct.EmailValidation(email, messages)
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
	if user.IsBanned() {
		return http.StatusUnauthorized, errors.New("Account banned")
	}
	if user.IsScraped() {
		return http.StatusUnauthorized, errors.New("Account need activation from Moderators, please contact us")
	}

	maxAge := getMaxAge()
	validUntil := timeHelper.FewDurationLater(time.Duration(maxAge) * time.Second)
	encoded, err := EncodeCookie(user.ID, validUntil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	cookie := &http.Cookie{
		Name:     CookieName,
		Domain:   getDomainName(),
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   maxAge,
	}
	http.SetCookie(w, cookie)
	// also set response header for convenience
	w.Header().Set("X-Auth-Token", encoded)
	return http.StatusOK, nil
}

// RegisterHanderFromForm sets cookie from a RegistrationForm.
func RegisterHanderFromForm(w http.ResponseWriter, r *http.Request, registrationForm formStruct.RegistrationForm) (int, error) {
	username := registrationForm.Username // email isn't set at this point
	pass := registrationForm.Password
	return SetCookieHandler(w, r, username, pass)
}

// RegisterHandler sets a cookie when user registered.
func RegisterHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	var registrationForm formStruct.RegistrationForm
	modelHelper.BindValueForm(&registrationForm, r)
	return RegisterHanderFromForm(w, r, registrationForm)
}

// CurrentUser determines the current user from the request or context
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
	userID, err := DecodeCookie(encoded)
	if err != nil {
		return user, err
	}

	userFromContext := getUserFromContext(r)

	if userFromContext.ID > 0 && userID == userFromContext.ID {
		user = userFromContext
	} else {
		if db.ORM.Preload("Notifications").Where("user_id = ?", userID).First(&user).RecordNotFound() { // We only load unread notifications
			return user, errors.New("User not found")
		}
		setUserToContext(r, user)
	}

	if user.IsBanned() {
		// recheck as user might've been banned in the meantime
		return user, errors.New("Account banned")
	}
	return user, nil
}

func getDomainName() string {
	domain := config.Conf.Cookies.DomainName
	if config.Conf.Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}

func getMaxAge() int {
	return config.Conf.Cookies.MaxAge
}

func getUserFromContext(r *http.Request) model.User {
	if rv := context.Get(r, UserContextKey); rv != nil {
		return rv.(model.User)
	}
	return model.User{}
}

func setUserToContext(r *http.Request, val model.User) {
	context.Set(r, UserContextKey, val)
}
