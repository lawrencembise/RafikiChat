package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetUser(c *gin.Context) {
	// Implement your logic to get user
	c.JSON(http.StatusOK, gin.H{"user": "user data"})
}

func CreateUser(c *gin.Context) {
	// Implement your logic to create user
	c.JSON(http.StatusOK, gin.H{"status": "user created"})
}