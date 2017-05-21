package router

import (
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
	"github.com/NyaaPantsu/nyaa/util"
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

	currentUser := GetUser(r)
	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(r)
		if !captcha.Authenticate(userCaptcha) {
			http.Error(w, "bad captcha", 403)
			return
		}
	}
	content := p.Sanitize(r.FormValue("comment"))

	if strings.TrimSpace(content) == "" {
		http.Error(w, "comment empty", 406)
		return
	}

	idNum, err := strconv.Atoi(id)

	userID := currentUser.ID
	comment := model.Comment{TorrentID: uint(idNum), UserID: userID, Content: content, CreatedAt: time.Now()}

	err = db.ORM.Create(&comment).Error
	if err != nil {
		util.SendError(w, err, 500)
		return
	}

	url, err := Router.Get("view_torrent").URL("id", id)
	if err == nil {
		http.Redirect(w, r, url.String(), 302)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ReportTorrentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	currentUser := GetUser(r)
	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(r)
		if !captcha.Authenticate(userCaptcha) {
			http.Error(w, "bad captcha", 403)
			return
		}
	}

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
		util.SendError(w, err, 500)
		return
	}

	url, err := Router.Get("view_torrent").URL("id", id)
	if err == nil {
		http.Redirect(w, r, url.String(), 302)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
