package http

import (
	"github.com/gin-gonic/gin"
	"rafikichat/internal/adapters/controllers"
)

// Routers is the function for declaring routes with controllers
func Routers(app *gin.Engine) {
	app.GET("/", controllers.ShowLoginPage)
	app.GET("/oauth/callback", controllers.HandleOAuthCallback)
	app.GET("/welcome", controllers.Welcome)
	app.GET("/logout", controllers.Logout)
}