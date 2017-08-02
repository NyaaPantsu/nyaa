package torrentValidator

import (
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/log"
	msg "github.com/NyaaPantsu/nyaa/utils/messages"
	"github.com/NyaaPantsu/nyaa/utils/metainfo"
	"github.com/NyaaPantsu/nyaa/utils/torrentLanguages"
	"github.com/NyaaPantsu/nyaa/utils/validator/tag"
	"github.com/gin-gonic/gin"
	"github.com/zeebo/bencode"
)

func (r *TorrentRequest) ValidateName() error {
	// then actually check that we have everything we need
	if len(r.Name) == 0 {
		return errTorrentNameInvalid
	}
	return nil
}

func (r *TorrentRequest) ValidateTags() error {
	// We need to parse it to json
	var tags []tagsValidator.CreateForm
	err := json.Unmarshal([]byte(r.Tags), &tags)
	if err != nil {
		r.Tags = ""
		return errTorrentTagsInvalid
	}
	// and filter out multiple tags with the same type (only keep the first one)
	var index config.ArrayString
	var filteredTags []tagsValidator.CreateForm
	for _, tag := range tags {
		if index.Contains(tag.Type) {
			continue
		}
		filteredTags = append(filteredTags, tag)
		index = append(index, tag.Type)
	}
	b, err := json.Marshal(filteredTags)
	if err != nil {
		r.Tags = ""
		log.Infof("Couldn't parse to json the tags %v", filteredTags)
		return errTorrentTagsInvalid
	}
	r.Tags = string(b)
	return nil
}

func (r *TorrentRequest) ValidateDescription() error {
	if len(r.Description) > config.Get().DescriptionLength {
		return errTorrentDescInvalid
	}
	return nil
}

func (r *TorrentRequest) ValidateMagnet() error {
	magnetURL, err := url.Parse(string(r.Magnet)) //?
	if err != nil {
		return err
	}
	xt := magnetURL.Query().Get("xt")
	if !strings.HasPrefix(xt, "urn:btih:") {
		return errTorrentMagnetInvalid
	}
	xt = strings.SplitAfter(xt, ":")[2]
	r.Infohash = strings.TrimSpace(strings.ToUpper(strings.Split(xt, "&")[0]))

	return nil
}

func (r *TorrentRequest) ValidateWebsiteLink() error {
	if r.WebsiteLink != "" {
		// WebsiteLink
		urlRegexp, _ := regexp.Compile(`^(https?:\/\/|ircs?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*(\/.*)?$`)
		if !urlRegexp.MatchString(r.WebsiteLink) {
			return errTorrentURIInvalid
		}
	}
	return nil
}

func (r *TorrentRequest) ValidateHash() error {
	isBase32, err := regexp.MatchString("^[2-7A-Z]{32}$", r.Infohash)
	if err != nil {
		return err
	}
	if !isBase32 {
		isBase16, err := regexp.MatchString("^[0-9A-F]{40}$", r.Infohash)
		if err != nil {
			return err
		}
		if !isBase16 {
			return errTorrentHashInvalid
		}
	} else {
		//convert to base16
		data, err := base32.StdEncoding.DecodeString(r.Infohash)
		if err != nil {
			return err
		}
		hash16 := make([]byte, hex.EncodedLen(len(data)))
		hex.Encode(hash16, data)
		r.Infohash = strings.ToUpper(string(hash16))
	}
	return nil
}

// ExtractCategory : takes an http request and computes category field for this form
func (r *TorrentRequest) ExtractCategory() error {
	catsSplit := strings.Split(r.Category, "_")
	// need this to prevent out of index panics
	if len(catsSplit) != 2 {
		return errTorrentCatInvalid
	}
	CatID, err := strconv.Atoi(catsSplit[0])
	if err != nil {
		return errTorrentCatInvalid
	}
	SubCatID, err := strconv.Atoi(catsSplit[1])
	if err != nil {
		return errTorrentCatInvalid
	}

	if !categories.Exists(r.Category) {
		return errTorrentCatInvalid
	}

	r.CategoryID = CatID
	r.SubCategoryID = SubCatID
	return nil
}

