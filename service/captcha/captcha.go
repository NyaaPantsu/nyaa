package captcha

import (
	"errors"
	"net/http"
	"time"

	"github.com/dchest/captcha"
)

const lifetime = time.Minute * 20

var (
	server            = captcha.Server(captcha.StdWidth, captcha.StdHeight)
	ErrInvalidCaptcha = errors.New("invalid captcha")
)

func init() {
	captcha.SetCustomStore(captcha.NewMemoryStore(1<<10, lifetime))
}

// Captcha is to be embedded into any form struct requiring a captcha
type Captcha struct {
	CaptchaID, Solution string
}

// GetID returns a new captcha id
func GetID() string {
	return captcha.New()
}

// Extract a Captcha struct from an HTML form
func Extract(r *http.Request) Captcha {
	return Captcha{
		CaptchaID: r.FormValue("captchaID"),
		Solution:  r.FormValue("solution"),
	}
}

// ServeFiles serves captcha images and audio
func ServeFiles(w http.ResponseWriter, r *http.Request) {
	server.ServeHTTP(w, r)
}

// Authenticate check's if a captcha solution is valid
func Authenticate(req Captcha) bool {
	return captcha.VerifyString(req.CaptchaID, req.Solution)
}
