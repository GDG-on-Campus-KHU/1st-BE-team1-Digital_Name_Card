package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"os"
)

// ExchangeCodeForToken: 로그인 코드로 액세스 토큰을 가져옵니다.
func ExchangeCodeForToken(code string) (string, error) {
	clientID := os.Getenv("FACEBOOK_CLIENT_ID")
	clientSecret := os.Getenv("FACEBOOK_CLIENT_SECRET")
	redirectURI := os.Getenv("FACEBOOK_REDIRECT_URI")

	tokenURL := "https://graph.facebook.com/v14.0/oauth/access_token"
	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)
	params.Add("redirect_uri", redirectURI)
	params.Add("code", code)

	resp, err := http.Get(tokenURL + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	if tokenResponse.AccessToken == "" {
		return "", errors.New("failed to get access token")
	}
	return tokenResponse.AccessToken, nil
}

// GetFacebookUserID: 액세스 토큰을 사용해 Facebook ID를 가져옵니다.
func GetFacebookUserID(accessToken string) (string, error) {
	resp, err := http.Get("https://graph.facebook.com/me?access_token=" + accessToken)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var userResponse struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userResponse); err != nil {
		return "", err
	}

	if userResponse.ID == "" {
		return "", errors.New("failed to get user ID")
	}
	return userResponse.ID, nil
}

// GetFacebookUserProfile: Facebook 프로필 가져오기
func GetFacebookUserProfile(accessToken string) (map[string]interface{}, error) {
	url := "https://graph.facebook.com/me?fields=id,name,link&access_token=" + accessToken

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}

	if profile["id"] == nil || profile["name"] == nil || profile["link"] == nil {
		return nil, errors.New("failed to fetch required user profile fields")
	}

	return profile, nil
}
