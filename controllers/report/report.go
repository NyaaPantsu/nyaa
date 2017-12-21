package reportController

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/models/reports"
	"github.com/NyaaPantsu/nyaa/models/torrents"
	"github.com/NyaaPantsu/nyaa/templates"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/gin-gonic/gin"
)

// ReportTorrentHandler : Controller for sending a torrent report
func ReportTorrentHandler(c *gin.Context) {
	fmt.Println("report")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	messages := msg.GetMessages(c)
	captchaError := "?reported"
	currentUser := router.GetUser(c)
	if currentUser.ID == 0 {
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
		reportMessage :=  c.PostForm("report_message")
		if len(reportMessage) > 60 {
			reportMessage = reportMessage[0:60]
		}
		_, err := reports.Create(c.PostForm("report_type"), reportMessage, torrent, currentUser)
		messages.AddInfoTf("infos", "report_msg", id)
		if err != nil {
			messages.ImportFromError("errors", err)
		}
		c.Redirect(http.StatusSeeOther, "/view/"+strconv.Itoa(int(torrent.ID))+captchaError)
	} else {
		ReportViewTorrentHandler(c)
	}
}

// ReportViewTorrentHandler : Controller for sending a torrent report
func ReportViewTorrentHandler(c *gin.Context) {

	type Report struct {
		ID        uint
		CaptchaID string
	}

	id, _ := strconv.ParseInt(c.Param("id"), 10, 32)
	messages := msg.GetMessages(c)
	currentUser := router.GetUser(c)

	torrent, err := torrents.FindByID(uint(id))
	if err != nil {
		messages.Error(err)
	}
	captchaID := ""
	if currentUser.ID == 0 {
		captchaID = captcha.GetID()
	}
	templates.Form(c, "site/torrents/report.jet.html", Report{torrent.ID, captchaID})
}
