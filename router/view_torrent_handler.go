package router

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"os"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/db"
	"github.com/NyaaPantsu/nyaa/model"
	"github.com/NyaaPantsu/nyaa/service"
	"github.com/NyaaPantsu/nyaa/service/activity"
	"github.com/NyaaPantsu/nyaa/service/api"
	"github.com/NyaaPantsu/nyaa/service/notifier"
	"github.com/NyaaPantsu/nyaa/service/report"
	"github.com/NyaaPantsu/nyaa/service/torrent"
	"github.com/NyaaPantsu/nyaa/service/user/permission"
	"github.com/NyaaPantsu/nyaa/util"
	"github.com/NyaaPantsu/nyaa/util/captcha"
	"github.com/NyaaPantsu/nyaa/util/filelist"
	msg "github.com/NyaaPantsu/nyaa/util/messages"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/gin-gonic/gin"
)

// ViewHandler : Controller for displaying a torrent
func ViewHandler(c *gin.Context) {
	id := c.Param("id")
	messages := msg.GetMessages(c)
	user := getUser(c)

	if c.Request.URL.Query()["success"] != nil {
		messages.AddInfo("infos", "Torrent uploaded successfully!")
	}

	torrent, err := torrentService.GetTorrentByID(id)

	if c.Request.URL.Query()["notif"] != nil {
		notifierService.ToggleReadNotification(torrent.Identifier(), user.ID)
	}

	if err != nil {
		NotFoundHandler(c)
		return
	}
	b := torrent.ToJSON()
	folder := filelist.FileListToFolder(torrent.FileList, "root")
	captchaID := ""
	if userPermission.NeedsCaptcha(user) {
		captchaID = captcha.GetID()
	}
	torrentTemplate(c, b, folder, captchaID)
}

// ViewHeadHandler : Controller for checking a torrent
func ViewHeadHandler(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 32)
	if err != nil {
		return
	}

	_, err = torrentService.GetRawTorrentByID(uint(id))

	if err != nil {
		NotFoundHandler(c)
		return
	}

	c.AbortWithStatus(http.StatusOK)
}

// PostCommentHandler : Controller for posting a comment
func PostCommentHandler(c *gin.Context) {
	id := c.Query("id")

	torrent, err := torrentService.GetTorrentByID(id)
	if err != nil {
		NotFoundHandler(c)
		return
	}

	currentUser := getUser(c)
	messages := msg.GetMessages(c)

	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(c)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddErrorT("errors", "bad_captcha")
		}
	}
	content := util.Sanitize(c.PostForm("comment"), "comment")

	if strings.TrimSpace(content) == "" {
		messages.AddErrorT("errors", "comment_empty")
	}
	if len(content) > config.Conf.CommentLength {
		messages.AddErrorT("errors", "comment_toolong")
	}
	if !messages.HasErrors() {
		userID := currentUser.ID

		comment := model.Comment{TorrentID: torrent.ID, UserID: userID, Content: content, CreatedAt: time.Now()}
		err := db.ORM.Create(&comment).Error
		if err != nil {
			messages.Error(err)
		}
		comment.Torrent = &torrent

		url := "/view/" + strconv.FormatUint(uint64(torrent.ID), 10) + "/" + torrent.Name
		torrent.Uploader.ParseSettings()
		if torrent.Uploader.Settings.Get("new_comment") {
			T, _, _ := publicSettings.TfuncAndLanguageWithFallback(torrent.Uploader.Language, torrent.Uploader.Language) // We need to send the notification to every user in their language
			notifierService.NotifyUser(torrent.Uploader, comment.Identifier(), fmt.Sprintf(T("new_comment_on_torrent"), torrent.Name), url, torrent.Uploader.Settings.Get("new_comment_email"))
		}

		if err != nil {
			messages.Error(err)
		}
	}
	ViewHandler(c)
}

// ReportTorrentHandler : Controller for sending a torrent report
func ReportTorrentHandler(c *gin.Context) {
	id := c.Query("id")
	messages := msg.GetMessages(c)
	currentUser := getUser(c)
	if userPermission.NeedsCaptcha(currentUser) {
		userCaptcha := captcha.Extract(c)
		if !captcha.Authenticate(userCaptcha) {
			messages.AddErrorT("errors", "bad_captcha")
		}
	}
	if !messages.HasErrors() {
		idNum, _ := strconv.Atoi(id)
		userID := currentUser.ID

		report := model.TorrentReport{
			Description: c.PostForm("report_type"),
			TorrentID:   uint(idNum),
			UserID:      userID,
			CreatedAt:   time.Now(),
		}

		err := db.ORM.Create(&report).Error
		messages.AddInfoTf("infos", "report_msg", id)
		if err != nil {
			messages.ImportFromError("errors", err)
		}
	}
	ViewHandler(c)
}

