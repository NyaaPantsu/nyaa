package messages

import (
	"errors"
	"fmt"

	"github.com/NyaaPantsu/nyaa/models/notifications"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/i18n"
)

// MessagesKey : use for context
const MessagesKey = "nyaapantsu.messages"

// Messages struct
type Messages struct {
	Errors map[string][]string
	Infos  map[string][]string
	c      *gin.Context
	T      i18n.TranslateFunc
}

// GetMessages : Initialize or return the messages object from context
func GetMessages(c *gin.Context) *Messages {
	if rv, ok := c.Get(MessagesKey); ok {
		mes := rv.(*Messages)
		T, _ := publicSettings.GetTfuncAndLanguageFromRequest(c)
		mes.T = T
		mes.c = c
		return mes
	}
	T, _ := publicSettings.GetTfuncAndLanguageFromRequest(c)
	mes := &Messages{make(map[string][]string), make(map[string][]string), c, T}
	announcements, _ := notifications.CheckAnnouncement()
	for _, announcement := range announcements {
		mes.AddInfo("system", announcement.Content)
	}
	c.Set(MessagesKey, mes)
	return mes
}

// AddError : Add an error in category name with message msg
func (mes *Messages) AddError(name string, msg string) error {
	if mes.Errors == nil {
		mes.Errors = make(map[string][]string)
	}
	mes.Errors[name] = append(mes.Errors[name], msg)
	mes.setMessagesInContext()
	return errors.New(msg)
}

// AddErrorf : Add an error in category name with message msg formatted with args
func (mes *Messages) AddErrorf(name string, msg string, args ...interface{}) error {
	return mes.AddError(name, fmt.Sprintf(msg, args...))
}

// AddErrorTf : Add an error in category name with translation string id formatted with args
func (mes *Messages) AddErrorTf(name string, id string, args ...interface{}) error {
	return mes.AddErrorf(name, mes.T(id), args...)
}

// AddErrorT : Add an error in category name with translation string id
func (mes *Messages) AddErrorT(name string, id string) error {
	return mes.AddError(name, mes.T(id))
}

// ImportFromError : Add an error in category name with message msg imported from type error
func (mes *Messages) ImportFromError(name string, err error) error {
	return mes.AddError(name, err.Error())
}

// ImportFromErrorT : Add an error in category name with message msg imported from type error
func (mes *Messages) ImportFromErrorT(name string, err error) error {
	return mes.AddError(name, mes.T(err.Error()))
}

// ImportFromErrorTf : Add an error in category name with message msg imported from type error
func (mes *Messages) ImportFromErrorTf(name string, err error, args ...interface{}) error {
	return mes.AddError(name, fmt.Sprintf(mes.T(err.Error()), args...))
}

// ImportFromError : Aliases to import directly an error in "errors" map index
func (mes *Messages) Error(err error) error {
	return mes.ImportFromError("errors", err)
}

// ErrorT : Aliases to import directly an error in "errors" map index and translate the error
func (mes *Messages) ErrorT(err error) error {
	return mes.ImportFromErrorT("errors", err)
}

// ErrorTf : Aliases to import directly an error in "errors" map index and translate the error with args
func (mes *Messages) ErrorTf(err error, args ...interface{}) error {
	return mes.ImportFromErrorTf("errors", err)
}

// AddInfo : Add an info in category name with message msg
func (mes *Messages) AddInfo(name string, msg string) {
	if mes.Infos == nil {
		mes.Infos = make(map[string][]string)
	}
	mes.Infos[name] = append(mes.Infos[name], msg)
	mes.setMessagesInContext()
}

// AddInfof : Add an info in category name with message msg formatted with args
func (mes *Messages) AddInfof(name string, msg string, args ...interface{}) {
	mes.AddInfo(name, fmt.Sprintf(msg, args...))
}

// AddInfoTf : Add an info in category name with translation string id formatted with args
func (mes *Messages) AddInfoTf(name string, id string, args ...interface{}) {
	mes.AddInfof(name, mes.T(id), args...)
}

// AddInfoT : Add an info in category name with translation string id
func (mes *Messages) AddInfoT(name string, id string) {
	mes.AddInfo(name, mes.T(id))
}

// ClearAllErrors : Erase all errors in messages
func (mes *Messages) ClearAllErrors() {
	mes.Errors = nil
	mes.setMessagesInContext()
}

// ClearAllInfos : Erase all infos in messages
func (mes *Messages) ClearAllInfos() {
	mes.Infos = nil
	mes.setMessagesInContext()
}

// ClearErrors : Erase all errors in messages
func (mes *Messages) ClearErrors(name string) {
	delete(mes.Errors, name)
	mes.setMessagesInContext()
}

// ClearInfos : Erase all infos in messages
func (mes *Messages) ClearInfos(name string) {
	delete(mes.Infos, name)
	mes.setMessagesInContext()
}

// GetAllErrors : Get all errors
func (mes *Messages) GetAllErrors() map[string][]string {
	mes = GetMessages(mes.c) // We need to look if any new errors from other functions has updated context
	return mes.Errors
}

// GetErrors : Get all errors in category name
func (mes *Messages) GetErrors(name string) []string {
	mes = GetMessages(mes.c) // We need to look if any new errors from other functions has updated context
	return mes.Errors[name]
}

// GetAllInfos : Get all infos
func (mes *Messages) GetAllInfos() map[string][]string {
	mes = GetMessages(mes.c) // We need to look if any new errors from other functions has updated context
	return mes.Infos
}

// GetInfos : Get all infos in category name
func (mes *Messages) GetInfos(name string) []string {
	mes = GetMessages(mes.c) // We need to look if any new errors from other functions has updated context
	return mes.Infos[name]
}

// HasErrors : Check if there are errors
func (mes *Messages) HasErrors() bool {
	mes = GetMessages(mes.c) // We need to look if any new errors from other functions has updated context
	return len(mes.Errors) > 0
}

// HasInfos : Check if there are infos
func (mes *Messages) HasInfos() bool {
	mes = GetMessages(mes.c) // We need to look if any new errors from other functions has updated context
	return len(mes.Infos) > 0
}

func (mes *Messages) setMessagesInContext() {
	mes.c.Set(MessagesKey, mes)
}
