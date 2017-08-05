package email

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/timeHelper"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
)

var verificationHandler = securecookie.New(config.EmailTokenHashKey, nil)

// SendEmailVerification sends an email verification token via email.
func SendEmailVerification(to string, token string) error {
	T, err := publicSettings.GetDefaultTfunc()
	if err != nil {
		return err
	}
	content := T("link") + " : " + config.Get().WebAddress.Nyaa + "/verify/email/" + token
	contentHTML := T("verify_email_content") + "<br/>" + "<a href=\"" + config.Get().WebAddress.Nyaa + "/verify/email/" + token + "\" target=\"_blank\">" + format.GetHostname(config.Get().WebAddress.Nyaa) + "/verify/email/" + token + "</a>"
	return SendEmailFromAdmin(to, T("verify_email_title"), content, contentHTML)
}

// SendVerificationToUser sends an email verification token to user.
func SendVerificationToUser(user *models.User, newEmail string) (int, error) {
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
func EmailVerification(token string, c *gin.Context) (int, error) {
	value := make(map[string]string)
	err := verificationHandler.Decode("", token, &value)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return http.StatusForbidden, errors.New("token_valid")
	}
	timeInt, _ := strconv.ParseInt(value["t"], 10, 0)
	if timeHelper.IsExpired(time.Unix(timeInt, 0)) {
		return http.StatusForbidden, errors.New("token_expired")
	}
	id, _ := strconv.Atoi(value["u"])
	user, _, err := users.FindByID(uint(id))
	if err != nil {
		return http.StatusNotFound, errors.New("user_not_found")
	}
	user.Email = value["e"]
	return user.Update()
}
