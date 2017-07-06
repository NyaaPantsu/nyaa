package controllers

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/models/comments"
	"github.com/NyaaPantsu/nyaa/models/reports"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/gin-gonic/gin"
)

// ReassignForm : Structure for reassign Form used by the reassign page
type ReassignForm struct {
	AssignTo uint
	By       string
	Data     string

	Torrents []uint
}

// ExtractInfo : Function to assign values from request to ReassignForm
func (f *ReassignForm) ExtractInfo(c *gin.Context) bool {
	f.By = c.PostForm("by")
	messages := msg.GetMessages(c)
	if f.By != "olduser" && f.By != "torrentid" {
		messages.AddErrorTf("errors", "no_action_exist", f.By)
		return false
	}

	f.Data = strings.Trim(c.PostForm("data"), " \r\n")
	if f.By == "olduser" {
		if f.Data == "" {
			messages.AddErrorT("errors", "user_not_found")
			return false
		} else if strings.Contains(f.Data, "\n") {
			messages.AddErrorT("errors", "multiple_username_error")
			return false
		}
	} else if f.By == "torrentid" {
		if f.Data == "" {
			messages.AddErrorT("errors", "no_id_given")
			return false
		}
		splitData := strings.Split(f.Data, "\n")
		for i, tmp := range splitData {
			tmp = strings.Trim(tmp, " \r")
			torrentID, err := strconv.ParseUint(tmp, 10, 0)
			if err != nil {
				messages.AddErrorTf("errors", "parse_error_line", i+1)
				return false // TODO: Shouldn't it continue to parse the rest and display the errored lines?
			}
			f.Torrents = append(f.Torrents, uint(torrentID))
		}
	}

	tmpID := c.PostForm("to")
	parsed, err := strconv.ParseUint(tmpID, 10, 32)
	if err != nil {
		messages.Error(err)
		return false
	}
	f.AssignTo = uint(parsed)
	_, _, _, _, err = cookies.RetrieveUserFromRequest(c, uint(parsed))
	if err != nil {
		messages.AddErrorTf("errors", "no_user_found_id", int(parsed))
		return false
	}

	return true
}

// ExecuteAction : Function for applying the changes from ReassignForm
func (f *ReassignForm) ExecuteAction() (int, error) {
	var toBeChanged []uint
	var err error
	if f.By == "olduser" {
		toBeChanged, err = users.FindOldUploadsByUsername(f.Data)
		if err != nil {
			return 0, err
		}
	} else if f.By == "torrentid" {
		toBeChanged = f.Torrents
	}

	num := 0
	for _, torrentID := range toBeChanged {
		torrent, err2 := torrents.FindRawByID(torrentID)
		if err2 == nil {
			torrent.UploaderID = f.AssignTo
			torrent.Update(true)
			num++
		}
	}
	return num, nil
}

// IndexModPanel : Controller for showing index page of Mod Panel
func IndexModPanel(c *gin.Context) {
	offset := 10
	torrents, _, _ := torrents.FindAll(offset, 0)
	users, _ := users.FindUsersForAdmin(offset, 0)
	comments, _ := comments.FindAll(offset, 0, "", "")
	torrentReports, _, _ := reports.GetAll(offset, 0)

	panelAdminTemplate(c, torrents, models.TorrentReportsToJSON(torrentReports), users, comments)
}

// TorrentsListPanel : Controller for listing torrents, can accept common search arguments
func TorrentsListPanel(c *gin.Context) {
	page := c.Param("page")
	messages := msg.GetMessages(c)
	deleted := c.Request.URL.Query()["deleted"]
	unblocked := c.Request.URL.Query()["unblocked"]
	blocked := c.Request.URL.Query()["blocked"]
	if deleted != nil {
		messages.AddInfoTf("infos", "torrent_deleted", "")
	}
	if blocked != nil {
		messages.AddInfoT("infos", "torrent_blocked")
	}
	if unblocked != nil {
		messages.AddInfoT("infos", "torrent_unblocked")
	}

	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	searchParam, torrents, count, err := search.ByQueryWithUser(c, pagenum)
	if err != nil {
		messages.Error(err)
	}
	category := ""
	if len(searchParam.Category) > 0 {
		category = searchParam.Category[0].String()
	}
	searchForm := searchForm{
		TorrentParam:     searchParam,
		Category:         category,
		ShowItemsPerPage: true,
	}

	nav := navigation{count, int(searchParam.Max), pagenum, "mod/torrents/p"}

	modelList(c, "admin/torrentlist.jet.html", torrents, nav, searchForm)
}

