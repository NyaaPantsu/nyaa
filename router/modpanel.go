package router

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/comment"
	"github.com/NyaaPantsu/nyaa/service/report"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/categories"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/gorilla/mux"
)

// ReassignForm : Structure for reassign Form used by the reassign page
type ReassignForm struct {
	AssignTo uint
	By       string
	Data     string

	Torrents []uint
}

// ExtractInfo : Function to assign values from request to ReassignForm
func (f *ReassignForm) ExtractInfo(r *http.Request) error {
	f.By = r.FormValue("by")
	if f.By != "olduser" && f.By != "torrentid" {
		return fmt.Errorf("what?")
	}

	f.Data = strings.Trim(r.FormValue("data"), " \r\n")
	if f.By == "olduser" {
		if f.Data == "" {
			return fmt.Errorf("No username given")
		} else if strings.Contains(f.Data, "\n") {
			return fmt.Errorf("More than one username given")
		}
	} else if f.By == "torrentid" {
		if f.Data == "" {
			return fmt.Errorf("No IDs given")
		}
		splitData := strings.Split(f.Data, "\n")
		for i, tmp := range splitData {
			tmp = strings.Trim(tmp, " \r")
			torrentID, err := strconv.ParseUint(tmp, 10, 0)
			if err != nil {
				return fmt.Errorf("Couldn't parse number on line %d", i+1)
			}
			f.Torrents = append(f.Torrents, uint(torrentID))
		}
	}

	tmp := r.FormValue("to")
	parsed, err := strconv.ParseUint(tmp, 10, 0)
	if err != nil {
		return err
	}
	f.AssignTo = uint(parsed)
	_, _, _, _, err = userService.RetrieveUser(r, tmp)
	if err != nil {
		return fmt.Errorf("User to assign to doesn't exist")
	}

	return nil
}

// ExecuteAction : Function for applying the changes from ReassignForm
func (f *ReassignForm) ExecuteAction() (int, error) {
	var toBeChanged []uint
	var err error
	if f.By == "olduser" {
		toBeChanged, err = userService.RetrieveOldUploadsByUsername(f.Data)
		if err != nil {
			return 0, err
		}
	} else if f.By == "torrentid" {
		toBeChanged = f.Torrents
	}

	num := 0
	for _, torrentID := range toBeChanged {
		torrent, err2 := torrentService.GetRawTorrentById(torrentID)
		if err2 == nil {
			torrent.UploaderID = f.AssignTo
			db.ORM.Model(&torrent).UpdateColumn(&torrent)
			num++
		}
	}
	return num, nil
}

// newPanelSearchForm : Helper that creates a search form without items/page field
// these need to be used when the templateVariables don't include `navigation`
func newPanelSearchForm() searchForm {
	form := newSearchForm()
	form.ShowItemsPerPage = false
	return form
}

//
func newPanelCommonVariables(r *http.Request) commonTemplateVariables {
	common := newCommonVariables(r)
	common.Search = newPanelSearchForm()
	return common
}

// IndexModPanel : Controller for showing index page of Mod Panel
func IndexModPanel(w http.ResponseWriter, r *http.Request) {
	offset := 10

	torrents, _, _ := torrentService.GetAllTorrents(offset, 0)
	users, _ := userService.RetrieveUsersForAdmin(offset, 0)
	comments, _ := commentService.GetAllComments(offset, 0, "", "")
	torrentReports, _, _ := reportService.GetAllTorrentReports(offset, 0)

	htv := panelIndexVbs{newPanelCommonVariables(r), torrents, model.TorrentReportsToJSON(torrentReports), users, comments}
	err := panelIndex.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

// TorrentsListPanel : Controller for listing torrents, can accept common search arguments
func TorrentsListPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	messages := msg.GetMessages(r)

	deleted := r.URL.Query()["deleted"]
	unblocked := r.URL.Query()["unblocked"]
	blocked := r.URL.Query()["blocked"]
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	searchParam, torrents, count, err := search.SearchByQueryWithUser(r, pagenum)
	searchForm := searchForm{
		SearchParam:      searchParam,
		Category:         searchParam.Category.String(),
		ShowItemsPerPage: true,
	}

	common := newCommonVariables(r)
	common.Navigation = navigation{count, int(searchParam.Max), pagenum, "mod_tlist_page"}
	common.Search = searchForm
	ptlv := modelListVbs{common, torrents, messages.GetAllErrors(), messages.GetAllInfos()}
	err = panelTorrentList.ExecuteTemplate(w, "admin_index.html", ptlv)
	log.CheckError(err)
}

