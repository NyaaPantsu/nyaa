package torrentLanguages

import (
	"strings"

	"sort"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
)

var torrentLanguages []string

func initTorrentLanguages() {
	languages := publicSettings.GetAvailableLanguages()

	for _, lang := range languages {
		torrentLanguages = append(torrentLanguages, lang.Code)
	}

	// Also support languages we don't have a translation
	torrentLanguages = append(torrentLanguages, config.Get().Torrents.AdditionalLanguages...)

	sort.Strings(torrentLanguages)
}

// GetTorrentLanguages returns a list of available torrent languages.
func GetTorrentLanguages() []string {
	if torrentLanguages == nil {
		initTorrentLanguages()
	}

	return torrentLanguages
}

// LanguageExists check if said language is available for torrents
func LanguageExists(lang string) bool {
	langs := GetTorrentLanguages()
	for _, code := range langs {
		if code == lang {
			return true
		}
	}

	return false
}

// FlagFromLanguage reads the language's country code.
func FlagFromLanguage(lang string) string {
	languageSplit := strings.Split(lang, "-")
	if len(languageSplit) > 1 {
		return languageSplit[1]
	}

	return ""
}
