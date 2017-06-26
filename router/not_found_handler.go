package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NotFoundHandler : Controller for displaying 404 error page
func NotFoundHandler(c *gin.Context) {
	httpError(c, http.StatusNotFound)
}