// TorrentEditUserPanel : Controller for editing a user torrent by a user, after GET request
func TorrentEditUserPanel(c *gin.Context) {
	id := c.Query("id")
	torrent, _ := torrentService.GetTorrentByID(id)
	currentUser := getUser(c)
	if userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
		uploadForm := apiService.NewTorrentRequest()
		uploadForm.Name = torrent.Name
		uploadForm.Category = strconv.Itoa(torrent.Category) + "_" + strconv.Itoa(torrent.SubCategory)
		uploadForm.Remake = torrent.Status == model.TorrentStatusRemake
		uploadForm.WebsiteLink = string(torrent.WebsiteLink)
		uploadForm.Description = string(torrent.Description)
		uploadForm.Hidden = torrent.Hidden
		uploadForm.Language = torrent.Language
		formTemplate(c, "user/torrent_edit.jet.html", uploadForm)
	} else {
		NotFoundHandler(c)
	}
}

// TorrentPostEditUserPanel : Controller for editing a user torrent by a user, after post request
func TorrentPostEditUserPanel(c *gin.Context) {
	var uploadForm apiService.TorrentRequest
	id := c.Query("id")
	messages := msg.GetMessages(c)
	torrent, _ := torrentService.GetTorrentByID(id)
	currentUser := getUser(c)
	if torrent.ID > 0 && userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
		errUp := uploadForm.ExtractEditInfo(c)
		if errUp != nil {
			messages.AddErrorT("errors", "fail_torrent_update")
		}
		if !messages.HasErrors() {
			status := model.TorrentStatusNormal
			if uploadForm.Remake { // overrides trusted
				status = model.TorrentStatusRemake
			} else if currentUser.IsTrusted() {
				status = model.TorrentStatusTrusted
			}
			// update some (but not all!) values
			torrent.Name = uploadForm.Name
			torrent.Category = uploadForm.CategoryID
			torrent.SubCategory = uploadForm.SubCategoryID
			torrent.Status = status
			torrent.Hidden = uploadForm.Hidden
			torrent.WebsiteLink = uploadForm.WebsiteLink
			torrent.Description = uploadForm.Description
			torrent.Language = uploadForm.Language
			db.ORM.Model(&torrent).UpdateColumn(&torrent)
			messages.AddInfoT("infos", "torrent_updated")
		}
		formTemplate(c, "user/torren_edit.jet.html", uploadForm)
	} else {
		NotFoundHandler(c)
	}
}

// TorrentDeleteUserPanel : Controller for deleting a user torrent by a user
func TorrentDeleteUserPanel(c *gin.Context) {
	id := c.Query("id")
	currentUser := getUser(c)
	torrent, _ := torrentService.GetTorrentByID(id)
	if userPermission.CurrentOrAdmin(currentUser, torrent.UploaderID) {
		_, _, err := torrentService.DeleteTorrent(id)
		if err == nil {
			_, username := torrentService.HideTorrentUser(torrent.UploaderID, torrent.Uploader.Username, torrent.Hidden)
			if userPermission.HasAdmin(currentUser) { // We hide username on log activity if user is not admin and torrent is hidden
				activity.Log(&model.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, currentUser.Username)
			} else {
				activity.Log(&model.User{}, torrent.Identifier(), "delete", "torrent_deleted_by", strconv.Itoa(int(torrent.ID)), username, username)
			}
			//delete reports of torrent
			whereParams := serviceBase.CreateWhereParams("torrent_id = ?", id)
			reports, _, _ := reportService.GetTorrentReportsOrderBy(&whereParams, "", 0, 0)
			for _, report := range reports {
				reportService.DeleteTorrentReport(report.ID)
			}
		}
		c.Redirect(http.StatusSeeOther, "/?deleted")
	} else {
		NotFoundHandler(c)
	}
}

// DownloadTorrent : Controller for downloading a torrent
func DownloadTorrent(c *gin.Context) {
	hash := c.Query("hash")

	if hash == "" && len(config.Conf.Torrents.FileStorage) == 0 {
		//File not found, send 404
		c.AbortWithError(http.StatusNotFound, errors.New("File not found"))
		return
	}

	//Check if file exists and open
	Openfile, err := os.Open(fmt.Sprintf("%s%c%s.torrent", config.Conf.Torrents.FileStorage, os.PathSeparator, hash))
	if err != nil {
		//File not found, send 404
		c.AbortWithError(http.StatusNotFound, errors.New("File not found"))
		return
	}
	defer Openfile.Close() //Close after function return

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	torrent, err := torrentService.GetRawTorrentByHash(hash)

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
