package controllers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"os"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/models/activities"
	"github.com/NyaaPantsu/nyaa/models/comments"
	"github.com/NyaaPantsu/nyaa/models/notifications"
	"github.com/NyaaPantsu/nyaa/models/reports"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	"github.com/NyaaPantsu/nyaa/utils/filelist"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/sanitize"
	"github.com/NyaaPantsu/nyaa/utils/search/structs"
	"github.com/NyaaPantsu/nyaa/utils/upload"
	"github.com/NyaaPantsu/nyaa/utils/validator/torrent"
	"github.com/gin-gonic/gin"
)

// ViewHandler : Controller for displaying a torrent
func ViewHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	messages := msg.GetMessages(c)
	user := getUser(c)

	if c.Request.URL.Query()["success"] != nil {
		messages.AddInfoT("infos", "torrent_uploaded")
	}
	if c.Request.URL.Query()["badcaptcha"] != nil {
		messages.AddErrorT("errors", "bad_captcha")
	}
	if c.Request.URL.Query()["reported"] != nil {
		messages.AddInfoTf("infos", "report_msg", id)
	}

	torrent, err := torrents.FindByID(uint(id))

	if c.Request.URL.Query()["notif"] != nil {
		notifications.ToggleReadNotification(torrent.Identifier(), user.ID)
	}

	if err != nil {
		NotFoundHandler(c)
		return
	}
	b := torrent.ToJSON()
	folder := filelist.FileListToFolder(torrent.FileList, "root")
	captchaID := ""
	if user.NeedsCaptcha() {
		captchaID = captcha.GetID()
	}
	torrentTemplate(c, b, folder, captchaID)
}

// ViewHeadHandler : Controller for checking a torrent
func ViewHeadHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		return
	}

	_, err = torrents.FindRawByID(uint(id))

	if err != nil {
		NotFoundHandler(c)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// PostCommentHandler : Controller for posting a comment
func PostCommentHandler(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)

	torrent, err := torrents.FindByID(uint(id))
	if err != nil {
		NotFoundHandler(c)
		return
	}

	currentUser := getUser(c)
	messages := msg.GetMessages(c)

	if currentUser.NeedsCaptcha() {
		userCaptcha := captcha.Extract(c)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddErrorT("errors", "bad_captcha")
		}
	}
	content := sanitize.Sanitize(c.PostForm("comment"), "comment")

	if strings.TrimSpace(content) == "" {
		messages.AddErrorT("errors", "comment_empty")
	}
	if len(content) > config.Get().CommentLength {
		messages.AddErrorT("errors", "comment_toolong")
	}
	if !messages.HasErrors() {

		_, err := comments.Create(content, torrent, currentUser)
		if err != nil {
			messages.Error(err)
		}
	}
	ViewHandler(c)
}

// ReportTorrentHandler : Controller for sending a torrent report
func ReportTorrentHandler(c *gin.Context) {
	fmt.Println("report")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	messages := msg.GetMessages(c)
	captchaError := "?reported"
	currentUser := getUser(c)
	if currentUser.NeedsCaptcha() {
		userCaptcha := captcha.Extract(c)
		if !captcha.Authenticate(userCaptcha) {
			captchaError = "?badcaptcha"
			messages.AddErrorT("errors", "bad_captcha")
		}
	}
	torrent, err := torrents.FindByID(uint(id))
	if err != nil {
		messages.Error(err)
	}
	if !messages.HasErrors() {
		_, err := reports.Create(c.PostForm("report_type"), torrent, currentUser)
		messages.AddInfoTf("infos", "report_msg", id)
		if err != nil {
			messages.ImportFromError("errors", err)
		}
		c.Redirect(http.StatusSeeOther, "/view/"+strconv.Itoa(int(torrent.ID))+captchaError)
	} else {
		ReportViewTorrentHandler(c)
	}
}

