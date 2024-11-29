package jwt

import (
	"fmt"
	"net/http"
	"sideproject/models"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

//go:embed cert/secret.pem
var secretKey []byte

//go:embed cert/public.pem
var publicKey []byte

type Claims struct {
	Accounts []models.User `json:"accounts"`
	jwt.RegisteredClaims
}

func GenerateToken(c *gin.Context) (string, error) {
	if secretKey == nil {
		return "", fmt.Errorf("secret key not found")
	}

	accounts, err := GetAccount(c)
	if err != nil {
		return "", err
	}

	claims := &Claims{
		Accounts: accounts,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "sideproject",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(secretKey)
}

func ValidateToken(tokenString string) (*Claims, error) {
	if publicKey == nil {
		return nil, fmt.Errorf("public key not found")
	}

	//JWT 복호화 및 파싱
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	//JWT 토큰이 유효한지 확인
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func FillContext(c *gin.Context) (*gin.Context, error) {
	//http request의 header에서 JWT 파싱
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
		c.Abort()
		return nil, fmt.Errorf("authorization header is required")
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
		c.Abort()
		return nil, fmt.Errorf("invalid authorization header format")
	}

	//JWT 토큰 검증
	claims, err := ValidateToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		c.Abort()
		return nil, fmt.Errorf("invalid token")
	}

	//JWT 토큰에서 accounts를 추출해서 context에 추가
	for _, account := range claims.Accounts {
		SetAccount(c, &account)
	}

	return c, nil
}

func GetAccount(c *gin.Context) ([]models.User, error) {
	//context에서 accounts를 찾아서 반환
	//accounts가 없으면 에러 반환
	result, err := c.Get("accounts")
	if !err {
		return nil, fmt.Errorf("accounts not found")
	}

	users, ok := result.([]models.User)
	if !ok {
		return nil, fmt.Errorf("accounts is not a user")
	}

	return users, nil
}

func SetAccount(c *gin.Context, user *models.User) (*gin.Context, error) {
	if user == nil {
		return c, fmt.Errorf("user is required")
	}

	//기존 context에 accounts가 있는지 확인
	//있으면 기존 accounts에 새로운 account를 추가
	//없으면 새로운 accounts를 생성해서 추가
	accounts, err := GetAccount(c)
	if err != nil {
		res := []models.User{*user}
		c.Set("accounts", res)
		return c, fmt.Errorf("failed to get accounts")
	}

	res := append(accounts, *user)
	c.Set("accounts", res)
	return c, nil
}
