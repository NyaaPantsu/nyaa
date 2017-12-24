package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
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
	extraParams := map[string]string{
		//Required
		"api_key":   apiKey,
		"subcat_id": subCat,
		"group_id":  "0",
		"lang_id":   lang,

		//Optional
		"description":  torrent.Description,
		"torrent_name": torrent.Name,
	}
	if config.IsSukebei() {
		extraParams["hentai"] = "1"
	}
	if torrent.IsRemake() {
		extraParams["reencode"] = "1"
	}
	if torrent.IsAnon() {
		extraParams["private"] = "1"
	}
	request, err := newfileUploadRequest("https://anidex.info/api/", extraParams, "file", torrent.GetPath())
	if err != nil {
		log.CheckError(err)
		return
	}
	client := &http.Client{}
	rsp, err := client.Do(request)
	if err != nil {
		log.CheckError(err)
		return
	}
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
func ToNyaasi(username string, password string, torrent *models.Torrent) {
	uploadMultiple := MultipleForm{PantsuID: torrent.ID, Nyaasi: service{Status: pendingState}}
	uploadMultiple.Nyaasi.Message = "Sorry u are not allowed"
	uploadMultiple.save(nyaasi)
	log.Info("Create NyaaSi instance")

	// If the torrent is posted as anonymous or apikey is not set, we set it with default value
	if username == "" || (torrent.Hidden && username != "") {
		username = config.Get().Upload.DefaultNyaasiUsername
		password = config.Get().Upload.DefaultNyaasiPassword
	}

	if username == "" || password == "" { // You need to check that username AND password are not empty even after config. Since they are left empty in config by default and are required
		log.Errorf("Username or Password is empty, we can't upload to Nyaa.si for torrent %d", torrent.ID)
		uploadMultiple.updateAndSave(nyaasi, errorState, "No valid account providen (required)")
		return
	}

	params := map[string]interface{}{
		"name":        torrent.Name,
		"category":    Category(nyaasi, torrent),
		"information": "",
		"description": torrent.Description,
		"anonymous":   torrent.IsAnon(),
		"hidden":      false,
		"remake":      torrent.IsRemake(),
		"trusted":     torrent.IsTrusted(),
	}
	torrentData, _ := json.Marshal(params)
	extraParams := map[string]string{
		"torrent_data": string(torrentData),
	}

	request, err := newfileUploadRequest("https://nyaa.si/api/upload", extraParams, "torrent", torrent.GetPath())
	if err != nil {
		log.CheckError(err)
		return
	}
	request.SetBasicAuth(username, password)
	client := &http.Client{}
	rsp, err := client.Do(request)
	if err != nil {
		log.CheckError(err)
		return
	}
	log.Info("Launch Nyaa.Si http request")

	if err != nil {
		uploadMultiple.updateAndSave(nyaasi, errorState, "Error during the HTTP POST request")
		log.CheckErrorWithMessage(err, "Error in request: %s")
		return
	}
	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		uploadMultiple.updateAndSave(nyaasi, errorState, "Unknown error")
		log.CheckErrorWithMessage(err, "Error in parsing request: %s")
		return
	}
	if uploadMultiple.Nyaasi.Status == pendingState {
		var data map[string]interface{}
		if err = json.Unmarshal(bodyByte, &data); err != nil {
			log.CheckErrorWithMessage(err, "Cannot unmarshal json Response after upload request to Nyaa.Si")
			uploadMultiple.Nyaasi.Status = errorState
			uploadMultiple.Nyaasi.Message = err.Error()
		}
		if _, ok := data["errors"]; ok {
			uploadMultiple.Nyaasi.Status = errorState
			uploadMultiple.Nyaasi.Message = string(bodyByte)
		} else {
			uploadMultiple.Nyaasi.Status = doneState
			uploadMultiple.Nyaasi.Message = fmt.Sprintf("%s", data["url"].(string))
		}
		uploadMultiple.save(nyaasi)
		log.Info("Nyaa.Si request done")
		fmt.Println(uploadMultiple)
	}
}

// ToTTosho : function to upload a torrent to TokyoTosho
func ToTTosho(apiKey string, torrent *models.Torrent) {
	uploadMultiple := MultipleForm{PantsuID: torrent.ID, TTosho: service{Status: pendingState}}
	uploadMultiple.save(ttosho)
	log.Info("Create TokyoTosho instance")

	// If the torrent is posted as anonymous or apikey is not set, we set it with default value
	if apiKey == "" || (torrent.Hidden && apiKey != "") {
		apiKey = config.Get().Upload.DefaultTokyoTToken
	}

	if apiKey == "" { // You need to check that apikey is not empty even after config. Since it is left empty in config by default and is required
		log.Errorf("ApiKey is empty, we can't upload to TokyoTosho for torrent %d", torrent.ID)
		uploadMultiple.updateAndSave(ttosho, errorState, "No ApiKey providen (required)")
		return
	}

	extraParams := map[string]string{
		//Required
		"apikey": apiKey,
		"url":    torrent.Download(),
		"type":   Category(ttosho, torrent),
		"send":   "true",

		//Optional
		"website": torrent.WebsiteLink,
		"comment": torrent.Description,
	}
	request, err := newUploadRequest("https://www.tokyotosho.info/new.php", extraParams)
	if err != nil {
		log.CheckError(err)
		return
	}
	client := &http.Client{}
	rsp, err := client.Do(request)
	if err != nil {
		log.CheckError(err)
		return
	}
	log.Info("Launch TokyoTosho http request")

	if err != nil {
		uploadMultiple.updateAndSave(ttosho, errorState, "Error during the HTTP POST request")
		log.CheckErrorWithMessage(err, "Error in request: %s")
		return
	}
	defer rsp.Body.Close()
	bodyByte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		uploadMultiple.updateAndSave(ttosho, errorState, "Unknown error")
		log.CheckErrorWithMessage(err, "Error in parsing request: %s")
		return
	}
	if uploadMultiple.TTosho.Status == pendingState {
		if strings.Contains(string(bodyByte), "OK,") {
			uploadMultiple.TTosho.Status = doneState
			idnumber := strings.Split(string(bodyByte), ",")
			uploadMultiple.TTosho.Message = fmt.Sprintf("Upload done! https://www.tokyotosho.info/details.php?id=%s", idnumber[1])
		} else {
			uploadMultiple.TTosho.Status = errorState
			uploadMultiple.TTosho.Message = string(bodyByte)
		}
		uploadMultiple.save(ttosho)
		log.Info("TokyoTosho request done")
		fmt.Println(uploadMultiple)
	}
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

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(part, file); err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Content-Type", writer.FormDataContentType())
	return request, nil
}

// Creates a new upload http request with optional extra params
func newUploadRequest(uri string, params map[string]string) (*http.Request, error) {
	var form url.Values
	for key, val := range params {
		form[key] = append(form[key], val)
	}

	body := bytes.NewBufferString(form.Encode())
	request, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}
	return request, nil
}
