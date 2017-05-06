package modelHelper

import (
	"reflect"

	"github.com/dorajistyle/goyangi/util/log"
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
