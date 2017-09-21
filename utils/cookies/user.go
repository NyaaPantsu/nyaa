package cookies

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/timeHelper"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

const (
	// CookieName : Name of cookie
	CookieName = "session"

	// UserContextKey : key for user context
	UserContextKey = "nyaapantsu.user"
)

// NewCurrentUserRetriever create CurrentUserRetriever Struct for languages
func NewCurrentUserRetriever() *CurrentUserRetriever {
	return &CurrentUserRetriever{}
}

// CurrentUserRetriever struct for languages
type CurrentUserRetriever struct{}

// RetrieveCurrentUser retrieve current user for languages
func (*CurrentUserRetriever) RetrieveCurrentUser(c *gin.Context) (*models.User, error) {
	user, _, err := CurrentUser(c)
	if user == nil {
		return &models.User{}, err
	}
	return user, err
}

// CreateUserAuthentication creates user authentication.
func CreateUserAuthentication(c *gin.Context, form *userValidator.LoginForm) (*models.User, int, error) {
	username := form.Username
	pass := form.Password
	user, status, err := users.Exists(username, pass)
	if err != nil {
		return user, status, err
	}
	status, err = SetLogin(c, user)
	return user, status, err
}

// If you want to keep login cookies between restarts you need to make these permanent
var cookieHandler = securecookie.New(
	getOrGenerateKey(config.Get().Cookies.HashKey, 64),
	getOrGenerateKey(config.Get().Cookies.EncryptionKey, 32))

func getOrGenerateKey(key string, requiredLen int) []byte {
	data := []byte(key)
	if len(data) == 0 {
		log.Infof("No cookie key '%s' is set in config files. The users won't be kept logged in during restart and accross websites.", key)
		data = securecookie.GenerateRandomKey(requiredLen)
	} else if len(data) != requiredLen {
		panic(fmt.Sprintf("failed to load cookie key. required key length is %d bytes and the provided key length is %d bytes.", requiredLen, len(data)))
	}
	return data
}

// Decode : Encoding & Decoding of the cookie value
func Decode(cookieValue string) (uint, error) {
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

// Encode : Encoding of the cookie value
func Encode(userID uint, validUntil time.Time) (string, error) {
	value := map[string]string{
		"u": strconv.FormatUint(uint64(userID), 10),
		"t": strconv.FormatInt(validUntil.Unix(), 10),
	}
	return cookieHandler.Encode(CookieName, value)
}

// Clear : Erase cookie session
func Clear(c *gin.Context) {
	c.SetCookie(CookieName, "", -1, "/", getDomainName(), false, true)
}

// SetLogin sets the authentication cookie
func SetLogin(c *gin.Context, user *models.User) (int, error) {
	maxAge := getMaxAge(false)
	if c.PostForm("remember_me") == "remember" {
		maxAge = getMaxAge(true)
	}
	validUntil := timeHelper.FewDurationLater(time.Duration(maxAge) * time.Second)
	encoded, err := Encode(user.ID, validUntil)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	c.SetCookie(CookieName, encoded, maxAge, "/", "", false, true)
	// also set response header for convenience
	c.Header("X-Auth-Token", encoded)
	return http.StatusOK, nil
}

// CurrentUser retrieves a current user.
func CurrentUser(c *gin.Context) (*models.User, int, error) {
	encoded := c.Request.Header.Get("X-Auth-Token")
	var user = &models.User{}
	if len(encoded) == 0 {
		// check cookie instead
		cookie, err := c.Cookie(CookieName)
		if err != nil {
			return user, http.StatusInternalServerError, err
		}
		encoded = cookie
	}
	userID, err := Decode(encoded)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}

	userFromContext := getUserFromContext(c)

	if userFromContext.ID > 0 && userID == userFromContext.ID {
		user = userFromContext
	} else {
		user, _, _ = users.SessionByID(userID)
		setUserToContext(c, user)
	}

	if user.IsBanned() {
		// recheck as user might've been banned in the meantime
		return user, http.StatusUnauthorized, errors.New("account_banned")
	}
	if err != nil {
		return user, http.StatusInternalServerError, err
	}
	return user, http.StatusOK, nil
}
func getUserFromContext(c *gin.Context) *models.User {
	if rv, ok := c.Get(UserContextKey); ok {
		return rv.(*models.User)
	}
	return &models.User{}
}

func setUserToContext(c *gin.Context, val *models.User) {
	c.Set(UserContextKey, val)
}

// RetrieveUserFromRequest retrieves a user.
func RetrieveUserFromRequest(c *gin.Context, id uint) (*models.User, bool, uint, int, error) {
	var user models.User
	var currentUserID uint
	var isAuthor bool

	if models.ORM.First(&user, id).RecordNotFound() {
		return nil, isAuthor, currentUserID, http.StatusNotFound, errors.New("user_not_found")
	}
	currentUser, _, err := CurrentUser(c)
	if err == nil {
		currentUserID = currentUser.ID
		isAuthor = currentUser.ID == user.ID
	}

	return &user, isAuthor, currentUserID, http.StatusOK, nil
}
