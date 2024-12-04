package jwt

import (
	"crypto/rsa"
	_ "embed"
	"fmt"
	"log"
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

var (
	privKey *rsa.PrivateKey
	pubKey  *rsa.PublicKey
)

// RSA 키 초기화
func InitializeKeys() error {
	sKey, err := jwt.ParseRSAPrivateKeyFromPEM(secretKey)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}
	privKey = sKey

	pKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		return fmt.Errorf("failed to parse public key: %v", err)
	}
	pubKey = pKey

	return nil
}

// JWT 토큰에 담을 정보
type Claims struct {
	Accounts []models.User `json:"accounts"`
	jwt.RegisteredClaims
}

// JWT 토큰 생성
func GenerateToken(c *gin.Context) (string, error) {
	if privKey == nil {
		return "", fmt.Errorf("secret key not found")
	}

	//context에서 accounts를 가져와서 JWT 토큰에 담기
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

	//JWT 토큰 생성
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privKey)
}

// JWT 토큰 검증
func ValidateToken(tokenString string) (*Claims, error) {
	if pubKey == nil {
		return nil, fmt.Errorf("public key not found")
	}

	//JWT 복호화 및 파싱
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return pubKey, nil
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

// JWT 토큰을 context에 채워넣기
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

// Context에서 accounts를 가져오기
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

// Context에 accounts를 추가
func SetAccount(c *gin.Context, user *models.User) (*gin.Context, error) {
	if user == nil {
		return c, fmt.Errorf("user is required")
	}

	//기존 context에 accounts가 있는지 확인
	//있으면 기존 accounts에 새로운 account를 추가
	//없으면 새로운 accounts를 생성해서 추가
	accounts, exist := c.Get("accounts")
	if !exist {
		res := []models.User{*user}
		c.Set("accounts", res)
		log.Printf("failed to get accounts")
		return c, nil
	}

	app := accounts.([]models.User)
	res := append(app, *user)
	c.Set("accounts", res)
	return c, nil
}
