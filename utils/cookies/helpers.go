package cookies

import (
	"github.com/NyaaPantsu/nyaa/config"
)

func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}

func getMaxAge(rememberMe bool) int {
	if rememberMe {
		return 365 * 24 * 3600
	}
	return config.Get().Cookies.MaxAge
}
