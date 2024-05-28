package controllers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"rafikichat/internal/infrastructure/config"
)

var cfg config.Config

func SetConfig(configuration config.Config) {
	cfg = configuration
}

func InitiateOAuth(c *gin.Context) {
	telegramOAuthURL := fmt.Sprintf("https://telegram.org/js/telegram-widget.js?2&data-telegram-login=%s&data-size=large&data-auth-url=%s",
		cfg.TelegramBotUsername, cfg.TelegramRedirectURI)
	c.Redirect(http.StatusTemporaryRedirect, telegramOAuthURL)
}

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
	hash := base64.StdEncoding.EncodeToString(h.Sum(nil))

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