package apiController

import (
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/controllers/feed"
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/models/users"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/crypto"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/search"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/NyaaPantsu/nyaa/utils/validator/user"
	"github.com/gin-gonic/gin"
)

/**
 * @apiDefine NotFoundError
 * @apiVersion 1.1.0
 * @apiError {String[]} errors List of errors messages with a 404 error message in it.
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 404 Not Found
 *     {
 *       "errors": [ "404_not_found", ... ]
 *     }
 */
/**
 * @apiDefine RequestError
 * @apiVersion 1.1.0
 * @apiError {Boolean} ok The request couldn't be done due to some errors.
 * @apiError {String[]} errors List of errors messages.
 * @apiError {Object[]} all_errors List of errors object messages for each wrong field
 *
 * @apiErrorExample Error-Response:
 *     HTTP/1.1 200 OK
 *     {
 *       "ok": false,
 *       "errors": [ ... ]
 *       "all_errors": {
 * 		 	"username": [ ... ],
 *        }
 *     }
 */

/**
 * @api {get} / Request Torrents index
 * @apiVersion 1.1.0
 * @apiName GetTorrents
 * @apiGroup Torrents
 *
 * @apiParam {Number} id Torrent unique ID.
 *
 * @apiSuccess {Object[]} torrents List of torrent object (see view for the properties).
 * @apiSuccess {Number} queryRecordCount Number of torrents given.
 * @apiSuccess {Number} totalRecordCount Total number of torrents.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *			"torrents": [...],
 *			"queryRecordCount": 50,
 *			"totalRecordCount": 798414
 *		}
 *
 * @apiUse NotFoundError
 */
// APIHandler : Controller for api request on torrent list
func APIHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	t := c.Query("t")
	if t != "" {
		feedController.RSSTorznabHandler(c)
	} else {
		APISearchHandler(c)
	}
}

/**
 * @api {get} /view/:id Request Torrent information
 * @apiVersion 1.1.0
 * @apiName GetTorrent
 * @apiGroup Torrents
 *
 * @apiParam {Number} id Torrent unique ID.
 *
 * @apiSuccess {Number} id ID of the torrent.
 * @apiSuccess {String} name Name of the torrent.
 * @apiSuccess {Number} status Status of the torrent.
 * @apiSuccess {String} hash Hash of the torrent.
 * @apiSuccess {Date} date Uploaded date of the torrent.
 * @apiSuccess {Number} filesize File size in Bytes of the torrent.
 * @apiSuccess {String} description Description of the torrent.
 * @apiSuccess {Object[]} comments Comments of the torrent.
 * @apiSuccess {String} sub_category Sub Category of the torrent.
 * @apiSuccess {String} category Category of the torrent.
 * @apiSuccess {String} anidb_id Anidb ID of the torrent.
 * @apiSuccess {Number} uploader_id ID of the torrent uploader.
 * @apiSuccess {String} uploader_name  Username of the torrent uploader.
 * @apiSuccess {String} uploader_old  Old username from nyaa of the torrent uploader.
 * @apiSuccess {String} website_link  External Link of the torrent.
 * @apiSuccess {String[]} languages  Languages of the torrent.
 * @apiSuccess {String} magnet  Magnet URI of the torrent.
 * @apiSuccess {String} torrent  Download URL of the torrent.
 * @apiSuccess {Number} seeders  Number of seeders of the torrent.
 * @apiSuccess {Number} leechers  Number of leechers of the torrent.
 * @apiSuccess {Number} completed  Downloads completed of the torrent.
 * @apiSuccess {Date} last_scrape  Last statistics update of the torrent.
 * @apiSuccess {Object[]} file_list  List of files in the torrent.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *	{
 *	"id": 952801,
 *	"name": "[HorribleSubs] Uchouten Kazoku S2 [720p]",
 *	"status": 1,
 *	"hash": "6E4D96F7A0B0456672E80B150CCB7C15868CD47D",
 *	"date": "2017-07-05T11:01:39Z",
 *	"filesize": 4056160259,
 *	"description": "<p>Unofficial batch</p>\n",
 *	"comments": [],
 *	"sub_category": "5",
 *	"category": "3",
 *	"anidb_id": "",
 *	"downloads": 0,
 *	"uploader_id": 7177,
 *	"uploader_name": "DarAR92",
 *	"uploader_old": "",
 *	"website_link": "http://horriblesubs.info/",
 *	"languages": [
 *	"en-us"
 *	],
 *	"magnet": "magnet:?xt=urn:btih:6E4D96F7A0B0456672E80B150CCB7C15868CD47D&dn=%5BHorribleSubs%5D+Uchouten+Kazoku+S2+%5B720p%5D&tr=http://nyaa.tracker.wf:7777/announce&tr=http://nyaa.tracker.wf:7777/announce&tr=udp://tracker.doko.moe:6969&tr=http://tracker.anirena.com:80/announce&tr=http://anidex.moe:6969/announce&tr=udp://tracker.opentrackr.org:1337&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://zer0day.ch:1337&tr=udp://9.rarbg.com:2710/announce&tr=udp://tracker2.christianbro.pw:6969/announce&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.leechers-paradise.org:6969&tr=udp://eddie4.nl:6969/announce&tr=udp://tracker.doko.moe:6969/announce",
 *	"torrent": "https://nyaa.pantsu.cat/download/6E4D96F7A0B0456672E80B150CCB7C15868CD47D",
 *	"seeders": 4,
 *	"leechers": 2,
 *	"completed": 28,
 *	"last_scrape": "2017-07-07T07:48:32.509635Z",
 *	"file_list": [
 *	{
 *	"path": "[HorribleSubs] Uchouten Kazoku S2 - 01[720p].mkv",
 *	"filesize": 338250895
 *	},
 *	{
 *	"path": "[HorribleSubs] Uchouten Kazoku S2 - 12 [720p].mkv",
 *	"filesize": 338556275
 *	}
 *	]
 *	}
 *
 * @apiUse NotFoundError
 */
