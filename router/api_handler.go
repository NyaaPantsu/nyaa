package router

import (
	"encoding/json"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/api"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/service/upload"
	"github.com/NyaaPantsu/nyaa/service/user"
	"github.com/NyaaPantsu/nyaa/service/user/form"
	"github.com/NyaaPantsu/nyaa/util/crypto"
	"github.com/NyaaPantsu/nyaa/util/log"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/modelHelper"
	"github.com/NyaaPantsu/nyaa/util/search"
	"github.com/gorilla/mux"
)

// APIHandler : Controller for api request on torrent list
func APIHandler(w http.ResponseWriter, r *http.Request) {
	t := r.URL.Query().Get("t")
	if t != "" {
		RSSTorznabHandler(w, r)
	} else {
		vars := mux.Vars(r)
		page := vars["page"]
		whereParams := serviceBase.WhereParams{}
		req := apiService.TorrentsRequest{}
		defer r.Body.Close()

		contentType := r.Header.Get("Content-Type")
		if contentType == "application/json" {
			d := json.NewDecoder(r.Body)
			if err := d.Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			if req.MaxPerPage == 0 {
				req.MaxPerPage = config.Conf.Navigation.TorrentsPerPage
			}
			if req.Page <= 0 {
				req.Page = 1
			}

			whereParams = req.ToParams()
		} else {
			var err error
			maxString := r.URL.Query().Get("max")
			if maxString != "" {
				req.MaxPerPage, err = strconv.Atoi(maxString)
				if !log.CheckError(err) {
					req.MaxPerPage = config.Conf.Navigation.TorrentsPerPage
				}
			} else {
				req.MaxPerPage = config.Conf.Navigation.TorrentsPerPage
			}

			req.Page = 1
			if page != "" {
				req.Page, err = strconv.Atoi(html.EscapeString(page))
				if !log.CheckError(err) {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				if req.Page <= 0 {
					NotFoundHandler(w, r)
					return
				}
			}
		}

		torrents, nbTorrents, err := torrentService.GetTorrents(whereParams, req.MaxPerPage, req.MaxPerPage*(req.Page-1))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b := model.APIResultJSON{
			Torrents: model.TorrentsToJSON(torrents),
		}
		b.QueryRecordCount = req.MaxPerPage
		b.TotalRecordCount = nbTorrents
		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(b)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// APIViewHandler : Controller for viewing a torrent by its ID
func APIViewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	torrent, err := torrentService.GetTorrentByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	b := torrent.ToJSON()
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
}

// APIViewHeadHandler : Controller for checking a torrent by its ID
func APIViewHeadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 32)
	defer r.Body.Close()
	if err != nil {
		return
	}

	_, err = torrentService.GetRawTorrentByID(uint(id))

	if err != nil {
		NotFoundHandler(w, r)
		return
	}

	w.Write(nil)
}