// TorrentReportListPanel : Controller for listing torrent reports, can accept pages
func TorrentReportListPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	messages := msg.GetMessages(r)
	pagenum := 1
	offset := 100
	var err error

	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	torrentReports, nbReports, _ := reportService.GetAllTorrentReports(offset, (pagenum-1)*offset)

	reportJSON := model.TorrentReportsToJSON(torrentReports)
	common := newCommonVariables(r)
	common.Navigation = navigation{nbReports, offset, pagenum, "mod_trlist_page"}
	ptrlv := modelListVbs{common, reportJSON, messages.GetAllErrors(), messages.GetAllInfos()}
	err = panelTorrentReportList.ExecuteTemplate(w, "admin_index.html", ptrlv)
	log.CheckError(err)
}

// UsersListPanel : Controller for listing users, can accept pages
func UsersListPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	pagenum := 1
	offset := 100
	var err error
	messages := msg.GetMessages(r)

	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	users, nbUsers := userService.RetrieveUsersForAdmin(offset, (pagenum-1)*offset)
	common := newCommonVariables(r)
	common.Navigation = navigation{nbUsers, offset, pagenum, "mod_ulist_page"}
	htv := modelListVbs{common, users, messages.GetAllErrors(), messages.GetAllInfos()}
	err = panelUserList.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

// CommentsListPanel : Controller for listing comments, can accept pages and userID
func CommentsListPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	pagenum := 1
	offset := 100
	userid := r.URL.Query().Get("userid")
	var err error
	messages := msg.GetMessages(r)

	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	var conditions string
	var values []interface{}
	if userid != "" {
		conditions = "user_id = ?"
		values = append(values, userid)
	}

	comments, nbComments := commentService.GetAllComments(offset, (pagenum-1)*offset, conditions, values...)
	common := newCommonVariables(r)
	common.Navigation = navigation{nbComments, offset, pagenum, "mod_clist_page"}
	htv := modelListVbs{common, comments, messages.GetAllErrors(), messages.GetAllInfos()}
	err = panelCommentList.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

// TorrentEditModPanel : Controller for editing a torrent after GET request
func TorrentEditModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	torrent, _ := torrentService.GetTorrentById(id)
	messages := msg.GetMessages(r)

	torrentJSON := torrent.ToJSON()
	uploadForm := newUploadForm()
	uploadForm.Name = torrentJSON.Name
	uploadForm.Category = torrentJSON.Category + "_" + torrentJSON.SubCategory
	uploadForm.Status = torrentJSON.Status
	uploadForm.WebsiteLink = string(torrentJSON.WebsiteLink)
	uploadForm.Description = string(torrentJSON.Description)
	htv := formTemplateVariables{newPanelCommonVariables(r), uploadForm, messages.GetAllErrors(), messages.GetAllInfos()}
	err := panelTorrentEd.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

