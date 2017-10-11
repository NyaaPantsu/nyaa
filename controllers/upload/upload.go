package uploadController

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"io/ioutil"

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
	
	if c.PostForm("anidex_api") != "" || c.PostForm("anidex_upload") != "" {
		AnidexUpload = true
	}

	if !messages.HasErrors() {
		// add to db and redirect
		torrent, err := torrents.Create(user, &uploadForm)
		log.CheckErrorWithMessage(err, "ERROR_TORRENT_CREATED: Error while creating entry in db")
		
		if AnidexUpload || c.PostForm("nyaasi_api") != "" || c.PostForm("tokyot_api") != "" {
			//User wants to upload to other websites too
			uploadMultiple := templates.NewUploadMultipleForm()
			uploadMultiple.PantsuID = torrent.ID
			
			if c.PostForm("anidex_api") != "" || user.AnidexAPIToken != "" {
				uploadMultiple.AnidexStatus = 1
				categoryId := "10"
				langId := "0"
				anonymous := "0"
				apiKey := c.PostForm("anidex_api")
				
				if c.PostForm("anidex_api") == "" {
					anonymous = "1"
				}
				
				postForm := url.Values{}
				//Required
				postForm.Set("api_key", apiKey)
				postForm.Set("subcat_id", categoryId)
				postForm.Set("file", "")
				postForm.Set("torrent_name", uploadForm.Name)
				postForm.Set("group_id", "0")
				postForm.Set("lang_id", langId)
				
				//Optional
				postForm.Set("description", "")
				postForm.Set("hentai", strconv.Itoa(int(config.IsSukebei)))
				postForm.Set("reencode", strconv.Itoa(int(uploadForm.Remake)))
				postForm.Set("private", anonymous)
				
				
				postForm.Set("debug", "1")
				
				body := bytes.NewBufferString(form.Encode())
				rsp, err := http.Post("https://anidex.info/api/", "application/x-www-form-urlencoded", body)
				
				if err != nil {
					uploadMultiple.AnidexStatus = 2
					uploadMultiple.AnidexMessage = "Error during the HTTP POST request"
				}
				defer rsp.Body.Close()
				body_byte, err := ioutil.ReadAll(rsp.Body)
				if err != nil {
					uploadMultiple.AnidexStatus = 2
					uploadMultiple.AnidexMessage = "Error"
				}
				if uploadMultiple.AnidexStatus == 1 {
					uploadMultiple.AnidexMessage = string(body_byte)
					uploadMultiple.AnidexStatus = 3
				}
			}
			
			if c.PostForm("nyaasi_api") != ""  {
				uploadMultiple.NyaasiStatus = 1
				uploadMultiple.NyaasiMessage = "Sorry u are not allowed"
			}
			
			if c.PostForm("tokyot_api") != ""  {
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
