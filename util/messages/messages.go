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
}

func GetMessages(r *http.Request) Messages {
	if rv := context.Get(r, MessagesKey); rv != nil {
        return rv.(Messages)
    } else {
    	context.Set(r, MessagesKey, Messages{})
    	return Messages{make(map[string][]string),make(map[string][]string)}
    }
}

func (mes Messages) AddError(r *http.Request, name string, msg string) {
	mes.Errors[name] = append(mes.Errors[name], msg)
	mes.setMessagesInContext(r)
}

func (mes Messages) ImportFromError(r *http.Request, name string, err error) {
	mes.AddError(r, name, err.Error())
}

func (mes Messages) AddInfo(r *http.Request, name string, msg string) {
	mes.Infos[name] = append(mes.Infos[name], msg)
	mes.setMessagesInContext(r)
}

func (mes Messages) AddErrorf(r *http.Request, name string, msg string, args ...interface{}) {
	mes.Errors[name] = append(mes.Errors[name], fmt.Sprintf(msg, args...))
	mes.setMessagesInContext(r)
}

func (mes Messages) AddInfof(r *http.Request, name string, msg string, args ...interface{}) {
	mes.Infos[name] = append(mes.Infos[name], fmt.Sprintf(msg, args...))
	mes.setMessagesInContext(r)
}

func (mes Messages) ClearInfos(r *http.Request) {
	mes.Errors = nil
	mes.setMessagesInContext(r)
}

func (mes Messages) ClearErrors(r *http.Request) {
	mes.Infos = nil
	mes.setMessagesInContext(r)
}

func (mes Messages) GetInfos(name string) []string {
	return mes.Infos[name]
}
func (mes Messages) GetErrors(name string) []string {
	return mes.Errors[name]
}
func (mes Messages) GetAllInfos() map[string][]string {
	return mes.Infos
}
func (mes Messages) GetAllErrors() map[string][]string {
	return mes.Errors
}

func (mes Messages) setMessagesInContext(r *http.Request) {
    context.Set(r, MessagesKey, mes)
}

func (mes Messages) HasErrors() bool {
	return len(mes.Errors) > 0
}
func (mes Messages) HasInfos() bool {
	return len(mes.Infos) > 0
}