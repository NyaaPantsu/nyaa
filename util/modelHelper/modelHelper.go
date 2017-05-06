package modelHelper

import (
	"reflect"
	"net/http"
	"github.com/ewhal/nyaa/util/log"
	"strconv"
)

func IsZeroOfUnderlyingType(x interface{}) bool {
	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

// AssignValue assign form values to model.
func AssignValue(model interface{}, form interface{}) {
	modelIndirect := reflect.Indirect(reflect.ValueOf(model))
	formElem := reflect.ValueOf(form).Elem()
	typeOfTForm := formElem.Type()
	for i := 0; i < formElem.NumField(); i++ {
		modelField := modelIndirect.FieldByName(typeOfTForm.Field(i).Name)
		if modelField.IsValid() {
			formField := formElem.Field(i)
			modelField.Set(formField)
		} else {
			log.Warnf("modelField : %s - %s", typeOfTForm.Field(i).Name, modelField)
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
			case "string" :
				formElem.Field(i).SetString(r.PostFormValue(tag.Get("form")))
			case "int" :
				nbr, _ := strconv.Atoi(r.PostFormValue(tag.Get("form")))
				formElem.Field(i).SetInt(nbr)
			case "float" :
				nbr, _ := strconv.Atoi(r.PostFormValue(tag.Get("form")))
				formElem.Field(i).SetFloat(float64(nbr))
		}
	}
}
