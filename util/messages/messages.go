package Messages
import (
	"github.com/gorilla/context"
	"fmt"
	"net/http"
)

const MessagesKey = "messages"

type Messages struct {
	Errors map[string][]string
	Infos map[string][]string
	r     *http.Request
}

func GetMessages(r *http.Request) Messages {
	if rv := context.Get(r, MessagesKey); rv != nil {
        return rv.(Messages)
    } else {
    	context.Set(r, MessagesKey, Messages{})
    	return Messages{make(map[string][]string),make(map[string][]string), r}
    }
}

func (mes Messages) AddError(name string, msg string) {
	mes.Errors[name] = append(mes.Errors[name], msg)
	mes.setMessagesInContext()
}
func (mes Messages) AddErrorf( name string, msg string, args ...interface{}) {
	mes.AddError(name, fmt.Sprintf(msg, args...))
}
func (mes Messages) ImportFromError(name string, err error) {
	mes.AddError(name, err.Error())
}

func (mes Messages) AddInfo(name string, msg string) {
	mes.Infos[name] = append(mes.Infos[name], msg)
	mes.setMessagesInContext()
}
func (mes Messages) AddInfof(name string, msg string, args ...interface{}) {
	mes.AddInfo(name, fmt.Sprintf(msg, args...))
}

func (mes Messages) ClearErrors() {
	mes.Infos = nil
	mes.setMessagesInContext()
}
func (mes Messages) ClearInfos() {
	mes.Errors = nil
	mes.setMessagesInContext()
}

func (mes Messages) GetAllErrors() map[string][]string {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return mes.Errors
}
func (mes Messages) GetErrors(name string) []string {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return mes.Errors[name]
}

func (mes Messages) GetAllInfos() map[string][]string {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return mes.Infos
}
func (mes Messages) GetInfos(name string) []string {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return mes.Infos[name]
}

func (mes Messages) HasErrors() bool {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return len(mes.Errors) > 0
}
func (mes Messages) HasInfos() bool {
	mes = GetMessages(mes.r) // We need to look if any new errors from other functions has updated context
	return len(mes.Infos) > 0
}

func (mes Messages) setMessagesInContext() {
    context.Set(mes.r, MessagesKey, mes)
}