// APIViewHandler : Controller for viewing a torrent by its ID
func APIViewHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	torrent, err := torrents.FindByID(uint(id))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	b := torrent.ToJSON()
	c.JSON(http.StatusOK, b)
}

/**
 * @api {get} /head/:id Request Torrent Head
 * @apiVersion 1.1.0
 * @apiName GetTorrentHead
 * @apiGroup Torrents
 *
 * @apiParam {Number} id Torrent unique ID.
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse NotFoundError
 */
// APIViewHeadHandler : Controller for checking a torrent by its ID
func APIViewHeadHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)

	if err != nil {
		return
	}

	_, err = torrents.FindRawByID(uint(id))

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.Writer.Write(nil)
}

/**
 * @api {post} /upload Upload a Torrent
 * @apiVersion 1.1.0
 * @apiName UploadTorrent
 * @apiGroup Torrents
 *
 * @apiParam {String} username Torrent uploader name.
 * @apiParam {String} name Torrent name.
 * @apiParam {String} magnet Torrent magnet URI.
 * @apiParam {String} category Torrent category.
 * @apiParam {Boolean} remake Torrent is a remake.
 * @apiParam {String} description Torrent description.
 * @apiParam {Number} status Torrent status.
 * @apiParam {Boolean} hidden Torrent hidden.
 * @apiParam {String} website_link Torrent website link.
 * @apiParam {String[]} languages Torrent languages.
 * @apiParam {File} torrent Torrent file to upload (you have to send a torrent file or a magnet, not both!).
 *
 * @apiSuccess {Boolean} ok The request is done without failing
 * @apiSuccess {String[]} infos Messages information relative to the request
 * @apiSuccess {Object} data The resulting torrent uploaded (see view for the properties)
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse RequestError
 */
// APIUploadHandler : Controller for uploading a torrent with api
func APIUploadHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
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
		uploadForm := upload.NewTorrentRequest()
		contentType := c.Request.Header.Get("Content-Type")
		if contentType != "application/json" && !strings.HasPrefix(contentType, "multipart/form-data") && contentType != "application/x-www-form-urlencoded" {
			// TODO What should we do here ? uploadForm is empty so we shouldn't
			// create a torrent from it
			messages.AddErrorT("errors", "error_content_type_post")
		}
		// As long as the right content-type is sent, formValue is smart enough to parse it
		err = upload.ExtractInfo(c, uploadForm)
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
				torrent, err := torrents.Create(user, uploadForm)
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

/**
 * @api {post} /update/ Update a Torrent
 * @apiVersion 1.1.0
 * @apiName UpdateTorrent
 * @apiGroup Torrents
 *
 * @apiParam {String} username Torrent uploader name.
 * @apiParam {Number} id Torrent ID.
 * @apiParam {String} name Torrent name.
 * @apiParam {String} category Torrent category.
 * @apiParam {Boolean} remake Torrent is a remake.
 * @apiParam {String} description Torrent description.
 * @apiParam {Number} status Torrent status.
 * @apiParam {Boolean} hidden Torrent hidden.
 * @apiParam {String} website_link Torrent website link.
 * @apiParam {String[]} languages Torrent languages.
 *
 * @apiSuccess {Boolean} ok The request is done without failing
 * @apiSuccess {String[]} infos Messages information relative to the request
 * @apiSuccess {Object} data The resulting torrent updated (see view for the properties)
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *
 * @apiUse RequestError
 */
