package validator

import (
	"github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/go-playground/validator"
)

func slicesErrors(fe validator.FieldError, mes *messages.Messages) error {
	switch fe.Tag() {
	case "len":
		return mes.AddErrorTf(fe.Field(), "error_len_array", fe.Field(), fe.Tag())
	case "min":
		return mes.AddErrorTf(fe.Field(), "error_min_array", fe.Field(), fe.Tag())
	case "max":
		return mes.AddErrorTf(fe.Field(), "error_max_array", fe.Field(), fe.Tag())
	case "gt":
		return mes.AddErrorTf(fe.Field(), "error_greater_array", fe.Field(), fe.Tag())
	case "gte":
		return mes.AddErrorTf(fe.Field(), "error_min_array", fe.Field(), fe.Tag())
	case "lt":
		return mes.AddErrorTf(fe.Field(), "error_less_array", fe.Field(), fe.Tag())
	case "lte":
		return mes.AddErrorTf(fe.Field(), "error_max_array", fe.Field(), fe.Tag())
	}
	return mes.AddErrorTf(fe.Field(), "error_field", fe.Field())
}
