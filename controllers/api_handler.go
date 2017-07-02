package controllers

import (
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/crypto"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
)

// APIHandler : Controller for api request on torrent list
func APIHandler(c *gin.Context) {
	t := c.Query("t")
	if t != "" {
		RSSTorznabHandler(c)
	} else {
		page := c.Param("page")
		whereParams := structs.WhereParams{}
		req := upload.TorrentsRequest{}

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

		torrents, nbTorrents, err := torrents.Find(whereParams, req.MaxPerPage, req.MaxPerPage*(req.Page-1))
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		b := upload.APIResultJSON{
			Torrents: models.APITorrentsToJSON(torrents),
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
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	torrent, err := torrents.FindByID(uint(id))

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

	_, err = torrents.FindRawByID(uint(id))

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
	user, _, _, _, err := users.FindByAPITokenAndName(token, username)
	messages := msg.GetMessages(c)

	if err != nil {
		messages.AddErrorT("errors", "error_api_token")
	}

	if !user.CanUpload() {
		messages.AddErrorT("errors", "uploads_disabled")
	}

	if user.ID == 0 {
		messages.AddErrorT("errors", "error_api_token")
	}

	if !messages.HasErrors() {
		uploadForm := torrentValidator.TorrentRequest{}
		contentType := c.Request.Header.Get("Content-Type")
		if contentType != "application/json" && !strings.HasPrefix(contentType, "multipart/form-data") && contentType != "application/x-www-form-urlencoded" {
			// TODO What should we do here ? uploadForm is empty so we shouldn't
			// create a torrent from it
			messages.AddErrorT("errors", "error_content_type_post")
		}
		// As long as the right content-type is sent, formValue is smart enough to parse it
		err = upload.ExtractInfo(c, &uploadForm)
		if err != nil {
			messages.Error(err)
		}

		if !messages.HasErrors() {
			uploadForm.Status = models.TorrentStatusNormal
			if uploadForm.Remake { // overrides trusted
				uploadForm.Status = models.TorrentStatusRemake
			} else if user.IsTrusted() {
				uploadForm.Status = models.TorrentStatusTrusted
			}
			err = torrents.ExistOrDelete(uploadForm.Infohash, user)
			if err != nil {
				messages.Error(err)
			}
			if !messages.HasErrors() {
				torrent, err := torrents.Create(user, &uploadForm)
				if err != nil {
					messages.Error(err)
				}
				messages.AddInfoT("infos", "torrent_uploaded")

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
	user, _, _, _, err := users.FindByAPITokenAndName(token, username)

	messages := msg.GetMessages(c)

	if err != nil {
		messages.AddErrorT("errors", "error_api_token")
	}

	if !user.CanUpload() {
		messages.AddErrorT("errors", "uploads_disabled")
	}

	if user.ID == 0 {
		messages.AddErrorT("errors", "error_api_token")
	}
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" && !strings.HasPrefix(contentType, "multipart/form-data") && contentType != "application/x-www-form-urlencoded" {
		// TODO What should we do here ? upload is empty so we shouldn't
		// create a torrent from it
		messages.AddErrorT("errors", "error_content_type_post")
	}
	update := torrentValidator.UpdateRequest{}
	err = upload.ExtractEditInfo(c, &update.Update)
	if err != nil {
		messages.Error(err)
	}
	if !messages.HasErrors() {
		c.Bind(&update)
		torrent, err := torrents.FindByID(update.ID)
		if err != nil {
			messages.AddErrorTf("errors", "torrent_not_exist", strconv.Itoa(int(update.ID)))
		}
		if torrent.UploaderID != 0 && torrent.UploaderID != user.ID { //&& user.Status != mod
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		upload.UpdateTorrent(&update, &torrent, user).Update(false)
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

	b := models.APITorrentsToJSON(torrents)
	c.JSON(http.StatusOK, b)
}

// APILoginHandler : Login with API
// This is not an OAuth api like and shouldn't be used for anything except getting the API Token in order to not store a password
func APILoginHandler(c *gin.Context) {
	b := userValidator.LoginForm{}
	messages := msg.GetMessages(c)
	contentType := c.Request.Header.Get("Content-type")
	if !strings.HasPrefix(contentType, "application/json") && !strings.HasPrefix(contentType, "multipart/form-data") && !strings.HasPrefix(contentType, "application/x-www-form-urlencoded") {
		// TODO What should we do here ? upload is empty so we shouldn't
		// create a torrent from it
		messages.AddErrorT("errors", "error_content_type_post")
	}
	c.Bind(&b)
	validator.ValidateForm(&b, messages)
	if !messages.HasErrors() {
		user, _, errorUser := cookies.CreateUserAuthentication(c, &b)
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
	user, _, _, _, err := users.FindByAPITokenAndName(token, username)

	messages := msg.GetMessages(c)
	if err != nil {
		messages.AddErrorT("errors", "error_api_token")
	}
	if !messages.HasErrors() {
		user.APIToken, _ = crypto.GenerateRandomToken32()
		user.APITokenExpiry = time.Unix(0, 0)
		_, errorUser := user.UpdateRaw()
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
	user, _, _, _, err := users.FindByAPITokenAndName(token, username)

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
