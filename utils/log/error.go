package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

// CheckError check error and return true if error is nil and return false if error is not nil.
func CheckError(err error) bool {
	return CheckErrorWithMessage(err, "")
}

// CheckErrorWithMessage check error with message and log messages with stack. And then return true if error is nil and return false if error is not nil.
func CheckErrorWithMessage(err error, msg string, args ...interface{}) bool {
	if err != nil {
		var stack [4096]byte
		runtime.Stack(stack[:], false)
		if len(args) == 0 {
			logrus.Error(msg + fmt.Sprintf("%q\n%s\n", err, stack[:]))
		} else {
			logrus.Error(fmt.Sprintf(msg, args...) + fmt.Sprintf("%q\n%s\n", err, stack[:]))
		}
		return false
	}
	return true
}
