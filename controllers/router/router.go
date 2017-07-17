package router

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var once sync.Once

// Get return a router signleton
func Get() *gin.Engine {
	once.Do(func() {
		router = gin.New()
		router.Use(gin.Logger())
		router.Use(gin.Recovery())
	})
	return router
}
