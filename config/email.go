package config

import "time"

// TODO: Perform email configuration at runtime
//       Future hosts shouldn't have to rebuild the binary to update a setting

const (
	// SendEmail : Enable Email
	SendEmail = true
	// EmailFrom : email address by default
	EmailFrom = "donotrespond@nyaa.pantsu.cat"
	// EmailTestTo : when testing to who send email
	EmailTestTo = ""
	// EmailHost : Host of mail server
	EmailHost = "localhost"
	// EmailUsername : Username needed for the connection
	EmailUsername = ""
	// EmailPassword : Password needed for the connection
	EmailPassword = ""
	// EmailPort : Mail Server port
	EmailPort = 465
	// EmailTimeout : Timeout for waiting server response
	EmailTimeout = 10 * time.Second
)