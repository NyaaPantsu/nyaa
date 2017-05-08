package router

import (
	"net/http"
	"strconv"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/service/torrent"
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
	b := torrent.ToJson()
	htv := ViewTemplateVariables{b, captcha.Captcha{CaptchaID: captcha.GetID()}, NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}

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
	content := p.Sanitize(r.FormValue("comment"))

	idNum, err := strconv.Atoi(id)
	comment := model.Comment{Username: "れんちょん", Content: content, TorrentId: idNum}
	db.ORM.Create(&comment)

	url, err := Router.Get("view_torrent").URL("id", id)
	if err == nil {
		http.Redirect(w, r, url.String(), 302)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
