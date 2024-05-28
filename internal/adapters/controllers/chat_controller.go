package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetChat(c *gin.Context) {
	// Implement your logic to get chat
	c.JSON(http.StatusOK, gin.H{"chat": "chat data"})
}

func CreateChat(c *gin.Context) {
	// Implement your logic to create chat
	c.JSON(http.StatusOK, gin.H{"status": "chat created"})
}