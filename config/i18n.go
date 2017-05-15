package config

type I18nConfig struct {
	TranslationsDirectory string `json:"translations_directory"`
	DefaultLanguage       string `json:"default_language"`
}

var DefaultI18nConfig = I18nConfig{
	TranslationsDirectory: "translations",
	DefaultLanguage:       "en-us",
}
