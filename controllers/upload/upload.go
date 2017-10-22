package uploadController

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"io/ioutil"
	"bytes"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/NyaaPantsu/nyaa/config"
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
	if c.PostForm("nyaasi_api") != "" || c.PostForm("nyaasi_upload") == "true"{
		NyaaSiUpload = true
	}
	if c.PostForm("tokyot_api") != "" || c.PostForm("tokyot_upload") == "true" {
		TokyoToshoUpload = true
	}

	if !messages.HasErrors() {
		// add to db and redirect
		torrent, err := torrents.Create(user, &uploadForm)
		log.CheckErrorWithMessage(err, "ERROR_TORRENT_CREATED: Error while creating entry in db")
		
		if AnidexUpload || NyaaSiUpload || TokyoToshoUpload {
			//User wants to upload to other websites too
			uploadMultiple := templates.NewUploadMultipleForm()
			uploadMultiple.PantsuID = torrent.ID
			
			if AnidexUpload {
				uploadMultiple.AnidexStatus = 1
	
				anonymous := false
				apiKey := c.PostForm("anidex_api")
				
				if apiKey == "" {
					anonymous = true
					apiKey = config.Get().Upload.DefaultAnidexToken
				}
				
				postForm := url.Values{}
				//Required
				postForm.Set("api_key", apiKey)
				postForm.Set("subcat_id", c.PostForm("anidex_form_category"))
				postForm.Set("file", "")
				postForm.Set("group_id", "0")
				postForm.Set("lang_id", c.PostForm("anidex_form_lang"))
				
				//Optional
				postForm.Set("description", "")
				if config.IsSukebei() {
					postForm.Set("hentai", "1")
				}
				if uploadForm.Remake {
					postForm.Set("reencode", "1")
				}
				if anonymous {
					postForm.Set("private", "1")
				}
				if c.PostForm("name") != "" {
					postForm.Set("torrent_name", c.PostForm("name"))
				}
					
				
				
				postForm.Set("debug", "1")
				
				rsp, err := http.Post("https://anidex.info/api/", "application/x-www-form-urlencoded", bytes.NewBufferString(postForm.Encode()))
				
				if err != nil {
					uploadMultiple.AnidexStatus = 2
					uploadMultiple.AnidexMessage = "Error during the HTTP POST request"
				}
				defer rsp.Body.Close()
				body_byte, err := ioutil.ReadAll(rsp.Body)
				if err != nil {
					uploadMultiple.AnidexStatus = 2
					uploadMultiple.AnidexMessage = "Unknown error"
				}
				if uploadMultiple.AnidexStatus == 1 {
					uploadMultiple.AnidexMessage = string(body_byte)
					if strings.Contains(uploadMultiple.AnidexMessage, "http://") {
						uploadMultiple.AnidexStatus = 3
					} else if strings.Contains(uploadMultiple.AnidexMessage, "error") {
						uploadMultiple.AnidexStatus = 2
					}
				}
			}
			
			if NyaaSiUpload {
				uploadMultiple.NyaasiStatus = 1
				uploadMultiple.NyaasiMessage = "Sorry u are not allowed"
			}
			
			if TokyoToshoUpload  {
				uploadMultiple.TToshoStatus = 1
			}
			
			variables := templates.Commonvariables(c)
			variables.Set("UploadMultiple", uploadMultiple)
			templates.Render(c, "site/torrents/upload_multiple.jet.html", variables)
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