// TorrentPostEditModPanel : Controller for editing a torrent after POST request
func TorrentPostEditModPanel(w http.ResponseWriter, r *http.Request) {
	var uploadForm uploadForm
	id := r.URL.Query().Get("id")
	messages := msg.GetMessages(r)
	torrent, _ := torrentService.GetTorrentById(id)
	if torrent.ID > 0 {
		errUp := uploadForm.ExtractEditInfo(r)
		if errUp != nil {
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		if !messages.HasErrors() {
			// update some (but not all!) values
			torrent.Name = uploadForm.Name
			torrent.Category = uploadForm.CategoryID
			torrent.SubCategory = uploadForm.SubCategoryID
			torrent.Status = uploadForm.Status
			torrent.WebsiteLink = uploadForm.WebsiteLink
			torrent.Description = uploadForm.Description
			// torrent.Uploader = nil // GORM will create a new user otherwise (wtf?!)
			db.ORM.Model(&torrent).UpdateColumn(&torrent)
			messages.AddInfoT("infos", "torrent_updated")
		}
	}
	htv := formTemplateVariables{newPanelCommonVariables(r), uploadForm, messages.GetAllErrors(), messages.GetAllInfos()}
	err := panelTorrentEd.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

// CommentDeleteModPanel : Controller for deleting a comment
func CommentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, _ = userService.DeleteComment(id)
	url, _ := Router.Get("mod_clist").URL()
	http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
}

// TorrentDeleteModPanel : Controller for deleting a torrent
func TorrentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	definitely := r.URL.Query()["definitely"]
	var returnRoute string
	if definitely != nil {
		_, _ = torrentService.DefinitelyDeleteTorrent(id)

		//delete reports of torrent
		whereParams := serviceBase.CreateWhereParams("torrent_id = ?", id)
		reports, _, _ := reportService.GetTorrentReportsOrderBy(&whereParams, "", 0, 0)
		for _, report := range reports {
			reportService.DeleteDefinitelyTorrentReport(report.ID)
		}
		returnRoute = "mod_tlist_deleted"
	} else {
		_, _ = torrentService.DeleteTorrent(id)

		//delete reports of torrent
		whereParams := serviceBase.CreateWhereParams("torrent_id = ?", id)
		reports, _, _ := reportService.GetTorrentReportsOrderBy(&whereParams, "", 0, 0)
		for _, report := range reports {
			reportService.DeleteTorrentReport(report.ID)
		}
		returnRoute = "mod_tlist"
	}
	url, _ := Router.Get(returnRoute).URL()
	http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
}

// TorrentReportDeleteModPanel : Controller for deleting a torrent report
func TorrentReportDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println(id)
	idNum, _ := strconv.ParseUint(id, 10, 64)
	_, _ = reportService.DeleteTorrentReport(uint(idNum))

	url, _ := Router.Get("mod_trlist").URL()
	http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
}

