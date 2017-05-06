package userService

import (
	"errors"
	"net/http"
	"time"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/modelHelper"
	"github.com/ewhal/nyaa/util/crypto"
	"github.com/ewhal/nyaa/util/email"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/timeHelper"

	"github.com/nicksnyder/go-i18n/i18n"
)

// SendEmailVerfication sends an email verification token via email.
func SendEmailVerfication(to string, token string, locale string) error {
	T, _ := i18n.Tfunc(locale)
	err := email.SendEmailFromAdmin(to,
		T("verify_email_title"),
		T("link")+" : "+config.WebAddress+"/verify/email/"+token,
		T("verify_email_content")+"<br/><a href=\""+config.WebAddress+"/verify/email/"+token+"\" target=\"_blank\">"+config.WebAddress+"/verify/email/"+token+"</a>")
	return err
}

// SendVerificationToUser sends an email verification token to user.
func SendVerificationToUser(user model.User) (int, error) {
	var status int
	var err error
	user.ActivateUntil = timeHelper.TwentyFourHoursLater()
	user.ActivationToken, err = crypto.GenerateRandomToken32()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	user.Activation = false
	log.Debugf("generated token : %s", user.ActivationToken)
	status, err = UpdateUserCore(&user)
	if err != nil {
		return status, err
	}
	err = SendEmailVerfication(user.Email, user.ActivationToken, "en-us")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}

// SendVerification sends an email verification token.
func SendVerification(r *http.Request) (int, error) {

	var user model.User
	currentUser, err := CurrentUser(r)
	if err != nil {
		return http.StatusUnauthorized, errors.New("Unauthorized.")
	}
	if db.ORM.First(&user, currentUser.Id).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	status, err := SendVerificationToUser(user)
	return status, err
}

// EmailVerification verifies an email of user.
func EmailVerification(w http.ResponseWriter, r *http.Request) (int, error) {
	var user model.User
	var verifyEmailForm VerifyEmailForm
	modelHelper.BindValueForm(&verifyEmailForm, r)
	log.Debugf("verifyEmailForm.ActivationToken : %s", verifyEmailForm.ActivationToken)
	if db.ORM.Where(&model.User{ActivationToken: verifyEmailForm.ActivationToken}).First(&user).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	isExpired := timeHelper.IsExpired(user.ActivateUntil)
	log.Debugf("passwordResetUntil : %s", user.ActivateUntil.UTC())
	log.Debugf("expired : %t", isExpired)
	if isExpired {
		return http.StatusForbidden, errors.New("token not valid.")
	}
	user.ActivationToken = ""
	user.ActivateUntil = time.Now()
	user.ActivatedAt = time.Now()
	user.Activation = true
	status, err := UpdateUserCore(&user)
	if err != nil {
		return status, err
	}
	status, err = SetCookie(<, user.Token)
	return status, err
}
