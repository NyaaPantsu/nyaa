package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service/captcha"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if config.UploadsDisabled {
		user := GetUser(r)
		if config.AdminsAreStillAllowedTo && user.Status != 2 && config.TrustedUsersAreStillAllowedTo && user.Status != 1 {
			http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
			return
		} else if !config.AdminsAreStillAllowedTo && user.Status == 2 {
			http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
			return
		} else if !config.TrustedUsersAreStillAllowedTo && user.Status == 1 {
			http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
			return
		}
	}

	if r.Method == "POST" {
		UploadPostHandler(w, r)
	} else if r.Method == "GET" {
		UploadGetHandler(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func UploadPostHandler(w http.ResponseWriter, r *http.Request) {
	var uploadForm UploadForm
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
	status := 1            // normal
	if uploadForm.Remake { // overrides trusted
		status = 2
	} else if user.Status == 1 {
		status = 3 // mark as trusted if user is trusted
	}

	var sameTorrents int
	db.ORM.Model(&model.Torrent{}).Table(config.TableName).Where("torrent_hash = ?", uploadForm.Infohash).Count(&sameTorrents)
	if sameTorrents == 0 {
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
		db.ORM.Table(config.TableName).Create(&torrent)

		// add filelist to files db, if we have one
		if len(uploadForm.FileList) > 0 {
			for _, uploadedFile := range uploadForm.FileList {
				file := model.File{TorrentID: torrent.ID, Filesize: uploadedFile.Filesize}
				err := file.SetPath(uploadedFile.Path)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				db.ORM.Create(&file)
			}
		}

		url, err := Router.Get("view_torrent").URL("id", strconv.FormatUint(uint64(torrent.ID), 10))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, url.String(), 302)
	} else {
		err = fmt.Errorf("Torrent already in database!")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func UploadGetHandler(w http.ResponseWriter, r *http.Request) {
	languages.SetTranslationFromRequest(uploadTemplate, r)

	var uploadForm UploadForm
	user := GetUser(r)
	if userPermission.NeedsCaptcha(user) {
		uploadForm.CaptchaID = captcha.GetID()
	} else {
		uploadForm.CaptchaID = ""
	}

	utv := UploadTemplateVariables{
		Upload:     uploadForm,
		Search:     NewSearchForm(),
		Navigation: NewNavigation(),
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
