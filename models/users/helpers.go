package users

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"golang.org/x/crypto/bcrypt"
)

// CheckEmail : check if email is in database
func CheckEmail(email string) bool {
	if len(email) == 0 {
		return false
	}
	var count int
	models.ORM.Model(models.User{}).Where("email = ?", email).Count(&count)

	return count != 0
}

// Exists : check if the users credentials match to a user in db
func Exists(email string, pass string) (user *models.User, status int, err error) {
	if email == "" || pass == "" {
		return user, http.StatusNotFound, errors.New("no_username_password")
	}
	var userExist = &models.User{}
	// search by email or username
	if userValidator.EmailValidation(email) {
		if models.ORM.Where("email = ?", email).First(userExist).RecordNotFound() {
			status, err = http.StatusNotFound, errors.New("user_not_found")
			return
		}
	} else if models.ORM.Where("username = ?", email).First(userExist).RecordNotFound() {
		status, err = http.StatusNotFound, errors.New("user_not_found")
		return
	}
	user = userExist
	err = bcrypt.CompareHashAndPassword([]byte(userExist.Password), []byte(pass))
	if err != nil {
		status, err = http.StatusUnauthorized, errors.New("incorrect_password")
		return
	}

	if userExist.IsScraped() {
		status, err = http.StatusUnauthorized, errors.New("account_need_activation")
		return
	}
	status, err = http.StatusOK, nil
	return
}

// SuggestUsername suggest user's name if user's name already occupied.
func SuggestUsername(username string) string {
	var count int
	var usernameCandidate string
	models.ORM.Model(models.User{}).Where(&models.User{Username: username}).Count(&count)
	log.Debugf("count Before : %d", count)
	if count == 0 {
		return username
	}
	var postfix int
	for {
		usernameCandidate = username + strconv.Itoa(postfix)
		log.Debugf("usernameCandidate: %s\n", usernameCandidate)
		models.ORM.Model(models.User{}).Where(&models.User{Username: usernameCandidate}).Count(&count)
		log.Debugf("count after : %d\n", count)
		postfix = postfix + 1
		if count == 0 {
			break
		}
	}
	return usernameCandidate
}
