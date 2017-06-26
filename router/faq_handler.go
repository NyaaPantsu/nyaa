package router

import (
	"github.com/gin-gonic/gin"
)

// FaqHandler : Controller for FAQ view page
func FaqHandler(c *gin.Context) {
	staticTemplate(c, "faq")
}
