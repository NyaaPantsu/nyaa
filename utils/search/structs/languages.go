package structs

import (
	"strings"

	"github.com/NyaaPantsu/nyaa/utils/publicSettings"
)

// ParseLanguages sets languages by string
func ParseLanguages(s []string) publicSettings.Languages {
	var languages publicSettings.Languages
	for _, lang := range s {
		lgSplit := splitsLanguages(lang)
		if len(lgSplit) > 0 {
			languages = append(languages, lgSplit...)
		}
	}
	return languages
}

func splitsLanguages(s string) publicSettings.Languages {
	var languages publicSettings.Languages
	if s != "" {
		parts := strings.Split(s, ",")
		for _, lang := range parts {
			if lang != "" {
				languages = append(languages, publicSettings.Language{Name: "", Code: lang}) // We just need the code
			}
		}
	}
	return languages
}
