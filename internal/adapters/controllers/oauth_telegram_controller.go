package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"rafikichat/internal/infrastructure/config"
)

var cfg config.Config

func SetConfig(configuration config.Config) {
	cfg = configuration
}

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

func HandleOAuthCallback(c *gin.Context) {
	authData := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		authData[key] = values[0]
	}

	fmt.Printf("Received auth data: %+v\n", authData)

	if err := checkTelegramAuthorization(authData); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	saveTelegramUserData(c, authData)
	c.Redirect(http.StatusTemporaryRedirect, "/welcome")
}

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

	fmt.Printf("Calculated hash: %s\nExpected hash: %s\n", hash, checkHash)

	if hash != checkHash {
		return fmt.Errorf("data is NOT from Telegram")
	}

	authDate, err := strconv.ParseInt(authData["auth_date"], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid auth_date")
	}

	if time.Since(time.Unix(authDate, 0)) > 24*time.Hour {
		return fmt.Errorf("data is outdated")
	}

	return nil
}

func saveTelegramUserData(c *gin.Context, authData map[string]string) {
	authDataJSON, _ := json.Marshal(authData)
	c.SetCookie("tg_user", string(authDataJSON), 86400, "/", "", false, true)
}

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

func Logout(c *gin.Context) {
	c.SetCookie("tg_user", "", -1, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/")
}
