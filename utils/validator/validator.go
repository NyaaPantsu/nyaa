package validator

import (
	"fmt"
	"reflect"
	"time"
	"unicode"

	"strconv"

	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("isutf", IsUTFLetterNumericValidator)
	validate.RegisterValidation("default", DefaultValidator)
}

// ValidateForm : Check if a form is valid according to its tags
func ValidateForm(form interface{}, mes *msg.Messages) {
	err := validate.Struct(form)
	if err != nil {
		for _, fieldError := range err.(validator.ValidationErrors) {
			switch fieldError.Tag() {
			case "required":
				mes.AddErrorTf(fieldError.Field(), "error_field_needed", fieldError.Field())
			case "eq":
				mes.AddErrorTf(fieldError.Field(), "error_equal", fieldError.Field(), fieldError.Param())
			case "ne", "nefield", "necsfield":
				mes.AddErrorTf(fieldError.Field(), "error_not_equal", fieldError.Field(), fieldError.Param())
			case "eqfield", "eqcsfield":
				mes.AddErrorTf(fieldError.Field(), "error_same_value", fieldError.Field(), fieldError.Param())
			case "gtfield", "gtcsfield":
				mes.AddErrorTf(fieldError.Field(), "error_greater_number", fieldError.Field(), fieldError.Param())
			case "ltfield", "ltcsfield":
				mes.AddErrorTf(fieldError.Field(), "error_less_number", fieldError.Field(), fieldError.Param())
			case "gtefield", "gtecsfield":
				mes.AddErrorTf(fieldError.Field(), "error_min_field", fieldError.Field(), fieldError.Param())
			case "ltefield", "ltecsfield":
				mes.AddErrorTf(fieldError.Field(), "error_max_field", fieldError.Field(), fieldError.Param())
			case "alpha":
				mes.AddErrorTf(fieldError.Field(), "error_alpha", fieldError.Field())
			case "alphanum":
				mes.AddErrorTf(fieldError.Field(), "error_alphanum", fieldError.Field())
			case "numeric":
				mes.AddErrorTf(fieldError.Field(), "error_numeric_valid", fieldError.Field())
			case "number":
				mes.AddErrorTf(fieldError.Field(), "error_number_valid", fieldError.Field())
			case "hexadecimal":
				mes.AddErrorTf(fieldError.Field(), "error_hexadecimal_valid", fieldError.Field())
			case "hexcolor":
				mes.AddErrorTf(fieldError.Field(), "error_hex_valid", fieldError.Field())
			case "rgb":
				mes.AddErrorTf(fieldError.Field(), "error_rgb_valid", fieldError.Field())
			case "rgba":
				mes.AddErrorTf(fieldError.Field(), "error_rgba_valid", fieldError.Field())
			case "hsl":
				mes.AddErrorTf(fieldError.Field(), "error_hsl_valid", fieldError.Field())
			case "hsla":
				mes.AddErrorTf(fieldError.Field(), "error_hsla_valid", fieldError.Field())
			case "url":
				mes.AddErrorTf(fieldError.Field(), "error_url_valid", fieldError.Field())
			case "uri":
				mes.AddErrorTf(fieldError.Field(), "error_uri_valid", fieldError.Field())
			case "base64":
				mes.AddErrorTf(fieldError.Field(), "error_base64_valid", fieldError.Field())
			case "contains":
				mes.AddErrorTf(fieldError.Field(), "error_contains", fieldError.Field(), fieldError.Param())
			case "containsany":
				mes.AddErrorTf(fieldError.Field(), "error_contains_any", fieldError.Field(), fieldError.Param())
			case "excludes":
				mes.AddErrorTf(fieldError.Field(), "error_excludes", fieldError.Field(), fieldError.Param())
			case "excludesall":
				mes.AddErrorTf(fieldError.Field(), "error_excludes_all", fieldError.Field(), fieldError.Param())
			case "excludesrune":
				mes.AddErrorTf(fieldError.Field(), "error_excludes_rune", fieldError.Field(), fieldError.Param())
			case "iscolor":
				mes.AddErrorTf(fieldError.Field(), "error_color_valid", fieldError.Field())
			default:
				switch fieldError.Kind() {
				case reflect.String:
					stringErrors(fieldError, mes)
				case reflect.Slice, reflect.Map, reflect.Array:
				case reflect.Struct:
					if fieldError.Type() != reflect.TypeOf(time.Time{}) {
						fmt.Printf("tag '%s' cannot be used on a struct type.", fieldError.Tag())
						mes.AddErrorTf(fieldError.Field(), "error_field", fieldError.Field())
					}
					dateErrors(fieldError, mes)
				default:
					numberErrors(fieldError, mes)
				}
			}
		}
	}
}

// Bind a validated form to a model, tag omit to not bind a field
func Bind(model interface{}, form interface{}) {
	modelIndirect := reflect.Indirect(reflect.ValueOf(model))
	formElem := reflect.ValueOf(form).Elem()
	typeOfTForm := formElem.Type()
	for i := 0; i < formElem.NumField(); i++ {
		tag := typeOfTForm.Field(i).Tag
		if tag.Get("omit") != "true" {
			modelField := modelIndirect.FieldByName(typeOfTForm.Field(i).Name)
			if modelField.IsValid() {
				formField := formElem.Field(i)
				modelField.Set(formField)
			} else {
				log.Warnf("modelField : %s - %s", typeOfTForm.Field(i).Name, modelField)
			}
		}
	}
}

// IsUTFLetterNumeric check if the string contains only unicode letters and numbers. Empty string is valid.
func IsUTFLetterNumeric(str string) bool {
	if len(str) == 0 {
		return true
	}
	for _, c := range str {
		if !unicode.IsLetter(c) && !unicode.IsNumber(c) { //letters && numbers are ok
			return false
		}
	}
	return true

}

// IsUTFLetterNumericValidator is an interface to IsUTFLetterNumeric from validator
func IsUTFLetterNumericValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String() // value
	return IsUTFLetterNumeric(value)
}

func DefaultValidator(fl validator.FieldLevel) bool {
	switch fl.Field().Kind() {
	case reflect.String:
		if len(fl.Field().String()) == 0 {
			fl.Field().SetString(fl.Param())
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if fl.Field().Int() == 0 {
			number, err := strconv.Atoi(fl.Param())
			if err != nil {
				fmt.Printf("Couldn't convert default value for field %s", fl.FieldName())
				return false
			}
			fl.Field().SetInt(int64(number))
		}
	case reflect.Float32, reflect.Float64:
		if fl.Field().Float() == 0 {
			number, err := strconv.ParseFloat(fl.Param(), 64)
			if err != nil {
				fmt.Printf("Couldn't convert default value for field %s", fl.FieldName())
				return false
			}
			fl.Field().SetFloat(number)
		}
	}
	return true
}
