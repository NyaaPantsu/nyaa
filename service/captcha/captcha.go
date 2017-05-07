package captcha

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/dchest/captcha"
)

const lifetime = time.Minute * 20

var (
	server   = captcha.Server(captcha.StdWidth, captcha.StdHeight)
	captchas = captchaMap{
		m: make(map[string]store, 64),
	}

	ErrInvalidCaptcha = errors.New("invalid captcha")
)

func init() {
	captcha.SetCustomStore(captcha.NewMemoryStore(1<<10, lifetime))

	go func() {
		t := time.Tick(time.Minute)
		for {
			<-t
			captchas.cleanUp()
		}
	}()
}

// Captcha is to be embedded into any form struct requiring a captcha
type Captcha struct {
	CaptchaID, Solution string
}

// Captchas are IP-specific and need eventual cleanup
type captchaMap struct {
	sync.Mutex
	m map[string]store
}

type store struct {
	id      string
	created time.Time
}

// Returns a captcha id by IP. If a captcha for this IP already exists, it is
// reloaded and returned. Otherwise, a new captcha is created.
func (n *captchaMap) get(ip string) string {
	n.Lock()
	defer n.Unlock()

	old, ok := n.m[ip]

	// No existing captcha, it expired or this IP already used the captcha
	if !ok || !captcha.Reload(old.id) {
		id := captcha.New()
		n.m[ip] = store{
			id:      id,
			created: time.Now(),
		}
		return id
	}

	old.created = time.Now()
	n.m[ip] = old
	return old.id
}

// Remove expired ip -> captchaID mappings
func (n *captchaMap) cleanUp() {
	n.Lock()
	defer n.Unlock()

	till := time.Now().Add(-lifetime)
	for ip, c := range n.m {
		if c.created.Before(till) {
			delete(n.m, ip)
		}
	}
}

// GetID returns a new or previous captcha id by IP
func GetID(ip string) string {
	return captchas.get(ip)
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
