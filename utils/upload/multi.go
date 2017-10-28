package upload

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cache"
	"github.com/gin-gonic/gin"
)

const (
	anidex = iota
	nyaasi
	ttosho
)

// Each service gives a status and a message when uploading to them
type service struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// MultipleForm is a struct used to follow the status of the uploads
type MultipleForm struct {
	PantsuID uint    `json:"id"`
	Anidex   service `json:"anidex"`
	Nyaasi   service `json:"nyaasi"`
	TTosho   service `json:"ttosho"`
}

// ToAnidex : function to upload a torrent to anidex
func ToAnidex(c *gin.Context, torrent *models.Torrent) {
	uploadMultiple := MultipleForm{PantsuID: torrent.ID, Anidex: service{Status: 1}}
	uploadMultiple.save(anidex)

	apiKey := c.PostForm("anidex_api")

	// If the torrent is posted as anonymous or apikey is not set, we set it with default value
	if apiKey == "" || torrent.IsAnon() {
		apiKey = config.Get().Upload.DefaultAnidexToken
	}

	if apiKey != "" { // You need to check that apikey is not empty even after config. Since it is left empty in config by default and is required
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
		if torrent.IsRemake() {
			postForm.Set("reencode", "1")
		}
		if torrent.IsAnon() {
			postForm.Set("private", "1")
		}
		if c.PostForm("name") != "" {
			postForm.Set("torrent_name", c.PostForm("name"))
		}

		postForm.Set("debug", "1")

		rsp, err := http.Post("https://anidex.info/api/", "application/x-www-form-urlencoded", bytes.NewBufferString(postForm.Encode()))

		if err != nil {
			uploadMultiple.updateAndSave(anidex, 2, "Error during the HTTP POST request")
			return
		}
		defer rsp.Body.Close()
		bodyByte, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			uploadMultiple.updateAndSave(anidex, 2, "Unknown error")
			return
		}
		if uploadMultiple.Anidex.Status == 1 {
			uploadMultiple.Anidex.Message = string(bodyByte)
			if strings.Contains(uploadMultiple.Anidex.Message, "http://") {
				uploadMultiple.Anidex.Status = 3
			} else if strings.Contains(uploadMultiple.Anidex.Message, "error") {
				uploadMultiple.Anidex.Status = 2
			}
			uploadMultiple.save(anidex)
		}
	}
}

// ToNyaasi : function to upload a torrent to anidex
func ToNyaasi(c *gin.Context, torrent *models.Torrent) {
	uploadMultiple := MultipleForm{PantsuID: torrent.ID, Nyaasi: service{Status: 1}}
	uploadMultiple.Nyaasi.Message = "Sorry u are not allowed"
	uploadMultiple.save(nyaasi)
}

// ToTTosho : function to upload a torrent to anidex
func ToTTosho(c *gin.Context, torrent *models.Torrent) {
	uploadMultiple := MultipleForm{PantsuID: torrent.ID, TTosho: service{Status: 1}}
	uploadMultiple.save(ttosho)
}

// Saves the multipleform in each go routines and share the state of each upload for 5 minutes
// After timeout, the multipleform is flushed from memory
func (u *MultipleForm) save(which int) {
	// We check if there is already a variable shared, if there is we update only the status/message of one service
	if found, ok := cache.C.Get(fmt.Sprintf("tstatus_%d", u.PantsuID)); ok {
		uploadStatus := found.(MultipleForm)
		switch which {
		case anidex:
			uploadStatus.Anidex = u.Anidex
			break
		case nyaasi:
			uploadStatus.Nyaasi = u.Nyaasi
			break
		case ttosho:
			uploadStatus.TTosho = u.TTosho
			break
		}
		u = &uploadStatus
	}
	// And then we save the variable in cache
	cache.C.Set(fmt.Sprintf("tstatus_%d", u.PantsuID), *u, 5*time.Minute)
}

// shortcut function to update and save a service
func (u *MultipleForm) updateAndSave(which int, code int, message string) {
	switch which {
	case anidex:
		u.Anidex.update(code, message)
		break
	case nyaasi:
		u.Nyaasi.update(code, message)
		break
	case ttosho:
		u.TTosho.update(code, message)
		break
	}
	u.save(which)
}

// shortcut function to update both code and message of a service
func (s *service) update(code int, message string) {
	s.Status = code
	s.Message = message
}
