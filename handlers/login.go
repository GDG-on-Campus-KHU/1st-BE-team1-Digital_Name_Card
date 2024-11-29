package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sideproject/config"
	"sideproject/models"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func GoogleForm(c *gin.Context) {
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
		"<html>"+
			"\n<head>\n    "+
			"<title>Go Oauth2.0 Test</title>\n"+
			"</head>\n"+
			"<body>\n<p>"+
			"<a href='./auth/google/login'>Google Login</a>"+
			"</p>\n"+
			"</body>\n"+
			"</html>"))
}

func GenerateStateOauthCookie(w http.ResponseWriter) string {
	expiration := time.Now().Add(1 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := &http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, cookie)
	return state
}

func GoogleLoginHandler(c *gin.Context) {

	state := GenerateStateOauthCookie(c.Writer)
	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleAuthCallback(c *gin.Context) {
	oauthstate, _ := c.Request.Cookie("oauthstate")

	if c.Request.FormValue("state") != oauthstate.Value {
		log.Printf("invalid google oauth state cookie:%s state:%s\n", oauthstate.Value, c.Request.FormValue("state"))
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := GetGoogleUserInfo(c.Request.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	log.Printf("User Info: %s\n", data)

	user := models.User{
		Uid:      uuid.NewString(),
		Email:    data["email"].(string),
		Nickname: data["name"].(string),
	}

	c.JSON(http.StatusOK, gin.H{
		"uid":      user.Uid,
		"email":    user.Email,
		"nickname": user.Nickname,
	})

	c.Redirect(http.StatusFound, "/profile")
}

func GetGoogleUserInfo(code string) (map[string]interface{}, error) {
	const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	token, err := config.AppConfig.GoogleLoginConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("Failed to Exchange %s\n", err.Error())
	}

	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("Failed to Get UserInfo %s\n", err.Error())
	}

	var data map[string]interface{}
	merr := json.NewDecoder(resp.Body).Decode(&data)
	if merr != nil {
		return nil, fmt.Errorf("Failed to decode JSON:", merr.Error())
	}

	defer resp.Body.Close()

	return data, err
}
