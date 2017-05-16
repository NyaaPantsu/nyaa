package router

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/log"
	"github.com/gorilla/mux"
)

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := torrentService.GetTorrentById(id)
	if err != nil {
		NotFoundHandler(w, r)
		return
	}
	b := torrent.ToJSON()
	captchaID := ""
	user := GetUser(r)
	if userPermission.NeedsCaptcha(user) {
		captchaID = captcha.GetID()
	}
	htv := ViewTemplateVariables{b, captchaID, NewSearchForm(), NewNavigation(), user, r.URL, mux.CurrentRoute(r)}

	languages.SetTranslationFromRequest(viewTemplate, r)
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
