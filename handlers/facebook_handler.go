package handlers

import (
	"net/http"
	"os"
	"sideproject/utils"

	"github.com/gin-gonic/gin"
)

// FacebookLoginHandler: Facebook 로그인 URL로 리디렉션
func FacebookLoginHandler(c *gin.Context) {
	clientID := os.Getenv("FACEBOOK_CLIENT_ID")
	redirectURI := os.Getenv("FACEBOOK_REDIRECT_URI")
	loginURL := "https://www.facebook.com/v14.0/dialog/oauth?client_id=" + clientID + "&redirect_uri=" + redirectURI + "&scope=email,public_profile"
	c.Redirect(http.StatusFound, loginURL)
}

// FacebookCallbackHandler: Facebook OAuth 후 사용자 정보 가져오기
func FacebookCallbackHandler(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing code"})
		return
	}

	// 액세스 토큰 요청
	accessToken, err := utils.ExchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 사용자 프로필 요청
	profile, err := utils.GetFacebookUserProfile(accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 필요한 필드만 반환
	c.JSON(http.StatusOK, gin.H{
		"id":   profile["id"],
		"name": profile["name"],
		"link": profile["link"],
	})
}
