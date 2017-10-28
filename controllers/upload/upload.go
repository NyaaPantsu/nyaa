package uploadController

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/cache"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/gin-gonic/gin"
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
	}

	err = torrents.ExistOrDelete(uploadForm.Infohash, user)
	if err != nil {
		messages.AddError("errors", err.Error())
	}

	AnidexUpload := false
	NyaaSiUpload := false
	TokyoToshoUpload := false

	if c.PostForm("anidex_api") != "" || c.PostForm("anidex_upload") == "true" {
		AnidexUpload = true
	}
	if c.PostForm("nyaasi_api") != "" || c.PostForm("nyaasi_upload") == "true" {
		NyaaSiUpload = true
	}
	if c.PostForm("tokyot_api") != "" || c.PostForm("tokyot_upload") == "true" {
		TokyoToshoUpload = true
	}

	if !messages.HasErrors() {
		// add to db
		torrent, err := torrents.Create(user, &uploadForm)
		log.CheckErrorWithMessage(err, "ERROR_TORRENT_CREATED: Error while creating entry in db")

		if AnidexUpload || NyaaSiUpload || TokyoToshoUpload {
			// User wants to upload to other websites too
			if AnidexUpload {
				go upload.ToAnidex(c, torrent)
			}

			if NyaaSiUpload {
				go upload.ToNyaasi(c, torrent)
			}

			if TokyoToshoUpload {
				go upload.ToTTosho(c, torrent)
			}
			// After that, we redirect to the page for upload status
			url := fmt.Sprintf("/upload/status/%d", torrent.ID)
			c.Redirect(302, url)
		} else {
			url := "/view/" + strconv.FormatUint(uint64(torrent.ID), 10)
			c.Redirect(302, url+"?success")
		}
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

// multiUploadStatus : controller to show the multi upload status
func multiUploadStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	if found, ok := cache.C.Get("tstatus_" + id); ok {
		uploadMultiple := found.(upload.MultipleForm)
		// if ?json we print the json format
		if _, ok = c.GetQuery("json"); ok {
			c.JSON(http.StatusFound, uploadMultiple)
			return
		}
		// else we send the upload multiple form (support of manual F5)
		variables := templates.Commonvariables(c)
		variables.Set("UploadMultiple", uploadMultiple)
		templates.Render(c, "site/torrents/upload_multiple.jet.html", variables)
	} else {
		// here it means the upload status is already flushed from memory
		c.AbortWithStatus(http.StatusNotFound)
	}
}
