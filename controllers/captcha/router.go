package captchaController

import (
	"github.com/NyaaPantsu/nyaa/controllers/router"
	"github.com/NyaaPantsu/nyaa/utils/captcha"
)

func init() {
	router.Get().Any("/captcha/*hash", captcha.ServeFiles)
}
