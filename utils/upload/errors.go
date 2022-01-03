package upload

import "errors"

var errPending = errors.New("Magnet being generated already")
var errConfig = errors.New("Magnet Empty or FileStorage not configured")