// TorrentReportListPanel : Controller for listing torrent reports, can accept pages
func TorrentReportListPanel(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	var err error

	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	torrentReports, nbReports, _ := reports.GetAll(offset, (pagenum-1)*offset)

	reportJSON := models.TorrentReportsToJSON(torrentReports)
	nav := navigation{nbReports, offset, pagenum, "mod/reports/p"}
	modelList(c, "admin/torrent_report.jet.html", reportJSON, nav, newSearchForm(c))
}

// UsersListPanel : Controller for listing users, can accept pages
func UsersListPanel(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	var err error

	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	users, nbUsers := users.FindUsersForAdmin(offset, (pagenum-1)*offset)
	nav := navigation{nbUsers, offset, pagenum, "mod/users/p"}
	modelList(c, "admin/userlist.jet.html", users, nav, newSearchForm(c))
}

// CommentsListPanel : Controller for listing comments, can accept pages and userID
func CommentsListPanel(c *gin.Context) {
	page := c.Param("page")
	pagenum := 1
	offset := 100
	userid := c.Query("userid")
	var err error

	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
	var conditions string
	var values []interface{}
	if userid != "" {
		conditions = "user_id = ?"
		values = append(values, userid)
	}

	comments, nbComments := comments.FindAll(offset, (pagenum-1)*offset, conditions, values...)
	nav := navigation{nbComments, offset, pagenum, "mod/comments/p"}
	modelList(c, "admin/commentlist.jet.html", comments, nav, newSearchForm(c))
}

// TorrentEditModPanel : Controller for editing a torrent after GET request
func TorrentEditModPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	torrent, _ := torrents.FindByID(uint(id))

	torrentJSON := torrent.ToJSON()
	uploadForm := upload.NewTorrentRequest()
	uploadForm.Name = torrentJSON.Name
	uploadForm.Category = torrentJSON.Category + "_" + torrentJSON.SubCategory
	uploadForm.Status = torrentJSON.Status
	uploadForm.Hidden = torrent.Hidden
	uploadForm.WebsiteLink = string(torrentJSON.WebsiteLink)
	uploadForm.Description = string(torrentJSON.Description)
	uploadForm.Languages = torrent.Languages

	formTemplate(c, "admin/paneltorrentedit.jet.html", uploadForm)
}

// TorrentPostEditModPanel : Controller for editing a torrent after POST request
func TorrentPostEditModPanel(c *gin.Context) {
	var uploadForm torrentValidator.TorrentRequest
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	messages := msg.GetMessages(c)
	torrent, _ := torrents.FindByID(uint(id))
	if torrent.ID > 0 {
		errUp := upload.ExtractEditInfo(c, &uploadForm)
		if errUp != nil {
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		if !messages.HasErrors() {
			// update some (but not all!) values
			torrent.Name = uploadForm.Name
			torrent.Category = uploadForm.CategoryID
			torrent.SubCategory = uploadForm.SubCategoryID
			torrent.Status = uploadForm.Status
			torrent.Hidden = uploadForm.Hidden
			torrent.WebsiteLink = uploadForm.WebsiteLink
			torrent.Description = uploadForm.Description
			torrent.Languages = uploadForm.Languages
			_, err := torrent.UpdateUnscope()
			messages.AddInfoT("infos", "torrent_updated")
			if err == nil { // We only log edit torrent for admins
				_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
				activities.Log(&models.User{}, torrent.Identifier(), "edit", "torrent_edited_by", strconv.Itoa(int(torrent.ID)), username, getUser(c).Username)
			}
		}
	}
	formTemplate(c, "admin/paneltorrentedit.jet.html", uploadForm)
}

// CommentDeleteModPanel : Controller for deleting a comment
func CommentDeleteModPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	comment, _, err := comments.Delete(uint(id))
	if err == nil {
		activities.Log(&models.User{}, comment.Identifier(), "delete", "comment_deleted_by", strconv.Itoa(int(comment.ID)), comment.User.Username, getUser(c).Username)
	}

	c.Redirect(http.StatusSeeOther, "/mod/comments?deleted")
}

// TorrentDeleteModPanel : Controller for deleting a torrent
func TorrentDeleteModPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	definitely := c.Request.URL.Query()["definitely"]

	var returnRoute = "/mod/torrents"
	torrent, errFind := torrents.FindByID(uint(id))
	if errFind == nil {
		var err error
		if definitely != nil {
			_, _, err = torrent.DefinitelyDelete()

			//delete reports of torrent
			whereParams := structs.CreateWhereParams("torrent_id = ?", id)
			reports, _, _ := reports.FindOrderBy(&whereParams, "", 0, 0)
			for _, report := range reports {
				report.Delete(true)
			}
			returnRoute = "/mod/torrents/deleted"
		} else {
			_, _, err = torrent.Delete(false)

			//delete reports of torrent
			whereParams := structs.CreateWhereParams("torrent_id = ?", id)
			reports, _, _ := reports.FindOrderBy(&whereParams, "", 0, 0)
			for _, report := range reports {
				report.Delete(false)
			}
		}
		if err == nil {
			_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
			activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, getUser(c).Username)
		}
	}

	c.Redirect(http.StatusSeeOther, returnRoute+"?deleted")
}

