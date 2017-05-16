package config

import "time"

// TODO: Perform email configuration at runtime
//       Future hosts shouldn't have to rebuild the binary to update a setting

const (
	SendEmail     = true
	EmailFrom     = "donotrespond@nyaa.pantsu.cat"
	EmailTestTo   = ""
	EmailHost     = "localhost"
	EmailUsername = ""
	EmailPassword = ""
	EmailPort     = 465
	EmailTimeout  = 10 * time.Second
)

var EmailTokenHashKey = []byte("CHANGE_THIS_BEFORE_DEPLOYING_YOU_GIT")