// ExtractLanguage : takes a http request, computes the torrent language from the form.
func (r *TorrentRequest) ExtractLanguage() error {
	isEnglishCategory := false
	for _, cat := range config.Get().Torrents.EnglishOnlyCategories {
		if cat == r.Category {
			isEnglishCategory = true
			break
		}
	}
	if len(r.Languages) == 0 {
		// If no language, but in an English category, set to en-us, else just stop the check.
		if !isEnglishCategory {
			return nil
		}
		r.Languages = append(r.Languages, "en")
		return nil
	}
	englishSelected := false
	for _, language := range r.Languages {
		if language == "en" {
			englishSelected = true
		}

		if language != "" && !torrentLanguages.LanguageExists(language) {
			return errTorrentLangInvalid
		}

		if strings.HasPrefix(language, "en") && isEnglishCategory {
			englishSelected = true
		}
	}

	// We shouldn't return an error for languages, just adding the right language is enough
	if !englishSelected && isEnglishCategory {
		r.Languages = append(r.Languages, "en")
		return nil
	}

	// We shouldn't return an error if someone has selected only english for languages and missed the right category. Just move the torrent in the right one
	// Multiple if conditions so we only do this for loop when needed
	if len(r.Languages) == 1 && strings.HasPrefix(r.Languages[0], "en") && !isEnglishCategory && r.CategoryID > 0 {
		for key, cat := range config.Get().Torrents.NonEnglishOnlyCategories {
			if cat == r.Category {
				r.Category = config.Get().Torrents.EnglishOnlyCategories[key]
				isEnglishCategory = true
				break
			}
		}
	}

	return nil
}

// ValidateMultipartUpload : Check if multipart upload is valid
func (r *TorrentRequest) ValidateMultipartUpload(c *gin.Context, uploadFormTorrent string) (multipart.File, error) {
	// first: parse torrent file (if any) to fill missing information
	tfile, _, err := c.Request.FormFile(uploadFormTorrent)
	if err == nil {
		var torrent metainfo.TorrentFile

		// decode torrent
		_, seekErr := tfile.Seek(0, io.SeekStart)
		if seekErr != nil {
			return tfile, seekErr
		}
		err = bencode.NewDecoder(tfile).Decode(&torrent)
		if err != nil {
			return tfile, metainfo.ErrInvalidTorrentFile
		}

		// check a few things
		if torrent.IsPrivate() {
			return tfile, errTorrentPrivate
		}
		trackers := torrent.GetAllAnnounceURLS()
		r.Trackers = CheckTrackers(trackers)
		if len(r.Trackers) == 0 {
			return tfile, errTorrentNoTrackers
		}

		// Name
		if len(r.Name) == 0 {
			r.Name = torrent.TorrentName()
		}

		// Magnet link: if a file is provided it should be empty
		if len(r.Magnet) != 0 {
			return tfile, errTorrentAndMagnet
		}

		_, seekErr = tfile.Seek(0, io.SeekStart)
		if seekErr != nil {
			return tfile, seekErr
		}
		infohash, err := metainfo.DecodeInfohash(tfile)
		if err != nil {
			return tfile, metainfo.ErrInvalidTorrentFile
		}
		r.Infohash = infohash
		r.Magnet = format.InfoHashToMagnet(infohash, r.Name, trackers...)

		// extract filesize
		r.Filesize = int64(torrent.TotalSize())

		// extract filelist
		fileInfos := torrent.Info.GetFiles()
		for _, fileInfo := range fileInfos {
			r.FileList = append(r.FileList, uploadedFile{
				Path:     fileInfo.Path,
				Filesize: int64(fileInfo.Length),
			})
		}
	} else {
		err = r.ValidateMagnet()
		if err != nil {
			return tfile, err
		}
		err = r.ValidateHash()
		if err != nil {
			return tfile, err
		}
		// TODO: Get Trackers from magnet URL
		r.Filesize = 0
		r.Filepath = ""

		return tfile, nil
	}
	return tfile, err
}

// ExtractInfo : Function to assign values from request to ReassignForm
func (f *ReassignForm) ExtractInfo(c *gin.Context) bool {
	f.By = c.PostForm("by")
	messages := msg.GetMessages(c)
	if f.By != "olduser" && f.By != "torrentid" {
		messages.AddErrorTf("errors", "no_action_exist", f.By)
		return false
	}

	f.Data = strings.Trim(c.PostForm("data"), " \r\n")
	if f.By == "olduser" {
		if f.Data == "" {
			messages.AddErrorT("errors", "user_not_found")
			return false
		} else if strings.Contains(f.Data, "\n") {
			messages.AddErrorT("errors", "multiple_username_error")
			return false
		}
	} else if f.By == "torrentid" {
		if f.Data == "" {
			messages.AddErrorT("errors", "no_id_given")
			return false
		}
		splitData := strings.Split(f.Data, "\n")
		for i, tmp := range splitData {
			tmp = strings.Trim(tmp, " \r")
			torrentID, err := strconv.ParseUint(tmp, 10, 0)
			if err != nil {
				messages.AddErrorTf("errors", "parse_error_line", i+1)
				return false // TODO: Shouldn't it continue to parse the rest and display the errored lines?
			}
			f.Torrents = append(f.Torrents, uint(torrentID))
		}
	}

	tmpID := c.PostForm("to")
	parsed, err := strconv.ParseUint(tmpID, 10, 32)
	if err != nil {
		messages.Error(err)
		return false
	}
	f.AssignTo = uint(parsed)
	_, _, _, _, err = cookies.RetrieveUserFromRequest(c, uint(parsed))
	if err != nil {
		messages.AddErrorTf("errors", "no_user_found_id", int(parsed))
		return false
	}

	return true
}
