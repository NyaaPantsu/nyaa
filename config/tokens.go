package config

// EmailTokenHashKey : /!\ Email hash for generating email activation token /!\
var EmailTokenHashKey = []byte("CHANGE_THIS_BEFORE_DEPLOYING_YOU_GIT")

// CSRFTokenHashKey : /!\ CSRF hash for generating CSRF tokens /!\
var CSRFTokenHashKey = []byte("CHANGE_THIS_BEFORE_DEPLOYING_YOU_GIT")
