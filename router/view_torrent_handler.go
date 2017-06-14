package router

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"os"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/api"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/notifier"
	"github.com/NyaaPantsu/nyaa/service/report"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/filelist"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/gorilla/mux"
)

// ViewHandler : Controller for displaying a torrent
func ViewHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	messages := msg.GetMessages(r)
	user := getUser(r)

	if r.URL.Query()["success"] != nil {
		messages.AddInfo("infos", "Torrent uploaded successfully!")
	}

	torrent, err := torrentService.GetTorrentByID(id)

	if r.URL.Query()["notif"] != nil {
		notifierService.ToggleReadNotification(torrent.Identifier(), user.ID)
	}

	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	b := torrent.ToJSON()
	folder := filelist.FileListToFolder(torrent.FileList, "root")
	captchaID := ""
	if userPermission.NeedsCaptcha(user) {
		captchaID = captcha.GetID()
	}
	htv := viewTemplateVariables{newCommonVariables(r), b, folder, captchaID, messages.GetAllErrors(), messages.GetAllInfos()}

	err = viewTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("ViewHandler(): %s", err)
	}
}

// ViewHeadHandler : Controller for checking a torrent
func ViewHeadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	if err != nil {
		return
	}

	_, err = torrentService.GetRawTorrentByID(uint(id))

	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	w.Write(nil)
}

// PostCommentHandler : Controller for posting a comment
func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := torrentService.GetTorrentByID(id)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	currentUser := getUser(r)
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
	if len(content) > 500 {
		messages.AddErrorT("errors", "comment_toolong")
	}
	if !messages.HasErrors() {
		userID := currentUser.ID

		comment := model.Comment{TorrentID: torrent.ID, UserID: userID, Content: content, CreatedAt: time.Now()}
		err := db.ORM.Create(&comment).Error
		if err != nil {
			messages.ImportFromError("errors", err)
		}
		comment.Torrent = &torrent

		url, err := Router.Get("view_torrent").URL("id", strconv.FormatUint(uint64(torrent.ID), 10))
		torrent.Uploader.ParseSettings()
		if torrent.Uploader.Settings.Get("new_comment") {
			T, _, _ := publicSettings.TfuncAndLanguageWithFallback(torrent.Uploader.Language, torrent.Uploader.Language) // We need to send the notification to every user in their language
			notifierService.NotifyUser(torrent.Uploader, comment.Identifier(), fmt.Sprintf(T("new_comment_on_torrent"), torrent.Name), url.String(), torrent.Uploader.Settings.Get("new_comment_email"))
		}

		if err != nil {
			messages.ImportFromError("errors", err)
		}
	}
	ViewHandler(w, r)
}

// ReportTorrentHandler : Controller for sending a torrent report
func ReportTorrentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	messages := msg.GetMessages(r)
	currentUser := getUser(r)
	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(r)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddErrorT("errors", "bad_captcha")
		}
	}
	if !messages.HasErrors() {
		idNum, _ := strconv.Atoi(id)
		userID := currentUser.ID

		report := model.TorrentReport{
			Description: r.FormValue("report_type"),
			TorrentID:   uint(idNum),
			UserID:      userID,
			CreatedAt:   time.Now(),
		}

		err := db.ORM.Create(&report).Error
		messages.AddInfoTf("infos", "report_msg", id)
		if err != nil {
			messages.ImportFromError("errors", err)
		}
	}
	ViewHandler(w, r)
}

// TorrentEditUserPanel : Controller for editing a user torrent by a user, after GET request
func TorrentEditUserPanel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := r.URL.Query().Get("id")
	torrent, _ := torrentService.GetTorrentByID(id)
	messages := msg.GetMessages(r)
	currentUser := getUser(r)
	if userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
		uploadForm := apiService.NewTorrentRequest()
		uploadForm.Name = torrent.Name
		uploadForm.Category = strconv.Itoa(torrent.Category) + "_" + strconv.Itoa(torrent.SubCategory)
		uploadForm.Remake = torrent.Status == model.TorrentStatusRemake
		uploadForm.WebsiteLink = string(torrent.WebsiteLink)
		uploadForm.Description = string(torrent.Description)
		uploadForm.Hidden = torrent.Hidden
		uploadForm.Language = torrent.Language
		htv := formTemplateVariables{newCommonVariables(r), uploadForm, messages.GetAllErrors(), messages.GetAllInfos()}
		err := userTorrentEd.ExecuteTemplate(w, "index.html", htv)
		log.CheckError(err)
	} else {
		NotFoundHandler(w, r)
	}
}

// TorrentPostEditUserPanel : Controller for editing a user torrent by a user, after post request
func TorrentPostEditUserPanel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var uploadForm apiService.TorrentRequest
	id := r.URL.Query().Get("id")
	messages := msg.GetMessages(r)
	torrent, _ := torrentService.GetTorrentByID(id)
	currentUser := getUser(r)
	if torrent.ID > 0 && userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
		errUp := uploadForm.ExtractEditInfo(r)
		if errUp != nil {
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		if !messages.HasErrors() {
			status := model.TorrentStatusNormal
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
			torrent.Hidden = uploadForm.Hidden
			torrent.WebsiteLink = uploadForm.WebsiteLink
			torrent.Description = uploadForm.Description
			torrent.Language = uploadForm.Language
			// torrent.Uploader = nil // GORM will create a new user otherwise (wtf?!)
			db.ORM.Model(&torrent).UpdateColumn(&torrent)
			messages.AddInfoT("infos", "torrent_updated")
		}
		htv := formTemplateVariables{newCommonVariables(r), uploadForm, messages.GetAllErrors(), messages.GetAllInfos()}
		err := userTorrentEd.ExecuteTemplate(w, "index.html", htv)
		log.CheckError(err)
	} else {
		NotFoundHandler(w, r)
	}
}

// TorrentDeleteUserPanel : Controller for deleting a user torrent by a user
func TorrentDeleteUserPanel(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id := r.URL.Query().Get("id")
	currentUser := getUser(r)
	torrent, _ := torrentService.GetTorrentByID(id)
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

// DownloadTorrent : Controller for downloading a torrent
func DownloadTorrent(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	hash := vars["hash"]

	if hash == "" && len(config.Conf.Torrents.FileStorage) == 0 {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}

	//Check if file exists and open
	Openfile, err := os.Open(fmt.Sprintf("%s%c%s.torrent", config.Conf.Torrents.FileStorage, os.PathSeparator, hash))
	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}
	defer Openfile.Close() //Close after function return

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	torrent, err := torrentService.GetRawTorrentByHash(hash)

	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found.", 404)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.torrent\"", torrent.Name))
	w.Header().Set("Content-Type", "application/x-bittorrent")
	w.Header().Set("Content-Length", FileSize)
	//Send the file
	// We reset the offset to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client
}
