package config

import "time"

const (
	SendEmail     = true
	EmailFrom     = "donotrespond@nyaa.pantsu.cat"
	EmailTestTo   = ""
	EmailHost     = "localhost"
	EmailUsername = ""
	EmailPassword = ""
	EmailPort     = 465
	// EmailTimeout  = 80 * time.Millisecond
	EmailTimeout = 10 * time.Second
)

var EmailTokenHashKey = []byte("CHANGE_THIS_BEFORE_DEPLOYING_YOU_RETARD")
