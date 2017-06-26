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
	"github.com/NyaaPantsu/nyaa/util/timeHelper"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/context"
	"github.com/gorilla/securecookie"
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
func ClearCookie(c *gin.Context) {
	c.SetCookie(CookieName, "", -1, "/", getDomainName(), false, true)
}

// SetCookieHandler sets the authentication cookie
func SetCookieHandler(c *gin.Context, user model.User) (int, error) {
	maxAge := getMaxAge()
	validUntil := timeHelper.FewDurationLater(time.Duration(maxAge) * time.Second)
	encoded, err := EncodeCookie(user.ID, validUntil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	c.SetCookie(CookieName, encoded, maxAge, "/", getDomainName(), false, true)
	// also set response header for convenience
	c.Header("X-Auth-Token", encoded)
	return http.StatusOK, nil
}

// RegisterHanderFromForm sets cookie from a RegistrationForm.
func RegisterHanderFromForm(c *gin.Context, registrationForm formStruct.RegistrationForm) (int, error) {
	username := registrationForm.Username // email isn't set at this point
	pass := registrationForm.Password
	user, status, err := checkAuth(c, username, pass)
	if err != nil {
		return status, err
	}
	return SetCookieHandler(c, user)
}

// RegisterHandler sets a cookie when user registered.
func RegisterHandler(c *gin.Context) (int, error) {
	var registrationForm formStruct.RegistrationForm
	c.Bind(&registrationForm)
	return RegisterHanderFromForm(c, registrationForm)
}

// CurrentUser determines the current user from the request or context
func CurrentUser(c *gin.Context) (model.User, error) {
	var user model.User
	encoded := c.Request.Header.Get("X-Auth-Token")
	if len(encoded) == 0 {
		// check cookie instead
		cookie, err := c.Cookie(CookieName)
		if err != nil {
			return user, err
		}
		encoded = cookie
	}
	userID, err := DecodeCookie(encoded)
	if err != nil {
		return user, err
	}

	userFromContext := getUserFromContext(c)

	if userFromContext.ID > 0 && userID == userFromContext.ID {
		user = userFromContext
	} else {
		if db.ORM.Preload("Notifications").Where("user_id = ?", userID).First(&user).RecordNotFound() { // We only load unread notifications
			return user, errors.New("user_not_found")
		}
		setUserToContext(c, user)
	}

	if user.IsBanned() {
		// recheck as user might've been banned in the meantime
		return user, errors.New("account_banned")
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

func getUserFromContext(c *gin.Context) model.User {
	if rv := context.Get(c.Request, UserContextKey); rv != nil {
		return rv.(model.User)
	}
	return model.User{}
}

func setUserToContext(c *gin.Context, val model.User) {
	context.Set(c.Request, UserContextKey, val)
}
