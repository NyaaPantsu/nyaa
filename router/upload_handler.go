package router

import (
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/upload"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/languages"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/gorilla/mux"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	user := GetUser(r)
	if !uploadService.IsUploadEnabled(*user) {
		http.Error(w, "Error uploads are disabled", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		UploadPostHandler(w, r)
	}

	UploadGetHandler(w, r)
}

func UploadPostHandler(w http.ResponseWriter, r *http.Request) {
	var uploadForm UploadForm
	defer r.Body.Close()
	user := GetUser(r)
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

	var sameTorrents int
	db.ORM.Model(&model.Torrent{}).Where("torrent_hash = ?", uploadForm.Infohash).Count(&sameTorrents)
	if sameTorrents > 0 {
		messages.AddError("errors", "Torrent already in database !")
	}

	if !messages.HasErrors() {
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
			WebsiteLink: uploadForm.WebsiteLink,
			UploaderID:  user.ID}
		db.ORM.Create(&torrent)

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

func UploadGetHandler(w http.ResponseWriter, r *http.Request) {
	messages := msg.GetMessages(r) // new util for errors and infos

	var uploadForm UploadForm
	_ = uploadForm.ExtractInfo(r)
	user := GetUser(r)
	if userPermission.NeedsCaptcha(user) {
		uploadForm.CaptchaID = captcha.GetID()
	} else {
		uploadForm.CaptchaID = ""
	}

	utv := UploadTemplateVariables{
		Upload:     uploadForm,
		FormErrors: messages.GetAllErrors(),
		Search:     NewSearchForm(),
		Navigation: NewNavigation(),
		T:          languages.GetTfuncFromRequest(r),
		User:       GetUser(r),
		URL:        r.URL,
		Route:      mux.CurrentRoute(r),
	}
	err := uploadTemplate.ExecuteTemplate(w, "index.html", utv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
