package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/notifier"
	"github.com/NyaaPantsu/nyaa/service/report"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/NyaaPantsu/nyaa/util/log"
	"github.com/NyaaPantsu/nyaa/util/filelist"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/gorilla/mux"
)

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	messages := msg.GetMessages(r)
	user := GetUser(r)

	if r.URL.Query()["success"] != nil {
		messages.AddInfo("infos", "Torrent uploaded successfully!")
	}

	torrent, err := torrentService.GetTorrentById(id)

	if r.URL.Query()["notif"] != nil {
		notifierService.ToggleReadNotification(torrent.Identifier(), user.ID)
	}

	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	b := torrent.ToJSON()
	folder := filelist.FileListToFolder(torrent.FileList)
	captchaID := ""
	if userPermission.NeedsCaptcha(user) {
		captchaID = captcha.GetID()
	}
	htv := ViewTemplateVariables{NewCommonVariables(r), b, folder, captchaID, messages.GetAllErrors(), messages.GetAllInfos()}

	err = viewTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("ViewHandler(): %s", err)
	}
}

func ViewHeadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		return
	}

	_, err = torrentService.GetRawTorrentById(uint(id))

	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	w.Write(nil)
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := torrentService.GetTorrentById(id)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	currentUser := GetUser(r)
	messages := msg.GetMessages(r)

	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(r)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddErrorT("errors", "bad_captcha")
		}
	}
	content := util.Sanitize(r.FormValue("comment"), "comment")

	if strings.TrimSpace(content) == "" {
		messages.AddErrorT("errors", "comment_empty")
	}
	if !messages.HasErrors() {
		userID := currentUser.ID
		comment := model.Comment{TorrentID: torrent.ID, UserID: userID, Content: content, CreatedAt: time.Now()}
		err := db.ORM.Create(&comment).Error
		comment.Torrent = &torrent

		url, err := Router.Get("view_torrent").URL("id", strconv.FormatUint(uint64(torrent.ID), 10))
		torrent.Uploader.ParseSettings()
		if torrent.Uploader.Settings.Get("new_comment") {
			T, _, _ := languages.TfuncAndLanguageWithFallback(torrent.Uploader.Language, torrent.Uploader.Language) // We need to send the notification to every user in their language
			notifierService.NotifyUser(torrent.Uploader, comment.Identifier(), fmt.Sprintf(T("new_comment_on_torrent"), torrent.Name), url.String(), torrent.Uploader.Settings.Get("new_comment_email"))
		}

		if err != nil {
			messages.ImportFromError("errors", err)
		}
	}
	ViewHandler(w, r)
}

func ReportTorrentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	messages := msg.GetMessages(r)
	currentUser := GetUser(r)
	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(r)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddErrorT("errors", "bad_captcha")
		}
	}
	if !messages.HasErrors() {
		idNum, err := strconv.Atoi(id)
		userID := currentUser.ID

		report := model.TorrentReport{
			Description: r.FormValue("report_type"),
			TorrentID:   uint(idNum),
			UserID:      userID,
			CreatedAt:   time.Now(),
		}

		err = db.ORM.Create(&report).Error
		if err != nil {
			messages.ImportFromError("errors", err)
		}
	}
	ViewHandler(w, r)
}

func TorrentEditUserPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	torrent, _ := torrentService.GetTorrentById(id)
	messages := msg.GetMessages(r)
	currentUser := GetUser(r)
	if userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
		uploadForm := NewUploadForm()
		uploadForm.Name = torrent.Name
		uploadForm.Category = strconv.Itoa(torrent.Category) + "_" + strconv.Itoa(torrent.SubCategory)
		uploadForm.Remake = torrent.Status == model.TorrentStatusRemake
		uploadForm.WebsiteLink = string(torrent.WebsiteLink)
		uploadForm.Description = string(torrent.Description)
		htv := UserTorrentEdVbs{NewCommonVariables(r), uploadForm, messages.GetAllErrors(), messages.GetAllInfos()}
		err := userTorrentEd.ExecuteTemplate(w, "index.html", htv)
		log.CheckError(err)
	} else {
		NotFoundHandler(w, r)
	}
}

func TorrentPostEditUserPanel(w http.ResponseWriter, r *http.Request) {
	var uploadForm UploadForm
	id := r.URL.Query().Get("id")
	messages := msg.GetMessages(r)
	torrent, _ := torrentService.GetTorrentById(id)
	currentUser := GetUser(r)
	if torrent.ID > 0 && userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
		errUp := uploadForm.ExtractEditInfo(r)
		if errUp != nil {
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		if !messages.HasErrors() {
			status := model.TorrentStatusNormal
			uploadForm.Remake = r.FormValue(UploadFormRemake) == "on"
			if uploadForm.Remake { // overrides trusted
				status = model.TorrentStatusRemake
			} else if currentUser.IsTrusted() {
				status = model.TorrentStatusTrusted
			}
			// update some (but not all!) values
			torrent.Name = uploadForm.Name
			torrent.Category = uploadForm.CategoryID
			torrent.SubCategory = uploadForm.SubCategoryID
			torrent.Status = status
			torrent.WebsiteLink = uploadForm.WebsiteLink
			torrent.Description = uploadForm.Description
			torrent.Uploader = nil // GORM will create a new user otherwise (wtf?!)
			db.ORM.Save(&torrent)
			messages.AddInfoT("infos", "torrent_updated")
		}
		htv := UserTorrentEdVbs{NewCommonVariables(r), uploadForm, messages.GetAllErrors(), messages.GetAllInfos()}
		err_ := userTorrentEd.ExecuteTemplate(w, "index.html", htv)
		log.CheckError(err_)
	} else {
		NotFoundHandler(w, r)
	}
}

func TorrentDeleteUserPanel(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	currentUser := GetUser(r)
	torrent, _ := torrentService.GetTorrentById(id)
	if userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
		_, err := torrentService.DeleteTorrent(id)
		if err == nil {
			//delete reports of torrent
			whereParams := serviceBase.CreateWhereParams("torrent_id = ?", id)
			reports, _, _ := reportService.GetTorrentReportsOrderBy(&whereParams, "", 0, 0)
			for _, report := range reports {
				reportService.DeleteTorrentReport(report.ID)
			}
		}
		url, _ := Router.Get("home").URL()
		http.Redirect(w, r, url.String()+"?deleted", http.StatusSeeOther)
	} else {
		NotFoundHandler(w, r)
	}
}
