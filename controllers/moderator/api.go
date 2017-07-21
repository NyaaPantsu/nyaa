package moderatorController

import (
	"net/http"

	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/gin-gonic/gin"
)

// APIMassMod : This function is used on the frontend for the mass
/* Query is: action=status|delete|owner|category|multiple
 * Needed: torrent_id[] Ids of torrents in checkboxes of name torrent_id
 *
 * Needed on context:
 * status=0|1|2|3|4 according to config/find.go (can be omitted if action=delete|owner|category|multiple)
 * owner is the User ID of the new owner of the torrents (can be omitted if action=delete|status|category|multiple)
 * category is the category string (eg. 1_3) of the new category of the torrents (can be omitted if action=delete|status|owner|multiple)
 *
 * withreport is the bool to enable torrent reports deletion (can be omitted)
 *
 * In case of action=multiple, torrents can be at the same time changed status, owner and category
 */
func APIMassMod(c *gin.Context) {
	torrentManyAction(c)
	messages := msg.GetMessages(c) // new utils for errors and infos
	c.Header("Content-Type", "application/json")

	var mapOk map[string]interface{}
	if !messages.HasErrors() {
		mapOk = map[string]interface{}{"ok": true, "infos": messages.GetAllInfos()["infos"]}
	} else { // We need to show error messages
		mapOk = map[string]interface{}{"ok": false, "errors": messages.GetAllErrors()["errors"]}
	}

	c.JSON(http.StatusOK, mapOk)
}
