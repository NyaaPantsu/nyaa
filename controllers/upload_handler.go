package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/service/api"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/gin-gonic/gin"
)

// UploadHandler : Main Controller for uploading a torrent
func UploadHandler(c *gin.Context) {
	user := getUser(c)
	if !uploadService.IsUploadEnabled(user) {
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
	var uploadForm apiService.TorrentRequest
	user := getUser(c)
	messages := msg.GetMessages(c) // new utils for errors and infos

	if userPermission.NeedsCaptcha(user) {
		userCaptcha := captcha.Extract(c)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddError("errors", captcha.ErrInvalidCaptcha.Error())
		}
	}

	// validation is done in ExtractInfo()
	err := uploadForm.ExtractInfo(c)
	if err != nil {
		messages.AddError("errors", err.Error())
	}
	status := models.TorrentStatusNormal
	if uploadForm.Remake { // overrides trusted
		status = models.TorrentStatusRemake
	} else if user.IsTrusted() {
		status = models.TorrentStatusTrusted
	}

	err = torrentService.ExistOrDelete(uploadForm.Infohash, user)
	if err != nil {
		messages.AddError("errors", err.Error())
	}

	if !messages.HasErrors() {
		// add to db and redirect
		torrent := models.Torrent{
			Name:        uploadForm.Name,
			Category:    uploadForm.CategoryID,
			SubCategory: uploadForm.SubCategoryID,
			Status:      status,
			Hidden:      uploadForm.Hidden,
			Hash:        uploadForm.Infohash,
			Date:        time.Now(),
			Filesize:    uploadForm.Filesize,
			Description: uploadForm.Description,
			WebsiteLink: uploadForm.WebsiteLink,
			UploaderID:  user.ID,
			Language:    uploadForm.Language}
		torrent.ParseTrackers(uploadForm.Trackers)
		db.ORM.Create(&torrent)

		if db.ElasticSearchClient != nil {
			err := torrent.AddToESIndex(db.ElasticSearchClient)
			if err == nil {
				log.Infof("Successfully added torrent to ES index.")
			} else {
				log.Errorf("Unable to add torrent to ES index: %s", err)
			}
		}

		torrentService.NewTorrentEvent(user, &torrent)

		// add filelist to files db, if we have one
		if len(uploadForm.FileList) > 0 {
			for _, uploadedFile := range uploadForm.FileList {
				file := models.File{TorrentID: torrent.ID, Filesize: uploadedFile.Filesize}
				err := file.SetPath(uploadedFile.Path)
				if err != nil {
					messages.AddError("errors", err.Error())
				}
				db.ORM.Create(&file)
			}
		}

		url := "/view/" + strconv.FormatUint(uint64(torrent.ID), 10)
		c.Redirect(302, url+"?success")
	}
}

// UploadGetHandler : Controller for uploading a torrent, after GET request or Failed Post request
func UploadGetHandler(c *gin.Context) {
	var uploadForm apiService.TorrentRequest
	_ = uploadForm.ExtractInfo(c)
	user := getUser(c)
	if userPermission.NeedsCaptcha(user) {
		uploadForm.CaptchaID = captcha.GetID()
	} else {
		uploadForm.CaptchaID = ""
	}
	formTemplate(c, "site/torrents/upload.jet.html", uploadForm)
}
