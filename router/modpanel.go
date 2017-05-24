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

type ReassignForm struct {
	AssignTo uint
	By       string
	Data     string

	Torrents []uint
}

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
			torrent_id, err := strconv.ParseUint(tmp, 10, 0)
			if err != nil {
				return fmt.Errorf("Couldn't parse number on line %d", i+1)
			}
			f.Torrents = append(f.Torrents, uint(torrent_id))
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
	for _, torrent_id := range toBeChanged {
		torrent, err2 := torrentService.GetRawTorrentById(torrent_id)
		if err2 == nil {
			torrent.UploaderID = f.AssignTo
			db.ORM.Save(&torrent)
			num += 1
		}
	}
	return num, nil
}

// Helper that creates a search form without items/page field
// these need to be used when the templateVariables don't include `Navigation`
func NewPanelSearchForm() SearchForm {
	form := NewSearchForm()
	form.ShowItemsPerPage = false
	return form
}

func NewPanelCommonVariables(r *http.Request) CommonTemplateVariables {
	common := NewCommonVariables(r)
	common.Search = NewPanelSearchForm()
	return common
}

func IndexModPanel(w http.ResponseWriter, r *http.Request) {
	offset := 10

	torrents, _, _ := torrentService.GetAllTorrents(offset, 0)
	users, _ := userService.RetrieveUsersForAdmin(offset, 0)
	comments, _ := commentService.GetAllComments(offset, 0, "", "")
	torrentReports, _, _ := reportService.GetAllTorrentReports(offset, 0)

	htv := PanelIndexVbs{NewPanelCommonVariables(r), torrents, model.TorrentReportsToJSON(torrentReports), users, comments}
	err := panelIndex.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

func TorrentsListPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

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
	searchForm := SearchForm{
		SearchParam:      searchParam,
		Category:         searchParam.Category.String(),
		ShowItemsPerPage: true,
	}

	messages := msg.GetMessages(r)
	common := NewCommonVariables(r)
	common.Navigation = Navigation{count, int(searchParam.Max), pagenum, "mod_tlist_page"}
	common.Search = searchForm
	ptlv := PanelTorrentListVbs{common, torrents, messages.GetAllErrors(), messages.GetAllInfos()}
	err = panelTorrentList.ExecuteTemplate(w, "admin_index.html", ptlv)
	log.CheckError(err)
}

func TorrentReportListPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	offset := 100

	torrentReports, nbReports, _ := reportService.GetAllTorrentReports(offset, (pagenum-1)*offset)

	reportJSON := model.TorrentReportsToJSON(torrentReports)
	common := NewCommonVariables(r)
	common.Navigation = Navigation{nbReports, offset, pagenum, "mod_trlist_page"}
	ptrlv := PanelTorrentReportListVbs{common, reportJSON}
	err = panelTorrentReportList.ExecuteTemplate(w, "admin_index.html", ptrlv)
	log.CheckError(err)
}

func UsersListPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	offset := 100

	users, nbUsers := userService.RetrieveUsersForAdmin(offset, (pagenum-1)*offset)
	common := NewCommonVariables(r)
	common.Navigation = Navigation{nbUsers, offset, pagenum, "mod_ulist_page"}
	htv := PanelUserListVbs{common, users}
	err = panelUserList.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

func CommentsListPanel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]

	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	offset := 100
	userid := r.URL.Query().Get("userid")
	var conditions string
	var values []interface{}
	if userid != "" {
		conditions = "user_id = ?"
		values = append(values, userid)
	}

	comments, nbComments := commentService.GetAllComments(offset, (pagenum-1)*offset, conditions, values...)
	common := NewCommonVariables(r)
	common.Navigation = Navigation{nbComments, offset, pagenum, "mod_clist_page"}
	htv := PanelCommentListVbs{common, comments}
	err = panelCommentList.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

func TorrentEditModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	torrent, _ := torrentService.GetTorrentById(id)
	messages := msg.GetMessages(r)

	torrentJson := torrent.ToJSON()
	uploadForm := NewUploadForm()
	uploadForm.Name = torrentJson.Name
	uploadForm.Category = torrentJson.Category + "_" + torrentJson.SubCategory
	uploadForm.Status = torrentJson.Status
	uploadForm.WebsiteLink = string(torrentJson.WebsiteLink)
	uploadForm.Description = string(torrentJson.Description)
	htv := PanelTorrentEdVbs{NewPanelCommonVariables(r), uploadForm, messages.GetAllErrors(), messages.GetAllInfos()}
	err := panelTorrentEd.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

func TorrentPostEditModPanel(w http.ResponseWriter, r *http.Request) {
	var uploadForm UploadForm
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
			torrent.Uploader = nil // GORM will create a new user otherwise (wtf?!)
			db.ORM.Save(&torrent)
			messages.AddInfoT("infos", "torrent_updated")
		}
	}
	htv := PanelTorrentEdVbs{NewPanelCommonVariables(r), uploadForm, messages.GetAllErrors(), messages.GetAllInfos()}
	err_ := panelTorrentEd.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err_)
}

func CommentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	_, _ = userService.DeleteComment(id)
	url, _ := Router.Get("mod_clist").URL()
	http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
}

func TorrentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_, _ = torrentService.DeleteTorrent(id)

	//delete reports of torrent
	whereParams := serviceBase.CreateWhereParams("torrent_id = ?", id)
	reports, _, _ := reportService.GetTorrentReportsOrderBy(&whereParams, "", 0, 0)
	for _, report := range reports {
		reportService.DeleteTorrentReport(report.ID)
	}
	url, _ := Router.Get("mod_tlist").URL()
	http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
}

func TorrentReportDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	fmt.Println(id)
	idNum, _ := strconv.ParseUint(id, 10, 64)
	_, _ = reportService.DeleteTorrentReport(uint(idNum))

	url, _ := Router.Get("mod_trlist").URL()
	http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
}

func TorrentReassignModPanel(w http.ResponseWriter, r *http.Request) {
	messages := msg.GetMessages(r)
	htv := PanelTorrentReassignVbs{NewPanelCommonVariables(r), ReassignForm{}, messages.GetAllErrors(), messages.GetAllInfos()}
	err := panelTorrentReassign.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

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

	htv := PanelTorrentReassignVbs{NewPanelCommonVariables(r), rForm, messages.GetAllErrors(), messages.GetAllInfos()}
	err_ := panelTorrentReassign.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err_)
}

func TorrentsPostListPanel(w http.ResponseWriter, r *http.Request) {
	torrentManyAction(r)
	TorrentsListPanel(w, r)
}

/*
 * This function is used on the frontend for the mass
 * Query is: action=status|delete|owner|category|multiple
 * Needed: torrent_id[] Ids of torrents in checkboxes of name torrent_id
 * 
 * Needed on context:
 * status=0|1|2|3|4 according to config/torrent.go (can be omitted if action=delete|owner|category|multiple)
 * owner is the User ID of the new owner of the torrents (can be omitted if action=delete|status|category|multiple)
 * category is the category string (eg. 1_3) of the new category of the torrents (can be omitted if action=delete|status|owner|multiple)
 * withreport is the bool to enable torrent reports deletion when action=delete (can be omitted if action=category|status|owner|multiple)
 *
 * In case of action=multiple, torrents can be at the same time changed status, owner and category  
 */
func ApiMassMod(w http.ResponseWriter, r *http.Request) {
	torrentManyAction(r)
	messages := msg.GetMessages(r) // new util for errors and infos
	var apiJson []byte
	w.Header().Set("Content-Type", "application/json")	
	
	if !messages.HasErrors() {
		mapOk := map[string]bool{"ok": true, "errors": false}
    	apiJson, _ = json.Marshal(mapOk)
	} else { // We need to show error messages
		mapNotOk := map[string]interface{}{"ok": false, "errors": messages.GetAllErrors()}
		apiJson, _ = json.Marshal(mapNotOk)
	}

	w.Write(apiJson)
}

/*
 * Controller to modify multiple torrents and can be used by the owner of the torrent or admin
 */

func torrentManyAction(r *http.Request) {
	currentUser := GetUser(r)
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

	if !userPermission.HasAdmin(currentUser)  {
		if r.FormValue("status") != "" { // Condition to check if a user try to change torrent status without having the right permission
			if (status == model.TorrentStatusTrusted && !currentUser.IsTrusted()) || status == model.TorrentStatusAPlus || status == 0 {
				status = model.TorrentStatusNormal
			}
			if !config.TorrentStatus[status] {
				messages.AddErrorTf("errors", "no_status_exist", status)
				status = -1
			}
		}
		if r.FormValue("owner") != "" { // Only admins can change owner of torrents
			owner = -1
		}
	}
	if r.FormValue("owner") != "" && userPermission.HasAdmin(currentUser) { // We check that the user given exist and if not we return an error
		_, _, errorUser := userService.RetrieveUserForAdmin(strconv.Itoa(owner))
		if errorUser != nil {
			owner = -1
			messages.AddErrorTf("errors", "no_user_found_id", owner)
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
		for _, torrent_id := range torrentsSelected {
			torrent, _ := torrentService.GetTorrentById(torrent_id)
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
					db.ORM.Save(&torrent)
				} else if action == "delete" {
					_, err = torrentService.DeleteTorrent(torrent_id)
					if err != nil {
						messages.ImportFromError("errors", err)
					} else {
						if withReport {
							whereParams := serviceBase.CreateWhereParams("torrent_id = ?", torrent_id)
							reports, _, _ := reportService.GetTorrentReportsOrderBy(&whereParams, "", 0, 0)
							for _, report := range reports {
								reportService.DeleteTorrentReport(report.ID)
							}
						}
						messages.AddInfoTf("infos", "torrent_deleted", torrent.Name)
					}
				} else {
					messages.AddErrorTf("errors", "no_action_exist", action)
				}
			} else {
				messages.AddErrorTf("errors", "torrent_not_exist", torrent_id)
			}
		}
	}
}
