package faqController

import "github.com/NyaaPantsu/nyaa/controllers/router"

func init() {
	router.Get().Any("/faq", FaqHandler)
}