// APIUpdateHandler : Controller for updating a torrent with api
func APIUpdateHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
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
		upload.UpdateTorrent(&update, torrent, user).Update(false)
	}
	apiResponseHandler(c)
}

/**
 * @api {get} /search/ Search Torrents
 * @apiVersion 1.1.0
 * @apiName FindTorrents
 * @apiGroup Torrents
 *
 * @apiParam {String[]} c In which categories to search.
 * @apiParam {String} q Query to search (torrent name).
 * @apiParam {String} limit Number of results per page.
 * @apiParam {String} userID Uploader ID owning the torrents.
 * @apiParam {String} fromID Show results with torrents ID superior to this.
 * @apiParam {String} s Torrent status.
 * @apiParam {String} maxage Torrents which have been uploaded the last x days.
 * @apiParam {String} toDate Torrents which have been uploaded since x <code>dateType</code>.
 * @apiParam {String} fromDate Torrents which have been uploaded the last x <code>dateType</code>.
 * @apiParam {String} dateType Which type of date (<code>d</code> for days, <code>m</code> for months, <code>y</code> for years).
 * @apiParam {String} minSize Filter by minimal size in <code>sizeType</code>.
 * @apiParam {String} maxSize Filter by maximal size in <code>sizeType</code>.
 * @apiParam {String} sizeType Which type of size (<code>b</code> for bytes, <code>k</code> for kilobytes, <code>m</code> for megabytes, <code>g</code> for gigabytes).
 * @apiParam {String} sort Torrent sorting type (0 = id, 1 = name, 2 = date, 3 = downloads, 4 = size, 5 = seeders, 6 = leechers, 7 = completed).
 * @apiParam {Boolean} order Order ascending or descending (true = ascending).
 * @apiParam {String[]} lang Filter the languages.
 * @apiParam {Number} page Search page.
 *
 * @apiSuccess {Object[]} torrents List of torrent object (see view for the properties).
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 *     {
 *			"torrents": [...],
 *			"queryRecordCount": 50,
 *			"totalRecordCount": 798414
 *		}
 *
 * @apiUse NotFoundError
 */
// APISearchHandler : Controller for searching with api
func APISearchHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	page := c.Param("page")
	currentUser := router.GetUser(c)

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
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
	}

	userID, err := strconv.ParseUint(c.Query("userID"), 10, 32)
	if err != nil {
		userID = 0
	}

	_, torrentSearch, nbTorrents, err := search.AuthorizedQuery(c, pagenum, currentUser.CurrentOrAdmin(uint(userID)))

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	maxQuery, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(config.Get().Navigation.TorrentsPerPage)))
	if err != nil {
		maxQuery = config.Get().Navigation.TorrentsPerPage
	} else if maxQuery > config.Get().Navigation.MaxTorrentsPerPage {
		maxQuery = config.Get().Navigation.MaxTorrentsPerPage
	}

	b := upload.APIResultJSON{
		TotalRecordCount: nbTorrents,
		Torrents:         torrents.APITorrentsToJSON(torrentSearch),
		QueryRecordCount: maxQuery,
	}
	c.JSON(http.StatusOK, b)
}

/**
 * @api {post} /login/ Login a user
 * @apiVersion 1.1.0
 * @apiName Login
 * @apiGroup Users
 *
 * @apiParam {String} username Username or Email.
 * @apiParam {String} password Password.
 *
 * @apiSuccess {Boolean} ok The request is done without failing
 * @apiSuccess {String[]} infos Messages information relative to the request
 * @apiSuccess {Object} data The connected user object
 *
 * @apiSuccessExample Success-Response:
 *     HTTP/1.1 200 OK
 * 		{
 * 			data:
 *       		[{
 * 					user_id:1,
 *					username:"username",
 *					status:1,
 *					token:"token",
 *					md5:"",
 *					created_at:"date",
 *					liking_count:0,
 *					liked_count:0
 *				}],
 *			infos: ["Logged", ... ],
 *			ok:true
 * 		}
 *
 * @apiUse RequestError
 */
// APILoginHandler : Login with API
// This is not an OAuth api like and shouldn't be used for anything except getting the API Token in order to not store a password
func APILoginHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")
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
	c.Header("Content-Type", "application/json")
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
	c.Header("Content-Type", "application/json")
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
