package router

import (
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
	"github.com/gin-gonic/gin"
)

// APIHandler : Controller for api request on torrent list
func APIHandler(c *gin.Context) {
	t := c.Query("t")
	if t != "" {
		RSSTorznabHandler(c)
	} else {
		page := c.Param("page")
		whereParams := serviceBase.WhereParams{}
		req := apiService.TorrentsRequest{}

		contentType := c.Request.Header.Get("Content-Type")
		if contentType == "application/json" {
			c.Bind(&req)

			if req.MaxPerPage == 0 {
				req.MaxPerPage = config.Conf.Navigation.TorrentsPerPage
			}
			if req.Page <= 0 {
				req.Page = 1
			}

			whereParams = req.ToParams()
		} else {
			var err error
			maxString := c.Query("max")
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
					c.AbortWithError(http.StatusInternalServerError, err)
					return
				}
				if req.Page <= 0 {
					NotFoundHandler(c)
					return
				}
			}
		}

		torrents, nbTorrents, err := torrentService.GetTorrents(whereParams, req.MaxPerPage, req.MaxPerPage*(req.Page-1))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		b := model.APIResultJSON{
			Torrents: model.APITorrentsToJSON(torrents),
		}
		b.QueryRecordCount = req.MaxPerPage
		b.TotalRecordCount = nbTorrents
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusOK, b)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

// APIViewHandler : Controller for viewing a torrent by its ID
func APIViewHandler(c *gin.Context) {
	id := c.Param("id")

	torrent, err := torrentService.GetTorrentByID(id)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	b := torrent.ToJSON()
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, b)
}

// APIViewHeadHandler : Controller for checking a torrent by its ID
func APIViewHeadHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		return
	}

	_, err = torrentService.GetRawTorrentByID(uint(id))

	if err != nil {
		NotFoundHandler(c)
		return
	}

	c.Writer.Write(nil)
}

// APIUploadHandler : Controller for uploading a torrent with api
func APIUploadHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	username := c.PostForm("username")
	user, _, _, _, err := userService.RetrieveUserByAPITokenAndName(token, username)
	messages := msg.GetMessages(c)

	if err != nil {
		messages.AddErrorT("errors", "error_api_token")
	}

	if !uploadService.IsUploadEnabled(user) {
		messages.AddErrorT("errors", "uploads_disabled")
	}

	if user.ID == 0 {
		messages.Error(apiService.ErrAPIKey)
	}

	if !messages.HasErrors() {
		upload := apiService.TorrentRequest{}
		contentType := c.Request.Header.Get("Content-Type")
		if contentType != "application/json" && !strings.HasPrefix(contentType, "multipart/form-data") && contentType != "application/x-www-form-urlencoded" {
			// TODO What should we do here ? upload is empty so we shouldn't
			// create a torrent from it
			messages.AddErrorT("errors", "error_content_type_post")
		}
		// As long as the right content-type is sent, formValue is smart enough to parse it
		err = upload.ExtractInfo(c)
		if err != nil {
			messages.Error(err)
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
				messages.Error(err)
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
				torrentService.NewTorrentEvent(user, &torrent)
				// add filelist to files db, if we have one
				if len(upload.FileList) > 0 {
					for _, uploadedFile := range upload.FileList {
						file := model.File{TorrentID: torrent.ID, Filesize: upload.Filesize}
						err := file.SetPath(uploadedFile.Path)
						if err != nil {
							c.AbortWithError(http.StatusBadRequest, err)
							return
						}
						db.ORM.Create(&file)
					}
				}
				apiResponseHandler(c, torrent.ToJSON())
				return
			}
		}
	}
	apiResponseHandler(c)
}

// APIUpdateHandler : Controller for updating a torrent with api
// FIXME Impossible to update a torrent uploaded by user 0
func APIUpdateHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	username := c.PostForm("username")
	user, _, _, _, err := userService.RetrieveUserByAPITokenAndName(token, username)

	messages := msg.GetMessages(c)

	if err != nil {
		messages.AddErrorT("errors", "error_api_token")
	}

	if !uploadService.IsUploadEnabled(user) {
		messages.AddErrorT("errors", "uploads_disabled")
	}

	if user.ID == 0 {
		messages.Error(apiService.ErrAPIKey)
	}
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" && !strings.HasPrefix(contentType, "multipart/form-data") && contentType != "application/x-www-form-urlencoded" {
		// TODO What should we do here ? upload is empty so we shouldn't
		// create a torrent from it
		messages.AddErrorT("errors", "error_content_type_post")
	}
	update := apiService.UpdateRequest{}
	err = update.Update.ExtractEditInfo(c)
	if err != nil {
		messages.Error(err)
	}
	if !messages.HasErrors() {
		torrent := model.Torrent{}
		db.ORM.Where("torrent_id = ?", c.PostForm("id")).First(&torrent)
		if torrent.ID == 0 {
			messages.Error(apiService.ErrTorrentID)
		}
		if torrent.UploaderID != 0 && torrent.UploaderID != user.ID { //&& user.Status != mod
			messages.Error(apiService.ErrRights)
		}
		update.UpdateTorrent(&torrent, user)
		torrentService.UpdateTorrent(&torrent)
	}
	apiResponseHandler(c)
}