// APIUploadHandler : Controller for uploading a torrent with api
func APIUploadHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username := r.FormValue("username")
	user, _, _, _, err := userService.RetrieveUserByAPITokenAndName(token, username)
	messages := msg.GetMessages(r)
	defer r.Body.Close()

	if err != nil {
		messages.AddError("errors", "Error API token doesn't exist")
	}

	if !uploadService.IsUploadEnabled(user) {
		messages.AddError("errors", "Error uploads are disabled")
	}

	if user.ID == 0 {
		messages.ImportFromError("errors", apiService.ErrAPIKey)
	}

	if !messages.HasErrors() {
		upload := apiService.TorrentRequest{}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" && !strings.HasPrefix(contentType, "multipart/form-data") && r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			// TODO What should we do here ? upload is empty so we shouldn't
			// create a torrent from it
			messages.AddError("errors", "Please provide either of Content-Type: application/json header or multipart/form-data")
		}
		// As long as the right content-type is sent, formValue is smart enough to parse it
		err = upload.ExtractInfo(r)
		if err != nil {
			messages.ImportFromError("errors", err)
		}

		if !messages.HasErrors() {
			status := model.TorrentStatusNormal
			if upload.Remake { // overrides trusted
				status = model.TorrentStatusRemake
			} else if user.IsTrusted() {
				status = model.TorrentStatusTrusted
			}
			err = torrentService.ExistOrDelete(upload.Infohash, user)
			if err != nil {
				messages.ImportFromError("errors", err)
			}
			if !messages.HasErrors() {
				torrent := model.Torrent{
					Name:        upload.Name,
					Category:    upload.CategoryID,
					SubCategory: upload.SubCategoryID,
					Status:      status,
					Hidden:      upload.Hidden,
					Hash:        upload.Infohash,
					Date:        time.Now(),
					Filesize:    upload.Filesize,
					Description: upload.Description,
					WebsiteLink: upload.WebsiteLink,
					UploaderID:  user.ID}
				torrent.ParseTrackers(upload.Trackers)
				db.ORM.Create(&torrent)

				if db.ElasticSearchClient != nil {
					err := torrent.AddToESIndex(db.ElasticSearchClient)
					if err == nil {
						log.Infof("Successfully added torrent to ES index.")
					} else {
						log.Errorf("Unable to add torrent to ES index: %s", err)
					}
				} else {
					log.Errorf("Unable to create elasticsearch client: %s", err)
				}
				messages.AddInfoT("infos", "torrent_uploaded")
				torrentService.NewTorrentEvent(Router, user, &torrent)
				// add filelist to files db, if we have one
				if len(upload.FileList) > 0 {
					for _, uploadedFile := range upload.FileList {
						file := model.File{TorrentID: torrent.ID, Filesize: upload.Filesize}
						err := file.SetPath(uploadedFile.Path)
						if err != nil {
							http.Error(w, err.Error(), http.StatusBadRequest)
							return
						}
						db.ORM.Create(&file)
					}
				}
				apiResponseHandler(w, r, torrent.ToJSON())
				return
			}
		}
	}
	apiResponseHandler(w, r)
}

// APIUpdateHandler : Controller for updating a torrent with api
// FIXME Impossible to update a torrent uploaded by user 0
func APIUpdateHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username := r.FormValue("username")
	user, _, _, _, err := userService.RetrieveUserByAPITokenAndName(token, username)
	defer r.Body.Close()
	messages := msg.GetMessages(r)

	if err != nil {
		messages.AddError("errors", "Error API token doesn't exist")
	}

	if !uploadService.IsUploadEnabled(user) {
		messages.AddError("errors", "Error uploads are disabled")
	}

	if user.ID == 0 {
		messages.ImportFromError("errors", apiService.ErrAPIKey)
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && !strings.HasPrefix(contentType, "multipart/form-data") && r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		// TODO What should we do here ? upload is empty so we shouldn't
		// create a torrent from it
		messages.AddError("errors", "Please provide either of Content-Type: application/json header or multipart/form-data")
	}
	update := apiService.UpdateRequest{}
	err = update.Update.ExtractEditInfo(r)
	if err != nil {
		messages.ImportFromError("errors", err)
	}
	if !messages.HasErrors() {
		torrent := model.Torrent{}
		db.ORM.Where("torrent_id = ?", r.FormValue("id")).First(&torrent)
		if torrent.ID == 0 {
			messages.ImportFromError("errors", apiService.ErrTorrentID)
		}
		if torrent.UploaderID != 0 && torrent.UploaderID != user.ID { //&& user.Status != mod
			messages.ImportFromError("errors", apiService.ErrRights)
		}
		update.UpdateTorrent(&torrent, user)
		torrentService.UpdateTorrent(&torrent)
	}
	apiResponseHandler(w, r)
}

