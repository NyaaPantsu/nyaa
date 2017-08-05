package captcha

import (
	"errors"
	"time"

	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
)

const lifetime = time.Minute * 20

var (
	server = captcha.Server(captcha.StdWidth, captcha.StdHeight)
	// ErrInvalidCaptcha : Error when captcha is invalid
	ErrInvalidCaptcha = errors.New("invalid captcha")
)

func init() {
	captcha.SetCustomStore(captcha.NewMemoryStore(1<<10, lifetime))
}

// Captcha is to be embedded into any form struct requiring a Captcha
type Captcha struct {
	CaptchaID, Solution string
}

// GetID returns a new Captcha ID
func GetID() string {
	return captcha.New()
}

// Extract a Captcha struct from an HTML form
func Extract(c *gin.Context) Captcha {
	return Captcha{
		CaptchaID: c.PostForm("captchaID"),
		Solution:  c.PostForm("solution"),
	}
}

// ServeFiles serves Captcha images and audio
func ServeFiles(c *gin.Context) {
	server.ServeHTTP(c.Writer, c.Request)
}

// Authenticate check's if a Captcha solution is valid
func Authenticate(req Captcha) bool {
	return captcha.VerifyString(req.CaptchaID, req.Solution)
}
