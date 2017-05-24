package form

import (
	"regexp"

	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
)

const EMAIL_REGEX = `(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`
const USERNAME_REGEX = `(\W)`

func EmailValidation(email string, mes *msg.Messages) bool {
	exp, errorRegex := regexp.Compile(EMAIL_REGEX)
	if regexpCompiled := log.CheckError(errorRegex); regexpCompiled {
		if exp.MatchString(email) {
			return true
		}
	}
	mes.AddError("email", "Email Address is not valid")
	return false
}

func ValidateUsername(username string, mes *msg.Messages) bool {
	exp, errorRegex := regexp.Compile(USERNAME_REGEX)
	if regexpCompiled := log.CheckError(errorRegex); regexpCompiled {
		if exp.MatchString(username) {
			mes.AddError("username", "Username contains illegal characters")
			return false
		}
	} else {
		return false
	}
	return true
}

func NewErrors() map[string][]string {
	err := make(map[string][]string)
	return err
}
func NewInfos() map[string][]string {
	infos := make(map[string][]string)
	return infos
}

func IsAgreed(termsAndConditions string) bool { // TODO: Inline function
	return termsAndConditions == "1"
}

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username           string `form:"username" needed:"true" len_min:"3" len_max:"20"`
	Email              string `form:"email"`
	Password           string `form:"password" needed:"true" len_min:"6" len_max:"72" equalInput:"ConfirmPassword"`
	ConfirmPassword    string `form:"password_confirmation" omit:"true" needed:"true"`
	CaptchaID          string `form:"captchaID" omit:"true" needed:"true"`
	TermsAndConditions bool   `form:"t_and_c" omit:"true" needed:"true" equal:"true" hum_name:"Terms and Conditions"`
}

// LoginForm is used when a user logs in.
type LoginForm struct {
	Username string `form:"username" needed:"true"`
	Password string `form:"password" needed:"true"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Username         string `form:"username" needed:"true" len_min:"3" len_max:"20"`
	Email            string `form:"email"`
	Language         string `form:"language" default:"en-us"`
	CurrentPassword  string `form:"current_password" len_min:"6" len_max:"72" omit:"true"`
	Password         string `form:"password" len_min:"6" len_max:"72" equalInput:"Confirm_Password"`
	Confirm_Password string `form:"password_confirmation" omit:"true"`
	Status           int    `form:"status" default:"0"`
}

// UserSettingsForm is used when updating a user.
type UserSettingsForm struct {
	NewTorrent        bool `form:"new_torrent" default:"true"`
	NewTorrentEmail   bool `form:"new_torrent_email" default:"true"`
	NewComment        bool `form:"new_comment" default:"true"`
	NewCommentEmail   bool `form:"new_comment_email" default:"false"`
	NewResponses      bool `form:"new_responses" default:"true"`
	NewResponsesEmail bool `form:"new_responses_email" default:"false"`
	NewFollower       bool `form:"new_follower" default:"true"`
	NewFollowerEmail  bool `form:"new_follower_email" default:"true"`
	Followed          bool `form:"followed" default:"false"`
	FollowedEmail     bool `form:"followed_email" default:"false"`
}

// PasswordForm is used when updating a user password.
type PasswordForm struct {
	CurrentPassword string `form:"currentPassword"`
	Password        string `form:"newPassword"`
}

// SendPasswordResetForm is used when sending a password reset token.
type SendPasswordResetForm struct {
	Email string `form:"email"`
}

// PasswordResetForm is used when reseting a password.
type PasswordResetForm struct {
	PasswordResetToken string `form:"token"`
	Password           string `form:"newPassword"`
}
