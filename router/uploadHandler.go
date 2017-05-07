package router

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/ewhal/nyaa/service/captcha"
	"github.com/gorilla/mux"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var uploadForm UploadForm
	if r.Method == "POST" {
		defer r.Body.Close()
		// validation is done in ExtractInfo()
		err = uploadForm.ExtractInfo(r)
		if err == nil {

			//add to db and redirect depending on result
			torrent := model.Torrents{
				Name:         uploadForm.Name,
				Category:     uploadForm.CategoryId,
				Sub_Category: uploadForm.SubCategoryId,
				Status:       1,
				Hash:         uploadForm.Infohash,
				Date:         time.Now().Unix(),
				Filesize:     uploadForm.Filesize, // FIXME: should set to NULL instead of 0
				Description:  uploadForm.Description,
				Comments:     []byte{}}
			db.ORM.Create(&torrent)
			fmt.Printf("%+v\n", torrent)
			url, err := Router.Get("view_torrent").URL("id", strconv.Itoa(torrent.Id))
			if err == nil {
				http.Redirect(w, r, url.String(), 302)
			}
		}
	} else if r.Method == "GET" {
		uploadForm.CaptchaID = captcha.GetID()
		htv := UploadTemplateVariables{uploadForm, NewSearchForm(), Navigation{}, r.URL, mux.CurrentRoute(r)}
		err = uploadTemplate.ExecuteTemplate(w, "index.html", htv)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
