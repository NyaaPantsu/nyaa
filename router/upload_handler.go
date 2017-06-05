package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/api"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/service/upload"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
)

// UploadHandler : Main Controller for uploading a torrent
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	user := getUser(r)
	if !uploadService.IsUploadEnabled(user) {
		http.Error(w, "Error uploads are disabled", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		UploadPostHandler(w, r)
	}

	UploadGetHandler(w, r)
}

// UploadPostHandler : Controller for uploading a torrent, after POST request, redirect or makes error in messages
func UploadPostHandler(w http.ResponseWriter, r *http.Request) {
	var uploadForm apiService.TorrentRequest
	defer r.Body.Close()
	user := getUser(r)
	messages := msg.GetMessages(r) // new util for errors and infos

	if userPermission.NeedsCaptcha(user) {
		userCaptcha := captcha.Extract(r)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddError("errors", captcha.ErrInvalidCaptcha.Error())
		}
	}

	// validation is done in ExtractInfo()
	err := uploadForm.ExtractInfo(r)
	if err != nil {
		messages.AddError("errors", err.Error())
	}
	status := model.TorrentStatusNormal
	if uploadForm.Remake { // overrides trusted
		status = model.TorrentStatusRemake
	} else if user.IsTrusted() {
		status = model.TorrentStatusTrusted
	}

	err = torrentService.ExistOrDelete(uploadForm.Infohash, user)
	if err != nil {
		messages.AddError("errors", err.Error())
	}

	if !messages.HasErrors() {
		// add to db and redirect
		torrent := model.Torrent{
			Name:        uploadForm.Name,
			Category:    uploadForm.CategoryID,
			SubCategory: uploadForm.SubCategoryID,
			Status:      status,
			Hidden:      uploadForm.Hidden,
			Hash:        uploadForm.Infohash,
			Date:        time.Now(),
			Filesize:    uploadForm.Filesize,
			Description: uploadForm.Description,
			WebsiteLink: uploadForm.WebsiteLink,
			UploaderID:  user.ID}
		torrent.ParseTrackers(uploadForm.Trackers)
		db.ORM.Create(&torrent)

		if db.ElasticSearchClient != nil {
			err := torrent.AddToESIndex(db.ElasticSearchClient)
			if err == nil {
				log.Infof("Successfully added torrent to ES index.")
			} else {
				log.Errorf("Unable to add torrent to ES index: %s", err)
			}
		}

		torrentService.NewTorrentEvent(Router, user, &torrent)

		// add filelist to files db, if we have one
		if len(uploadForm.FileList) > 0 {
			for _, uploadedFile := range uploadForm.FileList {
				file := model.File{TorrentID: torrent.ID, Filesize: uploadedFile.Filesize}
				err := file.SetPath(uploadedFile.Path)
				if err != nil {
					messages.AddError("errors", err.Error())
				}
				db.ORM.Create(&file)
			}
		}

		url, err := Router.Get("view_torrent").URL("id", strconv.FormatUint(uint64(torrent.ID), 10))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url.String()+"?success", 302)
	}
}

// UploadGetHandler : Controller for uploading a torrent, after GET request or Failed Post request
func UploadGetHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	messages := msg.GetMessages(r) // new util for errors and infos

	var uploadForm apiService.TorrentRequest
	_ = uploadForm.ExtractInfo(r)
	user := getUser(r)
	if userPermission.NeedsCaptcha(user) {
		uploadForm.CaptchaID = captcha.GetID()
	} else {
		uploadForm.CaptchaID = ""
	}

	utv := formTemplateVariables{
		commonTemplateVariables: newCommonVariables(r),
		Form:       uploadForm,
		FormErrors: messages.GetAllErrors(),
	}
	err := uploadTemplate.ExecuteTemplate(w, "index.html", utv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
