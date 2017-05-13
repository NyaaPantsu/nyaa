package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ewhal/nyaa/config"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/ewhal/nyaa/service/user/permission"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if config.UploadsDisabled {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	}
	var uploadForm UploadForm
	if r.Method == "POST" {
		defer r.Body.Close()
		user := GetUser(r)
		if userPermission.NeedsCaptcha(user) {
			userCaptcha := captcha.Extract(r)
			if !captcha.Authenticate(userCaptcha) {
				http.Error(w, captcha.ErrInvalidCaptcha.Error(), http.StatusInternalServerError)
				return
			}
		}

		// validation is done in ExtractInfo()
		err := uploadForm.ExtractInfo(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		status := 1 // normal
		if uploadForm.Remake { // overrides trusted
			status = 2
		} else if user.Status == 1 {
			status = 3 // mark as trusted if user is trusted
		}

		var sameTorrents int
		db.ORM.Model(&model.Torrent{}).Where("torrent_hash = ?", uploadForm.Infohash).Count(&sameTorrents)
		if (sameTorrents == 0) {
			// add to db and redirect
			torrent := model.Torrent{
				Name:        uploadForm.Name,
				Category:    uploadForm.CategoryID,
				SubCategory: uploadForm.SubCategoryID,
				Status:      status,
				Hash:        uploadForm.Infohash,
				Date:        time.Now(),
				Filesize:    uploadForm.Filesize,
				Description: uploadForm.Description,
				UploaderID:  user.ID}
			db.ORM.Create(&torrent)
			url, err := Router.Get("view_torrent").URL("id", strconv.FormatUint(uint64(torrent.ID), 10))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, url.String(), 302)
		} else {
			http.Error(w, fmt.Errorf("Torrent already in database!").Error(), http.StatusInternalServerError)
			return
		}
	} else if r.Method == "GET" {
		user := GetUser(r)
		if userPermission.NeedsCaptcha(user) {
			uploadForm.CaptchaID = captcha.GetID()
		} else {
			uploadForm.CaptchaID = ""
		}


		htv := UploadTemplateVariables{uploadForm, NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}
		languages.SetTranslationFromRequest(uploadTemplate, r, "en-us")
		err := uploadTemplate.ExecuteTemplate(w, "index.html", htv)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