// APISearchHandler : Controller for searching with api
func APISearchHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page := vars["page"]
	defer r.Body.Close()

	// db params url
	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if pagenum <= 0 {
			NotFoundHandler(w, r)
			return
		}
	}

	_, torrents, _, err := search.SearchByQueryWithUser(r, pagenum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b := model.TorrentsToJSON(torrents)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// APILoginHandler : Login with API
// This is not an OAuth api like and shouldn't be used for anything except getting the API Token in order to not store a password
func APILoginHandler(w http.ResponseWriter, r *http.Request) {
	b := form.LoginForm{}
	messages := msg.GetMessages(r)
	contentType := r.Header.Get("Content-type")
	if !strings.HasPrefix(contentType, "application/json") && !strings.HasPrefix(contentType, "multipart/form-data") && !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// TODO What should we do here ? upload is empty so we shouldn't
		// create a torrent from it
		messages.AddError("errors", "Please provide either of Content-Type: application/json header or multipart/form-data")
	}
	if strings.HasPrefix(contentType, "multipart/form-data") || strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		modelHelper.BindValueForm(&b, r)
	} else {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&b)
		if err != nil {
			messages.ImportFromError("errors", err)
		}
	}
	defer r.Body.Close()

	modelHelper.ValidateForm(&b, messages)
	if !messages.HasErrors() {
		user, _, errorUser := userService.CreateUserAuthenticationAPI(r, &b)
		if errorUser == nil {
			messages.AddInfo("infos", "Logged")
			apiResponseHandler(w, r, user.ToJSON())
			return
		}
		messages.ImportFromError("errors", errorUser)
	}
	apiResponseHandler(w, r)
}

// APIRefreshTokenHandler : Refresh Token with API
func APIRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username := r.FormValue("username")
	user, _, _, _, err := userService.RetrieveUserByAPITokenAndName(token, username)
	defer r.Body.Close()
	messages := msg.GetMessages(r)
	if err != nil {
		messages.AddError("errors", "Error API token doesn't exist")
	}
	if !messages.HasErrors() {
		user.APIToken, _ = crypto.GenerateRandomToken32()
		user.APITokenExpiry = time.Unix(0, 0)
		_, errorUser := userService.UpdateRawUser(user)
		if errorUser == nil {
			messages.AddInfoT("infos", "profile_updated")
			apiResponseHandler(w, r, user.ToJSON())
			return
		}
		messages.ImportFromError("errors", errorUser)
	}
	apiResponseHandler(w, r)
}

// APICheckTokenHandler : Check Token with API
func APICheckTokenHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	username := r.FormValue("username")
	user, _, _, _, err := userService.RetrieveUserByAPITokenAndName(token, username)
	defer r.Body.Close()
	messages := msg.GetMessages(r)
	if err != nil {
		messages.AddError("errors", "Error API token doesn't exist")
	} else {
		messages.AddInfo("infos", "Logged")
	}
	apiResponseHandler(w, r, user.ToJSON())
}

// This function is the global response for every simple Post Request API
// Please use it. Responses are of the type:
// {ok: bool, [errors | infos]: ArrayOfString [, data: ArrayOfObjects, all_errors: ArrayOfObjects]}
// To send errors or infos, you just need to use the Messages Util
func apiResponseHandler(w http.ResponseWriter, r *http.Request, obj ...interface{}) {
	messages := msg.GetMessages(r)
	var apiJSON []byte
	w.Header().Set("Content-Type", "application/json")

	if !messages.HasErrors() {
		mapOk := map[string]interface{}{"ok": true, "infos": messages.GetInfos("infos")}
		if len(obj) > 0 {
			mapOk["data"] = obj
			if len(obj) == 1 {
				mapOk["data"] = obj[0]
			}
		}
		apiJSON, _ = json.Marshal(mapOk)
	} else { // We need to show error messages
		mapNotOk := map[string]interface{}{"ok": false, "errors": messages.GetErrors("errors"), "all_errors": messages.GetAllErrors()}
		if len(obj) > 0 {
			mapNotOk["data"] = obj
			if len(obj) == 1 {
				mapNotOk["data"] = obj[0]
			}
		}
		if len(messages.GetAllErrors()) > 0 && len(messages.GetErrors("errors")) == 0 {
			mapNotOk["errors"] = "errors"
		}
		apiJSON, _ = json.Marshal(mapNotOk)
	}

	w.Write(apiJSON)
}
