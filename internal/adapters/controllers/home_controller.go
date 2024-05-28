package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HomeHandler(c *gin.Context) {
	// Implement your logic for home page
	c.JSON(http.StatusOK, gin.H{"status": "home page"})
}