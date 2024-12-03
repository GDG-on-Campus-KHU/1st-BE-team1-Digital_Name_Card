package main

import (
	"log"
	"sideproject/jwt"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {
	// 환경 변수 로드
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if err := jwt.InitializeKeys(); err != nil {
		log.Fatalf("Failed to initialize RSA keys: %v", err)
	}

	r := Setup()

	// CORS 설정 (모든 도메인 허용)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 정적 파일 서빙
	r.Static("/static", "./static")

	r.Run(":8080")
}
