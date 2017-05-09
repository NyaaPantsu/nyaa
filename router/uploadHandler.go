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
	"github.com/ewhal/nyaa/service/user"
	"github.com/ewhal/nyaa/util/languages"
	"github.com/gorilla/mux"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if config.UploadsDisabled {
		http.Error(w, "Error uploads are disabled", http.StatusInternalServerError)
		return
	}
	var err error
	var uploadForm UploadForm
	if r.Method == "POST" {
		defer r.Body.Close()
		// validation is done in ExtractInfo()
		err = uploadForm.ExtractInfo(r)
		if err == nil {
			user, _, err := userService.RetrieveCurrentUser(r)
			if err != nil {
				fmt.Printf("error %+v\n", err)
			}
			//add to db and redirect depending on result
			torrent := model.Torrent{
				Name:        uploadForm.Name,
				Category:    uploadForm.CategoryId,
				SubCategory: uploadForm.SubCategoryId,
				Status:      1,
				Hash:        uploadForm.Infohash,
				Date:        time.Now(),
				Filesize:    uploadForm.Filesize, // FIXME: should set to NULL instead of 0
				Description: uploadForm.Description,
				UploaderID:  user.ID}
			db.ORM.Create(&torrent)
			fmt.Printf("%+v\n", torrent)
			url, err := Router.Get("view_torrent").URL("id", strconv.FormatUint(uint64(torrent.ID), 10))
			if err == nil {
				http.Redirect(w, r, url.String(), 302)
			}
		}
	} else if r.Method == "GET" {
		uploadForm.CaptchaID = captcha.GetID()
		htv := UploadTemplateVariables{uploadForm, NewSearchForm(), Navigation{}, GetUser(r), r.URL, mux.CurrentRoute(r)}
		languages.SetTranslationFromRequest(uploadTemplate, r, "en-us")
		err = uploadTemplate.ExecuteTemplate(w, "index.html", htv)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