// TorrentReassignModPanel : Controller for reassigning a torrent, after GET request
func TorrentReassignModPanel(w http.ResponseWriter, r *http.Request) {
	messages := msg.GetMessages(r)
	htv := formTemplateVariables{newPanelCommonVariables(r), ReassignForm{}, messages.GetAllErrors(), messages.GetAllInfos()}
	err := panelTorrentReassign.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

// TorrentPostReassignModPanel : Controller for reassigning a torrent, after POST request
func TorrentPostReassignModPanel(w http.ResponseWriter, r *http.Request) {
	var rForm ReassignForm
	messages := msg.GetMessages(r)

	err2 := rForm.ExtractInfo(r)
	if err2 != nil {
		messages.ImportFromError("errors", err2)
	} else {
		count, err2 := rForm.ExecuteAction()
		if err2 != nil {
			messages.AddErrorT("errors", "something_went_wrong")
		} else {
			messages.AddInfoTf("infos", "nb_torrents_updated", count)
		}
	}

	htv := formTemplateVariables{newPanelCommonVariables(r), rForm, messages.GetAllErrors(), messages.GetAllInfos()}
	err := panelTorrentReassign.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

// TorrentsPostListPanel : Controller for listing torrents, after POST request when mass update
func TorrentsPostListPanel(w http.ResponseWriter, r *http.Request) {
	torrentManyAction(r)
	TorrentsListPanel(w, r)
}

// APIMassMod : This function is used on the frontend for the mass
/* Query is: action=status|delete|owner|category|multiple
 * Needed: torrent_id[] Ids of torrents in checkboxes of name torrent_id
 *
 * Needed on context:
 * status=0|1|2|3|4 according to config/torrent.go (can be omitted if action=delete|owner|category|multiple)
 * owner is the User ID of the new owner of the torrents (can be omitted if action=delete|status|category|multiple)
 * category is the category string (eg. 1_3) of the new category of the torrents (can be omitted if action=delete|status|owner|multiple)
 *
 * withreport is the bool to enable torrent reports deletion (can be omitted)
 *
 * In case of action=multiple, torrents can be at the same time changed status, owner and category
 */
func APIMassMod(w http.ResponseWriter, r *http.Request) {
	torrentManyAction(r)
	messages := msg.GetMessages(r) // new util for errors and infos
	var apiJSON []byte
	w.Header().Set("Content-Type", "application/json")

	if !messages.HasErrors() {
		mapOk := map[string]interface{}{"ok": true, "infos": messages.GetAllInfos()["infos"]}
		apiJSON, _ = json.Marshal(mapOk)
	} else { // We need to show error messages
		mapNotOk := map[string]interface{}{"ok": false, "errors": messages.GetAllErrors()["errors"]}
		apiJSON, _ = json.Marshal(mapNotOk)
	}

	w.Write(apiJSON)
}

// DeletedTorrentsModPanel : Controller for viewing deleted torrents, accept common search arguments
func DeletedTorrentsModPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	messages := msg.GetMessages(r) // new util for errors and infos
	deleted := r.URL.Query()["deleted"]
	unblocked := r.URL.Query()["unblocked"]
	blocked := r.URL.Query()["blocked"]
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	searchParam, torrents, count, err := search.SearchByQueryDeleted(r, pagenum)
	searchForm := searchForm{
		SearchParam:      searchParam,
		Category:         searchParam.Category.String(),
		ShowItemsPerPage: true,
	}

	common := newCommonVariables(r)
	common.Navigation = navigation{count, int(searchParam.Max), pagenum, "mod_tlist_page"}
	common.Search = searchForm
	ptlv := modelListVbs{common, torrents, messages.GetAllErrors(), messages.GetAllInfos()}
	err = panelTorrentList.ExecuteTemplate(w, "admin_index.html", ptlv)
	log.CheckError(err)
}

// DeletedTorrentsPostPanel : Controller for viewing deleted torrents after a mass update, accept common search arguments
func DeletedTorrentsPostPanel(w http.ResponseWriter, r *http.Request) {
	torrentManyAction(r)
	DeletedTorrentsModPanel(w, r)
}

// TorrentBlockModPanel : Controller to lock torrents, redirecting to previous page
func TorrentBlockModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	torrent, _, _ := torrentService.ToggleBlockTorrent(id)

	var returnRoute, action string
	if torrent.IsDeleted() {
		returnRoute = "mod_tlist_deleted"
	} else {
		returnRoute = "mod_tlist"
	}
	if torrent.IsBlocked() {
		action = "blocked"
	} else {
		action = "unblocked"
	}
	url, _ := Router.Get(returnRoute).URL()
	http.Redirect(w, r, url.String()+"?"+action, http.StatusSeeOther)
}

/*
 * Controller to modify multiple torrents and can be used by the owner of the torrent or admin
 */
