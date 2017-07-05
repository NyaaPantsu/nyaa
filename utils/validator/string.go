package validator

import (
	"github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/go-playground/validator"
)

func stringErrors(fe validator.FieldError, mes *messages.Messages) error {
	switch fe.Tag() {
	case "len":
		return mes.AddErrorTf(fe.Field(), "error_length", fe.Param(), fe.Field())
	case "min", "gte":
		return mes.AddErrorTf(fe.Field(), "error_min_length", fe.Param(), fe.Field())
	case "max", "lte":
		return mes.AddErrorTf(fe.Field(), "error_max_length", fe.Param(), fe.Field())
	case "lt":
		return mes.AddErrorTf(fe.Field(), "error_less_length", fe.Param(), fe.Field())
	case "gt":
		return mes.AddErrorTf(fe.Field(), "error_greater_length", fe.Param(), fe.Field())
	}
	return mes.AddErrorTf(fe.Field(), "error_field", fe.Field())
}
