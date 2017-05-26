package config

// I18nConfig : Config struct for translation
type I18nConfig struct {
	TranslationsDirectory string `json:"translations_directory"`
	DefaultLanguage       string `json:"default_language"`
}

// DefaultI18nConfig : Default configuration for translation
var DefaultI18nConfig = I18nConfig{
	TranslationsDirectory: "translations",
	DefaultLanguage:       "en-us",
}
