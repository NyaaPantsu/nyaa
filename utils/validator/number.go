package validator

import (
	"github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/go-playground/validator"
)

func numberErrors(fe validator.FieldError, mes *messages.Messages) error {
	switch fe.Tag() {
	case "len":
		return mes.AddErrorTf(fe.Field(), "error_equal", fe.Field(), fe.Tag())
	case "min":
		return mes.AddErrorTf(fe.Field(), "error_min_number", fe.Field(), fe.Tag())
	case "max":
		return mes.AddErrorTf(fe.Field(), "error_max_number", fe.Field(), fe.Tag())
	case "gt":
		return mes.AddErrorTf(fe.Field(), "error_greater_number", fe.Field(), fe.Tag())
	case "gte":
		return mes.AddErrorTf(fe.Field(), "error_min_number", fe.Field(), fe.Tag())
	case "lt":
		return mes.AddErrorTf(fe.Field(), "error_less_number", fe.Field(), fe.Tag())
	case "lte":
		return mes.AddErrorTf(fe.Field(), "error_max_number", fe.Field(), fe.Tag())
	}
	return mes.AddErrorTf(fe.Field(), "error_field", fe.Field())
}
