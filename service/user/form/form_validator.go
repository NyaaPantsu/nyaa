package form

import (
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/asaskevich/govalidator"
)

// Regex by: Philippe Verdy (in a comment somewhere on a website) - Valid every email RFC valid
 const emailRegex = `^((?:[-!#$%&'*+/=?^` + "`" + `{|}~\w]|\\.)+(?:\.(?:[-!#$%&'*+/=?^` + "`" + `{|}~\w]|\\.)+)*|"(?:[^\\"]|\\.)+")@(?:\[(?:((?:(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)\.){3}(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?))|IPv6:((?:[0-9A-F]{1,4}:){7}[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){6}:[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){5}:(?:[0-9A-F]{1,4}:)?[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){4}:(?:[0-9A-F]{1,4}:){0,2}[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){3}:(?:[0-9A-F]{1,4}:){0,3}[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){2}:(?:[0-9A-F]{1,4}:){0,4}[0-9A-F]{1,4}|[0-9A-F]{1,4}::(?:[0-9A-F]{1,4}:){0,5}[0-9A-F]{1,4}|::(?:[0-9A-F]{1,4}:){0,6}[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){1,7}:|(?:[0-9A-F]{1,4}:){6}(?:(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)\.){3}(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)|(?:[0-9A-F]{1,4}:){0,5}:(?:(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)\.){3}(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)|::(?:[0-9A-F]{1,4}:){0,5}(?:(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)\.){3}(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?))|([-a-z\d]{0,62}[a-z\d]:[^\[\\\]]+))\]|([a-z\d](?:[-a-z\d]{0,62}[a-z\d])?(?:\.[a-z\d](?:[-a-z\d]{0,62}[a-z\d])?)+))$`
 const usernameRegex = `(\W)`
 
  // EmailValidation : Check if an email is valid
  func EmailValidation(email string, mes *msg.Messages) bool {
 	exp, errorRegex := regexp.Compile(emailRegex)
 	if regexpCompiled := log.CheckError(errorRegex); regexpCompiled {
 		if exp.MatchString(email) {
 			return true
 		}
 		mes.AddErrorT("email", "email_not_valid")
 		return false
	}
	mes.AddErrorT("email", "email_not_valid")
	return false
}

// ValidateUsername : Check if a username is valid
func ValidateUsername(username string, mes *msg.Messages) bool {
	if govalidator.IsAlpha(username) {
		mes.AddErrorT("username", "username_illegal")
		return false
	}
	return true
}

// IsAgreed : Check if terms and conditions are valid
func IsAgreed(termsAndConditions string) bool { // TODO: Inline function
	return termsAndConditions == "1"
}

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username           string `form:"username" needed:"true" len_min:"3" len_max:"20"`
	Email              string `form:"email"`
	Password           string `form:"password" needed:"true" len_min:"6" len_max:"72" equalInput:"ConfirmPassword"`
	ConfirmPassword    string `form:"password_confirmation" hum_name:"Password Confirmation" omit:"true" needed:"true"`
	CaptchaID          string `form:"captchaID" omit:"true" needed:"true"`
	TermsAndConditions bool   `form:"t_and_c" omit:"true" needed:"true" equal:"true" hum_name:"Terms and Conditions"`
}

// LoginForm is used when a user logs in.
type LoginForm struct {
	Username string `form:"username" needed:"true" json:"username"`
	Password string `form:"password" needed:"true" json:"password"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Username        string `form:"username" json:"username" needed:"true" len_min:"3" len_max:"20"`
	Email           string `form:"email" json:"email"`
	Language        string `form:"language" json:"language" default:"en-us"`
	CurrentPassword string `form:"current_password" json:"current_password" len_min:"6" len_max:"72" omit:"true"`
	Password        string `form:"password" json:"password" len_min:"6" len_max:"72" equalInput:"ConfirmPassword"`
	ConfirmPassword string `form:"password_confirmation" json:"password_confirmation" hum_name:"Password Confirmation" omit:"true"`
	Status          int    `form:"status" json:"status" default:"0"`
	Theme           string `form:"theme" json:"theme"`
}

// UserSettingsForm is used when updating a user.
type UserSettingsForm struct {
	NewTorrent        bool `form:"new_torrent" json:"new_torrent" default:"true"`
	NewTorrentEmail   bool `form:"new_torrent_email" json:"new_torrent_email" default:"true"`
	NewComment        bool `form:"new_comment" json:"new_comment" default:"true"`
	NewCommentEmail   bool `form:"new_comment_email" json:"new_comment_email" default:"false"`
	NewResponses      bool `form:"new_responses" json:"new_responses" default:"true"`
	NewResponsesEmail bool `form:"new_responses_email" json:"new_responses_email" default:"false"`
	NewFollower       bool `form:"new_follower" json:"new_follower" default:"true"`
	NewFollowerEmail  bool `form:"new_follower_email" json:"new_follower_email" default:"true"`
	Followed          bool `form:"followed" json:"followed" default:"false"`
	FollowedEmail     bool `form:"followed_email" json:"followed_email" default:"false"`
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
