package controllers

import (
	"github.com/NyaaPantsu/nyaa/models"
	"github.com/NyaaPantsu/nyaa/utils/cookies"
	"github.com/gin-gonic/gin"
)

func getUser(c *gin.Context) *models.User {
	user, _, _ := cookies.CurrentUser(c)
	if user == nil {
		return &models.User{}
	}
	return user
}
