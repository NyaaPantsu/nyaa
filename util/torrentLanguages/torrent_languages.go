package torrentLanguages

import (
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/util/publicSettings"
)

var torrentLanguages []string

func initTorrentLanguages() {
	languages := publicSettings.GetAvailableLanguages()
	for code := range languages {
		torrentLanguages = append(torrentLanguages, code)
	}

	// Also support languages we don't have a translation
	for _, code := range config.Conf.Torrents.AdditionalLanguages {
		torrentLanguages = append(torrentLanguages, code)
	}
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
