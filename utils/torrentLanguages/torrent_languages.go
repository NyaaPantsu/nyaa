package torrentLanguages

import (
	"strings"

	"github.com/NyaaPantsu/nyaa/config"
	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
)

var torrentLanguages publicSettings.Languages

func initTorrentLanguages() {
	languages := publicSettings.GetAvailableLanguages()
	var langSort []string
	for _, lang := range languages {
		langSort = append(langSort, lang.Code)
	}

	// Also support languages we don't have a translation
	langSorted := publicSettings.ParseLanguages(append(langSort, config.Get().Torrents.AdditionalLanguages...))

	prevLang := ""
	for _, lang := range langSorted {
		if prevLang == lang.Code {
			last := len(torrentLanguages) - 1
			if last > 0 && !strings.Contains(torrentLanguages[last].Name, lang.Name) {
				torrentLanguages[last].Name += ", " + lang.Name
				torrentLanguages[last].Tag += ", " + lang.Tag
			}
		} else {
			prevLang = lang.Code
			torrentLanguages = append(torrentLanguages, lang)
		}
	}
}

// GetTorrentLanguages returns a list of available torrent languages.
func GetTorrentLanguages() publicSettings.Languages {
	if torrentLanguages == nil {
		initTorrentLanguages()
	}

	return torrentLanguages
}

// LanguageExists check if said language is available for torrents
func LanguageExists(languageCode string) bool {
	langs := GetTorrentLanguages()
	for _, lang := range langs {
		if lang.Code == publicSettings.GetParentTag(languageCode).String() {
			return true
		}
	}

	return false
}
