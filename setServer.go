package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sideproject/config"
	login "sideproject/handlers"
	"sideproject/utils"

	"github.com/gin-gonic/gin"
)

// User represents the user data we'll store in session
type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func Setup() *gin.Engine {
	config.Init()

	r := gin.Default()

	//r.Use(middleware.AuthMiddleware())
	r.GET("/", login.GoogleForm)
	r.GET("/auth/google/login", login.GoogleLoginHandler)
	r.GET("/auth/google/callback", login.GoogleAuthCallback)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	// Facebook OAuth 라우트
	r.GET("/auth/facebook/login", func(c *gin.Context) {
		loginURL := getFacebookLoginURL()
		c.Redirect(http.StatusFound, loginURL)
	})

	r.GET("/auth/facebook/callback", func(c *gin.Context) {
		code := c.DefaultQuery("code", "")
		if code == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		accessToken, err := utils.ExchangeCodeForToken(code)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token"})
			return
		}

		facebookID, err := utils.GetFacebookUserID(accessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user ID"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"FacebookID": facebookID})
	})

	return r
}

// Facebook 로그인 URL 생성
func getFacebookLoginURL() string {
	clientID := os.Getenv("FACEBOOK_CLIENT_ID")
	redirectURI := os.Getenv("FACEBOOK_REDIRECT_URI")
	if clientID == "" || redirectURI == "" {
		log.Fatal("FACEBOOK_CLIENT_ID or FACEBOOK_REDIRECT_URI is not set.")
	}
	return fmt.Sprintf("https://www.facebook.com/v14.0/dialog/oauth?client_id=%s&redirect_uri=%s&scope=email,public_profile,user_birthday,user_friends", clientID, redirectURI)
}
