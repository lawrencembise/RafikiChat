package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// HomeHandler is the function for handling home page
func HomeHandler(c *gin.Context) {
	// Implement your logic for home page
	c.JSON(http.StatusOK, gin.H{"status": "home page"})
}