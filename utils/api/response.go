package apiUtils

import (
	"net/http"

	msg "github.com/NyaaPantsu/nyaa/utils/messages"

	"github.com/gin-gonic/gin"
)

// ResponseHandler : This function is the global response for every simple Post Request API
// Please use it. Responses are of the type:
// {ok: bool, [errors | infos]: ArrayOfString [, data: ArrayOfObjects, all_errors: ArrayOfObjects]}
// To send errors or infos, you just need to use the Messages Util
func ResponseHandler(c *gin.Context, obj ...interface{}) {
	messages := msg.GetMessages(c)

	var mapOk map[string]interface{}
	if !messages.HasErrors() {
		mapOk = map[string]interface{}{"ok": true, "infos": messages.GetInfos("infos")}
		if len(obj) > 0 {
			mapOk["data"] = obj
			if len(obj) == 1 {
				mapOk["data"] = obj[0]
			}
		}
	} else { // We need to show error messages
		mapOk = map[string]interface{}{"ok": false, "errors": messages.GetErrors("errors"), "all_errors": messages.GetAllErrors()}
		if len(obj) > 0 {
			mapOk["data"] = obj
			if len(obj) == 1 {
				mapOk["data"] = obj[0]
			}
		}
		if len(messages.GetAllErrors()) > 0 && len(messages.GetErrors("errors")) == 0 {
			mapOk["errors"] = "errors"
		}
	}

	c.JSON(http.StatusOK, mapOk)
}
