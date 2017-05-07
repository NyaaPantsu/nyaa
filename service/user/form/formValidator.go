package form

import (
	"regexp"

	"github.com/ewhal/nyaa/util/log"
)

const EMAIL_REGEX = `(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`
const USERNAME_REGEX = `(\W)`

func EmailValidation(email string) bool {
	exp, err := regexp.Compile(EMAIL_REGEX)
	if regexpCompiled := log.CheckError(err); regexpCompiled {
		if exp.MatchString(email) {
			return true
		}
	}
	return false
}
func ValidateUsername(username string) bool {
	exp, err := regexp.Compile(USERNAME_REGEX)

	if username == "" {
		return false

	}
	if (len(username) < 3) || (len(username) > 15) {
		return false

	}
	if regexpCompiled := log.CheckError(err); regexpCompiled {
		if exp.MatchString(username) {
			return false
		}
	} else {
		return false
	}
	return true
}

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username  string `form:"registrationUsername"`
	Email     string `form:"registrationEmail"`
	Password  string `form:"registrationPassword"`
	CaptchaID string `form:"captchaID" inmodel:"false"`
}

// RegistrationForm is used when creating a user authentication.
type LoginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Email string `form:"email"`
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

// VerifyEmailForm is used when verifying an email.
type VerifyEmailForm struct {
	ActivationToken string `form:"token"`
}

// ActivateForm is used when activating user.
type ActivateForm struct {
	Activation bool `form:"activation"`
}

// UserRoleForm is used when adding or removing a role from a user.
type UserRoleForm struct {
	UserId int `form:"userId"`
	RoleId int `form:"roleId"`
}
