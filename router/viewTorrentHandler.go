package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/service/torrent"
	"github.com/ewhal/nyaa/util"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/ewhal/nyaa/util/log"
	"github.com/gorilla/mux"

	"fmt"
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
	htv := ViewTemplateVariables{b, captcha.Captcha{CaptchaID: captcha.GetID()}, NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}

	languages.SetTranslationFromRequest(viewTemplate, r, "en-us")
	err = viewTemplate.ExecuteTemplate(w, "index.html", htv)
	if err != nil {
		log.Errorf("ViewHandler(): %s", err)
	}
}

func PostCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	userCaptcha := captcha.Extract(r)
	if !captcha.Authenticate(userCaptcha) {
		http.Error(w, "bad captcha", 403)
	}
	currentUser := GetUser(r)
	content := p.Sanitize(r.FormValue("comment"))

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

	userCaptcha := captcha.Extract(r)
	if !captcha.Authenticate(userCaptcha) {
		http.Error(w, "bad captcha", 403)
	}
	currentUser := GetUser(r)

	idNum, err := strconv.Atoi(id)
	userID := currentUser.ID

	torrent, _ := torrentService.GetTorrentById(id)

	report := model.TorrentReport{
		Description: r.FormValue("report_type"),
		TorrentID:   uint(idNum),
		UserID:      userID,
		Torrent:     torrent,
		User:        *currentUser,
	}
	fmt.Println(report)

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
