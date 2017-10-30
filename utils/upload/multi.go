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
	"github.com/NyaaPantsu/nyaa/utils/log"
)

const (
	anidex = iota
	nyaasi
	ttosho
)

const (
	pendingState = iota + 1
	errorState
	doneState
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
// TODO: subCat and lang should be taken from torrent model and not asked to be typed again
// so making the conversion here would be better
func ToAnidex(torrent *models.Torrent, apiKey string, subCat string, lang string) {
	uploadMultiple := MultipleForm{PantsuID: torrent.ID, Anidex: service{Status: pendingState}}
	uploadMultiple.save(anidex)
	log.Info("Create anidex instance")

	// If the torrent is posted as anonymous or apikey is not set, we set it with default value
	if apiKey == "" || (torrent.Hidden && apiKey != "") {
		apiKey = config.Get().Upload.DefaultAnidexToken
	}

	if apiKey == "" { // You need to check that apikey is not empty even after config. Since it is left empty in config by default and is required
		log.Errorf("ApiKey is empty, we can't upload to anidex for torrent %d", torrent.ID)
		uploadMultiple.updateAndSave(anidex, errorState, "No ApiKey providen (required)")
		return
	}

	postForm := url.Values{}
	//Required
	postForm.Set("api_key", apiKey)
	postForm.Set("subcat_id", subCat)
	postForm.Set("file", "")
	postForm.Set("group_id", "0")
	postForm.Set("lang_id", lang)

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
	postForm.Set("torrent_name", torrent.Name)

	postForm.Set("debug", "1")

	rsp, err := http.Post("https://anidex.info/api/", "application/x-www-form-urlencoded", bytes.NewBufferString(postForm.Encode()))
	log.Info("Launch anidex http request")

	if err != nil {
		uploadMultiple.updateAndSave(anidex, errorState, "Error during the HTTP POST request")
		log.CheckErrorWithMessage(err, "Error in request: %s")
		return
	}
	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		uploadMultiple.updateAndSave(anidex, errorState, "Unknown error")
		log.CheckErrorWithMessage(err, "Error in parsing request: %s")
		return
	}
	if uploadMultiple.Anidex.Status == pendingState {
		uploadMultiple.Anidex.Message = string(bodyByte)
		if strings.Contains(uploadMultiple.Anidex.Message, "http://") {
			uploadMultiple.Anidex.Status = doneState
		} else {
			uploadMultiple.Anidex.Status = errorState
		}
		uploadMultiple.save(anidex)
		log.Info("Anidex request done")
		fmt.Println(uploadMultiple)
	}
}

// ToNyaasi : function to upload a torrent to anidex
func ToNyaasi(apiKey string, torrent *models.Torrent) {
	uploadMultiple := MultipleForm{PantsuID: torrent.ID, Nyaasi: service{Status: pendingState}}
	uploadMultiple.Nyaasi.Message = "Sorry u are not allowed"
	uploadMultiple.save(nyaasi)
}

// ToTTosho : function to upload a torrent to anidex
func ToTTosho(apiKey string, torrent *models.Torrent) {
	uploadMultiple := MultipleForm{PantsuID: torrent.ID, TTosho: service{Status: pendingState}}
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
