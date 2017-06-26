package modelHelper

import (
	"reflect"
	"strconv"

	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/gin-gonic/gin"
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
func BindValueForm(form interface{}, c *gin.Context) {
	formElem := reflect.ValueOf(form).Elem()
	for i := 0; i < formElem.NumField(); i++ {
		typeField := formElem.Type().Field(i)
		tag := typeField.Tag
		switch typeField.Type.Name() {
		case "string":
			formElem.Field(i).SetString(c.PostForm(tag.Get("form")))
		case "int":
			nbr, _ := strconv.Atoi(c.PostForm(tag.Get("form")))
			formElem.Field(i).SetInt(int64(nbr))
		case "float":
			nbr, _ := strconv.Atoi(c.PostForm(tag.Get("form")))
			formElem.Field(i).SetFloat(float64(nbr))
		case "bool":
			nbr, _ := strconv.ParseBool(c.PostForm(tag.Get("form")))
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
			checkMinLength(formElem.Field(i), tag, inputName, mes)
		}
		if tag.Get("len_max") != "" && (tag.Get("needed") != "" || formElem.Field(i).Len() > 0) { // Check maximum length
			checkMaxLength(formElem.Field(i), tag, inputName, mes)
		}
		if tag.Get("equalInput") != "" && (tag.Get("needed") != "" || formElem.Field(i).Len() > 0) {
			checkEqualValue(formElem, i, tag, inputName, mes)
		}
		checksOnFieldTypes(formElem.Field(i), inputName, typeField, tag, mes)
	}
}

func checkMinLength(fieldElem reflect.Value, tag reflect.StructTag, inputName string, mes *msg.Messages) {
	lenMin, _ := strconv.Atoi(tag.Get("len_min"))
	if fieldElem.Len() < lenMin {
		mes.AddErrorTf(tag.Get("form"), "error_min_length", strconv.Itoa(lenMin), inputName)
	}
}
func checkMaxLength(fieldElem reflect.Value, tag reflect.StructTag, inputName string, mes *msg.Messages) {
	lenMax, _ := strconv.Atoi(tag.Get("len_max"))
	if fieldElem.Len() > lenMax {
		mes.AddErrorTf(tag.Get("form"), "error_max_length", strconv.Itoa(lenMax), inputName)
	}
}

func checkEqualValue(formElem reflect.Value, i int, tag reflect.StructTag, inputName string, mes *msg.Messages) {
	otherInput := formElem.FieldByName(tag.Get("equalInput"))
	if formElem.Field(i).Interface() != otherInput.Interface() {
		mes.AddErrorTf(tag.Get("form"), "error_same_value", inputName)
	}
}

func checksOnFieldTypes(fieldElem reflect.Value, inputName string, typeField reflect.StructField, tag reflect.StructTag, mes *msg.Messages) {
	switch typeField.Type.Name() {
	case "string":
		if tag.Get("equal") != "" && fieldElem.String() != tag.Get("equal") {
			mes.AddErrorTf(tag.Get("form"), "error_wrong_value", inputName)
		}
		if tag.Get("needed") != "" && fieldElem.String() == "" {
			mes.AddErrorTf(tag.Get("form"), "error_field_needed", inputName)
		}
		if fieldElem.String() == "" && tag.Get("default") != "" {
			fieldElem.SetString(tag.Get("default"))
		}
	case "int":
		if tag.Get("equal") != "" { // Check minimum length
			equal, _ := strconv.Atoi(tag.Get("equal"))
			if fieldElem.Int() > int64(equal) {
				mes.AddErrorTf(tag.Get("form"), "error_wrong_value", inputName)
			}
		}
		if tag.Get("needed") != "" && fieldElem.Int() == 0 {
			mes.AddErrorTf(tag.Get("form"), "error_field_needed", inputName)
		}
		if fieldElem.Int() == 0 && tag.Get("default") != "" && tag.Get("notnull") != "" {
			defaultValue, _ := strconv.Atoi(tag.Get("default"))
			fieldElem.SetInt(int64(defaultValue))
		}
	case "float":
		if tag.Get("equal") != "" { // Check minimum length
			equal, _ := strconv.Atoi(tag.Get("equal"))
			if fieldElem.Float() != float64(equal) {
				mes.AddErrorTf(tag.Get("form"), "error_wrong_value", inputName)
			}
		}
		if tag.Get("needed") != "" && fieldElem.Float() == 0 {
			mes.AddErrorTf(tag.Get("form"), "error_field_needed", inputName)
		}
		if fieldElem.Float() == 0 && tag.Get("default") != "" && tag.Get("notnull") != "" {
			defaultValue, _ := strconv.Atoi(tag.Get("default"))
			fieldElem.SetFloat(float64(defaultValue))
		}
	case "bool":
		if tag.Get("equal") != "" { // Check minimum length
			equal, _ := strconv.ParseBool(tag.Get("equal"))
			if fieldElem.Bool() != equal {
				mes.AddErrorTf(tag.Get("form"), "error_wrong_value", inputName)
			}
		}
		if !fieldElem.Bool() && tag.Get("default") != "" && tag.Get("notnull") != "" {
			defaultValue, _ := strconv.ParseBool(tag.Get("default"))
			fieldElem.SetBool(defaultValue)
		}
	}
}
