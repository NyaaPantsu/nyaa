package uploadController

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/gin-gonic/gin"
	"github.com/NyaaPantsu/nyaa/utils/log"
)

// UploadHandler : Main Controller for uploading a torrent
func UploadHandler(c *gin.Context) {
	user := router.GetUser(c)
	if !user.CanUpload() {
		T := publicSettings.GetTfuncFromRequest(c)
		c.AbortWithError(http.StatusBadRequest, errors.New(string(T("uploads_disabled"))))
		return
	}

	if c.Request.Method == "POST" {
		UploadPostHandler(c)
	}

	UploadGetHandler(c)
}

// UploadPostHandler : Controller for uploading a torrent, after POST request, redirect or makes error in messages
func UploadPostHandler(c *gin.Context) {
	var uploadForm torrentValidator.TorrentRequest
	user := router.GetUser(c)
	messages := msg.GetMessages(c) // new utils for errors and infos

	if user.NeedsCaptcha() {
		userCaptcha := captcha.Extract(c)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddError("errors", captcha.ErrInvalidCaptcha.Error())
		}
	}

	// validation is done in ExtractInfo()
	err := upload.ExtractInfo(c, &uploadForm)
	if err != nil {
		messages.AddError("errors", err.Error())
	}

	uploadForm.Status = models.TorrentStatusNormal
	if uploadForm.Remake { // overrides trusted
		uploadForm.Status = models.TorrentStatusRemake
	} else if user.IsTrusted() {
		uploadForm.Status = models.TorrentStatusTrusted
	} else if user.IsBanned() {
		uploadForm.Status = models.TorrentStatusBlocked
	}

	err = torrents.ExistOrDelete(uploadForm.Infohash, user)
	if err != nil {
		messages.AddError("errors", err.Error())
	}

	if !messages.HasErrors() {
		// add to db and redirect
		torrent, err := torrents.Create(user, &uploadForm)
		log.CheckErrorWithMessage(err, "ERROR_TORRENT_CREATED: Error while creating entry in db")
		url := "/view/" + strconv.FormatUint(uint64(torrent.ID), 10)
		c.Redirect(302, url+"?success")
	}
}

// UploadGetHandler : Controller for uploading a torrent, after GET request or Failed Post request
func UploadGetHandler(c *gin.Context) {
	var uploadForm torrentValidator.TorrentRequest
	_ = upload.ExtractInfo(c, &uploadForm)
	user := router.GetUser(c)
	if user.NeedsCaptcha() {
		uploadForm.CaptchaID = captcha.GetID()
	} else {
		uploadForm.CaptchaID = ""
	}
	templates.Form(c, "site/torrents/upload.jet.html", uploadForm)
}
