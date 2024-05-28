package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"rafikichat/internal/infrastructure/config"
)

var cfg config.Config

// SetConfig is the function for setting up configuration variables
func SetConfig(configuration config.Config) {
	cfg = configuration
}
// ShowLoginPage show login page
func ShowLoginPage(c *gin.Context) {
	loginPageHTML := `
    <!DOCTYPE html>
    <html>
    <head>
        <meta charset="utf-8">
        <title>Login with Telegram</title>
    </head>
    <body>
        <h1>Login with Telegram</h1>
        <script async src="https://telegram.org/js/telegram-widget.js?2" data-telegram-login="%s" data-size="large" data-auth-url="/oauth/callback" data-request-access="write"></script>
    </body>
    </html>`
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, fmt.Sprintf(loginPageHTML, cfg.TelegramBotUsername))
}


// HandleOAuthCallback is the controller for handling Telegram OAuth
func HandleOAuthCallback(c *gin.Context) {
	authData := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		authData[key] = values[0]
	}

	if err := checkTelegramAuthorization(authData); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	saveTelegramUserData(c, authData)
	c.Redirect(http.StatusTemporaryRedirect, "/welcome")
}

// checkTelegramAuthorization is the function for checking telegram authorization
func checkTelegramAuthorization(authData map[string]string) error {
	checkHash := authData["hash"]
	delete(authData, "hash")

	var dataCheckArr []string
	for key, value := range authData {
		dataCheckArr = append(dataCheckArr, fmt.Sprintf("%s=%s", key, value))
	}

	sort.Strings(dataCheckArr)
	dataCheckString := string([]byte(fmt.Sprintf("%s\n", dataCheckArr)))

	secretKey := sha256.New()
	secretKey.Write([]byte(cfg.TelegramBotToken))

	h := hmac.New(sha256.New, secretKey.Sum(nil))
	h.Write([]byte(dataCheckString))
	hash := hex.EncodeToString(h.Sum(nil))

	fmt.Printf("Calculated hash: %s\n", hash)
	fmt.Printf("Check hash: %s\n", checkHash)

	if hash != checkHash {
		return fmt.Errorf("data is NOT from Telegram")
	}

	authDate, err := time.Parse("2006-01-02T15:04:05Z", authData["auth_date"])
	if err != nil {
		return err
	}

	if time.Since(authDate) > 24*time.Hour {
		return fmt.Errorf("data is outdated")
	}

	return nil
}

// saveTelegramUserData is the function that is called to save user data after authorization
func saveTelegramUserData(c *gin.Context, authData map[string]string) {
	authDataJSON, _ := json.Marshal(authData)
	c.SetCookie("tg_user", string(authDataJSON), 86400, "/", "", false, true)
}

// Welcome is the controller for showing welcome page after authorization
func Welcome(c *gin.Context) {
	tgUser, err := getTelegramUserData(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	html := fmt.Sprintf("<h1>Hello, %s %s!</h1>", tgUser["first_name"], tgUser["last_name"])
	if username, ok := tgUser["username"]; ok {
		html = fmt.Sprintf("<h1>Hello, <a href=\"https://t.me/%s\">%s %s</a>!</h1>", username, tgUser["first_name"], tgUser["last_name"])
	}
	if photoURL, ok := tgUser["photo_url"]; ok {
		html += fmt.Sprintf("<img src=\"%s\">", photoURL)
	}
	html += "<p><a href=\"/logout\">Log out</a></p>"

	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// getTelegramUserData is the function for getting user data from telegram
func getTelegramUserData(c *gin.Context) (map[string]string, error) {
	tgUserCookie, err := c.Cookie("tg_user")
	if err != nil {
		return nil, err
	}

	var tgUser map[string]string
	if err := json.Unmarshal([]byte(tgUserCookie), &tgUser); err != nil {
		return nil, err
	}

	return tgUser, nil
}

// Logout is the controller for user logout
func Logout(c *gin.Context) {
	c.SetCookie("tg_user", "", -1, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}