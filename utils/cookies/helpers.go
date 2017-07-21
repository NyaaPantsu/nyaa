package cookies

import "github.com/NyaaPantsu/nyaa/config"

func getDomainName() string {
	domain := config.Get().Cookies.DomainName
	if config.Get().Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}

func getMaxAge() int {
	return config.Get().Cookies.MaxAge
}