// APISearchHandler : Controller for searching with api
func APISearchHandler(c *gin.Context) {
	page := c.Param("page")

	// db params url
	var err error
	pagenum := 1
	if page != "" {
		pagenum, err = strconv.Atoi(html.EscapeString(page))
		if !log.CheckError(err) {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if pagenum <= 0 {
			NotFoundHandler(c)
			return
		}
	}

	_, torrents, _, err := search.SearchByQueryWithUser(c, pagenum)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	b := model.APITorrentsToJSON(torrents)
	c.JSON(http.StatusOK, b)
}

// APILoginHandler : Login with API
// This is not an OAuth api like and shouldn't be used for anything except getting the API Token in order to not store a password
func APILoginHandler(c *gin.Context) {
	b := form.LoginForm{}
	messages := msg.GetMessages(c)
	contentType := c.Request.Header.Get("Content-type")
	if !strings.HasPrefix(contentType, "application/json") && !strings.HasPrefix(contentType, "multipart/form-data") && !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// TODO What should we do here ? upload is empty so we shouldn't
		// create a torrent from it
		messages.AddErrorT("errors", "error_content_type_post")
	}
	c.Bind(&b)
	modelHelper.ValidateForm(&b, messages)
	if !messages.HasErrors() {
		user, _, errorUser := userService.CreateUserAuthenticationAPI(c, &b)
		if errorUser == nil {
			messages.AddInfo("infos", "Logged")
			apiResponseHandler(c, user.ToJSON())
			return
		}
		messages.Error(errorUser)
	}
	apiResponseHandler(c)
}

// APIRefreshTokenHandler : Refresh Token with API
func APIRefreshTokenHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	username := c.PostForm("username")
	user, _, _, _, err := userService.RetrieveUserByAPITokenAndName(token, username)

	messages := msg.GetMessages(c)
	if err != nil {
		messages.AddErrorT("errors", "error_api_token")
	}
	if !messages.HasErrors() {
		user.APIToken, _ = crypto.GenerateRandomToken32()
		user.APITokenExpiry = time.Unix(0, 0)
		_, errorUser := userService.UpdateRawUser(user)
		if errorUser == nil {
			messages.AddInfoT("infos", "profile_updated")
			apiResponseHandler(c, user.ToJSON())
			return
		}
		messages.Error(errorUser)
	}
	apiResponseHandler(c)
}

// APICheckTokenHandler : Check Token with API
func APICheckTokenHandler(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	username := c.PostForm("username")
	user, _, _, _, err := userService.RetrieveUserByAPITokenAndName(token, username)

	messages := msg.GetMessages(c)
	if err != nil {
		messages.AddErrorT("errors", "error_api_token")
	} else {
		messages.AddInfo("infos", "Logged")
	}
	apiResponseHandler(c, user.ToJSON())
}

// This function is the global response for every simple Post Request API
// Please use it. Responses are of the type:
// {ok: bool, [errors | infos]: ArrayOfString [, data: ArrayOfObjects, all_errors: ArrayOfObjects]}
// To send errors or infos, you just need to use the Messages Util
func apiResponseHandler(c *gin.Context, obj ...interface{}) {
	messages := msg.GetMessages(c)
	c.Header("Content-Type", "application/json")

	var mapOk map[string]interface{}
	if !messages.HasErrors() {
		mapOk = map[string]interface{}{"ok": true, "infos": messages.GetInfos("infos")}
		if len(obj) > 0 {
			mapOk["data"] = obj
			if len(obj) == 1 {
				mapOk["data"] = obj[0]
			}
		}
	} else { // We need to show error messages
		mapOk := map[string]interface{}{"ok": false, "errors": messages.GetErrors("errors"), "all_errors": messages.GetAllErrors()}
		if len(obj) > 0 {
			mapOk["data"] = obj
			if len(obj) == 1 {
				mapOk["data"] = obj[0]
			}
		}
		if len(messages.GetAllErrors()) > 0 && len(messages.GetErrors("errors")) == 0 {
			mapOk["errors"] = "errors"
		}
	}

	c.JSON(http.StatusOK, mapOk)
}
