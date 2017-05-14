package router

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service"
	"github.com/ewhal/nyaa/service/comment"
	"github.com/ewhal/nyaa/service/report"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/service/user"
	form "github.com/ewhal/nyaa/service/user/form"
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/log"
	"github.com/ewhal/nyaa/util/search"
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
	// TODO: clean shit from user_uploads_old if needed
	return num, nil
}

// Helper that creates a search form without items/page field
// these need to be used when the templateVariables don't include `Navigation`
func NewPanelSearchForm() SearchForm {
	form := NewSearchForm()
	form.ShowItemsPerPage = false
	return form
}


func IndexModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		offset := 10

		torrents, _, _ := torrentService.GetAllTorrents(offset, 0)
		users, _ := userService.RetrieveUsersForAdmin(offset, 0)
		comments, _ := commentService.GetAllComments(offset, 0, "", "")
		torrentReports, _, _ := reportService.GetAllTorrentReports(offset, 0)

		languages.SetTranslationFromRequest(panelIndex, r, "en-us")
		htv := PanelIndexVbs{torrents, model.TorrentReportsToJSON(torrentReports), users, comments, NewPanelSearchForm(), currentUser, r.URL}
		err := panelIndex.ExecuteTemplate(w, "admin_index.html", htv)
		log.CheckError(err)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}

func TorrentsListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
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

		searchParam, torrents, _, err := search.SearchByQuery(r, pagenum)
		searchForm := SearchForm{
			SearchParam:      searchParam,
			Category:         searchParam.Category.String(),
			ShowItemsPerPage: true,
		}

		languages.SetTranslationFromRequest(panelTorrentList, r, "en-us")
		htv := PanelTorrentListVbs{torrents, searchForm, Navigation{int(searchParam.Max), offset, pagenum, "mod_tlist_page"}, currentUser, r.URL}
		err = panelTorrentList.ExecuteTemplate(w, "admin_index.html", htv)
		log.CheckError(err)
	} else {

		http.Error(w, "admins only", http.StatusForbidden)
	}
}

func TorrentReportListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
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
		languages.SetTranslationFromRequest(panelTorrentReportList, r, "en-us")
		htv := PanelTorrentReportListVbs{reportJSON, NewSearchForm(), Navigation{nbReports, offset, pagenum, "mod_trlist_page"}, currentUser, r.URL}
		err = panelTorrentReportList.ExecuteTemplate(w, "admin_index.html", htv)
		log.CheckError(err)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}

func UsersListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
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
		languages.SetTranslationFromRequest(panelUserList, r, "en-us")
		htv := PanelUserListVbs{users, NewSearchForm(), Navigation{nbUsers, offset, pagenum, "mod_ulist_page"}, currentUser, r.URL}
		err = panelUserList.ExecuteTemplate(w, "admin_index.html", htv)
		log.CheckError(err)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}

func CommentsListPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
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
		languages.SetTranslationFromRequest(panelCommentList, r, "en-us")
		htv := PanelCommentListVbs{comments, NewSearchForm(), Navigation{nbComments, offset, pagenum, "mod_clist_page"}, currentUser, r.URL}
		err = panelCommentList.ExecuteTemplate(w, "admin_index.html", htv)
		log.CheckError(err)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}

}

func TorrentEditModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		id := r.URL.Query().Get("id")
		torrent, _ := torrentService.GetTorrentById(id)
		languages.SetTranslationFromRequest(panelTorrentEd, r, "en-us")

		torrentJson := torrent.ToJSON()
		uploadForm := NewUploadForm()
		uploadForm.Name = torrentJson.Name
		uploadForm.Category = torrentJson.Category + "_" + torrentJson.SubCategory
		uploadForm.Status = torrentJson.Status
		uploadForm.Description = string(torrentJson.Description)
		htv := PanelTorrentEdVbs{uploadForm, NewPanelSearchForm(), currentUser, form.NewErrors(), form.NewInfos(), r.URL}
		err := panelTorrentEd.ExecuteTemplate(w, "admin_index.html", htv)
		log.CheckError(err)

	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}

}

func TorrentPostEditModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if !userPermission.HasAdmin(currentUser) {
		http.Error(w, "admins only", http.StatusForbidden)
		return
	}
	var uploadForm UploadForm
	id := r.URL.Query().Get("id")
	err := form.NewErrors()
	infos := form.NewInfos()
	torrent, _ := torrentService.GetTorrentById(id)
	if torrent.ID > 0 {
		errUp := uploadForm.ExtractEditInfo(r)
		if errUp != nil {
			err["errors"] = append(err["errors"], "Failed to update torrent!")
		}
		if len(err) == 0 {
			// update some (but not all!) values
			torrent.Name = uploadForm.Name
			torrent.Category = uploadForm.CategoryID
			torrent.SubCategory = uploadForm.SubCategoryID
			torrent.Status = uploadForm.Status
			torrent.Description = uploadForm.Description
			torrent.Uploader = nil // GORM will create a new user otherwise (wtf?!)
			db.ORM.Save(&torrent)
			infos["infos"] = append(infos["infos"], "Torrent details updated.")
		}
	}
	languages.SetTranslationFromRequest(panelTorrentEd, r, "en-us")
	htv := PanelTorrentEdVbs{uploadForm, NewPanelSearchForm(), currentUser, err, infos, r.URL}
	err_ := panelTorrentEd.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err_)
}

func CommentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	id := r.URL.Query().Get("id")

	if userPermission.HasAdmin(currentUser) {
		_ = form.NewErrors()
		_, _ = userService.DeleteComment(id)
		url, _ := Router.Get("mod_clist").URL()
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}

func TorrentDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	id := r.URL.Query().Get("id")
	if userPermission.HasAdmin(currentUser) {
		_ = form.NewErrors()
		_, _ = torrentService.DeleteTorrent(id)

		//delete reports of torrent
		whereParams := serviceBase.CreateWhereParams("torrent_id = ?", id)
		reports, _, _ := reportService.GetTorrentReportsOrderBy(&whereParams, "", 0, 0)
		for _, report := range reports {
			reportService.DeleteTorrentReport(report.ID)
		}
		url, _ := Router.Get("mod_tlist").URL()
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}

func TorrentReportDeleteModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if userPermission.HasAdmin(currentUser) {
		id := r.URL.Query().Get("id")
		fmt.Println(id)
		idNum, _ := strconv.ParseUint(id, 10, 64)
		_ = form.NewErrors()
		_, _ = reportService.DeleteTorrentReport(uint(idNum))

		url, _ := Router.Get("mod_trlist").URL()
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	} else {
		http.Error(w, "admins only", http.StatusForbidden)
	}
}

func TorrentReassignModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if !userPermission.HasAdmin(currentUser) {
		http.Error(w, "admins only", http.StatusForbidden)
		return
	}
	languages.SetTranslationFromRequest(panelTorrentReassign, r, "en-us")

	htv := PanelTorrentReassignVbs{ReassignForm{}, NewPanelSearchForm(), currentUser, form.NewErrors(), form.NewInfos(), r.URL}
	err := panelTorrentReassign.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err)
}

func TorrentPostReassignModPanel(w http.ResponseWriter, r *http.Request) {
	currentUser := GetUser(r)
	if !userPermission.HasAdmin(currentUser) {
		http.Error(w, "admins only", http.StatusForbidden)
		return
	}
	var rForm ReassignForm
	err := form.NewErrors()
	infos := form.NewInfos()

	err2 := rForm.ExtractInfo(r)
	if err2 != nil {
		err["errors"] = append(err["errors"], err2.Error())
	} else {
		count, err2 := rForm.ExecuteAction()
		if err2 != nil {
			err["errors"] = append(err["errors"], "Something went wrong")
		} else {
			infos["infos"] = append(infos["infos"], fmt.Sprintf("%d torrents updated.", count))
		}
	}

	htv := PanelTorrentReassignVbs{rForm, NewPanelSearchForm(), currentUser, err, infos, r.URL}
	err_ := panelTorrentReassign.ExecuteTemplate(w, "admin_index.html", htv)
	log.CheckError(err_)
}
