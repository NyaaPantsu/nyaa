package userValidator

import (
	"net/http"
	"path"
	"testing"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/gin-gonic/gin"
)

// run before config/parse.go:init()
var _ = func() (_ struct{}) {
	config.Configpaths[1] = path.Join("..", "..", "..", config.Configpaths[1])
	config.Configpaths[0] = path.Join("..", "..", "..", config.Configpaths[0])
	config.Reload()
	config.Get().DBType = models.SqliteType
	config.Get().DBParams = ":memory:?cache=shared&mode=memory"

	models.ORM, _ = models.GormInit(models.DefaultLogger)
	return
}()

func TestForms(t *testing.T) {
	t.Parallel()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &gin.Context{Request: req}
	messages := msg.GetMessages(c)
	registration := &RegistrationForm{
		"lol", "", "testing", "testing", "xxx", "1",
	}
	login := &LoginForm{"lol", "testing", "/", "false"}
	user := &UserForm{"lol", "", "", "testing", "testing", "testing", 0, "", "test", "test", "test", "test"}
	userSettings := &UserSettingsForm{}
	password := &PasswordForm{"testing", "testing"}
	passwordReset := &SendPasswordResetForm{"lol@gt.com"}
	passwordResetForm := &PasswordResetForm{"testing", "testing"}

	validator.ValidateForm(registration, messages)
	if messages.HasErrors() {
		t.Errorf("Error on RegistrationForm struct, please check validation arguments: %v", messages.GetAllErrors())
	}
	validator.ValidateForm(login, messages)
	if messages.HasErrors() {
		t.Errorf("Error on LoginForm struct, please check validation arguments: %v", messages.GetAllErrors())
	}
	validator.ValidateForm(user, messages)
	if messages.HasErrors() {
		t.Errorf("Error on User struct, please check validation arguments: %v", messages.GetAllErrors())
	}
	validator.ValidateForm(userSettings, messages)
	if messages.HasErrors() {
		t.Errorf("Error on UserSettingsForm struct, please check validation arguments: %v", messages.GetAllErrors())
	}
	validator.ValidateForm(password, messages)
	if messages.HasErrors() {
		t.Errorf("Error on PasswordForm struct, please check validation arguments: %v", messages.GetAllErrors())
	}
	validator.ValidateForm(passwordReset, messages)
	if messages.HasErrors() {
		t.Errorf("Error on SendPasswordResetForm struct, please check validation arguments: %v", messages.GetAllErrors())
	}
	validator.ValidateForm(passwordResetForm, messages)
	if messages.HasErrors() {
		t.Errorf("Error on PasswordResetForm struct, please check validation arguments: %v", messages.GetAllErrors())
	}
}
