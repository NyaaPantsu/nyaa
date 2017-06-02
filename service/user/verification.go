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
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/email"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/NyaaPantsu/nyaa/util/timeHelper"
	"github.com/gorilla/securecookie"
)

var verificationHandler = securecookie.New(config.EmailTokenHashKey, nil)

// SendEmailVerification sends an email verification token via email.
func SendEmailVerification(to string, token string) error {
	T, err := publicSettings.GetDefaultTfunc()
	if err != nil {
		return err
	}
	content := T("link") + " : " + config.Conf.WebAddress.Nyaa + "/verify/email/" + token
	contentHTML := T("verify_email_content") + "<br/>" + "<a href=\"" + config.Conf.WebAddress.Nyaa + "/verify/email/" + token + "\" target=\"_blank\">" + util.GetHostname(config.Conf.WebAddress.Nyaa) + "/verify/email/" + token + "</a>"
	return email.SendEmailFromAdmin(to, T("verify_email_title"), content, contentHTML)
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
	err = SendEmailVerification(newEmail, encoded)
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
		return http.StatusForbidden, errors.New("Token is not valid")
	}
	timeInt, _ := strconv.ParseInt(value["t"], 10, 0)
	if timeHelper.IsExpired(time.Unix(timeInt, 0)) {
		return http.StatusForbidden, errors.New("Token has expired")
	}
	var user model.User
	if db.ORM.Where("user_id = ?", value["u"]).First(&user).RecordNotFound() {
		return http.StatusNotFound, errors.New("User is not found")
	}
	user.Email = value["e"]
	return UpdateUserCore(&user)
}