// ReportTorrentHandler : Controller for sending a torrent report
func ReportViewTorrentHandler(c *gin.Context) {

	type Report struct {
		ID        uint
		CaptchaID string
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	messages := msg.GetMessages(c)
	currentUser := getUser(c)
	if currentUser.ID > 0 {
		torrent, err := torrents.FindByID(uint(id))
		if err != nil {
			messages.Error(err)
		}
		captchaID := ""
		if currentUser.NeedsCaptcha() {
			captchaID = captcha.GetID()
		}
		formTemplate(c, "site/torrents/report.jet.html", Report{torrent.ID, captchaID})
	} else {
		c.Status(404)
	}
}

// TorrentEditUserPanel : Controller for editing a user torrent by a user, after GET request
func TorrentEditUserPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	torrent, _ := torrents.FindByID(uint(id))
	currentUser := getUser(c)
	if currentUser.CurrentOrAdmin(torrent.UploaderID) {
		uploadForm := torrentValidator.TorrentRequest{}
		uploadForm.Name = torrent.Name
		uploadForm.Category = strconv.Itoa(torrent.Category) + "_" + strconv.Itoa(torrent.SubCategory)
		uploadForm.Remake = torrent.Status == models.TorrentStatusRemake
		uploadForm.WebsiteLink = string(torrent.WebsiteLink)
		uploadForm.Description = string(torrent.Description)
		uploadForm.Hidden = torrent.Hidden
		uploadForm.Languages = torrent.Languages
		formTemplate(c, "site/torrents/edit.jet.html", uploadForm)
	} else {
		NotFoundHandler(c)
	}
}

// TorrentPostEditUserPanel : Controller for editing a user torrent by a user, after post request
func TorrentPostEditUserPanel(c *gin.Context) {
	var uploadForm torrentValidator.UpdateRequest
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	uploadForm.ID = uint(id)
	messages := msg.GetMessages(c)
	torrent, _ := torrents.FindByID(uint(id))
	currentUser := getUser(c)
	if torrent.ID > 0 && currentUser.CurrentOrAdmin(torrent.UploaderID) {
		errUp := upload.ExtractEditInfo(c, &uploadForm.Update)
		if errUp != nil {
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		if !messages.HasErrors() {
			upload.UpdateTorrent(&uploadForm, torrent, currentUser).Update(currentUser.HasAdmin())
			messages.AddInfoT("infos", "torrent_updated")
		}
		formTemplate(c, "site/torrents/edit.jet.html", uploadForm.Update)
	} else {
		NotFoundHandler(c)
	}
}

// TorrentDeleteUserPanel : Controller for deleting a user torrent by a user
func TorrentDeleteUserPanel(c *gin.Context) {
	id, _ := strconv.ParseInt(c.Query("id"), 10, 32)
	currentUser := getUser(c)
	torrent, _ := torrents.FindByID(uint(id))
	if currentUser.CurrentOrAdmin(torrent.UploaderID) {
		_, _, err := torrent.Delete(false)
		if err == nil {
			if torrent.Uploader == nil {
				torrent.Uploader = &models.User{}
			}
			_, username := torrents.HideUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
			if currentUser.HasAdmin() { // We hide username on log activity if user is not admin and torrent is hidden
				activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, currentUser.Username)
			} else {
				activities.Log(&models.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, username)
			}
			//delete reports of torrent
			whereParams := structs.CreateWhereParams("torrent_id = ?", id)
			torrentReports, _, _ := reports.FindOrderBy(&whereParams, "", 0, 0)
			for _, report := range torrentReports {
				report.Delete(false)
			}
		}
		c.Redirect(http.StatusSeeOther, "/?deleted")
	} else {
		NotFoundHandler(c)
	}
}

// DownloadTorrent : Controller for downloading a torrent
func DownloadTorrent(c *gin.Context) {
	hash := c.Param("hash")

	if hash == "" && len(config.Get().Torrents.FileStorage) == 0 {
		//File not found, send 404
		c.AbortWithError(http.StatusNotFound, errors.New("File not found"))
		return
	}

	//Check if file exists and open
	Openfile, err := os.Open(fmt.Sprintf("%s%c%s.torrent", config.Get().Torrents.FileStorage, os.PathSeparator, hash))
	if err != nil {
		//File not found, send 404
		c.AbortWithError(http.StatusNotFound, errors.New("File not found"))
		return
	}
	defer Openfile.Close() //Close after function return

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	torrent, err := torrents.FindRawByHash(hash)

	if err != nil {
		//File not found, send 404
		c.AbortWithError(http.StatusNotFound, errors.New("File not found"))
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.torrent\"", torrent.Name))
	c.Header("Content-Type", "application/x-bittorrent")
	c.Header("Content-Length", FileSize)
	//Send the file
	// We reset the offset to 0
	Openfile.Seek(0, 0)
	io.Copy(c.Writer, Openfile) //'Copy' the file to the client
}
