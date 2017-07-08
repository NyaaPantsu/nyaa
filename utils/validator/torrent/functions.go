package torrentValidator

import (
	"encoding/base32"
	"encoding/hex"
	"errors"
	"io"
	"mime/multipart"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/categories"
	"github.com/NyaaPantsu/nyaa/utils/format"
	"github.com/NyaaPantsu/nyaa/utils/metainfo"
	"github.com/NyaaPantsu/nyaa/utils/torrentLanguages"
	"github.com/gin-gonic/gin"
	"github.com/zeebo/bencode"
)

func (r *TorrentRequest) ValidateName() error {
	// then actually check that we have everything we need
	if len(r.Name) == 0 {
		return errors.New("torrent_name_invalid")
	}
	return nil
}

func (r *TorrentRequest) ValidateDescription() error {
	if len(r.Description) > config.Get().DescriptionLength {
		return errors.New("torrent_desc_invalid")
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
		return errors.New("torrent_magnet_invalid")
	}
	xt = strings.SplitAfter(xt, ":")[2]
	r.Infohash = strings.TrimSpace(strings.ToUpper(strings.Split(xt, "&")[0]))

	return nil
}

func (r *TorrentRequest) ValidateWebsiteLink() error {
	if r.WebsiteLink != "" {
		// WebsiteLink
		urlRegexp, _ := regexp.Compile(`^(https?:\/\/|ircs?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`)
		if !urlRegexp.MatchString(r.WebsiteLink) {
			return errors.New("torrent_uri_invalid")
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
			return errors.New("torrent_hash_invalid")
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
		return errors.New("torrent_cat_invalid")
	}
	CatID, err := strconv.Atoi(catsSplit[0])
	if err != nil {
		return errors.New("torrent_cat_invalid")
	}
	SubCatID, err := strconv.Atoi(catsSplit[1])
	if err != nil {
		return errors.New("torrent_cat_invalid")
	}

	if !categories.Exists(r.Category) {
		return errors.New("torrent_cat_invalid")
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
		r.Languages = append(r.Languages, "en-us")
		return nil
	}
	englishSelected := false
	for _, language := range r.Languages {
		if language == "en-us" {
			englishSelected = true
		}

		if language != "" && !torrentLanguages.LanguageExists(language) {
			return errors.New("torrent_lang_invalid")
		}

		if strings.HasPrefix(language, "en") && isEnglishCategory {
			englishSelected = true
		}
	}

	// We shouldn't return an error for languages, just adding the right language is enough
	if !englishSelected && isEnglishCategory {
		r.Languages = append(r.Languages, "en-us")
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
			return tfile, errors.New("torrent_private")
		}
		trackers := torrent.GetAllAnnounceURLS()
		r.Trackers = CheckTrackers(trackers)
		if len(r.Trackers) == 0 {
			return tfile, errors.New("torrent_no_working_trackers")
		}

		// Name
		if len(r.Name) == 0 {
			r.Name = torrent.TorrentName()
		}

		// Magnet link: if a file is provided it should be empty
		if len(r.Magnet) != 0 {
			return tfile, errors.New("torrent_plus_magnet")
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
