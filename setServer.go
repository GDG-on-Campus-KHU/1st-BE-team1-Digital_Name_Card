package main

import (
	"net/http"
	"sideproject/config"
	login "sideproject/handlers"

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

	return r
}