// TorrentReportDeleteModPanel : Controller for deleting a torrent report
func TorrentReportDeleteModPanel(c *gin.Context) {
	id := c.Query("id")

	fmt.Println(id)
	idNum, _ := strconv.ParseUint(id, 10, 64)
	_, _, _ = reports.Delete(uint(idNum))
	/* If we need to log report delete activity
	if err == nil {
		activity.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_report_deleted_by", strconv.Itoa(int(report.ID)), getUser(c).Username)
	}
	*/
	c.Redirect(http.StatusSeeOther, "/mod/reports?deleted")
}

// TorrentReassignModPanel : Controller for reassigning a torrent, after GET request
func TorrentReassignModPanel(c *gin.Context) {
	formTemplate(c, "admin/reassign.jet.html", ReassignForm{})
}

// TorrentPostReassignModPanel : Controller for reassigning a torrent, after POST request
func TorrentPostReassignModPanel(c *gin.Context) {
	var rForm ReassignForm
	messages := msg.GetMessages(c)

	if rForm.ExtractInfo(c) {
		count, err2 := rForm.ExecuteAction()
		if err2 != nil {
			messages.AddErrorT("errors", "something_went_wrong")
		} else {
			messages.AddInfoTf("infos", "nb_torrents_updated", count)
		}
	}
	formTemplate(c, "admin/reassign.jet.html", rForm)
}

// TorrentsPostListPanel : Controller for listing torrents, after POST request when mass update
func TorrentsPostListPanel(c *gin.Context) {
	torrentManyAction(c)
	TorrentsListPanel(c)
}

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

// DeletedTorrentsModPanel : Controller for viewing deleted torrents, accept common search arguments
func DeletedTorrentsModPanel(c *gin.Context) {
	page := c.Param("page")
	messages := msg.GetMessages(c) // new utils for errors and infos
	deleted := c.Request.URL.Query()["deleted"]
	unblocked := c.Request.URL.Query()["unblocked"]
	blocked := c.Request.URL.Query()["blocked"]
	if deleted != nil {
		messages.AddInfoT("infos", "torrent_deleted_definitely")
	}
	if blocked != nil {
		messages.AddInfoT("infos", "torrent_blocked")
	}
	if unblocked != nil {
		messages.AddInfoT("infos", "torrent_unblocked")
	}
	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	searchParam, torrents, count, err := search.ByQueryDeleted(c, pagenum)
	if err != nil {
		messages.Error(err)
	}
	category := ""
	if len(searchParam.Category) > 0 {
		category = searchParam.Category[0].String()
	}
	searchForm := searchForm{
		TorrentParam:     searchParam,
		Category:         category,
		ShowItemsPerPage: true,
	}

	nav := navigation{count, int(searchParam.Max), pagenum, "mod/torrents/deleted/p"}
	search := searchForm
	modelList(c, "admin/torrentlist.jet.html", torrents, nav, search)
}

// DeletedTorrentsPostPanel : Controller for viewing deleted torrents after a mass update, accept common search arguments
func DeletedTorrentsPostPanel(c *gin.Context) {
	torrentManyAction(c)
	DeletedTorrentsModPanel(c)
}

// TorrentBlockModPanel : Controller to lock torrents, redirecting to previous page
func TorrentBlockModPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	torrent, _, err := torrents.ToggleBlock(uint(id))
	var returnRoute, action string
	if torrent.IsDeleted() {
		returnRoute = "/mod/torrents/deleted"
	} else {
		returnRoute = "/mod/torrents"
	}
	if torrent.IsBlocked() {
		action = "blocked"
	} else {
		action = "unblocked"
	}
	if err == nil {
		_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
		activities.Log(&models.User{}, torrent.Identifier(), action, "torrent_"+action+"_by", strconv.Itoa(int(torrent.ID)), username, getUser(c).Username)
	}

	c.Redirect(http.StatusSeeOther, returnRoute+"?"+action)
}

/*
 * Controller to modify multiple torrents and can be used by the owner of the torrent or admin
 */
func torrentManyAction(c *gin.Context) {
	currentUser := getUser(c)

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

	if !config.Conf.Torrents.Status[status] { // Check if the status exist
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
					_, err := torrent.UpdateUnscope()
					if err == nil {
						_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
						activities.Log(&models.User{}, torrent.Identifier(), "edited", "torrent_edited_by", strconv.Itoa(int(torrent.ID)), username, getUser(c).Username)
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
						_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
						activities.Log(&models.User{}, torrent.Identifier(), "deleted", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, getUser(c).Username)
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
