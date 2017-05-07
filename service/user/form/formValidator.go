package form

import (
	"regexp"

	"github.com/ewhal/nyaa/util/log"
)

const EMAIL_REGEX = `(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`
const USERNAME_REGEX = `(\W)`

func EmailValidation(email string, err map[string][]string) (bool, map[string][]string) {
	exp, errorRegex := regexp.Compile(EMAIL_REGEX)
	if regexpCompiled := log.CheckError(errorRegex); regexpCompiled {
		if exp.MatchString(email) {
			return true, err
		}
	} 
	err["email"] = append(err["email"], "Email Address is not valid")
	return false, err
}
func ValidateUsername(username string, err map[string][]string) (bool, map[string][]string)  {
	exp, errorRegex := regexp.Compile(USERNAME_REGEX)
	if regexpCompiled := log.CheckError(errorRegex); regexpCompiled {
		if exp.MatchString(username) {
			err["username"] = append(err["username"], "Username contains illegal characters")
			return false, err
		}
	} else {
		return false, err
	}
	return true, err
}

func NewErrors() map[string][]string {
	err := make(map[string][]string)
	return err
}
func IsAgreed(t_and_c string) bool {
	if t_and_c == "1" {
		return true
	}
	return false
}

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username  string `form:"username" needed:"true" min_len:"3" max_len:"20"`
	Email     string `form:"email" needed:"true"`
	Password  string `form:"password" needed:"true" min_len:"6" max_len:"25"`
	Confirm_Password string `form:"password_confirmation" omit:"true" needed:"true"`
	CaptchaID string `form:"captchaID" omit:"true" needed:"true"`
	T_and_C   bool   `form:"t_and_c" omit:"true" needed:"true" equal:"true" hum_name:"Terms and Conditions"`
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
