package modelHelper

import (
	"fmt"
	"github.com/ewhal/nyaa/util/log"
	"net/http"
	"reflect"
	"strconv"
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

// AssignValue assign form values to model.
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

func ValidateForm(form interface{}, errorForm map[string][]string) map[string][]string {
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
				errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Minimal length of %s required for the input: %s", strconv.Itoa(lenMin), inputName))
			}
		}
		if tag.Get("len_max") != "" && (tag.Get("needed") != "" || formElem.Field(i).Len() > 0) { // Check maximum length
			lenMax, _ := strconv.Atoi(tag.Get("len_max"))
			if formElem.Field(i).Len() > lenMax {
				errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Maximal length of %s required for the input: %s", strconv.Itoa(lenMax), inputName))
			}
		}
		if tag.Get("equalInput") != "" && (tag.Get("needed") != "" || formElem.Field(i).Len() > 0) {
			otherInput := formElem.FieldByName(tag.Get("equalInput"))
			if formElem.Field(i).Interface() != otherInput.Interface() {
				errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Must be same %s", inputName))
			}
		}
		switch typeField.Type.Name() {
		case "string":
			if tag.Get("equal") != "" && formElem.Field(i).String() != tag.Get("equal") {
				errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Wrong value for the input: %s", inputName))
			}
			if tag.Get("needed") != "" && formElem.Field(i).String() == "" {
				errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Field needed: %s", inputName))
			}
			if formElem.Field(i).String() == "" && tag.Get("default") != "" {
				formElem.Field(i).SetString(tag.Get("default"))
			}
		case "int":
			if tag.Get("equal") != "" { // Check minimum length
				equal, _ := strconv.Atoi(tag.Get("equal"))
				if formElem.Field(i).Int() > int64(equal) {
					errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Wrong value for the input: %s", inputName))
				}
			}
			if tag.Get("needed") != "" && formElem.Field(i).Int() == 0 {
				errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Field needed: %s", inputName))
			}
			if formElem.Field(i).Interface == nil && tag.Get("default") != "" {
				defaultValue, _ := strconv.Atoi(tag.Get("default"))
				formElem.Field(i).SetInt(int64(defaultValue))
			}
		case "float":
			if tag.Get("equal") != "" { // Check minimum length
				equal, _ := strconv.Atoi(tag.Get("equal"))
				if formElem.Field(i).Float() != float64(equal) {
					errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Wrong value for the input: %s", inputName))
				}
			}
			if tag.Get("needed") != "" && formElem.Field(i).Float() == 0 {
				errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Field needed: %s", inputName))
			}
			if formElem.Field(i).Interface == nil && tag.Get("default") != "" {
				defaultValue, _ := strconv.Atoi(tag.Get("default"))
				formElem.Field(i).SetFloat(float64(defaultValue))
			}
		case "bool":
			if tag.Get("equal") != "" { // Check minimum length
				equal, _ := strconv.ParseBool(tag.Get("equal"))
				if formElem.Field(i).Bool() != equal {
					errorForm[tag.Get("form")] = append(errorForm[tag.Get("form")], fmt.Sprintf("Wrong value for the input: %s", inputName))
				}
			}
		}
	}
	return errorForm
}
