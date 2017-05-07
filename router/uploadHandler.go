package router

import (
	"fmt"
	"github.com/ewhal/nyaa/db"
	"github.com/ewhal/nyaa/model"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var uploadTemplate = template.Must(template.New("upload").Funcs(FuncMap).ParseFiles("templates/index.html", "templates/upload.html"))

func init() {
	template.Must(uploadTemplate.ParseGlob("templates/_*.html")) // common
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var uploadForm UploadForm
	if r.Method == "POST" {
		defer r.Body.Close()
		err = uploadForm.ExtractInfo(r)
		if err == nil {
			//validate name + hash
			//add to db and redirect depending on result
			torrent := model.Torrents{
				Name:         uploadForm.Name,
				Category:     uploadForm.CategoryId,
				Sub_Category: uploadForm.SubCategoryId,
				Status:       1,
				Hash:         uploadForm.Infohash,
				Date:         time.Now().Unix(),
				Description:  []byte{},
				Comments:     []byte{}}
			fmt.Printf("%+v\n", torrent)
			db.ORM.Create(&torrent)
			fmt.Printf("%+v\n", torrent)
			url, err := Router.Get("view_torrent").URL("id", strconv.Itoa(torrent.Id))
			if err == nil {
				http.Redirect(w, r, url.String(), 302)
			}
		}
		fmt.Printf("%+v\n", uploadForm)
	} else if r.Method == "GET" {
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
