package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/notifier"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/gorilla/mux"
)

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	messages := msg.GetMessages(r)
	user := GetUser(r)

	if (r.URL.Query()["success"] != nil) {
		messages.AddInfo("infos", "Torrent uploaded successfully!")
	}

	torrent, err := torrentService.GetTorrentById(id)
	
	if (r.URL.Query()["notif"] != nil) {
		notifierService.ToggleReadNotification(torrent.Identifier(), user.ID)
	}

	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	b := torrent.ToJSON()
	captchaID := ""
	if userPermission.NeedsCaptcha(user) {
		captchaID = captcha.GetID()
	}
	htv := ViewTemplateVariables{NewCommonVariables(r), b, captchaID, messages.GetAllErrors(), messages.GetAllInfos()}

	err = viewTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("ViewHandler(): %s", err)
	}
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := torrentService.GetTorrentById(id)
	if (err != nil) {
		NotFoundHandler(w, r)
		return
	}

	currentUser := GetUser(r)
	messages := msg.GetMessages(r)

	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(r)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddError("errors", "Bad captcham!")
		}
	}
	content := p.Sanitize(r.FormValue("comment"))

	if strings.TrimSpace(content) == "" {
		messages.AddError("errors", "Comment empty!")
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
	ViewHandler(w,r)
}

func ReportTorrentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	messages := msg.GetMessages(r)
	currentUser := GetUser(r)
	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(r)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddError("errors", "Bad captchae!")
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
	ViewHandler(w,r)
}
