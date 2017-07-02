package cookies

func getDomainName() string {
	domain := config.Conf.Cookies.DomainName
	if config.Conf.Environment == "DEVELOPMENT" {
		domain = ""
	}
	return domain
}

func getMaxAge() int {
	return config.Conf.Cookies.MaxAge
}
