package config

// TODO: Perform environment configuration at runtime
//       Future hosts shouldn't have to rebuild the binary to update a setting

const (
	// Environment should be one of: DEVELOPMENT, TEST, PRODUCTION
	Environment = "DEVELOPMENT"
	// WebAddress : url of the website
	WebAddress = "nyaa.pantsu.cat"
	// AuthTokenExpirationDay : Number of Days for token expiration when logged in
	AuthTokenExpirationDay = 1000
)
