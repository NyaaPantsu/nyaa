package moderatorController

import (
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/models/reports"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/gin-gonic/gin"
)

/*
 * Controller to modify multiple torrents and can be used by the owner of the torrent or admin
 */
func torrentManyAction(c *gin.Context) {
	currentUser := router.GetUser(c)

	torrentsSelected := c.PostFormArray("torrent_id") // should be []string
	action := c.PostForm("action")
	status, _ := strconv.Atoi(c.PostForm("status"))
	owner, _ := strconv.Atoi(c.PostForm("owner"))
	category := c.PostForm("category")
	withReport, _ := strconv.ParseBool(c.DefaultPostForm("withreport", "false"))
	messages := msg.GetMessages(c) // new utils for errors and infos
	catID, subCatID := -1, -1
	var err error

	if action == "" {
		messages.AddErrorT("errors", "no_action_selected")
	}
	if action == "status" && c.PostForm("status") == "" { // We need to check the form value, not the int one because hidden is 0
		messages.AddErrorT("errors", "no_move_location_selected")
	}
	if action == "owner" && c.PostForm("owner") == "" { // We need to check the form value, not the int one because renchon is 0
		messages.AddErrorT("errors", "no_owner_selected")
	}
	if action == "category" && category == "" {
		messages.AddErrorT("errors", "no_category_selected")
	}
	if len(torrentsSelected) == 0 {
		messages.AddErrorT("errors", "select_one_element")
	}

	if !config.Get().Torrents.Status[status] { // Check if the status exist
		messages.AddErrorTf("errors", "no_status_exist", status)
		status = -1
	}
	if !currentUser.HasAdmin() {
		if c.PostForm("status") != "" { // Condition to check if a user try to change torrent status without having the right permission
			if (status == models.TorrentStatusTrusted && !currentUser.IsTrusted()) || status == models.TorrentStatusAPlus || status == 0 {
				status = models.TorrentStatusNormal
			}
		}
		if c.PostForm("owner") != "" { // Only admins can change owner of torrents
			owner = -1
		}
		withReport = false // Users should not be able to remove reports
	}
	if c.PostForm("owner") != "" && currentUser.HasAdmin() { // We check that the user given exist and if not we return an error
		_, _, errorUser := users.FindForAdmin(uint(owner))
		if errorUser != nil {
			messages.AddErrorTf("errors", "no_user_found_id", owner)
			owner = -1
		}
	}
	if category != "" {
		catsSplit := strings.Split(category, "_")
		// need this to prevent out of index panics
		if len(catsSplit) == 2 {
			catID, err = strconv.Atoi(catsSplit[0])
			if err != nil {
				messages.AddErrorT("errors", "invalid_torrent_category")
			}
			subCatID, err = strconv.Atoi(catsSplit[1])
			if err != nil {
				messages.AddErrorT("errors", "invalid_torrent_category")
			}

			if !categories.Exists(category) {
				messages.AddErrorT("errors", "invalid_torrent_category")
			}
		}
	}

	if !messages.HasErrors() {
		for _, torrentID := range torrentsSelected {
			id, _ := strconv.Atoi(torrentID)
			torrent, _ := torrents.FindByID(uint(id))
			if torrent.ID > 0 && currentUser.CurrentOrAdmin(torrent.UploaderID) {
				if action == "status" || action == "multiple" || action == "category" || action == "owner" {

					/* If we don't delete, we make changes according to the form posted and we save at the end */
					if c.PostForm("status") != "" && status != -1 {
						torrent.Status = status
						messages.AddInfoTf("infos", "torrent_moved", torrent.Name)
					}
					if c.PostForm("owner") != "" && owner != -1 {
						torrent.UploaderID = uint(owner)
						messages.AddInfoTf("infos", "torrent_owner_changed", torrent.Name)
					}
					if category != "" && catID != -1 && subCatID != -1 {
						torrent.Category = catID
						torrent.SubCategory = subCatID
						messages.AddInfoTf("infos", "torrent_category_changed", torrent.Name)
					}

					/* Changes are done, we save */
					_, err = torrent.UpdateUnscope()
					if err == nil {
						if torrent.Uploader == nil {
							torrent.Uploader = &models.User{}
						}
						_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
						activities.Log(&models.User{}, torrent.Identifier(), "edited", "torrent_edited_by", strconv.Itoa(int(torrent.ID)), username, router.GetUser(c).Username)
					}
				} else if action == "delete" {
					if status == models.TorrentStatusBlocked { // Then we should lock torrents before deleting them
						torrent.Status = status
						messages.AddInfoTf("infos", "torrent_moved", torrent.Name)
						torrent.UpdateUnscope()
					}
					_, _, err = torrent.Delete(false)
					if err != nil {
						messages.ImportFromError("errors", err)
					} else {
						messages.AddInfoTf("infos", "torrent_deleted", torrent.Name)
						if torrent.Uploader == nil {
							torrent.Uploader = &models.User{}
						}
						_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
						activities.Log(&models.User{}, torrent.Identifier(), "deleted", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, router.GetUser(c).Username)
					}
				} else {
					messages.AddErrorTf("errors", "no_action_exist", action)
				}
				if withReport {
					whereParams := structs.CreateWhereParams("torrent_id = ?", torrentID)
					reports, _, _ := reports.FindOrderBy(&whereParams, "", 0, 0)
					for _, report := range reports {
						report.Delete(false)
					}
					messages.AddInfoTf("infos", "torrent_reports_deleted", torrent.Name)
				}
			} else {
				messages.AddErrorTf("errors", "torrent_not_exist", torrentID)
			}
		}
	}
}
