package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type Config struct {
	GoogleLoginConfig oauth2.Config
	KakaoLoginConfig  KakaoConfig
}

type KakaoConfig struct {
	RestApiKey  string
	RedirectURL string
}

var (
	AppConfig *Config = &Config{}
)

func GoogleConfig() oauth2.Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	googleLoginConfig := oauth2.Config{
		RedirectURL:  "http://localhost:5000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint: google.Endpoint,
	}

	return googleLoginConfig
}

func KakaoConfigInit() KakaoConfig {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	kakaoLoginConfig := KakaoConfig{
		RestApiKey:  os.Getenv("KAKAO_REST_API_KEY"),
		RedirectURL: "http://localhost:5000/auth/kakao/callback",
	}

	return kakaoLoginConfig
}

func Init() {
	googleLoginConfig := GoogleConfig()

	AppConfig = &Config{
		GoogleLoginConfig: googleLoginConfig,
	}

}
