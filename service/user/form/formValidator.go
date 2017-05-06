package form

import (
	"regexp"
	"github.com/ewhal/nyaa/util/log"
)

const EMAIL_REGEX = `(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`

func EmailValidation(email string) bool {
	exp, err := regexp.Compile(EMAIL_REGEX)
	if regexpCompiled := log.CheckError(err); regexpCompiled {
		if exp.MatchString(email) {
			return true
		}
	}
	return false
}

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username string `form:"registrationUsername" binding:"required"`
	Email    string `form:"registrationEmail" binding:"required"`
	Password string `form:"registrationPassword" binding:"required"`
}

// RegistrationForm is used when creating a user authentication.
type LoginForm struct {
	Email    string `form:"loginEmail" binding:"required"`
	Password string `form:"loginPassword" binding:"required"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Email string `form:"email" binding:"required"`
}

// PasswordForm is used when updating a user password.
type PasswordForm struct {
	CurrentPassword string `form:"currentPassword" binding:"required"`
	Password        string `form:"newPassword" binding:"required"`
}

// SendPasswordResetForm is used when sending a password reset token.
type SendPasswordResetForm struct {
	Email string `form:"email" binding:"required"`
}

// PasswordResetForm is used when reseting a password.
type PasswordResetForm struct {
	PasswordResetToken string `form:"token" binding:"required"`
	Password           string `form:"newPassword" binding:"required"`
}

// VerifyEmailForm is used when verifying an email.
type VerifyEmailForm struct {
	ActivationToken string `form:"token" binding:"required"`
}

// ActivateForm is used when activating user.
type ActivateForm struct {
	Activation bool `form:"activation" binding:"required"`
}

// UserRoleForm is used when adding or removing a role from a user.
type UserRoleForm struct {
	UserId int `form:"userId" binding:"required"`
	RoleId int `form:"roleId" binding:"required"`
}
