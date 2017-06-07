package torrentLanguages

import (
	"strings"

	"github.com/NyaaPantsu/nyaa/util/publicSettings"
	"github.com/nicksnyder/go-i18n/i18n"
)

// TorrentLanguage stores info about a language (code, flag), and defines the
// translation ID to be used when translating its name. The language code is read
// automatically from the available translations, but the flag should be defined
// inside them.
type TorrentLanguage struct {
	Code string
	Flag string
	/* We need to translate every language name in each language,
	so you see eg. Spanish and not Español */
	NameTranslationID string
}

var torrentLanguages map[string]TorrentLanguage

func initTorrentLanguages() {
	torrentLanguages = make(map[string]TorrentLanguage)
	languages := publicSettings.GetAvailableLanguages()
	for code := range languages {
		T, _, _ := i18n.TfuncAndLanguage(code)

		// Read the flag from the translation file, if it has one.
		flag := T("flag")
		if flag == "flag" {
			// Try using the flag from the second part of the language code (en-us would give "us", for example)
			split := strings.Split(code, "-")
			if len(split) > 1 {
				flag = split[1]
			}
		}

		if flag != "flag" {
			torrentLanguages[code] = TorrentLanguage{code, flag, "language_" + code + "_name"}
		}
	}
}

// GetTorrentLanguages returns a map of available torrent languages.
func GetTorrentLanguages() map[string]TorrentLanguage {
	if torrentLanguages == nil {
		initTorrentLanguages()
	}

	return torrentLanguages
}
