package config

type I18nConfig struct {
	TranslationsDirectory string `json:"translations_directory"`
	DefaultLanguage       string `json:"default_language"`
}

var DefaultI18nConfig = I18nConfig{
	TranslationsDirectory: "translations",
	DefaultLanguage:       "en-us", // TODO: Remove refs to "en-us" from the code and templates
}
