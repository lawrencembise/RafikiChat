package http

import (
	"github.com/gin-gonic/gin"
	"rafikichat/internal/adapters/controllers"
)

func Routers(app *gin.Engine) {
	app.GET("/oauth/telegram", controllers.InitiateOAuth)
	app.GET("/oauth/callback", controllers.HandleOAuthCallback)
	app.GET("/welcome", controllers.Welcome)
	app.GET("/logout", controllers.Logout)
}