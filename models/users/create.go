package users

import (
	"errors"
	"net/http"
	"time"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/crypto"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// CreateUserFromForm creates a user from a registration form.
func CreateUserFromRequest(registrationForm *userValidator.RegistrationForm) (*models.User, error) {
	var user = &models.User{}
	log.Debugf("registrationForm %+v\n", registrationForm)
	validator.Bind(user, registrationForm)
	if user.Email == "" {
		user.MD5 = ""
	} else {
		// Despite the email not being verified yet we calculate this for convenience reasons
		var err error
		user.MD5, err = crypto.GenerateMD5Hash(user.Email)
		if err != nil {
			return user, err
		}
	}
	user.Email = "" // unset email because it will be verified later
	user.CreatedAt = time.Now()
	// User settings to default
	user.Settings.ToDefault()
	user.SaveSettings()
	// currently unused but needs to be set:
	user.APIToken, _ = crypto.GenerateRandomToken32()
	user.APITokenExpiry = time.Unix(0, 0)
	// set default number of PP
	user.Pantsu = 1

	if models.ORM.Create(user).Error != nil {
		return user, errors.New("user not created")
	}

	return user, nil
}

// CreateUser creates a user.
func CreateUser(c *gin.Context) (*models.User, int) {
	var user = &models.User{}
	var registrationForm userValidator.RegistrationForm
	var err error
	messages := msg.GetMessages(c)
	c.Bind(&registrationForm)
	usernameCandidate := SuggestUsername(registrationForm.Username)
	if usernameCandidate != registrationForm.Username {
		messages.AddErrorTf("username", "username_taken", usernameCandidate)
		return user, http.StatusInternalServerError
	}
	if registrationForm.Email != "" && CheckEmail(registrationForm.Email) {
		messages.AddErrorT("email", "email_in_db")
		return user, http.StatusInternalServerError
	}
	password, err := bcrypt.GenerateFromPassword([]byte(registrationForm.Password), 10)
	if err != nil {
		messages.Error(err)
		return user, http.StatusInternalServerError
	}
	registrationForm.Password = string(password)
	user, err = CreateUserFromRequest(&registrationForm)
	if err != nil {
		messages.Error(err)
		return user, http.StatusInternalServerError
	}
	return user, http.StatusOK
}