func torrentManyAction(r *http.Request) {
	currentUser := getUser(r)
	r.ParseForm()
	torrentsSelected := r.Form["torrent_id"] // should be []string
	action := r.FormValue("action")
	status, _ := strconv.Atoi(r.FormValue("status"))
	owner, _ := strconv.Atoi(r.FormValue("owner"))
	category := r.FormValue("category")
	withReport, _ := strconv.ParseBool(r.FormValue("withreport"))
	messages := msg.GetMessages(r) // new util for errors and infos
	catID, subCatID := -1, -1
	var err error

	if action == "" {
		messages.AddErrorT("errors", "no_action_selected")
	}
	if action == "status" && r.FormValue("status") == "" { // We need to check the form value, not the int one because hidden is 0
		messages.AddErrorT("errors", "no_move_location_selected")
	}
	if action == "owner" && r.FormValue("owner") == "" { // We need to check the form value, not the int one because renchon is 0
		messages.AddErrorT("errors", "no_owner_selected")
	}
	if action == "category" && category == "" {
		messages.AddErrorT("errors", "no_category_selected")
	}
	if len(torrentsSelected) == 0 {
		messages.AddErrorT("errors", "select_one_element")
	}

	if r.FormValue("withreport") == "" { // Default behavior for withreport
		withReport = false
	}
	if !config.TorrentStatus[status] { // Check if the status exist
		messages.AddErrorTf("errors", "no_status_exist", status)
		status = -1
	}
	if !userPermission.HasAdmin(currentUser) {
		if r.FormValue("status") != "" { // Condition to check if a user try to change torrent status without having the right permission
			if (status == model.TorrentStatusTrusted && !currentUser.IsTrusted()) || status == model.TorrentStatusAPlus || status == 0 {
				status = model.TorrentStatusNormal
			}
		}
		if r.FormValue("owner") != "" { // Only admins can change owner of torrents
			owner = -1
		}
		withReport = false // Users should not be able to remove reports
	}
	if r.FormValue("owner") != "" && userPermission.HasAdmin(currentUser) { // We check that the user given exist and if not we return an error
		_, _, errorUser := userService.RetrieveUserForAdmin(strconv.Itoa(owner))
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

			if !categories.CategoryExists(category) {
				messages.AddErrorT("errors", "invalid_torrent_category")
			}
		}
	}

	if !messages.HasErrors() {
		for _, torrentID := range torrentsSelected {
			torrent, _ := torrentService.GetTorrentById(torrentID)
			if torrent.ID > 0 && userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
				if action == "status" || action == "multiple" || action == "category" || action == "owner" {

					/* If we don't delete, we make changes according to the form posted and we save at the end */
					if r.FormValue("status") != "" && status != -1 {
						torrent.Status = status
						messages.AddInfoTf("infos", "torrent_moved", torrent.Name)
					}
					if r.FormValue("owner") != "" && owner != -1 {
						torrent.UploaderID = uint(owner)
						messages.AddInfoTf("infos", "torrent_owner_changed", torrent.Name)
					}
					if category != "" && catID != -1 && subCatID != -1 {
						torrent.Category = catID
						torrent.SubCategory = subCatID
						messages.AddInfoTf("infos", "torrent_category_changed", torrent.Name)
					}

					/* Changes are done, we save */
					db.ORM.Unscoped().Model(&torrent).UpdateColumn(&torrent)
				} else if action == "delete" {
					_, err = torrentService.DeleteTorrent(torrentID)
					if err != nil {
						messages.ImportFromError("errors", err)
					} else {
						messages.AddInfoTf("infos", "torrent_deleted", torrent.Name)
					}
				} else {
					messages.AddErrorTf("errors", "no_action_exist", action)
				}
				if withReport {
					whereParams := serviceBase.CreateWhereParams("torrent_id = ?", torrentID)
					reports, _, _ := reportService.GetTorrentReportsOrderBy(&whereParams, "", 0, 0)
					for _, report := range reports {
						reportService.DeleteTorrentReport(report.ID)
					}
					messages.AddInfoTf("infos", "torrent_reports_deleted", torrent.Name)
				}
			} else {
				messages.AddErrorTf("errors", "torrent_not_exist", torrentID)
			}
		}
	}
}
