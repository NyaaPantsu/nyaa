package userValidator

import (
	"regexp"

	"github.com/NyaaPantsu/nyaa/utils/log"
	"github.com/NyaaPantsu/nyaa/utils/validator"
)

// Regex by: Philippe Verdy (in a comment somewhere on a website) - Valid every email RFC valid
const emailRegex = `^((?:[-!#$%&'*+/=?^` + "`" + `{|}~\w]|\\.)+(?:\.(?:[-!#$%&'*+/=?^` + "`" + `{|}~\w]|\\.)+)*|"(?:[^\\"]|\\.)+")@(?:\[(?:((?:(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)\.){3}(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?))|IPv6:((?:[0-9A-F]{1,4}:){7}[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){6}:[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){5}:(?:[0-9A-F]{1,4}:)?[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){4}:(?:[0-9A-F]{1,4}:){0,2}[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){3}:(?:[0-9A-F]{1,4}:){0,3}[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){2}:(?:[0-9A-F]{1,4}:){0,4}[0-9A-F]{1,4}|[0-9A-F]{1,4}::(?:[0-9A-F]{1,4}:){0,5}[0-9A-F]{1,4}|::(?:[0-9A-F]{1,4}:){0,6}[0-9A-F]{1,4}|(?:[0-9A-F]{1,4}:){1,7}:|(?:[0-9A-F]{1,4}:){6}(?:(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)\.){3}(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)|(?:[0-9A-F]{1,4}:){0,5}:(?:(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)\.){3}(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)|::(?:[0-9A-F]{1,4}:){0,5}(?:(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?)\.){3}(?:[01][\d]{0,2}|2(?:[0-4]\d?|5[0-5]?|[6-9])?|[3-9]\d?))|([-a-z\d]{0,62}[a-z\d]:[^\[\\\]]+))\]|([a-z\d](?:[-a-z\d]{0,62}[a-z\d])?(?:\.[a-z\d](?:[-a-z\d]{0,62}[a-z\d])?)+))$`
const usernameRegex = `(\S)`

// EmailValidation : Check if an email is valid
func EmailValidation(email string) bool {
	exp, errorRegex := regexp.Compile(emailRegex)
	if regexpCompiled := log.CheckError(errorRegex); regexpCompiled {
		if exp.MatchString(email) {
			return true
		}
		return false
	}
	return false
}

// ValidateUsername : Check if a username is valid
func ValidateUsername(username string) bool {
	return validator.IsUTFLetterNumeric(username) && len(username) > 0
}

// IsAgreed : Check if terms and conditions are valid
func IsAgreed(termsAndConditions string) bool { // TODO: Inline function
	return termsAndConditions == "1"
}
