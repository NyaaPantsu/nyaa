package userService

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/util/email"
	"github.com/ewhal/nyaa/util/timeHelper"
	"github.com/gorilla/securecookie"
	"github.com/nicksnyder/go-i18n/i18n"
)

var verificationHandler = securecookie.New(config.EmailTokenHashKey, nil)

// SendEmailVerfication sends an email verification token via email.
func SendEmailVerification(to string, token string, locale string) error {
	T, err := i18n.Tfunc(locale)
	if err != nil {
		return err
	}
	content := T("link") + " : https://" + config.WebAddress + "/verify/email/" + token
	content_html := T("verify_email_content") + "<br/>" + "<a href=\"https://" + config.WebAddress + "/verify/email/" + token + "\" target=\"_blank\">" + config.WebAddress + "/verify/email/" + token + "</a>"
	return email.SendEmailFromAdmin(to, T("verify_email_title"), content, content_html)
	/* // debug code DO NOT LEAVE THIS ENABLED
	fmt.Printf("sending email to %s\n----\n%s\n%s\n----\n", to, content, content_html)
	return nil*/
}

// SendVerificationToUser sends an email verification token to user.
func SendVerificationToUser(user model.User, newEmail string) (int, error) {
	validUntil := timeHelper.TwentyFourHoursLater() // TODO: longer duration?
	value := map[string]string{
		"t": strconv.FormatInt(validUntil.Unix(), 10),
		"u": strconv.FormatUint(uint64(user.ID), 10),
		"e": newEmail,
	}
	encoded, err := verificationHandler.Encode("", value)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	err = SendEmailVerification(newEmail, encoded, "en-us")
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// EmailVerification verifies the token used for email verification
func EmailVerification(token string, w http.ResponseWriter) (int, error) {
	value := make(map[string]string)
	err := verificationHandler.Decode("", token, &value)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return http.StatusForbidden, errors.New("Token is not valid.")
	}
	time_int, _ := strconv.ParseInt(value["t"], 10, 0)
	if timeHelper.IsExpired(time.Unix(time_int, 0)) {
		return http.StatusForbidden, errors.New("Token has expired.")
	}
	var user model.User
	if db.ORM.Where("user_id = ?", value["u"]).First(&user).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found.")
	}
	user.Email = value["e"]
	return UpdateUserCore(&user)
}
