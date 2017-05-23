package Messages
import (
	"github.com/gorilla/context"
	"fmt"
	"net/http"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/NyaaPantsu/nyaa/util/languages"
)

const MessagesKey = "messages"

type Messages struct {
	Errors map[string][]string
	Infos map[string][]string
	r     *http.Request
	T     i18n.TranslateFunc
}

func GetMessages(r *http.Request) *Messages {
	if rv := context.Get(r, MessagesKey); rv != nil {
        mes := rv.(*Messages)
        T, _ := languages.GetTfuncAndLanguageFromRequest(r)
        mes.T = T
        mes.r = r
        return mes
    } else {
    	context.Set(r, MessagesKey, &Messages{})
    	T, _ := languages.GetTfuncAndLanguageFromRequest(r)
    	return &Messages{make(map[string][]string),make(map[string][]string), r, T}
    }
}

func (mes *Messages) AddError(name string, msg string) {
	if (mes.Errors == nil) {
		mes.Errors = make(map[string][]string)
	}
	mes.Errors[name] = append(mes.Errors[name], msg)
	mes.setMessagesInContext()
}
func (mes *Messages) AddErrorf( name string, msg string, args ...interface{}) {
	mes.AddError(name, fmt.Sprintf(msg, args...))
}
func (mes *Messages) AddErrorTf( name string, id string, args ...interface{}) {
	mes.AddErrorf(name, mes.T(id), args...)
}
func (mes *Messages) AddErrorT( name string, id string) {
	mes.AddError(name, mes.T(id))
}
func (mes *Messages) ImportFromError(name string, err error) {
	mes.AddError(name, err.Error())
}

func (mes *Messages) AddInfo(name string, msg string) {
	if (mes.Infos == nil) {
		mes.Infos = make(map[string][]string)
	}
	mes.Infos[name] = append(mes.Infos[name], msg)
	mes.setMessagesInContext()
}
func (mes *Messages) AddInfof(name string, msg string, args ...interface{}) {
	mes.AddInfo(name, fmt.Sprintf(msg, args...))
}
func (mes *Messages) AddInfoTf(name string, id string, args ...interface{}) {
	mes.AddInfof(name, mes.T(id), args...)
}
func (mes *Messages) AddInfoT(name string, id string) {
	mes.AddInfo(name, mes.T(id))
}

func (mes *Messages) ClearErrors() {
	mes.Infos = nil
	mes.setMessagesInContext()
}
func (mes *Messages) ClearInfos() {
	mes.Errors = nil
	mes.setMessagesInContext()
}

func (mes *Messages) GetAllErrors() map[string][]string {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return mes.Errors
}
func (mes *Messages) GetErrors(name string) []string {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return mes.Errors[name]
}

func (mes *Messages) GetAllInfos() map[string][]string {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return mes.Infos
}
func (mes *Messages) GetInfos(name string) []string {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return mes.Infos[name]
}

func (mes *Messages) HasErrors() bool {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return len(mes.Errors) > 0
}
func (mes *Messages) HasInfos() bool {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return len(mes.Infos) > 0
}

func (mes *Messages) setMessagesInContext() {
    context.Set(mes.r, MessagesKey, mes)
}