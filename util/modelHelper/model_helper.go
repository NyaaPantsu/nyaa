package modelHelper

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
)

// TODO: Rewrite module
//       Functions are highly complex and require a lot of additional error handling

func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// AssignValue assign form values to model.
func AssignValue(model interface{}, form interface{}) {
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

// BindValueForm assign populate form from a request
func BindValueForm(form interface{}, r *http.Request) {
	r.ParseForm()
	formElem := reflect.ValueOf(form).Elem()
	for i := 0; i < formElem.NumField(); i++ {
		typeField := formElem.Type().Field(i)
		tag := typeField.Tag
		switch typeField.Type.Name() {
		case "string":
			formElem.Field(i).SetString(r.PostFormValue(tag.Get("form")))
		case "int":
			nbr, _ := strconv.Atoi(r.PostFormValue(tag.Get("form")))
			formElem.Field(i).SetInt(int64(nbr))
		case "float":
			nbr, _ := strconv.Atoi(r.PostFormValue(tag.Get("form")))
			formElem.Field(i).SetFloat(float64(nbr))
		case "bool":
			nbr, _ := strconv.ParseBool(r.PostFormValue(tag.Get("form")))
			formElem.Field(i).SetBool(nbr)
		}
	}
}

// ValidateForm : Check if a form is valid according to its tags
func ValidateForm(form interface{}, mes *msg.Messages) {
	formElem := reflect.ValueOf(form).Elem()
	for i := 0; i < formElem.NumField(); i++ {
		typeField := formElem.Type().Field(i)
		tag := typeField.Tag
		inputName := typeField.Name
		if tag.Get("hum_name") != "" { // For more human input name than gibberish
			inputName = tag.Get("hum_name")
		}
		if tag.Get("len_min") != "" && (tag.Get("needed") != "" || formElem.Field(i).Len() > 0) { // Check minimum length
			lenMin, _ := strconv.Atoi(tag.Get("len_min"))
			if formElem.Field(i).Len() < lenMin {
				mes.AddErrorf(tag.Get("form"), "Minimal length of %s required for the input: %s", strconv.Itoa(lenMin), inputName)
			}
		}
		if tag.Get("len_max") != "" && (tag.Get("needed") != "" || formElem.Field(i).Len() > 0) { // Check maximum length
			lenMax, _ := strconv.Atoi(tag.Get("len_max"))
			if formElem.Field(i).Len() > lenMax {
				mes.AddErrorf(tag.Get("form"), "Maximal length of %s required for the input: %s", strconv.Itoa(lenMax), inputName)
			}
		}
		if tag.Get("equalInput") != "" && (tag.Get("needed") != "" || formElem.Field(i).Len() > 0) {
			otherInput := formElem.FieldByName(tag.Get("equalInput"))
			if formElem.Field(i).Interface() != otherInput.Interface() {
				mes.AddErrorf(tag.Get("form"), "Must be same %s", inputName)
			}
		}
		switch typeField.Type.Name() {
		case "string":
			if tag.Get("equal") != "" && formElem.Field(i).String() != tag.Get("equal") {
				mes.AddErrorf(tag.Get("form"), "Wrong value for the input: %s", inputName)
			}
			if tag.Get("needed") != "" && formElem.Field(i).String() == "" {
				mes.AddErrorf(tag.Get("form"), "Field needed: %s", inputName)
			}
			if formElem.Field(i).String() == "" && tag.Get("default") != "" {
				formElem.Field(i).SetString(tag.Get("default"))
			}
		case "int":
			if tag.Get("equal") != "" { // Check minimum length
				equal, _ := strconv.Atoi(tag.Get("equal"))
				if formElem.Field(i).Int() > int64(equal) {
					mes.AddErrorf(tag.Get("form"), "Wrong value for the input: %s", inputName)
				}
			}
			if tag.Get("needed") != "" && formElem.Field(i).Int() == 0 {
				mes.AddErrorf(tag.Get("form"), "Field needed: %s", inputName)
			}
			if formElem.Field(i).Interface == nil && tag.Get("default") != "" { // FIXME: always false :'(
				defaultValue, _ := strconv.Atoi(tag.Get("default"))
				formElem.Field(i).SetInt(int64(defaultValue))
			}
		case "float":
			if tag.Get("equal") != "" { // Check minimum length
				equal, _ := strconv.Atoi(tag.Get("equal"))
				if formElem.Field(i).Float() != float64(equal) {
					mes.AddErrorf(tag.Get("form"), "Wrong value for the input: %s", inputName)
				}
			}
			if tag.Get("needed") != "" && formElem.Field(i).Float() == 0 {
				mes.AddErrorf(tag.Get("form"), "Field needed: %s", inputName)
			}
			if formElem.Field(i).Interface == nil && tag.Get("default") != "" { // FIXME: always false :'(
				defaultValue, _ := strconv.Atoi(tag.Get("default"))
				formElem.Field(i).SetFloat(float64(defaultValue))
			}
		case "bool":
			if tag.Get("equal") != "" { // Check minimum length
				equal, _ := strconv.ParseBool(tag.Get("equal"))
				if formElem.Field(i).Bool() != equal {
					mes.AddErrorf(tag.Get("form"), "Wrong value for the input: %s", inputName)
				}
			}
			if formElem.Field(i).Interface == nil && tag.Get("default") != "" { // FIXME: always false :'(
				defaultValue, _ := strconv.ParseBool(tag.Get("default"))
				formElem.Field(i).SetBool(defaultValue)
			}
		}
	}
}
