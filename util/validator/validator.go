package validator

import (
	"reflect"
	"strconv"

	"github.com/asaskevich/govalidator"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
)
/

// ValidateForm : Check if a form is valid according to its tags
func ValidateForm(form interface{}, mes *msg.Messages) {
	result, err := govalidator.ValidateStruct(form)
	if err != nil {
		println("error: " + err.Error())
	}
}
